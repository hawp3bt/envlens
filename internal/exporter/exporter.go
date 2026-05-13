// Package exporter provides functionality to export parsed env data
// into various formats such as shell export scripts, Docker env-file
// format, and JSON key-value pairs.
package exporter

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strings"
)

// Format represents the target export format.
type Format string

const (
	FormatShell  Format = "shell"
	FormatDocker Format = "docker"
	FormatJSON   Format = "json"
)

// Options controls export behaviour.
type Options struct {
	Format  Format
	Masked  bool // if true, sensitive values are replaced with "***"
	Sorted  bool // if true, keys are emitted in alphabetical order
}

// DefaultOptions returns sensible export defaults.
func DefaultOptions() Options {
	return Options{
		Format: FormatShell,
		Sorted: true,
	}
}

// Export writes the env map to w using the specified options.
func Export(env map[string]string, opts Options, w io.Writer) error {
	keys := sortedKeys(env, opts.Sorted)

	switch opts.Format {
	case FormatShell:
		return writeShell(env, keys, w)
	case FormatDocker:
		return writeDocker(env, keys, w)
	case FormatJSON:
		return writeJSON(env, keys, w)
	default:
		return fmt.Errorf("exporter: unknown format %q", opts.Format)
	}
}

func writeShell(env map[string]string, keys []string, w io.Writer) error {
	for _, k := range keys {
		v := shellEscape(env[k])
		if _, err := fmt.Fprintf(w, "export %s=%s\n", k, v); err != nil {
			return err
		}
	}
	return nil
}

func writeDocker(env map[string]string, keys []string, w io.Writer) error {
	for _, k := range keys {
		if _, err := fmt.Fprintf(w, "%s=%s\n", k, env[k]); err != nil {
			return err
		}
	}
	return nil
}

func writeJSON(env map[string]string, keys []string, w io.Writer) error {
	ordered := make([]map[string]string, 0, len(keys))
	for _, k := range keys {
		ordered = append(ordered, map[string]string{"key": k, "value": env[k]})
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(ordered)
}

func shellEscape(v string) string {
	if strings.ContainsAny(v, " \t\n'\"\\$`") {
		v = strings.ReplaceAll(v, "'", "'\\''")
		return "'" + v + "'"
	}
	return v
}

func sortedKeys(env map[string]string, sorted bool) []string {
	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}
	if sorted {
		sort.Strings(keys)
	}
	return keys
}
