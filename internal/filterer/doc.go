// Package filterer provides flexible key-based filtering of environment variable
// maps. It supports glob patterns (e.g. "DB_*"), anchored regex patterns
// (e.g. "^AWS"), and both include and exclude modes.
//
// Example usage:
//
//	env := map[string]string{"DB_HOST": "localhost", "PORT": "8080"}
//
//	out, err := filterer.Filter(env, filterer.Options{
//		Mode:     filterer.ModeInclude,
//		Patterns: []string{"DB_*"},
//	})
//
// Patterns that begin with '^' or end with '$' are compiled as regular
// expressions. All other patterns are treated as path-style globs via
// path.Match.
package filterer
