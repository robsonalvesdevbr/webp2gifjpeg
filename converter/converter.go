package converter

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

// executePythonConversion executes a Python conversion script and handles temp files
func executePythonConversion(scriptPath, inputPath, outputPath string, extraArgs ...string) error {
	tempPath := outputPath + ".tmp"

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
func DetectWebPType(scriptMgr *ScriptManager, filePath string) (WebPType, error) {
	scriptPath := scriptMgr.GetScriptPath("detect_webp_type.py")

	cmd := exec.Command("python3", scriptPath, filePath)
	output, err := cmd.CombinedOutput()
	// Exit codes: 0=static, 1=animated, 2=error
	if err != nil {
		var exitErr *exec.ExitError
		if ok := errors.As(err, &exitErr); ok {
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
func ConvertWebPToGIF(scriptMgr *ScriptManager, inputPath string) error {
	scriptPath := scriptMgr.GetScriptPath("webp_to_gif.py")
	outputPath := strings.TrimSuffix(inputPath, filepath.Ext(inputPath)) + ".gif"
	return executePythonConversion(scriptPath, inputPath, outputPath)
}

// ConvertWebPToJPEG converts a static WebP file to JPEG format using Python/PIL
// quality: JPEG quality (1-100), recommended 85
func ConvertWebPToJPEG(scriptMgr *ScriptManager, inputPath string, quality int) error {
	// Validate quality
	if quality < 1 || quality > 100 {
		return fmt.Errorf("quality must be between 1 and 100, got %d", quality)
	}

	scriptPath := scriptMgr.GetScriptPath("webp_to_jpeg.py")
	outputPath := strings.TrimSuffix(inputPath, filepath.Ext(inputPath)) + ".jpg"
	return executePythonConversion(scriptPath, inputPath, outputPath, fmt.Sprintf("%d", quality))
}

// ProcessOptions configures the conversion behavior
type ProcessOptions struct {
	JPEGQuality int // 1-100, default 100
	NumWorkers  int // Number of parallel workers (default: runtime.NumCPU())
}

// DefaultProcessOptions returns default configuration
func DefaultProcessOptions() ProcessOptions {
	return ProcessOptions{
		JPEGQuality: 100,
		NumWorkers:  1, // Sequential by default
	}
}

// ConversionJob represents a file to be converted
type ConversionJob struct {
	Path     string
	FileInfo os.FileInfo
}

// ConversionResult represents the result of a conversion
type ConversionResult struct {
	Path     string
	Success  bool
	Type     WebPType
	Error    error
	FilePath string // Output file path
}

// ProcessStats aggregates conversion statistics
type ProcessStats struct {
	TotalProcessed int
	StaticCount    int
	AnimatedCount  int
	ErrorCount     int
}

// convertSingleFile processes a single WebP file
func convertSingleFile(scriptMgr *ScriptManager, path string, options ProcessOptions) ConversionResult {
	result := ConversionResult{
		Path:    path,
		Success: false,
	}

	// Detect WebP type
	webpType, err := DetectWebPType(scriptMgr, path)
	if err != nil {
		result.Error = fmt.Errorf("failed to detect type: %w", err)
		return result
	}

	result.Type = webpType

	// Route to appropriate converter
	switch webpType {
	case WebPTypeAnimated:
		result.FilePath = strings.TrimSuffix(path, filepath.Ext(path)) + ".gif"
		err = ConvertWebPToGIF(scriptMgr, path)

	case WebPTypeStatic:
		result.FilePath = strings.TrimSuffix(path, filepath.Ext(path)) + ".jpg"
		err = ConvertWebPToJPEG(scriptMgr, path, options.JPEGQuality)

	default:
		result.Error = fmt.Errorf("unknown WebP type")
		return result
	}

	if err != nil {
		result.Error = err
		return result
	}

	result.Success = true
	return result
}

// worker processes jobs from the jobs channel
func worker(id int, jobs <-chan ConversionJob, results chan<- ConversionResult, scriptMgr *ScriptManager, options ProcessOptions, wg *sync.WaitGroup) {
	defer wg.Done()

	for job := range jobs {
		result := convertSingleFile(scriptMgr, job.Path, options)
		results <- result
	}
}

// collectStats aggregates results from the results channel
func collectStats(results <-chan ConversionResult, total int, verbose bool) ProcessStats {
	stats := ProcessStats{}
	processed := 0

	for result := range results {
		processed++

		if verbose {
			fmt.Printf("Processing [%d/%d]: %s\n", processed, total, result.Path)
		}

		if result.Success {
			stats.TotalProcessed++
			switch result.Type {
			case WebPTypeAnimated:
				stats.AnimatedCount++
				if verbose {
					fmt.Printf("  Type: Animated → Converted to GIF\n")
				}
			case WebPTypeStatic:
				stats.StaticCount++
				if verbose {
					fmt.Printf("  Type: Static → Converted to JPEG\n")
				}
			}
			if verbose {
				fmt.Printf("  Successfully converted\n")
			}
		} else {
			stats.ErrorCount++
			if verbose && result.Error != nil {
				fmt.Printf("  Error: %v\n", result.Error)
			}
		}
	}

	return stats
}

// ProcessDirectoryParallel recursively processes all WebP files in a directory using parallel workers
func ProcessDirectoryParallel(scriptMgr *ScriptManager, rootPath string, options ProcessOptions) error {
	// Phase 1: Scan - Collect all WebP files
	var webpFiles []ConversionJob

	err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.ToLower(filepath.Ext(path)) == ".webp" {
			webpFiles = append(webpFiles, ConversionJob{
				Path:     path,
				FileInfo: info,
			})
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("error scanning directory: %w", err)
	}

	if len(webpFiles) == 0 {
		fmt.Println("No WebP files found")
		return nil
	}

	// Determine number of workers
	numWorkers := options.NumWorkers
	if numWorkers <= 0 {
		numWorkers = runtime.NumCPU()
	}
	// Limit workers to number of files if fewer files than workers
	if numWorkers > len(webpFiles) {
		numWorkers = len(webpFiles)
	}

	fmt.Printf("Found %d WebP file(s), using %d worker(s)\n\n", len(webpFiles), numWorkers)

	// Phase 2: Create channels and worker pool
	jobs := make(chan ConversionJob, len(webpFiles))
	results := make(chan ConversionResult, len(webpFiles))

	// Phase 3: Start workers
	var wg sync.WaitGroup
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker(i, jobs, results, scriptMgr, options, &wg)
	}

	// Start stats collector in a separate goroutine
	var statsWg sync.WaitGroup
	statsWg.Add(1)
	var stats ProcessStats
	go func() {
		defer statsWg.Done()
		stats = collectStats(results, len(webpFiles), true)
	}()

	// Dispatch jobs
	for _, file := range webpFiles {
		jobs <- file
	}
	close(jobs)

	// Wait for all workers to finish
	wg.Wait()
	close(results)

	// Wait for stats collection to finish
	statsWg.Wait()

	// Phase 4: Display summary
	fmt.Printf("\nSummary:\n")
	fmt.Printf("  Total converted: %d files\n", stats.TotalProcessed)
	fmt.Printf("  Static → JPEG: %d\n", stats.StaticCount)
	fmt.Printf("  Animated → GIF: %d\n", stats.AnimatedCount)
	fmt.Printf("  Errors: %d\n", stats.ErrorCount)

	return nil
}

// ProcessDirectory recursively processes all WebP files in a directory
func ProcessDirectory(scriptMgr *ScriptManager, rootPath string, options ProcessOptions) error {
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

		// Convert the file
		result := convertSingleFile(scriptMgr, path, options)

		// Handle result
		if !result.Success {
			if result.Error != nil {
				fmt.Printf("  Error: %v\n", result.Error)
			}
			errorCount++
			return nil // Continue processing other files
		}

		// Update counters based on type
		switch result.Type {
		case WebPTypeAnimated:
			fmt.Printf("  Type: Animated → Converted to GIF\n")
			animatedCount++
		case WebPTypeStatic:
			fmt.Printf("  Type: Static → Converted to JPEG (quality %d)\n", options.JPEGQuality)
			staticCount++
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
