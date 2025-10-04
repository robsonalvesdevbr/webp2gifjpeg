package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/robson/webp2gifjpeg/converter"
)

func main() {
	// Define command line flags
	dirPtr := flag.String("dir", ".", "Directory to process (default: current directory)")
	qualityPtr := flag.Int("quality", 85, "JPEG quality for static WebP (1-100, default: 85)")
	flag.Parse()

	// Validate quality
	if *qualityPtr < 1 || *qualityPtr > 100 {
		fmt.Fprintf(os.Stderr, "Error: quality must be between 1 and 100\n")
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
	fmt.Printf("JPEG Quality: %d\n\n", *qualityPtr)

	// Process all WebP files in directory with options
	options := converter.ProcessOptions{
		JPEGQuality: *qualityPtr,
	}

	if err := converter.ProcessDirectory(absPath, options); err != nil {
		fmt.Fprintf(os.Stderr, "Error processing directory: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("\nConversion completed!")
}
