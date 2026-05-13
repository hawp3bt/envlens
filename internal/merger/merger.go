// Package merger provides functionality to merge multiple .env files
// with configurable conflict resolution strategies.
package merger

import "fmt"

// Strategy defines how key conflicts are resolved during a merge.
type Strategy int

const (
	// StrategyFirst keeps the value from the first file that defines the key.
	StrategyFirst Strategy = iota
	// StrategyLast keeps the value from the last file that defines the key.
	StrategyLast
	// StrategyError returns an error if the same key appears with different values.
	StrategyError
)

// Conflict records a key that appeared in more than one source with differing values.
type Conflict struct {
	Key    string
	Values []string // one entry per source file, in order
}

// Result holds the merged environment map and any conflicts detected.
type Result struct {
	Env       map[string]string
	Conflicts []Conflict
}

// Merge combines multiple env maps according to the given Strategy.
// sources are processed left-to-right.
func Merge(sources []map[string]string, strategy Strategy) (*Result, error) {
	result := &Result{
		Env: make(map[string]string),
	}

	// Track per-key values across sources for conflict detection.
	seen := make(map[string][]string)

	for _, src := range sources {
		for k, v := range src {
			seen[k] = append(seen[k], v)
		}
	}

	for k, vals := range seen {
		if len(vals) == 1 {
			result.Env[k] = vals[0]
			continue
		}

		// Check if all values are identical.
		allSame := true
		for _, v := range vals[1:] {
			if v != vals[0] {
				allSame = false
				break
			}
		}

		if allSame {
			result.Env[k] = vals[0]
			continue
		}

		result.Conflicts = append(result.Conflicts, Conflict{Key: k, Values: vals})

		switch strategy {
		case StrategyFirst:
			result.Env[k] = vals[0]
		case StrategyLast:
			result.Env[k] = vals[len(vals)-1]
		case StrategyError:
			return nil, fmt.Errorf("merge conflict on key %q: values %v", k, vals)
		}
	}

	return result, nil
}
