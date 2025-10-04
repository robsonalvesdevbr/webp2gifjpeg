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
	flag.Parse()

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

	fmt.Printf("Processing WebP files in: %s\n\n", absPath)

	// Process all WebP files in directory
	if err := converter.ProcessDirectory(absPath); err != nil {
		fmt.Fprintf(os.Stderr, "Error processing directory: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("\nConversion completed!")
}
