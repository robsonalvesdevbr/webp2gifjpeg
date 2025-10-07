package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/robsonalvesdevbr/webpconvert/converter"
)

func main() {
	// Define command line flags
	dirPtr := flag.String("dir", ".", "Directory to process (default: current directory)")
	qualityPtr := flag.Int("quality", 100, "JPEG quality for static WebP (1-100, default: 100)")
	workersPtr := flag.Int("workers", runtime.NumCPU(), "Number of parallel workers (default: CPU count)")
	keepOriginalPtr := flag.Bool("keep-original", false, "Keep original WebP files after conversion (default: false)")
	flag.Parse()

	// Validate quality
	if *qualityPtr < 1 || *qualityPtr > 100 {
		fmt.Fprintf(os.Stderr, "Error: quality must be between 1 and 100\n")
		os.Exit(1)
	}

	// Validate workers
	if *workersPtr < 1 {
		fmt.Fprintf(os.Stderr, "Error: workers must be at least 1\n")
		os.Exit(1)
	}

	// Get absolute path
	absPath, err := filepath.Abs(*dirPtr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error resolving path: %v\n", err)
		os.Exit(1)
	}

	// Check if directory exists
	info, err := os.Stat(absPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error accessing directory: %v\n", err)
		os.Exit(1)
	}

	if !info.IsDir() {
		fmt.Fprintf(os.Stderr, "Path is not a directory: %s\n", absPath)
		os.Exit(1)
	}

	fmt.Printf("Processing WebP files in: %s\n", absPath)
	fmt.Printf("JPEG Quality: %d\n", *qualityPtr)
	fmt.Printf("Parallel Workers: %d\n", *workersPtr)
	fmt.Printf("Keep Original: %v\n\n", *keepOriginalPtr)

	// Process all WebP files in directory with options
	options := converter.ProcessOptions{
		JPEGQuality:  *qualityPtr,
		NumWorkers:   *workersPtr,
		KeepOriginal: *keepOriginalPtr,
	}

	// Use parallel processing if more than 1 worker is specified
	if *workersPtr > 1 {
		if err := converter.ProcessDirectoryParallel(absPath, options); err != nil {
			fmt.Fprintf(os.Stderr, "Error processing directory: %v\n", err)
			os.Exit(1)
		}
	} else {
		if err := converter.ProcessDirectory(absPath, options); err != nil {
			fmt.Fprintf(os.Stderr, "Error processing directory: %v\n", err)
			os.Exit(1)
		}
	}

	fmt.Println("\nConversion completed!")
}
