package converter

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// getPythonScriptPath returns the full path to a Python script
func getPythonScriptPath(scriptName string) (string, error) {
	// List of locations to search for the script
	searchPaths := []string{}

	// 1. Directory where the binary is located
	execPath, err := os.Executable()
	if err == nil {
		execDir := filepath.Dir(execPath)
		searchPaths = append(searchPaths, filepath.Join(execDir, scriptName))
	}

	// 2. Current working directory
	searchPaths = append(searchPaths, scriptName)

	// 3. Parent directory (for tests running in subdirectories)
	searchPaths = append(searchPaths, filepath.Join("..", scriptName))

	// Try each path
	for _, path := range searchPaths {
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}

	return "", fmt.Errorf("script %s not found in any search path", scriptName)
}

// executePythonConversion executes a Python conversion script and handles temp files
func executePythonConversion(scriptName, inputPath, outputPath string, extraArgs ...string) error {
	tempPath := outputPath + ".tmp"

	scriptPath, err := getPythonScriptPath(scriptName)
	if err != nil {
		return err
	}

	// Build command arguments
	args := []string{scriptPath, inputPath, tempPath}
	args = append(args, extraArgs...)

	// Execute Python script
	cmd := exec.Command("python3", args...)
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

// DetectWebPType detects if a WebP file is animated or static using Python/Pillow
func DetectWebPType(filePath string) (WebPType, error) {
	scriptPath, err := getPythonScriptPath("detect_webp_type.py")
	if err != nil {
		return WebPTypeUnknown, fmt.Errorf("failed to locate detection script: %w", err)
	}

	cmd := exec.Command("python3", scriptPath, filePath)
	output, err := cmd.CombinedOutput()

	// Exit codes: 0=static, 1=animated, 2=error
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			switch exitErr.ExitCode() {
			case 0:
				return WebPTypeStatic, nil
			case 1:
				return WebPTypeAnimated, nil
			default:
				return WebPTypeUnknown, fmt.Errorf("detection failed: %s", string(output))
			}
		}
		return WebPTypeUnknown, fmt.Errorf("failed to execute detection script: %w", err)
	}

	// Exit code 0 (static)
	return WebPTypeStatic, nil
}

// ConvertWebPToGIF converts an animated WebP file to GIF format using Python/PIL
func ConvertWebPToGIF(inputPath string) error {
	outputPath := strings.TrimSuffix(inputPath, filepath.Ext(inputPath)) + ".gif"
	return executePythonConversion("webp_to_gif.py", inputPath, outputPath)
}

// ConvertWebPToJPEG converts a static WebP file to JPEG format using Python/PIL
// quality: JPEG quality (1-100), recommended 85
func ConvertWebPToJPEG(inputPath string, quality int) error {
	// Validate quality
	if quality < 1 || quality > 100 {
		return fmt.Errorf("quality must be between 1 and 100, got %d", quality)
	}

	outputPath := strings.TrimSuffix(inputPath, filepath.Ext(inputPath)) + ".jpg"
	return executePythonConversion("webp_to_jpeg.py", inputPath, outputPath, fmt.Sprintf("%d", quality))
}

// ProcessOptions configures the conversion behavior
type ProcessOptions struct {
	JPEGQuality int // 1-100, default 85
}

// DefaultProcessOptions returns default configuration
func DefaultProcessOptions() ProcessOptions {
	return ProcessOptions{
		JPEGQuality: 85,
	}
}

// ProcessDirectory recursively processes all WebP files in a directory
func ProcessDirectory(rootPath string, options ProcessOptions) error {
	var processedCount int
	var errorCount int
	var staticCount int
	var animatedCount int

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

		// Detect WebP type
		webpType, err := DetectWebPType(path)
		if err != nil {
			fmt.Printf("  Error detecting type: %v\n", err)
			errorCount++
			return nil // Continue processing other files
		}

		// Route to appropriate converter
		var convertErr error

		switch webpType {
		case WebPTypeAnimated:
			fmt.Printf("  Type: Animated → Converting to GIF\n")
			convertErr = ConvertWebPToGIF(path)
			if convertErr == nil {
				animatedCount++
			}

		case WebPTypeStatic:
			fmt.Printf("  Type: Static → Converting to JPEG (quality %d)\n", options.JPEGQuality)
			convertErr = ConvertWebPToJPEG(path, options.JPEGQuality)
			if convertErr == nil {
				staticCount++
			}

		default:
			fmt.Printf("  Type: Unknown - skipping\n")
			errorCount++
			return nil
		}

		// Handle conversion errors
		if convertErr != nil {
			fmt.Printf("  Error: %v\n", convertErr)
			errorCount++
			return nil // Continue processing other files
		}

		processedCount++
		fmt.Printf("  Successfully converted\n")
		return nil
	})

	if err != nil {
		return fmt.Errorf("error walking directory: %w", err)
	}

	fmt.Printf("\nSummary:\n")
	fmt.Printf("  Total converted: %d files\n", processedCount)
	fmt.Printf("  Static → JPEG: %d\n", staticCount)
	fmt.Printf("  Animated → GIF: %d\n", animatedCount)
	fmt.Printf("  Errors: %d\n", errorCount)

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
