package converter

import (
	"image/gif"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// TestConvertWebPToGIF tests the basic conversion functionality
func TestConvertWebPToGIF(t *testing.T) {
	// Skip if ffmpeg is not available
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		t.Skip("ffmpeg not found, skipping test")
	}

	// Create a temporary directory for test files
	tmpDir := t.TempDir()

	// Create a simple test WebP using ffmpeg
	webpPath := filepath.Join(tmpDir, "test.webp")

	// Create a simple solid color WebP using ffmpeg
	cmd := exec.Command("ffmpeg", "-f", "lavfi", "-i", "color=c=red:s=100x100:d=1", "-y", webpPath)
	if err := cmd.Run(); err != nil {
		t.Skipf("Failed to create test WebP file: %v", err)
	}

	// Convert to GIF
	if err := ConvertWebPToGIF(webpPath); err != nil {
		t.Fatalf("ConvertWebPToGIF failed: %v", err)
	}

	// Check that WebP was removed
	if _, err := os.Stat(webpPath); !os.IsNotExist(err) {
		t.Error("Original WebP file was not removed")
	}

	// Check that GIF was created
	gifPath := filepath.Join(tmpDir, "test.gif")
	if _, err := os.Stat(gifPath); os.IsNotExist(err) {
		t.Error("GIF file was not created")
	}

	// Verify GIF can be decoded
	gifFile, err := os.Open(gifPath)
	if err != nil {
		t.Fatalf("Failed to open GIF file: %v", err)
	}
	defer gifFile.Close()

	_, err = gif.DecodeAll(gifFile)
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

// TestIsAnimatedWebP tests the animated WebP detection
func TestIsAnimatedWebP(t *testing.T) {
	// Skip if ffmpeg is not available
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		t.Skip("ffmpeg not found, skipping test")
	}

	tmpDir := t.TempDir()

	// Create a simple WebP file using ffmpeg
	webpPath := filepath.Join(tmpDir, "test.webp")
	cmd := exec.Command("ffmpeg", "-f", "lavfi", "-i", "color=c=green:s=10x10:d=1", "-y", webpPath)
	if err := cmd.Run(); err != nil {
		t.Skipf("Failed to create test file: %v", err)
	}

	// Test detection
	isAnimated, err := IsAnimatedWebP(webpPath)
	if err != nil {
		t.Errorf("IsAnimatedWebP returned error: %v", err)
	}

	// For this simple test, we just verify it returns without error
	// The actual animation detection would require proper animated WebP files
	t.Logf("IsAnimatedWebP result: %v", isAnimated)
}

// TestConvertWebPToGIF_NonExistentFile tests error handling
func TestConvertWebPToGIF_NonExistentFile(t *testing.T) {
	err := ConvertWebPToGIF("/nonexistent/file.webp")
	if err == nil {
		t.Error("Expected error for non-existent file, got nil")
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
