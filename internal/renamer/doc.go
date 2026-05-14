// Package renamer provides utilities for renaming keys inside an env map
// according to a list of rules.
//
// Rules can target exact key names or key prefixes, making it straightforward
// to migrate naming conventions across environments (e.g. renaming APP_* to
// SVC_*) without manually editing each variable.
//
// Collision detection ensures that a rename operation never silently overwrites
// an existing key — an error is returned instead so the caller can decide how
// to resolve the conflict.
package renamer
