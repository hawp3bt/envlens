package differ

import (
	"fmt"
	"io"
	"strings"
)

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorGray   = "\033[90m"
)

// FormatOptions controls how a diff result is rendered.
type FormatOptions struct {
	Colorized    bool
	ShowUnchanged bool
}

// Format writes a human-readable diff to the given writer.
func Format(w io.Writer, result *Result, opts FormatOptions) {
	for _, entry := range result.Entries {
		if entry.Type == Unchanged && !opts.ShowUnchanged {
			continue
		}
		line := formatEntry(entry, opts.Colorized)
		fmt.Fprintln(w, line)
	}
	fmt.Fprintln(w)
	summary := result.Summary()
	if opts.Colorized {
		summary = colorYellow + summary + colorReset
	}
	fmt.Fprintln(w, summary)
}

func formatEntry(e Entry, colorized bool) string {
	var sb strings.Builder

	switch e.Type {
	case Added:
		prefix := "+ "
		line := fmt.Sprintf("%s%s=%s", prefix, e.Key, e.NewValue)
		if colorized {
			line = colorGreen + line + colorReset
		}
		sb.WriteString(line)
	case Removed:
		prefix := "- "
		line := fmt.Sprintf("%s%s=%s", prefix, e.Key, e.OldValue)
		if colorized {
			line = colorRed + line + colorReset
		}
		sb.WriteString(line)
	case Changed:
		old := fmt.Sprintf("~ %s: %s -> %s", e.Key, e.OldValue, e.NewValue)
		if colorized {
			old = colorYellow + old + colorReset
		}
		sb.WriteString(old)
	case Unchanged:
		line := fmt.Sprintf("  %s=%s", e.Key, e.OldValue)
		if colorized {
			line = colorGray + line + colorReset
		}
		sb.WriteString(line)
	}

	return sb.String()
}
