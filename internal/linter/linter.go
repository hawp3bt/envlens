// Package linter provides key naming convention checks for .env files.
package linter

import (
	"fmt"
	"regexp"
	"strings"
)

// Severity represents the level of a lint finding.
type Severity string

const (
	SeverityError   Severity = "error"
	SeverityWarning Severity = "warning"
)

// Finding represents a single lint result for a key.
type Finding struct {
	Key      string
	Message  string
	Severity Severity
}

var (
	validKeyPattern   = regexp.MustCompile(`^[A-Z][A-Z0-9_]*$`)
	doubleUnder       = regexp.MustCompile(`__`)
	leadingDigit      = regexp.MustCompile(`^[0-9]`)
)

// Lint checks a map of env key-value pairs for naming convention violations.
// It returns a slice of Findings describing any issues found.
func Lint(env map[string]string) []Finding {
	var findings []Finding

	for key := range env {
		if key == "" {
			findings = append(findings, Finding{
				Key:      key,
				Message:  "key must not be empty",
				Severity: SeverityError,
			})
			continue
		}

		if leadingDigit.MatchString(key) {
			findings = append(findings, Finding{
				Key:      key,
				Message:  "key must not start with a digit",
				Severity: SeverityError,
			})
		}

		if strings.ToUpper(key) != key {
			findings = append(findings, Finding{
				Key:      key,
				Message:  "key should be UPPER_SNAKE_CASE",
				Severity: SeverityWarning,
			})
		}

		if doubleUnder.MatchString(key) {
			findings = append(findings, Finding{
				Key:      key,
				Message:  "key contains consecutive underscores",
				Severity: SeverityWarning,
			})
		}

		if strings.HasSuffix(key, "_") || strings.HasPrefix(key, "_") {
			findings = append(findings, Finding{
				Key:      key,
				Message:  "key must not start or end with an underscore",
				Severity: SeverityWarning,
			})
		}
	}

	return findings
}

// Summary returns a human-readable summary string for a slice of findings.
func Summary(findings []Finding) string {
	if len(findings) == 0 {
		return "lint: no issues found"
	}
	var sb strings.Builder
	for _, f := range findings {
		fmt.Fprintf(&sb, "[%s] %s: %s\n", f.Severity, f.Key, f.Message)
	}
	return strings.TrimRight(sb.String(), "\n")
}
