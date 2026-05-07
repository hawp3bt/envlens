package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSnapshotCmd_SavesSnapshot(t *testing.T) {
	envContent := "APP_ENV=production\nSECRET_KEY=supersecret\nPORT=8080\n"
	envFile := writeTempEnv(t, envContent)
	snapshotDir := filepath.Join(t.TempDir(), "snaps")

	err := runSnapshot([]string{
		"--file", envFile,
		"--env", "prod",
		"--dir", snapshotDir,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	entries, _ := os.ReadDir(snapshotDir)
	if len(entries) != 1 {
		t.Errorf("expected 1 snapshot file, got %d", len(entries))
	}
}

func TestSnapshotCmd_MissingFlags_ReturnsError(t *testing.T) {
	err := runSnapshot([]string{"--file", "some.env"})
	if err == nil {
		t.Error("expected error when --env is missing")
	}
}

func TestSnapshotCmd_ListSnapshots(t *testing.T) {
	envContent := "KEY=value\n"
	envFile := writeTempEnv(t, envContent)
	snapshotDir := filepath.Join(t.TempDir(), "snaps")

	_ = runSnapshot([]string{"--file", envFile, "--env", "staging", "--dir", snapshotDir})
	_ = runSnapshot([]string{"--file", envFile, "--env", "staging", "--dir", snapshotDir})

	// Capture stdout by redirecting
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := runSnapshot([]string{"--list", "staging", "--dir", snapshotDir})

	w.Close()
	os.Stdout = old

	if err != nil {
		t.Fatalf("list failed: %v", err)
	}

	buf := make([]byte, 1024)
	n, _ := r.Read(buf)
	output := string(buf[:n])

	if strings.Count(output, "staging_") < 2 {
		t.Errorf("expected at least 2 snapshot lines, got: %s", output)
	}
}

func TestSnapshotCmd_ListEmpty(t *testing.T) {
	snapshotDir := t.TempDir()

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := runSnapshot([]string{"--list", "ghost", "--dir", snapshotDir})

	w.Close()
	os.Stdout = old

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	buf := make([]byte, 256)
	n, _ := r.Read(buf)
	if !strings.Contains(string(buf[:n]), "no snapshots") {
		t.Errorf("expected 'no snapshots' message")
	}
}
