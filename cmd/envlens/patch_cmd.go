package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/yourorg/envlens/internal/parser"
	"github.com/yourorg/envlens/internal/patcher"
)

func runPatch(args []string) error {
	fs := flag.NewFlagSet("patch", flag.ContinueOnError)
	file := fs.String("file", "", "path to .env file to patch")
	ops := fs.String("ops", "", "JSON array of patch ops, e.g. '[{\"action\":\"set\",\"key\":\"FOO\",\"value\":\"bar\"}]'")
	output := fs.String("output", "", "write patched env to this file (default: overwrite input)")
	dryRun := fs.Bool("dry-run", false, "print results without writing")

	if err := fs.Parse(args); err != nil {
		return err
	}

	if *file == "" {
		return fmt.Errorf("--file is required")
	}
	if *ops == "" {
		return fmt.Errorf("--ops is required")
	}

	env, err := parser.ParseFile(*file)
	if err != nil {
		return fmt.Errorf("parse: %w", err)
	}

	var patchOps []patcher.Op
	if err := json.Unmarshal([]byte(*ops), &patchOps); err != nil {
		return fmt.Errorf("invalid ops JSON: %w", err)
	}

	patched, results, err := patcher.Patch(env.ToMap(), patchOps)
	if err != nil {
		return fmt.Errorf("patch: %w", err)
	}

	for _, r := range results {
		status := "ok"
		if !r.Applied {
			status = "skip"
		}
		fmt.Fprintf(os.Stderr, "[%s] %s %s -> %s\n", status, r.Op.Action, r.Op.Key, r.Note)
	}

	if *dryRun {
		return nil
	}

	dest := *output
	if dest == "" {
		dest = *file
	}
	return patcher.WritePatched(dest, patched)
}
