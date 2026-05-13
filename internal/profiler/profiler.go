// Package profiler analyzes .env files and produces a usage profile,
// summarizing key counts, secret density, and potential issues.
package profiler

import (
	"fmt"
	"sort"
	"strings"

	"github.com/user/envlens/internal/masker"
)

// Profile holds the analysis result for a set of env entries.
type Profile struct {
	TotalKeys    int
	SecretKeys   int
	EmptyValues  int
	UniqueValues int
	TopKeys      []string // top 5 longest keys
	SecretRatio  float64
}

// String returns a human-readable summary of the profile.
func (p Profile) String() string {
	return fmt.Sprintf(
		"Total keys: %d | Secrets: %d (%.0f%%) | Empty: %d | Unique values: %d",
		p.TotalKeys, p.SecretKeys, p.SecretRatio*100, p.EmptyValues, p.UniqueValues,
	)
}

// Analyze builds a Profile from the provided env map.
func Analyze(env map[string]string) Profile {
	m := masker.New()

	seen := make(map[string]struct{})
	var keys []string

	p := Profile{}

	for k, v := range env {
		p.TotalKeys++
		keys = append(keys, k)

		if m.IsSensitive(k) {
			p.SecretKeys++
		}

		norm := strings.TrimSpace(v)
		if norm == "" {
			p.EmptyValues++
		}

		if _, exists := seen[norm]; !exists {
			seen[norm] = struct{}{}
			p.UniqueValues++
		}
	}

	if p.TotalKeys > 0 {
		p.SecretRatio = float64(p.SecretKeys) / float64(p.TotalKeys)
	}

	// top 5 longest keys
	sort.Slice(keys, func(i, j int) bool {
		return len(keys[i]) > len(keys[j])
	})
	if len(keys) > 5 {
		keys = keys[:5]
	}
	p.TopKeys = keys

	return p
}
