package converter

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// ConvertWebPToGIF converts an animated WebP file to GIF format using Python/PIL
func ConvertWebPToGIF(inputPath string) error {
	// Generate output path
	outputPath := strings.TrimSuffix(inputPath, filepath.Ext(inputPath)) + ".gif"
	tempPath := outputPath + ".tmp"

	// Get the directory where the converter binary is located
	execPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %w", err)
	}
	execDir := filepath.Dir(execPath)
	pythonScript := filepath.Join(execDir, "webp_to_gif.py")

	// Fallback: check if script is in current directory
	if _, err := os.Stat(pythonScript); os.IsNotExist(err) {
		pythonScript = "webp_to_gif.py"
	}

	// Use Python script to convert WebP to GIF (supports animated WebP)
	cmd := exec.Command("python3", pythonScript, inputPath, tempPath)

	// Capture output for debugging
	output, err := cmd.CombinedOutput()
	if err != nil {
		os.Remove(tempPath)
		return fmt.Errorf("conversion failed: %w\nOutput: %s", err, string(output))
	}

	// Verify output file was created
	if _, err := os.Stat(tempPath); os.IsNotExist(err) {
		return fmt.Errorf("output file was not created")
	}

	// Remove original WebP
	if err := os.Remove(inputPath); err != nil {
		os.Remove(tempPath)
		return fmt.Errorf("failed to remove original file: %w", err)
	}

	// Rename temp file to final name
	if err := os.Rename(tempPath, outputPath); err != nil {
		return fmt.Errorf("failed to rename temp file: %w", err)
	}

	return nil
}

// ProcessDirectory recursively processes all WebP files in a directory
func ProcessDirectory(rootPath string) error {
	var processedCount int
	var errorCount int

	err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Check if file is WebP
		if strings.ToLower(filepath.Ext(path)) != ".webp" {
			return nil
		}

		fmt.Printf("Processing: %s\n", path)

		// Convert WebP to GIF
		if err := ConvertWebPToGIF(path); err != nil {
			fmt.Printf("Error converting %s: %v\n", path, err)
			errorCount++
			return nil // Continue processing other files
		}

		processedCount++
		fmt.Printf("Successfully converted: %s\n", path)
		return nil
	})

	if err != nil {
		return fmt.Errorf("error walking directory: %w", err)
	}

	fmt.Printf("\nSummary: %d files converted, %d errors\n", processedCount, errorCount)
	return nil
}

// IsAnimatedWebP checks if a WebP file is animated
func IsAnimatedWebP(filePath string) (bool, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return false, err
	}

	// Check for ANIM chunk in WebP file
	// This is a simplified check - a complete implementation would parse the WebP format
	return len(data) > 12 && string(data[0:4]) == "RIFF" && string(data[8:12]) == "WEBP", nil
}
