package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestExportCmd_ShellFormat_ContainsExport(t *testing.T) {
	bin := buildBinary(t)
	f := writeTempEnv(t, "APP_NAME=envlens\nPORT=9000\n")

	out, err := exec.Command(bin, "export", "--file", f).CombinedOutput()
	if err != nil {
		t.Fatalf("command failed: %v\n%s", err, out)
	}
	if !strings.Contains(string(out), "export APP_NAME=envlens") {
		t.Errorf("expected shell export line, got:\n%s", out)
	}
}

func TestExportCmd_DockerFormat_NoExportKeyword(t *testing.T) {
	bin := buildBinary(t)
	f := writeTempEnv(t, "DB_HOST=localhost\nDB_PORT=5432\n")

	out, err := exec.Command(bin, "export", "--file", f, "--format", "docker").CombinedOutput()
	if err != nil {
		t.Fatalf("command failed: %v\n%s", err, out)
	}
	s := string(out)
	if strings.Contains(s, "export ") {
		t.Errorf("docker format must not contain 'export', got:\n%s", s)
	}
	if !strings.Contains(s, "DB_HOST=localhost") {
		t.Errorf("expected DB_HOST=localhost, got:\n%s", s)
	}
}

func TestExportCmd_JSONFormat_IsValidStructure(t *testing.T) {
	bin := buildBinary(t)
	f := writeTempEnv(t, "REGION=us-east-1\n")

	out, err := exec.Command(bin, "export", "--file", f, "--format", "json").CombinedOutput()
	if err != nil {
		t.Fatalf("command failed: %v\n%s", err, out)
	}
	s := string(out)
	if !strings.Contains(s, "\"key\"") || !strings.Contains(s, "REGION") {
		t.Errorf("expected JSON output with key field, got:\n%s", s)
	}
}

func TestExportCmd_MissingFileFlag_ReturnsError(t *testing.T) {
	bin := buildBinary(t)
	out, err := exec.Command(bin, "export").CombinedOutput()
	if err == nil {
		t.Fatalf("expected error when --file is missing, got nil\n%s", out)
	}
	if !strings.Contains(string(out), "--file") {
		t.Errorf("error message should mention --file, got:\n%s", out)
	}
}

func TestExportCmd_NonExistentFile_ReturnsError(t *testing.T) {
	bin := buildBinary(t)
	ghost := filepath.Join(os.TempDir(), "ghost_does_not_exist.env")
	out, err := exec.Command(bin, "export", "--file", ghost).CombinedOutput()
	if err == nil {
		t.Fatalf("expected error for missing file, got nil\n%s", out)
	}
}
