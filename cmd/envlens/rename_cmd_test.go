package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRenameCmd_ExactKey_RenamesInOutput(t *testing.T) {
	dir := t.TempDir()
	env := writeTempEnv(t, dir, "OLD_KEY=hello\nSAFE=world\n")
	out := filepath.Join(dir, "out.env")

	err := runRename([]string{
		"-file", env,
		"-rules", `[{"from":"OLD_KEY","to":"NEW_KEY"}]`,
		"-out", out,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, _ := os.ReadFile(out)
	if !strings.Contains(string(data), "NEW_KEY=hello") {
		t.Errorf("expected NEW_KEY=hello in output, got:\n%s", data)
	}
	if strings.Contains(string(data), "OLD_KEY") {
		t.Errorf("OLD_KEY should not appear in output")
	}
}

func TestRenameCmd_DryRun_DoesNotWrite(t *testing.T) {
	dir := t.TempDir()
	env := writeTempEnv(t, dir, "APP_HOST=localhost\n")
	out := filepath.Join(dir, "out.env")

	err := runRename([]string{
		"-file", env,
		"-rules", `[{"from":"APP_HOST","to":"SVC_HOST"}]`,
		"-out", out,
		"-dry-run",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := os.Stat(out); !os.IsNotExist(err) {
		t.Error("dry-run should not create output file")
	}
}

func TestRenameCmd_MissingFlags_ReturnsError(t *testing.T) {
	err := runRename([]string{"-file", "some.env"})
	if err == nil {
		t.Fatal("expected error when -rules flag is missing")
	}
}

func TestRenameCmd_InvalidRulesJSON_ReturnsError(t *testing.T) {
	dir := t.TempDir()
	env := writeTempEnv(t, dir, "KEY=val\n")
	err := runRename([]string{
		"-file", env,
		"-rules", `not-json`,
	})
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestRenameCmd_PrefixRule_RenamesMatchingKeys(t *testing.T) {
	dir := t.TempDir()
	env := writeTempEnv(t, dir, "APP_HOST=h\nAPP_PORT=9\nOTHER=x\n")
	out := filepath.Join(dir, "out.env")

	err := runRename([]string{
		"-file", env,
		"-rules", `[{"from":"APP_","to":"SVC_","is_prefix":true}]`,
		"-out", out,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, _ := os.ReadFile(out)
	if !strings.Contains(string(data), "SVC_HOST=h") {
		t.Errorf("expected SVC_HOST in output, got:\n%s", data)
	}
	if !strings.Contains(string(data), "SVC_PORT=9") {
		t.Errorf("expected SVC_PORT in output, got:\n%s", data)
	}
	if !strings.Contains(string(data), "OTHER=x") {
		t.Errorf("expected OTHER unchanged in output, got:\n%s", data)
	}
}
