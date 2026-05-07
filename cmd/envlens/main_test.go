package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0600); err != nil {
		t.Fatalf("failed to write temp env file: %v", err)
	}
	return p
}

func buildBinary(t *testing.T) string {
	t.Helper()
	bin := filepath.Join(t.TempDir(), "envlens")
	cmd := exec.Command("go", "build", "-o", bin, ".")
	cmd.Dir = "."
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("build failed: %v\n%s", err, out)
	}
	return bin
}

func TestCLI_BasicDiff(t *testing.T) {
	bin := buildBinary(t)
	base := writeTempEnv(t, "APP_ENV=production\nDB_HOST=localhost\n")
	target := writeTempEnv(t, "APP_ENV=staging\nDB_HOST=localhost\nNEW_KEY=hello\n")

	cmd := exec.Command(bin, "--base", base, "--target", target)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("unexpected error: %v\n%s", err, out)
	}

	output := string(out)
	if !strings.Contains(output, "APP_ENV") {
		t.Errorf("expected APP_ENV in output, got:\n%s", output)
	}
	if !strings.Contains(output, "NEW_KEY") {
		t.Errorf("expected NEW_KEY in output, got:\n%s", output)
	}
}

func TestCLI_MissingFlags(t *testing.T) {
	bin := buildBinary(t)
	cmd := exec.Command(bin)
	out, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatal("expected non-zero exit when flags are missing")
	}
	if !strings.Contains(string(out), "--base and --target are required") {
		t.Errorf("expected usage error message, got:\n%s", out)
	}
}

func TestCLI_JSONFormat(t *testing.T) {
	bin := buildBinary(t)
	base := writeTempEnv(t, "SECRET_KEY=abc123\n")
	target := writeTempEnv(t, "SECRET_KEY=xyz789\n")

	cmd := exec.Command(bin, "--base", base, "--target", target, "--format", "json", "--mask=true")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("unexpected error: %v\n%s", err, out)
	}

	output := string(out)
	if !strings.Contains(output, "SECRET_KEY") {
		t.Errorf("expected SECRET_KEY in JSON output, got:\n%s", output)
	}
	if strings.Contains(output, "abc123") || strings.Contains(output, "xyz789") {
		t.Errorf("expected secret values to be masked, got:\n%s", output)
	}
}
