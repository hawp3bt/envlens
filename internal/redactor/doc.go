// Package redactor provides a Redactor type that replaces sensitive
// environment variable values with a configurable placeholder before
// the data is surfaced in reports, exports, or log output.
//
// Sensitivity is determined by the masker package heuristics (key
// substrings such as PASSWORD, SECRET, TOKEN, KEY, etc.) and can be
// extended with caller-supplied extra key substrings via Options.
//
// Basic usage:
//
//	r := redactor.New(redactor.DefaultOptions())
//	safe := r.Redact(envMap)
//
// To extend with custom sensitive key patterns:
//
//	opts := redactor.DefaultOptions()
//	opts.ExtraKeys = []string{"INTERNAL", "PRIV"}
//	r := redactor.New(opts)
package redactor
