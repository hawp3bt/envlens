package renamer

import (
	"fmt"
	"strings"
)

// Rule describes a single rename operation.
type Rule struct {
	From string // exact key name or prefix (when IsPrefix is true)
	To   string
	IsPrefix bool
}

// Options controls Rename behaviour.
type Options struct {
	Rules     []Rule
	DryRun    bool
	IgnoreCase bool
}

// Result holds the outcome of a rename operation.
type Result struct {
	Original string
	Renamed  string
	Applied  bool
}

// Rename applies the given rules to the supplied env map and returns a new map
// together with a per-key result log. The original map is never mutated.
func Rename(env map[string]string, opts Options) (map[string]string, []Result, error) {
	out := make(map[string]string, len(env))
	results := make([]Result, 0, len(env))

	for k, v := range env {
		newKey, applied := applyRules(k, opts)
		if applied {
			if _, exists := out[newKey]; exists {
				return nil, nil, fmt.Errorf("rename collision: key %q would overwrite existing key %q", k, newKey)
			}
		}
		results = append(results, Result{Original: k, Renamed: newKey, Applied: applied})
		out[newKey] = v
	}
	return out, results, nil
}

func applyRules(key string, opts Options) (string, bool) {
	cmp := key
	if opts.IgnoreCase {
		cmp = strings.ToLower(key)
	}
	for _, r := range opts.Rules {
		from := r.From
		if opts.IgnoreCase {
			from = strings.ToLower(from)
		}
		if r.IsPrefix {
			if strings.HasPrefix(cmp, from) {
				return r.To + key[len(r.From):], true
			}
		} else {
			if cmp == from {
				return r.To, true
			}
		}
	}
	return key, false
}
