package reporter

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/yourorg/envlens/internal/auditor"
	"github.com/yourorg/envlens/internal/differ"
	"github.com/yourorg/envlens/internal/masker"
)

// Options controls report generation.
type Options struct {
	Format      string // "text" or "json"
	ShowSame    bool
	MaskSecrets bool
	Audit       bool
	Out         io.Writer
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{
		Format:      "text",
		ShowSame:    false,
		MaskSecrets: true,
		Audit:       false,
		Out:         os.Stdout,
	}
}

// Generate produces a diff report between two env maps.
func Generate(base, target map[string]string, opts Options) error {
	m := masker.New()
	if opts.MaskSecrets {
		base = m.MaskMap(base)
		target = m.MaskMap(target)
	}

	entries := differ.Diff(base, target)

	var auditReport *auditor.Report
	if opts.Audit {
		auditReport = auditor.Audit(target)
	}

	switch opts.Format {
	case "json":
		return writeJSON(opts.Out, entries, auditReport, opts)
	default:
		return writeText(opts.Out, entries, auditReport, opts)
	}
}

func writeText(w io.Writer, entries []differ.Entry, audit *auditor.Report, opts Options) error {
	fmt.Fprint(w, differ.Format(entries, opts.ShowSame))
	if audit != nil && len(audit.Issues) > 0 {
		fmt.Fprintln(w, "\n--- Audit Issues ---")
		fmt.Fprint(w, audit.Summary())
	}
	return nil
}

type jsonOutput struct {
	Diff  []differ.Entry  `json:"diff"`
	Audit []auditor.Issue `json:"audit,omitempty"`
}

func writeJSON(w io.Writer, entries []differ.Entry, audit *auditor.Report, opts Options) error {
	out := jsonOutput{Diff: entries}
	if !opts.ShowSame {
		filtered := entries[:0]
		for _, e := range entries {
			if e.Status != "unchanged" {
				filtered = append(filtered, e)
			}
		}
		out.Diff = filtered
	}
	if audit != nil {
		out.Audit = audit.Issues
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}
