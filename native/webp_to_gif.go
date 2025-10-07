package native

/*
#cgo pkg-config: libwebp libwebpdemux
#cgo LDFLAGS: -lgif
#include <stdlib.h>
#include <string.h>
#include <webp/decode.h>
#include <webp/demux.h>
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

	// Set GIF screen descriptor
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
		// Decode frame to RGBA
		var frameWidth, frameHeight C.int
		rgba := C.WebPDecodeRGBA(
			iter.fragment.bytes,
			iter.fragment.size,
			&frameWidth,
			&frameHeight,
		)
		if rgba == nil {
			return fmt.Errorf("failed to decode frame %d", iter.frame_num)
		}

		// Convert RGBA to indexed color (GIF format)
		w := int(frameWidth)
		h := int(frameHeight)
		indexedData, colorMap, err := rgbaToIndexed(rgba, w, h)
		if err != nil {
			C.free(unsafe.Pointer(rgba))
			return fmt.Errorf("failed to convert frame %d to indexed: %w", iter.frame_num, err)
		}
		C.free(unsafe.Pointer(rgba))

		// Create GIF color map
		gifColorMap := C.GifMakeMapObject(256, nil)
		if gifColorMap == nil {
			return fmt.Errorf("failed to create color map for frame %d", iter.frame_num)
		}

		// Copy color map
		colors := unsafe.Slice(gifColorMap.Colors, len(colorMap))
		for i := 0; i < len(colorMap); i++ {
			colors[i].Red = C.GifByteType(colorMap[i].R)
			colors[i].Green = C.GifByteType(colorMap[i].G)
			colors[i].Blue = C.GifByteType(colorMap[i].B)
		}

		// Add graphics control extension (for timing)
		duration := int(iter.duration) / 10 // Convert ms to centiseconds
		if duration < 1 {
			duration = 10 // Default 100ms
		}

		var gce [4]C.GifByteType
		gce[0] = 0 // No disposal method, no transparency
		gce[1] = C.GifByteType(duration & 0xff)
		gce[2] = C.GifByteType((duration >> 8) & 0xff)
		gce[3] = 0 // No transparent color

		if C.EGifPutExtension(gifFile, C.GRAPHICS_EXT_FUNC_CODE, 4, unsafe.Pointer(&gce[0])) == C.GIF_ERROR {
			C.GifFreeMapObject(gifColorMap)
			return fmt.Errorf("failed to write graphics control extension")
		}

		// Write frame
		if C.EGifPutImageDesc(gifFile, 0, 0, C.int(w), C.int(h), C.bool(false), gifColorMap) == C.GIF_ERROR {
			C.GifFreeMapObject(gifColorMap)
			return fmt.Errorf("failed to write image descriptor for frame %d", iter.frame_num)
		}

		// Write scanlines
		for y := 0; y < h; y++ {
			line := (*C.GifByteType)(unsafe.Pointer(&indexedData[y*w]))
			if C.EGifPutLine(gifFile, line, C.int(w)) == C.GIF_ERROR {
				C.GifFreeMapObject(gifColorMap)
				return fmt.Errorf("failed to write scanline %d in frame %d", y, iter.frame_num)
			}
		}

		C.GifFreeMapObject(gifColorMap)

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

// RGB represents an RGB color
type RGB struct {
	R, G, B byte
}

// rgbaToIndexed converts RGBA data to indexed color with palette (simple quantization)
func rgbaToIndexed(rgba *C.uint8_t, width, height int) ([]byte, []RGB, error) {
	size := width * height
	indexed := make([]byte, size)
	rgbaSlice := unsafe.Slice(rgba, size*4)

	// Simple color quantization - use 256 color palette
	// This is a basic implementation; for better quality, use a proper quantization algorithm
	colorMap := make(map[uint32]byte)
	palette := make([]RGB, 0, 256)

	for i := 0; i < size; i++ {
		r := rgbaSlice[i*4]
		g := rgbaSlice[i*4+1]
		b := rgbaSlice[i*4+2]
		a := rgbaSlice[i*4+3]

		// Handle transparency by compositing on white
		if a < 255 {
			alpha := float32(a) / 255.0
			invAlpha := 1.0 - alpha
			r = C.uint8_t(float32(r)*alpha + 255*invAlpha)
			g = C.uint8_t(float32(g)*alpha + 255*invAlpha)
			b = C.uint8_t(float32(b)*alpha + 255*invAlpha)
		}

		// Quantize to 6-bit per channel (6x6x6 = 216 colors)
		qr := (uint32(r) * 5 / 255)
		qg := (uint32(g) * 5 / 255)
		qb := (uint32(b) * 5 / 255)

		colorKey := (qr << 16) | (qg << 8) | qb

		idx, exists := colorMap[colorKey]
		if !exists {
			if len(palette) >= 256 {
				// Palette full, use closest existing color
				idx = findClosestColor(palette, byte(r), byte(g), byte(b))
			} else {
				idx = byte(len(palette))
				palette = append(palette, RGB{
					R: byte(qr * 255 / 5),
					G: byte(qg * 255 / 5),
					B: byte(qb * 255 / 5),
				})
				colorMap[colorKey] = idx
			}
		}
		indexed[i] = idx
	}

	// Fill remaining palette entries with black if needed
	for len(palette) < 256 {
		palette = append(palette, RGB{0, 0, 0})
	}

	return indexed, palette, nil
}

// findClosestColor finds the closest color in palette
func findClosestColor(palette []RGB, r, g, b byte) byte {
	minDist := uint32(0xFFFFFFFF)
	closest := byte(0)

	for i, c := range palette {
		dr := int(r) - int(c.R)
		dg := int(g) - int(c.G)
		db := int(b) - int(c.B)
		dist := uint32(dr*dr + dg*dg + db*db)

		if dist < minDist {
			minDist = dist
			closest = byte(i)
		}
	}

	return closest
}
