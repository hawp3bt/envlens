// Package interpolator resolves variable references within .env file values.
// It supports ${VAR} and $VAR syntax, expanding references from the same map
// or from a provided base environment.
package interpolator

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

var refPattern = regexp.MustCompile(`\$\{([A-Za-z_][A-Za-z0-9_]*)\}|\$([A-Za-z_][A-Za-z0-9_]*)`)

// Options controls interpolation behaviour.
type Options struct {
	// FallbackToOS allows unresolved references to be looked up in os.Environ.
	FallbackToOS bool
	// ErrorOnMissing causes Interpolate to return an error for unresolved refs.
	ErrorOnMissing bool
}

// DefaultOptions returns a sensible default Options value.
func DefaultOptions() Options {
	return Options{
		FallbackToOS:   true,
		ErrorOnMissing: false,
	}
}

// Interpolate resolves variable references in every value of env using the
// values already present in env (and optionally os.Environ).
// The input map is not modified; a new map is returned.
func Interpolate(env map[string]string, opts Options) (map[string]string, error) {
	out := make(map[string]string, len(env))
	for k, v := range env {
		resolved, err := resolve(v, env, opts)
		if err != nil {
			return nil, fmt.Errorf("key %q: %w", k, err)
		}
		out[k] = resolved
	}
	return out, nil
}

func resolve(value string, env map[string]string, opts Options) (string, error) {
	var resolveErr error
	result := refPattern.ReplaceAllStringFunc(value, func(match string) string {
		if resolveErr != nil {
			return match
		}
		name := extractName(match)
		if v, ok := env[name]; ok {
			return v
		}
		if opts.FallbackToOS {
			if v, ok := os.LookupEnv(name); ok {
				return v
			}
		}
		if opts.ErrorOnMissing {
			resolveErr = fmt.Errorf("unresolved variable reference: %s", name)
			return match
		}
		return ""
	})
	if resolveErr != nil {
		return "", resolveErr
	}
	return result, nil
}

func extractName(match string) string {
	match = strings.TrimPrefix(match, "$")
	match = strings.TrimPrefix(match, "{")
	match = strings.TrimSuffix(match, "}")
	return match
}
