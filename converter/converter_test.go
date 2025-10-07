package converter

import (
	"image"
	_ "image/gif"
	_ "image/jpeg"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/robsonalvesdevbr/webp2gifjpeg/native"
)

// TestDetectWebPType tests WebP type detection
func TestDetectWebPType(t *testing.T) {
	// Skip if ffmpeg is not available
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		t.Skip("ffmpeg not found, skipping test")
	}

	tmpDir := t.TempDir()

	// Create a simple animated WebP file using ffmpeg
	webpPath := filepath.Join(tmpDir, "animated.webp")
	cmd := exec.Command("ffmpeg", "-f", "lavfi", "-i", "color=c=red:s=100x100:d=1", "-y", webpPath)
	if err := cmd.Run(); err != nil {
		t.Skipf("Failed to create test WebP file: %v", err)
	}

	// Test detection
	webpType, err := native.DetectWebPType(webpPath)
	if err != nil {
		t.Fatalf("DetectWebPType failed: %v", err)
	}

	// ffmpeg with duration creates animated WebP
	if webpType != native.WebPTypeAnimated && webpType != native.WebPTypeStatic {
		t.Errorf("Unexpected WebP type: %v", webpType)
	}
	t.Logf("Detected WebP type: %s", webpType)
}

// TestConvertWebPToJPEG tests static WebP to JPEG conversion
func TestConvertWebPToJPEG(t *testing.T) {
	// Skip if ffmpeg is not available
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		t.Skip("ffmpeg not found, skipping test")
	}

	tmpDir := t.TempDir()

	// Create a static WebP (using -frames:v 1 to force single frame)
	webpPath := filepath.Join(tmpDir, "static.webp")
	cmd := exec.Command("ffmpeg", "-f", "lavfi", "-i", "color=c=blue:s=100x100", "-frames:v", "1", "-y", webpPath)
	if err := cmd.Run(); err != nil {
		t.Skipf("Failed to create test WebP file: %v", err)
	}

	// Convert to JPEG
	jpegPath := filepath.Join(tmpDir, "static.jpg")
	if err := native.ConvertWebPToJPEG(webpPath, jpegPath, 95); err != nil {
		t.Fatalf("ConvertWebPToJPEG failed: %v", err)
	}

	// Check that JPEG was created
	if _, err := os.Stat(jpegPath); os.IsNotExist(err) {
		t.Fatal("JPEG file was not created")
	}

	// Verify JPEG can be decoded
	jpegFile, err := os.Open(jpegPath)
	if err != nil {
		t.Fatalf("Failed to open JPEG file: %v", err)
	}
	defer jpegFile.Close()

	_, _, err = image.Decode(jpegFile)
	if err != nil {
		t.Errorf("Failed to decode JPEG file: %v", err)
	}
}

// TestConvertWebPToGIF tests animated WebP to GIF conversion
func TestConvertWebPToGIF(t *testing.T) {
	// Skip if ffmpeg is not available
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		t.Skip("ffmpeg not found, skipping test")
	}

	tmpDir := t.TempDir()

	// Create an animated WebP using ffmpeg
	webpPath := filepath.Join(tmpDir, "animated.webp")
	cmd := exec.Command("ffmpeg", "-f", "lavfi", "-i", "color=c=red:s=100x100:d=1", "-y", webpPath)
	if err := cmd.Run(); err != nil {
		t.Skipf("Failed to create test WebP file: %v", err)
	}

	// Convert to GIF
	gifPath := filepath.Join(tmpDir, "animated.gif")
	if err := native.ConvertWebPToGIF(webpPath, gifPath); err != nil {
		t.Fatalf("ConvertWebPToGIF failed: %v", err)
	}

	// Check that GIF was created
	if _, err := os.Stat(gifPath); os.IsNotExist(err) {
		t.Fatal("GIF file was not created")
	}

	// Verify GIF can be decoded
	gifFile, err := os.Open(gifPath)
	if err != nil {
		t.Fatalf("Failed to open GIF file: %v", err)
	}
	defer gifFile.Close()

	_, _, err = image.Decode(gifFile)
	if err != nil {
		t.Errorf("Failed to decode GIF file: %v", err)
	}
}

// TestProcessDirectory tests directory processing
func TestProcessDirectory(t *testing.T) {
	// Skip if ffmpeg is not available
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		t.Skip("ffmpeg not found, skipping test")
	}

	// Create a temporary directory structure
	tmpDir := t.TempDir()
	subDir := filepath.Join(tmpDir, "subdir")
	if err := os.Mkdir(subDir, 0755); err != nil {
		t.Fatalf("Failed to create subdirectory: %v", err)
	}

	// Create WebP files in different directories using ffmpeg
	testFiles := []string{
		filepath.Join(tmpDir, "test1.webp"),
		filepath.Join(tmpDir, "test2.webp"),
		filepath.Join(subDir, "test3.webp"),
	}

	for _, path := range testFiles {
		cmd := exec.Command("ffmpeg", "-f", "lavfi", "-i", "color=c=blue:s=50x50:d=1", "-y", path)
		if err := cmd.Run(); err != nil {
			t.Skipf("Failed to create test WebP file %s: %v", path, err)
		}
	}

	// Process directory
	options := DefaultProcessOptions()
	if err := ProcessDirectory(tmpDir, options); err != nil {
		t.Fatalf("ProcessDirectory failed: %v", err)
	}

	// Check that all WebP files were converted
	for _, webpPath := range testFiles {
		// WebP should be removed
		if _, err := os.Stat(webpPath); !os.IsNotExist(err) {
			t.Errorf("WebP file %s was not removed", webpPath)
		}

		// Check if GIF or JPEG was created (depends on if WebP is animated or static)
		basePath := webpPath[:len(webpPath)-5]
		gifPath := basePath + ".gif"
		jpegPath := basePath + ".jpg"

		gifExists := true
		jpegExists := true

		if _, err := os.Stat(gifPath); os.IsNotExist(err) {
			gifExists = false
		}
		if _, err := os.Stat(jpegPath); os.IsNotExist(err) {
			jpegExists = false
		}

		// At least one output format should exist
		if !gifExists && !jpegExists {
			t.Errorf("Neither GIF nor JPEG was created for %s", webpPath)
		}
	}
}

// TestProcessDirectoryParallel tests parallel directory processing
func TestProcessDirectoryParallel(t *testing.T) {
	// Skip if ffmpeg is not available
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		t.Skip("ffmpeg not found, skipping test")
	}

	// Create a temporary directory structure
	tmpDir := t.TempDir()
	subDir := filepath.Join(tmpDir, "subdir")
	if err := os.Mkdir(subDir, 0755); err != nil {
		t.Fatalf("Failed to create subdirectory: %v", err)
	}

	// Create multiple WebP files in different directories using ffmpeg
	testFiles := []string{
		filepath.Join(tmpDir, "test1.webp"),
		filepath.Join(tmpDir, "test2.webp"),
		filepath.Join(tmpDir, "test3.webp"),
		filepath.Join(tmpDir, "test4.webp"),
		filepath.Join(subDir, "test5.webp"),
		filepath.Join(subDir, "test6.webp"),
	}

	for _, path := range testFiles {
		cmd := exec.Command("ffmpeg", "-f", "lavfi", "-i", "color=c=blue:s=50x50:d=1", "-y", path)
		if err := cmd.Run(); err != nil {
			t.Skipf("Failed to create test WebP file %s: %v", path, err)
		}
	}

	// Process directory with parallel workers
	options := ProcessOptions{
		JPEGQuality: 85,
		NumWorkers:  4,
	}
	if err := ProcessDirectoryParallel(tmpDir, options); err != nil {
		t.Fatalf("ProcessDirectoryParallel failed: %v", err)
	}

	// Check that all WebP files were converted
	for _, webpPath := range testFiles {
		// WebP should be removed
		if _, err := os.Stat(webpPath); !os.IsNotExist(err) {
			t.Errorf("WebP file %s was not removed", webpPath)
		}

		// Check if GIF or JPEG was created
		basePath := webpPath[:len(webpPath)-5]
		gifPath := basePath + ".gif"
		jpegPath := basePath + ".jpg"

		gifExists := true
		jpegExists := true

		if _, err := os.Stat(gifPath); os.IsNotExist(err) {
			gifExists = false
		}
		if _, err := os.Stat(jpegPath); os.IsNotExist(err) {
			jpegExists = false
		}

		// At least one output format should exist
		if !gifExists && !jpegExists {
			t.Errorf("Neither GIF nor JPEG was created for %s", webpPath)
		}
	}
}

// TestProcessDirectoryParallel_EmptyDir tests parallel processing with no files
func TestProcessDirectoryParallel_EmptyDir(t *testing.T) {
	tmpDir := t.TempDir()

	options := ProcessOptions{
		JPEGQuality: 85,
		NumWorkers:  4,
	}
	if err := ProcessDirectoryParallel(tmpDir, options); err != nil {
		t.Fatalf("ProcessDirectoryParallel on empty dir failed: %v", err)
	}
}

// TestProcessDirectory_NonExistentDir tests error handling for invalid directory
func TestProcessDirectory_NonExistentDir(t *testing.T) {
	options := DefaultProcessOptions()
	err := ProcessDirectory("/nonexistent/directory", options)
	if err == nil {
		t.Error("Expected error for non-existent directory, got nil")
	}
}

// TestQualityValidation tests JPEG quality validation
func TestQualityValidation(t *testing.T) {
	tmpDir := t.TempDir()
	webpPath := filepath.Join(tmpDir, "test.webp")
	jpegPath := filepath.Join(tmpDir, "test.jpg")

	// Create a dummy file
	if err := os.WriteFile(webpPath, []byte("dummy"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Test invalid quality (below 1)
	err := native.ConvertWebPToJPEG(webpPath, jpegPath, 0)
	if err == nil {
		t.Error("Expected error for quality=0, got nil")
	}

	// Test invalid quality (above 100)
	err = native.ConvertWebPToJPEG(webpPath, jpegPath, 101)
	if err == nil {
		t.Error("Expected error for quality=101, got nil")
	}
}
