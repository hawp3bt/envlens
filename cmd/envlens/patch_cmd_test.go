package main_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestPatchCmd_SetKey(t *testing.T) {
	bin := buildBinary(t)
	dir := t.TempDir()
	env := filepath.Join(dir, "test.env")
	writeTempEnv(t, env, "APP_PORT=8080\nAPP_NAME=myapp\n")

	out, err := runCmd(t, bin, "patch",
		"--file", env,
		"--ops", `[{"action":"set","key":"APP_PORT","value":"9090"}]`,
	)
	if err != nil {
		t.Fatalf("unexpected error: %v\noutput: %s", err, out)
	}
	data, _ := os.ReadFile(env)
	if !strings.Contains(string(data), "APP_PORT=9090") {
		t.Errorf("expected patched value in file, got:\n%s", data)
	}
}

func TestPatchCmd_DeleteKey(t *testing.T) {
	bin := buildBinary(t)
	dir := t.TempDir()
	env := filepath.Join(dir, "test.env")
	writeTempEnv(t, env, "REMOVE_ME=yes\nKEEP=no\n")

	_, err := runCmd(t, bin, "patch",
		"--file", env,
		"--ops", `[{"action":"delete","key":"REMOVE_ME"}]`,
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, _ := os.ReadFile(env)
	if strings.Contains(string(data), "REMOVE_ME") {
		t.Errorf("expected key deleted, got:\n%s", data)
	}
}

func TestPatchCmd_DryRun_DoesNotWrite(t *testing.T) {
	bin := buildBinary(t)
	dir := t.TempDir()
	env := filepath.Join(dir, "test.env")
	original := "FOO=original\n"
	writeTempEnv(t, env, original)

	_, err := runCmd(t, bin, "patch",
		"--file", env,
		"--ops", `[{"action":"set","key":"FOO","value":"changed"}]`,
		"--dry-run",
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, _ := os.ReadFile(env)
	if string(data) != original {
		t.Errorf("dry-run should not modify file, got:\n%s", data)
	}
}

func TestPatchCmd_MissingFileFlag_ReturnsError(t *testing.T) {
	bin := buildBinary(t)
	_, err := runCmd(t, bin, "patch", "--ops", `[]`)
	if err == nil {
		t.Error("expected error when --file is missing")
	}
}

func TestPatchCmd_InvalidOpsJSON_ReturnsError(t *testing.T) {
	bin := buildBinary(t)
	dir := t.TempDir()
	env := filepath.Join(dir, "test.env")
	writeTempEnv(t, env, "X=1\n")

	_, err := runCmd(t, bin, "patch",
		"--file", env,
		"--ops", `not-json`,
	)
	if err == nil {
		t.Error("expected error for invalid JSON ops")
	}
}
