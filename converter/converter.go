package converter

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/robsonalvesdevbr/webpconvert/native"
)

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
	Type     native.WebPType
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
func convertSingleFile(path string, options ProcessOptions) ConversionResult {
	result := ConversionResult{
		Path:    path,
		Success: false,
	}

	// Detect WebP type using native implementation
	webpType, err := native.DetectWebPType(path)
	if err != nil {
		result.Error = fmt.Errorf("failed to detect type: %w", err)
		return result
	}

	result.Type = webpType

	// Create temp output path
	baseWithoutExt := strings.TrimSuffix(path, filepath.Ext(path))
	var outputPath string
	var tempPath string

	// Route to appropriate converter
	switch webpType {
	case native.WebPTypeAnimated:
		outputPath = baseWithoutExt + ".gif"
		tempPath = outputPath + ".tmp"
		err = native.ConvertWebPToGIF(path, tempPath)

	case native.WebPTypeStatic:
		outputPath = baseWithoutExt + ".jpg"
		tempPath = outputPath + ".tmp"
		err = native.ConvertWebPToJPEG(path, tempPath, options.JPEGQuality)

	default:
		result.Error = fmt.Errorf("unknown WebP type")
		return result
	}

	if err != nil {
		os.Remove(tempPath)
		result.Error = err
		return result
	}

	// Verify temp file was created
	if _, err := os.Stat(tempPath); os.IsNotExist(err) {
		result.Error = fmt.Errorf("output file was not created")
		return result
	}

	// Remove original WebP
	if err := os.Remove(path); err != nil {
		os.Remove(tempPath)
		result.Error = fmt.Errorf("failed to remove original file: %w", err)
		return result
	}

	// Rename temp file to final name
	if err := os.Rename(tempPath, outputPath); err != nil {
		result.Error = fmt.Errorf("failed to rename temp file: %w", err)
		return result
	}

	result.Success = true
	result.FilePath = outputPath
	return result
}

// worker processes jobs from the jobs channel
func worker(id int, jobs <-chan ConversionJob, results chan<- ConversionResult, options ProcessOptions, wg *sync.WaitGroup) {
	defer wg.Done()

	for job := range jobs {
		result := convertSingleFile(job.Path, options)
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
			case native.WebPTypeAnimated:
				stats.AnimatedCount++
				if verbose {
					fmt.Printf("  Type: Animated → Converted to GIF\n")
				}
			case native.WebPTypeStatic:
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
func ProcessDirectoryParallel(rootPath string, options ProcessOptions) error {
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
		go worker(i, jobs, results, options, &wg)
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

		// Convert the file
		result := convertSingleFile(path, options)

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
		case native.WebPTypeAnimated:
			fmt.Printf("  Type: Animated → Converted to GIF\n")
			animatedCount++
		case native.WebPTypeStatic:
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
