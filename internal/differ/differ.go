package differ

import (
	"fmt"
	"sort"
)

// DiffType represents the kind of difference between two env files.
type DiffType string

const (
	Added    DiffType = "added"
	Removed  DiffType = "removed"
	Changed  DiffType = "changed"
	Unchanged DiffType = "unchanged"
)

// Entry represents a single diff result for a key.
type Entry struct {
	Key      string
	Type     DiffType
	OldValue string
	NewValue string
}

// Result holds the full diff between two env maps.
type Result struct {
	Entries []Entry
}

// HasChanges returns true if any entries are not unchanged.
func (r *Result) HasChanges() bool {
	for _, e := range r.Entries {
		if e.Type != Unchanged {
			return true
		}
	}
	return false
}

// Summary returns a human-readable summary line.
func (r *Result) Summary() string {
	var added, removed, changed int
	for _, e := range r.Entries {
		switch e.Type {
		case Added:
			added++
		case Removed:
			removed++
		case Changed:
			changed++
		}
	}
	return fmt.Sprintf("+%d added, -%d removed, ~%d changed", added, removed, changed)
}

// Diff computes the difference between a base env map and a target env map.
// Keys present in both maps are compared; missing or extra keys are flagged.
func Diff(base, target map[string]string) *Result {
	result := &Result{}

	allKeys := make(map[string]struct{})
	for k := range base {
		allKeys[k] = struct{}{}
	}
	for k := range target {
		allKeys[k] = struct{}{}
	}

	keys := make([]string, 0, len(allKeys))
	for k := range allKeys {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		baseVal, inBase := base[k]
		targetVal, inTarget := target[k]

		switch {
		case inBase && !inTarget:
			result.Entries = append(result.Entries, Entry{Key: k, Type: Removed, OldValue: baseVal})
		case !inBase && inTarget:
			result.Entries = append(result.Entries, Entry{Key: k, Type: Added, NewValue: targetVal})
		case baseVal != targetVal:
			result.Entries = append(result.Entries, Entry{Key: k, Type: Changed, OldValue: baseVal, NewValue: targetVal})
		default:
			result.Entries = append(result.Entries, Entry{Key: k, Type: Unchanged, OldValue: baseVal, NewValue: targetVal})
		}
	}

	return result
}
