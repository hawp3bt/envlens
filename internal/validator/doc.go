// Package validator provides rule-based validation for .env key-value pairs.
//
// It supports the following rule types:
//
//   - nonempty: the value must not be blank or whitespace-only
//   - bool:     the value must be one of true, false, 1, or 0
//   - int:      the value must be a valid integer
//   - float:    the value must be a valid floating-point number
//   - url:      the value must be a well-formed absolute URL
//
// Rules are applied against a map[string]string (typically produced by the
// parser package) and any violations are returned as a slice of Finding values.
//
// A Rule may also carry an optional Required flag. When Required is true and
// the key is absent from the map entirely, a Finding is emitted regardless of
// the rule type.
//
// Example:
//
//	rules := []validator.Rule{
//		{Key: "PORT",    Type: "int",      Required: true},
//		{Key: "API_URL", Type: "url",      Required: true},
//		{Key: "DEBUG",   Type: "bool"},
//		{Key: "RATIO",   Type: "float"},
//	}
//	findings := validator.Validate(envMap, rules)
package validator
