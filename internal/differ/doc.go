// Package differ provides functionality for computing and displaying the
// difference between two parsed .env file maps.
//
// Usage:
//
//	base, _ := parser.ParseFile("staging.env")
//	prod, _  := parser.ParseFile("production.env")
//
//	result := differ.Diff(base.ToMap(), prod.ToMap())
//	if result.HasChanges() {
//		differ.Format(os.Stdout, result, differ.FormatOptions{
//			Colorized:     true,
//			ShowUnchanged: false,
//		})
//	}
//
// The diff result categorises each key as one of: Added, Removed, Changed,
// or Unchanged. The formatter renders the result with optional ANSI colour
// highlighting suitable for terminal output.
package differ
