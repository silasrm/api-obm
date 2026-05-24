package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/silasrm/api-obm/internal/importer"
)

func main() {
	inputFlag := flag.String("input", "", "Path to MySQL dump file")
	outputFlag := flag.String("output", "migrations/postgres/001_obm_schema.sql", "Path to output PostgreSQL file")
	flag.Parse()

	if *inputFlag == "" {
		fmt.Fprintln(os.Stderr, "Error: -input flag is required")
		flag.Usage()
		os.Exit(1)
	}

	inputFile, err := os.Open(*inputFlag)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening input file: %v\n", err)
		os.Exit(1)
	}
	defer inputFile.Close()

	outDir := filepath.Dir(*outputFlag)
	if err := os.MkdirAll(outDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating output directory: %v\n", err)
		os.Exit(1)
	}

	outputFile, err := os.Create(*outputFlag)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating output file: %v\n", err)
		os.Exit(1)
	}
	defer outputFile.Close()

	if err := importer.Convert(inputFile, outputFile, importer.Options{AppendMetadata: true}); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
		os.Exit(1)
	}
}
