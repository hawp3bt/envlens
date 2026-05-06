// Package reporter provides functionality for generating diff reports
// between two sets of environment variables.
//
// It combines the differ and masker packages to produce human-readable
// or machine-parseable output showing what has changed between two
// .env files or environments.
//
// Supported output formats:
//   - text: coloured, line-oriented diff output (default)
//   - json: structured JSON array of diff entries
//
// Secret masking is enabled by default to prevent accidental exposure
// of sensitive values such as passwords, tokens, and API keys.
//
// Example usage:
//
//	opts := reporter.DefaultOptions()
//	opts.Format = reporter.FormatJSON
//	reporter.Generate(baseEnv, targetEnv, opts)
package reporter
