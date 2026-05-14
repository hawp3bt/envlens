package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/user/envlens/internal/parser"
	"github.com/user/envlens/internal/renamer"
)

// runRename implements the `envlens rename` sub-command.
// Flags:
//   -file      path to the .env file
//   -rules     JSON array of rename rules, e.g. '[{"from":"OLD","to":"NEW"}]'
//   -out       output path (defaults to stdout)
//   -dry-run   print what would change without writing
func runRename(args []string) error {
	fs := flag.NewFlagSet("rename", flag.ContinueOnError)
	filePath := fs.String("file", "", "path to .env file (required)")
	rulesJSON := fs.String("rules", "", "JSON rename rules (required)")
	outPath := fs.String("out", "", "output file path (default: stdout)")
	dryRun := fs.Bool("dry-run", false, "preview changes without writing")

	if err := fs.Parse(args); err != nil {
		return err
	}
	if *filePath == "" || *rulesJSON == "" {
		return fmt.Errorf("flags -file and -rules are required")
	}

	var rawRules []struct {
		From     string `json:"from"`
		To       string `json:"to"`
		IsPrefix bool   `json:"is_prefix"`
	}
	if err := json.Unmarshal([]byte(*rulesJSON), &rawRules); err != nil {
		return fmt.Errorf("invalid rules JSON: %w", err)
	}
	rules := make([]renamer.Rule, len(rawRules))
	for i, r := range rawRules {
		rules[i] = renamer.Rule{From: r.From, To: r.To, IsPrefix: r.IsPrefix}
	}

	entries, err := parser.ParseFile(*filePath)
	if err != nil {
		return fmt.Errorf("parse error: %w", err)
	}
	envMap := entries.ToMap()

	opts := renamer.Options{Rules: rules, DryRun: *dryRun}
	renamed, results, err := renamer.Rename(envMap, opts)
	if err != nil {
		return fmt.Errorf("rename error: %w", err)
	}

	if *dryRun {
		for _, r := range results {
			if r.Applied {
				fmt.Printf("  rename: %s -> %s\n", r.Original, r.Renamed)
			}
		}
		return nil
	}

	w := os.Stdout
	if *outPath != "" {
		f, err := os.Create(*outPath)
		if err != nil {
			return fmt.Errorf("cannot create output file: %w", err)
		}
		defer f.Close()
		w = f
	}
	for k, v := range renamed {
		fmt.Fprintf(w, "%s=%s\n", k, v)
	}
	_ = results
	return nil
}
