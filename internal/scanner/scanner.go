// Package scanner detects duplicate keys and conflicting definitions
// within a single .env file or across multiple env maps.
package scanner

import (
	"fmt"
	"strings"
)

// Finding represents a single scan result.
type Finding struct {
	Key      string
	Severity string // "error" | "warning"
	Message  string
}

// Result holds all findings from a scan operation.
type Result struct {
	Findings []Finding
}

// HasErrors returns true if any finding has severity "error".
func (r Result) HasErrors() bool {
	for _, f := range r.Findings {
		if f.Severity == "error" {
			return true
		}
	}
	return false
}

// Summary returns a human-readable summary line.
func (r Result) Summary() string {
	if len(r.Findings) == 0 {
		return "no issues found"
	}
	return fmt.Sprintf("%d issue(s) found", len(r.Findings))
}

// ScanLines inspects raw lines from a single .env file for duplicate keys.
// It returns a Result containing any duplicate-key findings.
func ScanLines(lines []string) Result {
	seen := make(map[string]int) // key -> first occurrence line number (1-based)
	var findings []Finding

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}
		parts := strings.SplitN(trimmed, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		if key == "" {
			continue
		}
		lineNum := i + 1
		if first, exists := seen[key]; exists {
			findings = append(findings, Finding{
				Key:      key,
				Severity: "error",
				Message:  fmt.Sprintf("duplicate key %q: first defined on line %d, redefined on line %d", key, first, lineNum),
			})
		} else {
			seen[key] = lineNum
		}
	}
	return Result{Findings: findings}
}

// ScanMaps checks for keys present in base that have conflicting casing variants
// across the provided maps (e.g. DB_HOST vs db_host).
func ScanMaps(maps ...map[string]string) Result {
	normalized := make(map[string]string) // lower-key -> original key
	var findings []Finding

	for _, m := range maps {
		for k := range m {
			lower := strings.ToLower(k)
			if orig, exists := normalized[lower]; exists && orig != k {
				findings = append(findings, Finding{
					Key:      k,
					Severity: "warning",
					Message:  fmt.Sprintf("case conflict: %q and %q differ only in casing", orig, k),
				})
			} else {
				normalized[lower] = k
			}
		}
	}
	return Result{Findings: findings}
}
