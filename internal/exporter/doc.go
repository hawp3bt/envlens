// Package exporter converts a parsed env map into various serialisation
// formats suitable for use in shell scripts, Docker Compose env-files,
// or JSON pipelines.
//
// Supported formats:
//
//   - shell  — emits "export KEY=VALUE" lines with proper shell quoting.
//   - docker — emits "KEY=VALUE" lines compatible with Docker --env-file.
//   - json   — emits a JSON array of {"key": ..., "value": ...} objects.
//
// Usage:
//
//	opts := exporter.DefaultOptions()
//	opts.Format = exporter.FormatDocker
//	err := exporter.Export(envMap, opts, os.Stdout)
package exporter
