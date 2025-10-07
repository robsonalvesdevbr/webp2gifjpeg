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

	// Decode WebP to RGBA
	var width, height C.int
	decoded := C.WebPDecodeRGBA(
		(*C.uint8_t)(unsafe.Pointer(&data[0])),
		C.size_t(len(data)),
		&width,
		&height,
	)

	if decoded == nil {
		return fmt.Errorf("failed to decode WebP image")
	}
	defer C.free(unsafe.Pointer(decoded))

	w := int(width)
	h := int(height)

	// Convert RGBA to RGB (remove alpha channel, composite on white background)
	rgbData := make([]byte, w*h*3)
	rgbaData := unsafe.Slice(decoded, w*h*4)

	for i := 0; i < w*h; i++ {
		r := rgbaData[i*4]
		g := rgbaData[i*4+1]
		b := rgbaData[i*4+2]
		a := rgbaData[i*4+3]

		// Composite on white background if there's transparency
		if a < 255 {
			alpha := float32(a) / 255.0
			invAlpha := 1.0 - alpha
			r = C.uint8_t(float32(r)*alpha + 255*invAlpha)
			g = C.uint8_t(float32(g)*alpha + 255*invAlpha)
			b = C.uint8_t(float32(b)*alpha + 255*invAlpha)
		}

		rgbData[i*3] = byte(r)
		rgbData[i*3+1] = byte(g)
		rgbData[i*3+2] = byte(b)
	}

	// Encode to JPEG
	if err := encodeJPEG(outputPath, rgbData, w, h, quality); err != nil {
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
