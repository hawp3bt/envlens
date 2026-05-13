// Package transformer provides utilities for transforming .env key-value maps
// by applying rename, prefix, suffix, and uppercase/lowercase operations.
package transformer

import (
	"strings"
)

// Options controls which transformations are applied.
type Options struct {
	Prefix      string
	Suffix      string
	RenameMap   map[string]string
	UppercaseKeys bool
	LowercaseKeys bool
}

// DefaultOptions returns an Options with no transformations applied.
func DefaultOptions() Options {
	return Options{
		RenameMap: make(map[string]string),
	}
}

// Transform applies the configured transformations to the input map and returns
// a new map. The original map is not modified.
func Transform(env map[string]string, opts Options) (map[string]string, error) {
	out := make(map[string]string, len(env))

	for k, v := range env {
		newKey := k

		// Apply rename first
		if renamed, ok := opts.RenameMap[k]; ok {
			newKey = renamed
		}

		// Apply case transformation
		if opts.UppercaseKeys {
			newKey = strings.ToUpper(newKey)
		} else if opts.LowercaseKeys {
			newKey = strings.ToLower(newKey)
		}

		// Apply prefix and suffix
		if opts.Prefix != "" {
			newKey = opts.Prefix + newKey
		}
		if opts.Suffix != "" {
			newKey = newKey + opts.Suffix
		}

		out[newKey] = v
	}

	return out, nil
}
