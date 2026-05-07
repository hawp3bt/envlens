package snapshot

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Snapshot represents a saved state of an environment file.
type Snapshot struct {
	Environment string            `json:"environment"`
	CapturedAt  time.Time         `json:"captured_at"`
	Entries     map[string]string `json:"entries"`
}

// Save writes a snapshot to the given directory as a JSON file.
// The filename is derived from the environment name and timestamp.
func Save(dir, environment string, entries map[string]string) (string, error) {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("snapshot: create dir: %w", err)
	}

	snap := Snapshot{
		Environment: environment,
		CapturedAt:  time.Now().UTC(),
		Entries:     entries,
	}

	filename := fmt.Sprintf("%s_%s.json", environment, snap.CapturedAt.Format("20060102T150405Z"))
	path := filepath.Join(dir, filename)

	f, err := os.Create(path)
	if err != nil {
		return "", fmt.Errorf("snapshot: create file: %w", err)
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(snap); err != nil {
		return "", fmt.Errorf("snapshot: encode: %w", err)
	}

	return path, nil
}

// Load reads a snapshot from the given file path.
func Load(path string) (*Snapshot, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("snapshot: open: %w", err)
	}
	defer f.Close()

	var snap Snapshot
	if err := json.NewDecoder(f).Decode(&snap); err != nil {
		return nil, fmt.Errorf("snapshot: decode: %w", err)
	}
	return &snap, nil
}

// List returns all snapshot file paths in the given directory for an environment.
func List(dir, environment string) ([]string, error) {
	pattern := filepath.Join(dir, environment+"_*.json")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return nil, fmt.Errorf("snapshot: list: %w", err)
	}
	return matches, nil
}
