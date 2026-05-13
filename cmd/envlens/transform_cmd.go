package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/user/envlens/internal/parser"
	"github.com/user/envlens/internal/transformer"
)

func runTransform(args []string) error {
	fs := flag.NewFlagSet("transform", flag.ContinueOnError)
	file := fs.String("file", "", "path to .env file (required)")
	prefix := fs.String("prefix", "", "prefix to prepend to all keys")
	suffix := fs.String("suffix", "", "suffix to append to all keys")
	upper := fs.Bool("uppercase", false, "convert all keys to uppercase")
	lower := fs.Bool("lowercase", false, "convert all keys to lowercase")
	renameRaw := fs.String("rename", "", "comma-separated old=new pairs, e.g. OLD=NEW,FOO=BAR")
	format := fs.String("format", "env", "output format: env or json")

	if err := fs.Parse(args); err != nil {
		return err
	}
	if *file == "" {
		return fmt.Errorf("--file flag is required")
	}

	entries, err := parser.ParseFile(*file)
	if err != nil {
		return fmt.Errorf("parsing file: %w", err)
	}

	envMap := make(map[string]string, len(entries))
	for _, e := range entries {
		envMap[e.Key] = e.Value
	}

	opts := transformer.DefaultOptions()
	opts.Prefix = *prefix
	opts.Suffix = *suffix
	opts.UppercaseKeys = *upper
	opts.LowercaseKeys = *lower

	if *renameRaw != "" {
		for _, pair := range strings.Split(*renameRaw, ",") {
			parts := strings.SplitN(pair, "=", 2)
			if len(parts) == 2 {
				opts.RenameMap[parts[0]] = parts[1]
			}
		}
	}

	out, err := transformer.Transform(envMap, opts)
	if err != nil {
		return fmt.Errorf("transforming: %w", err)
	}

	switch *format {
	case "json":
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(out)
	default:
		for k, v := range out {
			fmt.Fprintf(os.Stdout, "%s=%s\n", k, v)
		}
	}
	return nil
}
