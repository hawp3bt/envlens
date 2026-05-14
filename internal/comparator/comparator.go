// Package comparator provides multi-file .env comparison across more than two
// environments, summarising which keys are present, missing, or divergent.
package comparator

import "sort"

// Presence describes whether a key exists in a given environment file.
type Presence string

const (
	Present Presence = "present"
	Missing Presence = "missing"
)

// KeyStatus holds the per-environment status of a single key.
type KeyStatus struct {
	Key    string
	Values map[string]string   // env label -> value (empty string when missing)
	Status map[string]Presence // env label -> presence
	Diverges bool              // true when present in >1 env with differing values
}

// Result is the full comparison output.
type Result struct {
	Envs     []string    // ordered env labels
	Statuses []KeyStatus // one entry per unique key, sorted alphabetically
}

// Compare accepts a map of label -> parsed env map and returns a Result
// describing every key across all environments.
func Compare(envs map[string]map[string]string) Result {
	labels := make([]string, 0, len(envs))
	for l := range envs {
		labels = append(labels, l)
	}
	sort.Strings(labels)

	// collect all unique keys
	keySet := map[string]struct{}{}
	for _, m := range envs {
		for k := range m {
			keySet[k] = struct{}{}
		}
	}
	keys := make([]string, 0, len(keySet))
	for k := range keySet {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	statuses := make([]KeyStatus, 0, len(keys))
	for _, k := range keys {
		ks := KeyStatus{
			Key:    k,
			Values: make(map[string]string, len(labels)),
			Status: make(map[string]Presence, len(labels)),
		}
		seenValues := map[string]struct{}{}
		for _, lbl := range labels {
			if v, ok := envs[lbl][k]; ok {
				ks.Values[lbl] = v
				ks.Status[lbl] = Present
				seenValues[v] = struct{}{}
			} else {
				ks.Values[lbl] = ""
				ks.Status[lbl] = Missing
			}
		}
		ks.Diverges = len(seenValues) > 1
		statuses = append(statuses, ks)
	}

	return Result{Envs: labels, Statuses: statuses}
}
