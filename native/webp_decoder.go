package native

/*
#cgo pkg-config: libwebp
#include <stdlib.h>
#include <string.h>
#include <webp/decode.h>

// DecodedImage holds decoded WebP image data
typedef struct {
    uint8_t* data;
    int width;
    int height;
    int has_alpha;
    int stride;
} DecodedImage;

// decode_webp_advanced decodes WebP with maximum quality settings
DecodedImage* decode_webp_advanced(const uint8_t* webp_data, size_t webp_size) {
    if (!webp_data || webp_size == 0) {
        return NULL;
    }

    // Initialize decoder configuration
    WebPDecoderConfig config;
    if (!WebPInitDecoderConfig(&config)) {
        return NULL;
    }

    // Get image features first
    WebPBitstreamFeatures features;
    VP8StatusCode status = WebPGetFeatures(webp_data, webp_size, &features);
    if (status != VP8_STATUS_OK) {
        return NULL;
    }

    // Configure decoder for maximum quality
    config.options.bypass_filtering = 0;           // Apply deblocking filters
    config.options.no_fancy_upsampling = 0;        // Use high-quality upsampling
    config.options.use_threads = 1;                // Enable multi-threading
    config.options.dithering_strength = 100;       // Maximum dithering
    config.options.flip = 0;                       // Don't flip
    config.options.alpha_dithering_strength = 100; // Maximum alpha dithering
    config.options.use_scaling = 0;                // No scaling
    config.options.scaled_width = 0;
    config.options.scaled_height = 0;

    // Configure output format - ALWAYS use RGBA for consistency and correctness
    // This prevents buffer mismatch issues and simplifies code
    config.output.colorspace = MODE_RGBA;

    // Decode with advanced configuration
    status = WebPDecode(webp_data, webp_size, &config);
    if (status != VP8_STATUS_OK) {
        WebPFreeDecBuffer(&config.output);
        return NULL;
    }

    // Allocate result structure
    DecodedImage* result = (DecodedImage*)malloc(sizeof(DecodedImage));
    if (!result) {
        WebPFreeDecBuffer(&config.output);
        return NULL;
    }

    // Set image properties
    result->width = config.output.width;
    result->height = config.output.height;
    result->has_alpha = features.has_alpha;

    // Calculate data size and stride - always 4 channels (RGBA)
    int stride = result->width * 4;
    size_t data_size = stride * result->height;
    result->stride = stride;

    // Allocate and copy image data
    result->data = (uint8_t*)malloc(data_size);
    if (!result->data) {
        free(result);
        WebPFreeDecBuffer(&config.output);
        return NULL;
    }

    // Copy decoded data - always from RGBA buffer
    memcpy(result->data, config.output.u.RGBA.rgba, data_size);

    // Free WebP decoder buffer
    WebPFreeDecBuffer(&config.output);

    return result;
}

// free_decoded_image frees memory allocated for DecodedImage
void free_decoded_image(DecodedImage* img) {
    if (img) {
        if (img->data) {
            free(img->data);
        }
        free(img);
    }
}
*/
import "C"
import (
	"fmt"
	"unsafe"
)

// DecodedWebPImage represents a decoded WebP image with high quality settings
type DecodedWebPImage struct {
	Data     []byte
	Width    int
	Height   int
	HasAlpha bool
	Stride   int
}

// DecodeWebPAdvanced decodes a WebP image with maximum quality settings
// This function uses WebPDecoderConfig with optimized settings for best quality
func DecodeWebPAdvanced(webpData []byte) (*DecodedWebPImage, error) {
	if len(webpData) == 0 {
		return nil, fmt.Errorf("empty WebP data")
	}

	// Call C function to decode with advanced settings
	cDecoded := C.decode_webp_advanced(
		(*C.uint8_t)(unsafe.Pointer(&webpData[0])),
		C.size_t(len(webpData)),
	)

	if cDecoded == nil {
		return nil, fmt.Errorf("failed to decode WebP image with advanced decoder")
	}
	defer C.free_decoded_image(cDecoded)

	// Extract image properties
	width := int(cDecoded.width)
	height := int(cDecoded.height)
	hasAlpha := int(cDecoded.has_alpha) != 0
	stride := int(cDecoded.stride)

	// Calculate data size
	dataSize := stride * height

	// Copy image data to Go slice
	data := make([]byte, dataSize)
	C.memcpy(
		unsafe.Pointer(&data[0]),
		unsafe.Pointer(cDecoded.data),
		C.size_t(dataSize),
	)

	return &DecodedWebPImage{
		Data:     data,
		Width:    width,
		Height:   height,
		HasAlpha: hasAlpha,
		Stride:   stride,
	}, nil
}

// ToRGB converts the decoded image to RGB format (compositing alpha on white if needed)
func (img *DecodedWebPImage) ToRGB() []byte {
	// Convert RGBA to RGB by compositing on white background
	// Note: decoder always returns RGBA now (4 channels)
	rgbSize := img.Width * img.Height * 3
	rgbData := make([]byte, rgbSize)

	for i := 0; i < img.Width*img.Height; i++ {
		r := float32(img.Data[i*4])
		g := float32(img.Data[i*4+1])
		b := float32(img.Data[i*4+2])
		a := float32(img.Data[i*4+3])

		// Composite on white background using float arithmetic for precision
		// This matches what Pillow/PIL does and prevents color precision loss
		if a < 255 {
			alpha := a / 255.0
			invAlpha := 1.0 - alpha
			r = r*alpha + 255*invAlpha
			g = g*alpha + 255*invAlpha
			b = b*alpha + 255*invAlpha
		}

		rgbData[i*3] = byte(r + 0.5)   // Round to nearest
		rgbData[i*3+1] = byte(g + 0.5)
		rgbData[i*3+2] = byte(b + 0.5)
	}

	return rgbData
}
