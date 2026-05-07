package auditor

import (
	"fmt"
	"sort"
	"strings"
)

// Issue represents a single audit finding.
type Issue struct {
	Key      string
	Severity string // "warn" or "error"
	Message  string
}

// Report holds all audit results for a parsed env map.
type Report struct {
	Issues []Issue
}

// HasErrors returns true if any issue has severity "error".
func (r *Report) HasErrors() bool {
	for _, i := range r.Issues {
		if i.Severity == "error" {
			return true
		}
	}
	return false
}

// Summary returns a human-readable summary string.
func (r *Report) Summary() string {
	if len(r.Issues) == 0 {
		return "No issues found."
	}
	var sb strings.Builder
	for _, issue := range r.Issues {
		fmt.Fprintf(&sb, "[%s] %s: %s\n", strings.ToUpper(issue.Severity), issue.Key, issue.Message)
	}
	return sb.String()
}

// Audit inspects an env map and returns a Report of findings.
func Audit(env map[string]string) *Report {
	report := &Report{}

	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, key := range keys {
		val := env[key]

		if val == "" {
			report.Issues = append(report.Issues, Issue{
				Key:      key,
				Severity: "warn",
				Message:  "empty value",
			})
		}

		if strings.ToLower(val) == "changeme" || strings.ToLower(val) == "todo" || strings.ToLower(val) == "fixme" {
			report.Issues = append(report.Issues, Issue{
				Key:      key,
				Severity: "error",
				Message:  fmt.Sprintf("placeholder value detected: %q", val),
			})
		}

		if strings.Contains(val, " ") && !strings.HasPrefix(val, "\"") {
			report.Issues = append(report.Issues, Issue{
				Key:      key,
				Severity: "warn",
				Message:  "value contains spaces but is not quoted",
			})
		}
	}

	return report
}
