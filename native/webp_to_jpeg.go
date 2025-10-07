package native

/*
#cgo pkg-config: libwebp
#cgo LDFLAGS: -ljpeg
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <webp/decode.h>
#include <jpeglib.h>
#include <setjmp.h>

// Complete JPEG encoding in C to avoid CGO pointer issues
int encode_jpeg_to_file(const char *filename, unsigned char *rgb_data, int width, int height, int quality) {
	FILE *outfile = fopen(filename, "wb");
	if (!outfile) {
		return -1;
	}

	struct jpeg_compress_struct cinfo;
	struct jpeg_error_mgr jerr;

	cinfo.err = jpeg_std_error(&jerr);
	jpeg_create_compress(&cinfo);
	jpeg_stdio_dest(&cinfo, outfile);

	cinfo.image_width = width;
	cinfo.image_height = height;
	cinfo.input_components = 3;
	cinfo.in_color_space = JCS_RGB;

	jpeg_set_defaults(&cinfo);
	jpeg_set_quality(&cinfo, quality, TRUE);

	// Optimize for quality
	cinfo.dct_method = JDCT_ISLOW;  // Highest quality DCT method
	cinfo.optimize_coding = TRUE;   // Optimize Huffman tables

	// Enable progressive encoding FIRST (before setting chroma)
	// jpeg_simple_progression() modifies the scan script and may reset component info
	jpeg_simple_progression(&cinfo);

	// CRITICAL: Disable chroma subsampling for maximum quality (4:4:4)
	// This MUST be set AFTER jpeg_simple_progression() to prevent reset
	// Default 4:2:0 loses 75% of color resolution - 4:4:4 keeps 100%
	cinfo.comp_info[0].h_samp_factor = 1;
	cinfo.comp_info[0].v_samp_factor = 1;
	cinfo.comp_info[1].h_samp_factor = 1;
	cinfo.comp_info[1].v_samp_factor = 1;
	cinfo.comp_info[2].h_samp_factor = 1;
	cinfo.comp_info[2].v_samp_factor = 1;

	jpeg_start_compress(&cinfo, TRUE);

	int row_stride = width * 3;
	while (cinfo.next_scanline < cinfo.image_height) {
		JSAMPROW row_pointer = &rgb_data[cinfo.next_scanline * row_stride];
		jpeg_write_scanlines(&cinfo, &row_pointer, 1);
	}

	jpeg_finish_compress(&cinfo);
	jpeg_destroy_compress(&cinfo);
	fclose(outfile);

	return 0;
}

// Check if WebP has alpha channel
int webp_has_alpha(const uint8_t *data, size_t data_size) {
	WebPBitstreamFeatures features;
	if (WebPGetFeatures(data, data_size, &features) != VP8_STATUS_OK) {
		return -1;  // Error
	}
	return features.has_alpha ? 1 : 0;
}
*/
import "C"
import (
	"fmt"
	"os"
	"unsafe"
)

// ConvertWebPToJPEG converts a static WebP file to JPEG format
// quality: JPEG quality (1-100)
func ConvertWebPToJPEG(inputPath, outputPath string, quality int) error {
	// Validate quality
	if quality < 1 || quality > 100 {
		return fmt.Errorf("quality must be between 1 and 100, got %d", quality)
	}

	// Read WebP file
	data, err := os.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("failed to read WebP file: %w", err)
	}

	if len(data) == 0 {
		return fmt.Errorf("WebP file is empty")
	}

	// Decode WebP using advanced decoder with maximum quality settings
	decoded, err := DecodeWebPAdvanced(data)
	if err != nil {
		return fmt.Errorf("failed to decode WebP with advanced decoder: %w", err)
	}

	// Convert to RGB (compositing alpha on white if needed)
	rgbData := decoded.ToRGB()

	// Encode to JPEG
	if err := encodeJPEG(outputPath, rgbData, decoded.Width, decoded.Height, quality); err != nil {
		return fmt.Errorf("failed to encode JPEG: %w", err)
	}

	return nil
}

// encodeJPEG encodes RGB data to JPEG file using libjpeg
func encodeJPEG(filename string, rgbData []byte, width, height, quality int) error {
	cFilename := C.CString(filename)
	defer C.free(unsafe.Pointer(cFilename))

	// Copy RGB data to C memory to avoid CGO pointer issues
	cRGBData := C.malloc(C.size_t(len(rgbData)))
	if cRGBData == nil {
		return fmt.Errorf("failed to allocate memory for RGB data")
	}
	defer C.free(cRGBData)

	C.memcpy(cRGBData, unsafe.Pointer(&rgbData[0]), C.size_t(len(rgbData)))

	// Call C function to encode JPEG
	result := C.encode_jpeg_to_file(cFilename, (*C.uchar)(cRGBData), C.int(width), C.int(height), C.int(quality))
	if result != 0 {
		return fmt.Errorf("JPEG encoding failed")
	}

	return nil
}
