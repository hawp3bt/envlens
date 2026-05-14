// Package patcher provides functionality for applying declarative patch
// operations to an env map.
//
// Supported operations:
//
//   - set    – create or overwrite a key with a new value
//   - delete – remove a key from the map
//   - rename – move a key to a new name, preserving its value
//
// Patch never mutates the input map; it always returns a fresh copy.
// WritePatched serialises the resulting map back to a .env file.
package patcher
