package main

import (
	"fmt"
	"os"

	"github.com/yourorg/envlens/internal/exporter"
	"github.com/yourorg/envlens/internal/parser"
)

// runExport handles the `envlens export` sub-command.
// Flags:
//
//	--file   path to the .env file (required)
//	--format shell|docker|json  (default: shell)
//	--sorted sort keys alphabetically (default: true)
func runExport(args []string) error {
	var filePath, format string
	var sorted bool

	fs := newFlagSet("export")
	fs.StringVar(&filePath, "file", "", "path to .env file (required)")
	fs.StringVar(&format, "format", "shell", "output format: shell, docker, json")
	fs.BoolVar(&sorted, "sorted", true, "sort keys alphabetically")

	if err := fs.Parse(args); err != nil {
		return err
	}
	if filePath == "" {
		return fmt.Errorf("--file is required")
	}

	entries, err := parser.ParseFile(filePath)
	if err != nil {
		return fmt.Errorf("parse error: %w", err)
	}

	envMap := make(map[string]string, len(entries))
	for _, e := range entries {
		envMap[e.Key] = e.Value
	}

	opts := exporter.Options{
		Format: exporter.Format(format),
		Sorted: sorted,
	}

	if err := exporter.Export(envMap, opts, os.Stdout); err != nil {
		return fmt.Errorf("export error: %w", err)
	}
	return nil
}
