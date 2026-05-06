package reporter

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/envlens/internal/differ"
	"github.com/envlens/internal/masker"
)

// Format represents the output format for reports.
type Format string

const (
	FormatText Format = "text"
	FormatJSON Format = "json"
)

// Options configures report generation.
type Options struct {
	Format       Format
	ShowUnchanged bool
	MaskSecrets  bool
	Output       io.Writer
}

// DefaultOptions returns sensible defaults for report generation.
func DefaultOptions() Options {
	return Options{
		Format:      FormatText,
		MaskSecrets: true,
		Output:      os.Stdout,
	}
}

// Generate produces a diff report between two env maps.
func Generate(base, target map[string]string, opts Options) error {
	if opts.Output == nil {
		opts.Output = os.Stdout
	}

	m := masker.New()
	if opts.MaskSecrets {
		base = m.MaskMap(base)
		target = m.MaskMap(target)
	}

	results := differ.Diff(base, target)

	switch opts.Format {
	case FormatJSON:
		return writeJSON(opts.Output, results, opts.ShowUnchanged)
	default:
		return writeText(opts.Output, results, opts.ShowUnchanged)
	}
}

func writeText(w io.Writer, results []differ.Result, showUnchanged bool) error {
	output := differ.Format(results, showUnchanged)
	_, err := fmt.Fprint(w, output)
	return err
}

func writeJSON(w io.Writer, results []differ.Result, showUnchanged bool) error {
	var sb strings.Builder
	sb.WriteString("[\n")
	first := true
	for _, r := range results {
		if !showUnchanged && r.Status == differ.StatusUnchanged {
			continue
		}
		if !first {
			sb.WriteString(",\n")
		}
		first = false
		sb.WriteString(fmt.Sprintf(
			`  {"key": %q, "status": %q, "base_value": %q, "target_value": %q}`,
			r.Key, r.Status, r.BaseValue, r.TargetValue,
		))
	}
	sb.WriteString("\n]\n")
	_, err := fmt.Fprint(w, sb.String())
	return err
}
