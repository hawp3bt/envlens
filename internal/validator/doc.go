// Package validator provides rule-based validation for .env key-value pairs.
//
// It supports the following rule types:
//
//   - nonempty: the value must not be blank or whitespace-only
//   - bool:     the value must be one of true, false, 1, or 0
//   - int:      the value must be a valid integer
//   - url:      the value must be a well-formed absolute URL
//
// Rules are applied against a map[string]string (typically produced by the
// parser package) and any violations are returned as a slice of Finding values.
//
// Example:
//
//	rules := []validator.Rule{
//		{Key: "PORT",    Type: "int"},
//		{Key: "API_URL", Type: "url"},
//		{Key: "DEBUG",   Type: "bool"},
//	}
//	findings := validator.Validate(envMap, rules)
package validator
