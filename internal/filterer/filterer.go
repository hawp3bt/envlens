// Package filterer provides key-based filtering of env maps using glob patterns,
// prefix matching, and regex rules.
package filterer

import (
	"fmt"
	"path"
	"regexp"
	"strings"
)

// Mode controls how the filter is applied.
type Mode string

const (
	ModeInclude Mode = "include"
	ModeExclude Mode = "exclude"
)

// Options configures the filtering behaviour.
type Options struct {
	Mode     Mode
	Patterns []string // glob or prefix patterns (e.g. "DB_*", "^AWS")
}

// DefaultOptions returns sensible defaults: include everything.
func DefaultOptions() Options {
	return Options{
		Mode:     ModeInclude,
		Patterns: []string{"*"},
	}
}

// Filter applies the options to src and returns a filtered copy.
func Filter(src map[string]string, opts Options) (map[string]string, error) {
	matchers, err := compilePatterns(opts.Patterns)
	if err != nil {
		return nil, err
	}

	out := make(map[string]string, len(src))
	for k, v := range src {
		matched := anyMatch(matchers, k)
		if opts.Mode == ModeInclude && matched {
			out[k] = v
		} else if opts.Mode == ModeExclude && !matched {
			out[k] = v
		}
	}
	return out, nil
}

type matcher interface {
	matches(key string) bool
}

type globMatcher struct{ pattern string }

func (g globMatcher) matches(key string) bool {
	ok, _ := path.Match(g.pattern, key)
	return ok
}

type regexMatcher struct{ re *regexp.Regexp }

func (r regexMatcher) matches(key string) bool { return r.re.MatchString(key) }

func compilePatterns(patterns []string) ([]matcher, error) {
	var matchers []matcher
	for _, p := range patterns {
		if strings.HasPrefix(p, "^") || strings.HasSuffix(p, "$") {
			re, err := regexp.Compile(p)
			if err != nil {
				return nil, fmt.Errorf("filterer: invalid regex %q: %w", p, err)
			}
			matchers = append(matchers, regexMatcher{re})
		} else {
			matchers = append(matchers, globMatcher{p})
		}
	}
	return matchers, nil
}

func anyMatch(matchers []matcher, key string) bool {
	for _, m := range matchers {
		if m.matches(key) {
			return true
		}
	}
	return false
}
