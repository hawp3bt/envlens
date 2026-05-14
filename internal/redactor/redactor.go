// Package redactor provides functionality to redact sensitive values
// from env maps before they are written to logs, reports, or output files.
package redactor

import (
	"strings"

	"github.com/yourorg/envlens/internal/masker"
)

// Options controls redaction behaviour.
type Options struct {
	// Replacement is the string used in place of a redacted value.
	// Defaults to "[REDACTED]".
	Replacement string

	// ExtraKeys lists additional key substrings to treat as sensitive,
	// beyond those detected by the masker.
	ExtraKeys []string
}

// DefaultOptions returns a sensible default Options.
func DefaultOptions() Options {
	return Options{
		Replacement: "[REDACTED]",
	}
}

// Redactor applies secret redaction to env maps.
type Redactor struct {
	opts   Options
	masker *masker.Masker
}

// New creates a Redactor with the given options.
func New(opts Options) *Redactor {
	if opts.Replacement == "" {
		opts.Replacement = "[REDACTED]"
	}
	return &Redactor{
		opts:   opts,
		masker: masker.New(),
	}
}

// Redact returns a copy of env where all sensitive values are replaced
// with the configured replacement string.
func (r *Redactor) Redact(env map[string]string) map[string]string {
	out := make(map[string]string, len(env))
	for k, v := range env {
		if r.isSensitive(k) {
			out[k] = r.opts.Replacement
		} else {
			out[k] = v
		}
	}
	return out
}

// RedactValue returns the replacement string if the key is sensitive,
// otherwise it returns the original value unchanged.
func (r *Redactor) RedactValue(key, value string) string {
	if r.isSensitive(key) {
		return r.opts.Replacement
	}
	return value
}

func (r *Redactor) isSensitive(key string) bool {
	if r.masker.IsSensitive(key) {
		return true
	}
	upper := strings.ToUpper(key)
	for _, extra := range r.opts.ExtraKeys {
		if strings.Contains(upper, strings.ToUpper(extra)) {
			return true
		}
	}
	return false
}
