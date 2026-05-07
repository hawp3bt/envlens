package snapshot_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/yourorg/envlens/internal/snapshot"
)

func TestSave_CreatesFile(t *testing.T) {
	dir := t.TempDir()
	entries := map[string]string{"APP_ENV": "production", "PORT": "8080"}

	path, err := snapshot.Save(dir, "prod", entries)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Errorf("expected file to exist at %s", path)
	}

	if !strings.Contains(filepath.Base(path), "prod_") {
		t.Errorf("expected filename to contain environment prefix, got %s", path)
	}
}

func TestLoad_ReturnsCorrectData(t *testing.T) {
	dir := t.TempDir()
	entries := map[string]string{"DB_HOST": "localhost", "DB_PORT": "5432"}

	path, err := snapshot.Save(dir, "staging", entries)
	if err != nil {
		t.Fatalf("save failed: %v", err)
	}

	snap, err := snapshot.Load(path)
	if err != nil {
		t.Fatalf("load failed: %v", err)
	}

	if snap.Environment != "staging" {
		t.Errorf("expected environment 'staging', got %q", snap.Environment)
	}
	if snap.Entries["DB_HOST"] != "localhost" {
		t.Errorf("expected DB_HOST=localhost, got %q", snap.Entries["DB_HOST"])
	}
	if snap.CapturedAt.IsZero() {
		t.Error("expected non-zero CapturedAt")
	}
}

func TestList_ReturnsMatchingFiles(t *testing.T) {
	dir := t.TempDir()

	_, _ = snapshot.Save(dir, "prod", map[string]string{"A": "1"})
	_, _ = snapshot.Save(dir, "prod", map[string]string{"A": "2"})
	_, _ = snapshot.Save(dir, "dev", map[string]string{"A": "3"})

	files, err := snapshot.List(dir, "prod")
	if err != nil {
		t.Fatalf("list failed: %v", err)
	}

	if len(files) != 2 {
		t.Errorf("expected 2 prod snapshots, got %d", len(files))
	}
}

func TestLoad_InvalidPath_ReturnsError(t *testing.T) {
	_, err := snapshot.Load("/nonexistent/path/snap.json")
	if err == nil {
		t.Error("expected error for missing file")
	}
}
