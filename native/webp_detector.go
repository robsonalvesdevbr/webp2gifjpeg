package native

/*
#cgo pkg-config: libwebp libwebpdemux
#include <stdlib.h>
#include <string.h>
#include <webp/decode.h>
#include <webp/demux.h>
*/
import "C"
import (
	"fmt"
	"os"
	"unsafe"
)

// WebPType represents the type of a WebP file
type WebPType int

const (
	WebPTypeUnknown WebPType = iota
	WebPTypeStatic
	WebPTypeAnimated
)

func (t WebPType) String() string {
	switch t {
	case WebPTypeStatic:
		return "static"
	case WebPTypeAnimated:
		return "animated"
	default:
		return "unknown"
	}
}

// DetectWebPType detects if a WebP file is animated or static using libwebp
func DetectWebPType(filePath string) (WebPType, error) {
	// Read file into memory
	data, err := os.ReadFile(filePath)
	if err != nil {
		return WebPTypeUnknown, fmt.Errorf("failed to read file: %w", err)
	}

	if len(data) == 0 {
		return WebPTypeUnknown, fmt.Errorf("file is empty")
	}

	// Allocate C memory and copy data to avoid CGO pointer issues
	cData := C.malloc(C.size_t(len(data)))
	if cData == nil {
		return WebPTypeUnknown, fmt.Errorf("failed to allocate memory")
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
		return WebPTypeUnknown, fmt.Errorf("failed to create WebP demuxer")
	}
	defer C.WebPDemuxDelete(demux)

	// Get canvas information
	flags := C.WebPDemuxGetI(demux, C.WEBP_FF_FORMAT_FLAGS)

	// Check if file has animation flag
	if (flags & C.ANIMATION_FLAG) != 0 {
		return WebPTypeAnimated, nil
	}

	return WebPTypeStatic, nil
}

// GetWebPInfo returns basic information about a WebP file
func GetWebPInfo(filePath string) (width, height, frameCount int, err error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("failed to read file: %w", err)
	}

	if len(data) == 0 {
		return 0, 0, 0, fmt.Errorf("file is empty")
	}

	// Allocate C memory and copy data
	cData := C.malloc(C.size_t(len(data)))
	if cData == nil {
		return 0, 0, 0, fmt.Errorf("failed to allocate memory")
	}
	defer C.free(cData)

	C.memcpy(cData, unsafe.Pointer(&data[0]), C.size_t(len(data)))

	webpData := C.WebPData{
		bytes: (*C.uint8_t)(cData),
		size:  C.size_t(len(data)),
	}

	demux := C.WebPDemux(&webpData)
	if demux == nil {
		return 0, 0, 0, fmt.Errorf("failed to create WebP demuxer")
	}
	defer C.WebPDemuxDelete(demux)

	width = int(C.WebPDemuxGetI(demux, C.WEBP_FF_CANVAS_WIDTH))
	height = int(C.WebPDemuxGetI(demux, C.WEBP_FF_CANVAS_HEIGHT))
	frameCount = int(C.WebPDemuxGetI(demux, C.WEBP_FF_FRAME_COUNT))

	return width, height, frameCount, nil
}
