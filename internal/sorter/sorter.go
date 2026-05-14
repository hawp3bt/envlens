// Package sorter provides utilities for sorting and ordering .env file entries.
// Keys can be sorted alphabetically, by group prefix, or by a custom priority list.
package sorter

import (
	"sort"
	"strings"
)

// Strategy defines how keys should be ordered.
type Strategy string

const (
	StrategyAlpha    Strategy = "alpha"    // alphabetical A→Z
	StrategyAlphaDesc Strategy = "alpha_desc" // alphabetical Z→A
	StrategyGroup    Strategy = "group"    // grouped by prefix (e.g. DB_, AWS_)
	StrategyPriority Strategy = "priority" // user-supplied key order first, rest alpha
)

// Options controls sorting behaviour.
type Options struct {
	Strategy     Strategy
	PriorityKeys []string // used with StrategyPriority
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{Strategy: StrategyAlpha}
}

// Sort returns the keys of env in the order dictated by opts.
func Sort(env map[string]string, opts Options) []string {
	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}

	switch opts.Strategy {
	case StrategyAlphaDesc:
		sort.Slice(keys, func(i, j int) bool { return keys[i] > keys[j] })
	case StrategyGroup:
		sort.Slice(keys, func(i, j int) bool {
			gi, gj := groupPrefix(keys[i]), groupPrefix(keys[j])
			if gi != gj {
				return gi < gj
			}
			return keys[i] < keys[j]
		})
	case StrategyPriority:
		priority := make(map[string]int, len(opts.PriorityKeys))
		for idx, k := range opts.PriorityKeys {
			priority[k] = idx
		}
		sort.Slice(keys, func(i, j int) bool {
			pi, iOk := priority[keys[i]]
			pj, jOk := priority[keys[j]]
			switch {
			case iOk && jOk:
				return pi < pj
			case iOk:
				return true
			case jOk:
				return false
			default:
				return keys[i] < keys[j]
			}
		})
	default: // StrategyAlpha
		sort.Strings(keys)
	}

	return keys
}

// groupPrefix returns the portion of a key up to (and including) the first
// underscore, or the whole key if no underscore is present.
func groupPrefix(key string) string {
	if idx := strings.Index(key, "_"); idx != -1 {
		return key[:idx+1]
	}
	return key
}
