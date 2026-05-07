// Package auditor inspects parsed .env maps for common issues such as
// empty values, placeholder strings (e.g. "changeme", "TODO"), and
// unquoted values that contain whitespace.
//
// Usage:
//
//	report := auditor.Audit(envMap)
//	if report.HasErrors() {
//		fmt.Println(report.Summary())
//	}
//
// Severity levels:
//   - "warn"  — suspicious but not necessarily broken
//   - "error" — likely misconfiguration that should block deployment
package auditor
