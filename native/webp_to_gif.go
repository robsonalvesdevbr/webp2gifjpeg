package native

/*
#cgo pkg-config: libwebp libwebpdemux
#cgo LDFLAGS: -lgif
#include <stdlib.h>
#include <string.h>
#include <webp/decode.h>
#include <webp/demux.h>
#include <webp/mux_types.h>
#include <gif_lib.h>
*/
import "C"
import (
	"fmt"
	"os"
	"unsafe"
)

// ConvertWebPToGIF converts an animated WebP file to GIF format
func ConvertWebPToGIF(inputPath, outputPath string) error {
	// Read WebP file
	data, err := os.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("failed to read WebP file: %w", err)
	}

	if len(data) == 0 {
		return fmt.Errorf("WebP file is empty")
	}

	// Allocate C memory and copy data to avoid CGO pointer issues
	cData := C.malloc(C.size_t(len(data)))
	if cData == nil {
		return fmt.Errorf("failed to allocate memory")
	}
	defer C.free(cData)

	C.memcpy(cData, unsafe.Pointer(&data[0]), C.size_t(len(data)))

	// Create WebPData structure
	webpData := C.WebPData{
		bytes: (*C.uint8_t)(cData),
		size:  C.size_t(len(data)),
	}

	// Create demuxer
	demux := C.WebPDemux(&webpData)
	if demux == nil {
		return fmt.Errorf("failed to create WebP demuxer")
	}
	defer C.WebPDemuxDelete(demux)

	// Get canvas info
	width := int(C.WebPDemuxGetI(demux, C.WEBP_FF_CANVAS_WIDTH))
	height := int(C.WebPDemuxGetI(demux, C.WEBP_FF_CANVAS_HEIGHT))
	frameCount := int(C.WebPDemuxGetI(demux, C.WEBP_FF_FRAME_COUNT))

	if frameCount == 0 {
		return fmt.Errorf("no frames found in WebP file")
	}

	// Open GIF file for writing
	cFilename := C.CString(outputPath)
	defer C.free(unsafe.Pointer(cFilename))

	var errCode C.int
	gifFile := C.EGifOpenFileName(cFilename, C.bool(false), &errCode)
	if gifFile == nil {
		return fmt.Errorf("failed to create GIF file: error code %d", errCode)
	}
	defer C.EGifCloseFile(gifFile, &errCode)

	// Set GIF screen descriptor WITHOUT global color map (use local per frame)
	// This allows each frame to have its own optimized 256-color palette
	if C.EGifPutScreenDesc(gifFile, C.int(width), C.int(height), 8, 0, nil) == C.GIF_ERROR {
		return fmt.Errorf("failed to write GIF screen descriptor")
	}

	// Add Netscape 2.0 extension for looping
	if err := addLoopingExtension(gifFile); err != nil {
		return fmt.Errorf("failed to add looping extension: %w", err)
	}

	// Iterate through frames
	var iter C.WebPIterator
	if C.WebPDemuxGetFrame(demux, 1, &iter) == 0 {
		return fmt.Errorf("failed to get first frame")
	}
	defer C.WebPDemuxReleaseIterator(&iter)

	for {
		// Decode frame with proper blending/compositing
		// Use WebPDecodeRGBA on the full fragment to get properly composited frame
		fragmentSize := int(iter.fragment.size)

		var outWidth, outHeight C.int
		rgbaData := C.WebPDecodeRGBA(
			iter.fragment.bytes,
			C.size_t(fragmentSize),
			&outWidth,
			&outHeight,
		)

		if rgbaData == nil {
			return fmt.Errorf("failed to decode frame %d", iter.frame_num)
		}
		defer C.WebPFree(unsafe.Pointer(rgbaData))

		// Convert RGBA to RGB pixels
		frameWidth := int(outWidth)
		frameHeight := int(outHeight)
		pixelCount := frameWidth * frameHeight
		rgbPixels := make([]RGB, pixelCount)

		rgbaSlice := unsafe.Slice((*byte)(rgbaData), pixelCount*4)
		for i := 0; i < pixelCount; i++ {
			rgbPixels[i] = RGB{
				R: rgbaSlice[i*4],
				G: rgbaSlice[i*4+1],
				B: rgbaSlice[i*4+2],
				// Skip alpha channel (i*4+3)
			}
		}

		// Quantize frame to 256 colors using Octree (like Pillow does)
		// Use simple Octree without dithering to match Python/Pillow behavior
		indexedData, framePalette := QuantizeImageOctreeWithDimensions(rgbPixels, 256, frameWidth, frameHeight)

		// Create local color map for this frame
		localColorMap := C.GifMakeMapObject(256, nil)
		if localColorMap == nil {
			return fmt.Errorf("failed to create local color map for frame %d", iter.frame_num)
		}

		// Copy frame palette to local color map
		localColors := unsafe.Slice(localColorMap.Colors, len(framePalette))
		for i := 0; i < len(framePalette); i++ {
			localColors[i].Red = C.GifByteType(framePalette[i].R)
			localColors[i].Green = C.GifByteType(framePalette[i].G)
			localColors[i].Blue = C.GifByteType(framePalette[i].B)
		}

		// Add graphics control extension (for timing)
		duration := int(iter.duration) / 10 // Convert ms to centiseconds
		if duration < 1 {
			duration = 10 // Default 100ms
		}

		var gce [4]C.GifByteType
		// Disposal method: 0 = unspecified (let decoder decide, like Pillow)
		gce[0] = 0x00 // No disposal method specified
		gce[1] = C.GifByteType(duration & 0xff)
		gce[2] = C.GifByteType((duration >> 8) & 0xff)
		gce[3] = 0 // No transparent color

		if C.EGifPutExtension(gifFile, C.GRAPHICS_EXT_FUNC_CODE, 4, unsafe.Pointer(&gce[0])) == C.GIF_ERROR {
			C.GifFreeMapObject(localColorMap)
			return fmt.Errorf("failed to write graphics control extension")
		}

		// Write frame WITH local color map
		if C.EGifPutImageDesc(gifFile, 0, 0, C.int(frameWidth), C.int(frameHeight), C.bool(false), localColorMap) == C.GIF_ERROR {
			C.GifFreeMapObject(localColorMap)
			return fmt.Errorf("failed to write image descriptor for frame %d", iter.frame_num)
		}

		// Write scanlines
		for y := 0; y < frameHeight; y++ {
			line := (*C.GifByteType)(unsafe.Pointer(&indexedData[y*frameWidth]))
			if C.EGifPutLine(gifFile, line, C.int(frameWidth)) == C.GIF_ERROR {
				return fmt.Errorf("failed to write scanline %d in frame %d", y, iter.frame_num)
			}
		}

		// Move to next frame
		if C.WebPDemuxNextFrame(&iter) == 0 {
			break
		}
	}

	return nil
}

// addLoopingExtension adds Netscape 2.0 extension for infinite looping
func addLoopingExtension(gifFile *C.GifFileType) error {
	// Netscape 2.0 application extension
	appExt := []byte("NETSCAPE2.0")
	if C.EGifPutExtensionLeader(gifFile, C.APPLICATION_EXT_FUNC_CODE) == C.GIF_ERROR {
		return fmt.Errorf("failed to write extension leader")
	}

	if C.EGifPutExtensionBlock(gifFile, C.int(len(appExt)), unsafe.Pointer(&appExt[0])) == C.GIF_ERROR {
		return fmt.Errorf("failed to write application extension")
	}

	// Loop count sub-block (0 = infinite)
	loopBlock := []byte{1, 0, 0} // sub-block id=1, loop count=0 (infinite)
	if C.EGifPutExtensionBlock(gifFile, 3, unsafe.Pointer(&loopBlock[0])) == C.GIF_ERROR {
		return fmt.Errorf("failed to write loop sub-block")
	}

	if C.EGifPutExtensionTrailer(gifFile) == C.GIF_ERROR {
		return fmt.Errorf("failed to write extension trailer")
	}

	return nil
}

// analyzeAllFramesForGlobalPalette analyzes all frames to create a global color palette
// This prevents color flickering between frames in the output GIF
func analyzeAllFramesForGlobalPalette(demux *C.WebPDemuxer, width, height, frameCount int) ([]RGB, error) {
	// Collect colors from all frames
	allColors := make(map[uint32]int) // color -> frequency

	var iter C.WebPIterator
	if C.WebPDemuxGetFrame(demux, 1, &iter) == 0 {
		return nil, fmt.Errorf("failed to get first frame for analysis")
	}
	defer C.WebPDemuxReleaseIterator(&iter)

	frameNum := 0
	for {
		frameNum++

		// Decode frame properly
		fragmentSize := int(iter.fragment.size)

		var outWidth, outHeight C.int
		rgbaData := C.WebPDecodeRGBA(
			iter.fragment.bytes,
			C.size_t(fragmentSize),
			&outWidth,
			&outHeight,
		)

		if rgbaData == nil {
			return nil, fmt.Errorf("failed to decode frame %d during analysis", frameNum)
		}

		// Collect colors from RGBA data
		frameWidth := int(outWidth)
		frameHeight := int(outHeight)
		pixelCount := frameWidth * frameHeight

		rgbaSlice := unsafe.Slice((*byte)(rgbaData), pixelCount*4)
		for i := 0; i < pixelCount; i++ {
			r := uint32(rgbaSlice[i*4])
			g := uint32(rgbaSlice[i*4+1])
			b := uint32(rgbaSlice[i*4+2])
			colorKey := (r << 16) | (g << 8) | b
			allColors[colorKey]++
		}

		C.WebPFree(unsafe.Pointer(rgbaData))

		// Move to next frame
		if C.WebPDemuxNextFrame(&iter) == 0 {
			break
		}
	}

	// Convert color histogram to RGB slice
	colorList := make([]RGB, 0, len(allColors))
	for colorKey := range allColors {
		r := byte(colorKey >> 16)
		g := byte(colorKey >> 8)
		b := byte(colorKey)
		colorList = append(colorList, RGB{R: r, G: g, B: b})
	}

	// Use Octree to quantize to 256 colors
	// Use reasonable dimensions for dithering (not needed for palette generation)
	_, globalPalette := QuantizeImageOctreeWithDimensions(colorList, 256, width, height)

	return globalPalette, nil
}

// mapPixelsToGlobalPalette maps RGB pixels to the global palette without dithering
func mapPixelsToGlobalPalette(pixels []RGB, globalPalette []RGB, width, height int) []byte {
	indexed := make([]byte, len(pixels))

	// Build color lookup cache for performance
	colorCache := make(map[uint32]byte, len(pixels)/4)

	// Direct nearest-color matching without dithering
	for i, p := range pixels {
		// Try cache first
		colorKey := (uint32(p.R) << 16) | (uint32(p.G) << 8) | uint32(p.B)
		paletteIdx, inCache := colorCache[colorKey]

		if !inCache {
			// Find closest color using perceptual distance
			paletteIdx = findClosestColorInPalette(globalPalette, p.R, p.G, p.B)
			colorCache[colorKey] = paletteIdx
		}

		indexed[i] = paletteIdx
	}

	return indexed
}

// clampByte clamps an integer to byte range
func clampByte(val int) byte {
	if val < 0 {
		return 0
	}
	if val > 255 {
		return 255
	}
	return byte(val)
}

// findClosestColorInPalette finds the closest color in palette using weighted Euclidean distance
// Weights are based on human color perception (green is more important)
func findClosestColorInPalette(palette []RGB, r, g, b byte) byte {
	minDist := uint32(0xFFFFFFFF)
	closest := byte(0)

	for i, c := range palette {
		dr := int(r) - int(c.R)
		dg := int(g) - int(c.G)
		db := int(b) - int(c.B)

		// Weighted distance for better perceptual matching
		// Human eye is more sensitive to green, then red, then blue
		dist := uint32(2*dr*dr + 4*dg*dg + 3*db*db)

		if dist < minDist {
			minDist = dist
			closest = byte(i)
		}
		if dist == 0 {
			break // Exact match
		}
	}

	return closest
}

