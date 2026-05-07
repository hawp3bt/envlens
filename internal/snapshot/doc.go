// Package snapshot provides functionality for saving and loading point-in-time
// captures of environment variable maps.
//
// Snapshots are stored as JSON files in a configurable directory, named by
// environment and timestamp. They can be loaded later and diffed against
// current or other snapshots to track environment drift over time.
//
// Example usage:
//
//	path, err := snapshot.Save(".envlens/snapshots", "production", entries)
//	snap, err := snapshot.Load(path)
//	files, err := snapshot.List(".envlens/snapshots", "production")
package snapshot
