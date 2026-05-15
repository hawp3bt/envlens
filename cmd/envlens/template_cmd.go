package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/yourorg/envlens/internal/parser"
	"github.com/yourorg/envlens/internal/templater"
)

func runTemplate(args []string) error {
	fs := flag.NewFlagSet("template", flag.ContinueOnError)
	envFile := fs.String("file", "", "path to .env file")
	tmplFile := fs.String("tmpl", "", "path to template file")
	outFile := fs.String("out", "", "output file (default: stdout)")
	missing := fs.String("missing", "zero", "behaviour for missing keys: zero|error")

	if err := fs.Parse(args); err != nil {
		return err
	}

	if *envFile == "" {
		return fmt.Errorf("--file is required")
	}
	if *tmplFile == "" {
		return fmt.Errorf("--tmpl is required")
	}

	entries, err := parser.ParseFile(*envFile)
	if err != nil {
		return fmt.Errorf("parse env: %w", err)
	}

	envMap := make(map[string]string, len(entries))
	for _, e := range entries {
		envMap[e.Key] = e.Value
	}

	tmplBytes, err := os.ReadFile(*tmplFile)
	if err != nil {
		return fmt.Errorf("read template: %w", err)
	}

	opts := templater.Options{MissingKey: *missing}
	res, err := templater.Render(string(tmplBytes), envMap, opts)
	if err != nil {
		return fmt.Errorf("render: %w", err)
	}

	if *outFile == "" {
		fmt.Print(res.Output)
		return nil
	}

	if err := os.WriteFile(*outFile, []byte(res.Output), 0644); err != nil {
		return fmt.Errorf("write output: %w", err)
	}
	fmt.Fprintf(os.Stderr, "wrote %s\n", *outFile)
	return nil
}
