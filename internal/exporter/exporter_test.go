package exporter_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/yourorg/envlens/internal/exporter"
)

var sampleEnv = map[string]string{
	"APP_NAME": "envlens",
	"DB_PASS":  "s3cr3t word",
	"PORT":     "8080",
}

func TestExport_ShellFormat_ContainsExport(t *testing.T) {
	var buf bytes.Buffer
	opts := exporter.DefaultOptions()
	if err := exporter.Export(sampleEnv, opts, &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "export APP_NAME=envlens") {
		t.Errorf("expected shell export line, got:\n%s", out)
	}
}

func TestExport_ShellFormat_EscapesSpaces(t *testing.T) {
	var buf bytes.Buffer
	opts := exporter.DefaultOptions()
	if err := exporter.Export(sampleEnv, opts, &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "export DB_PASS='s3cr3t word'") {
		t.Errorf("expected quoted value, got:\n%s", out)
	}
}

func TestExport_DockerFormat_NoExportKeyword(t *testing.T) {
	var buf bytes.Buffer
	opts := exporter.Options{Format: exporter.FormatDocker, Sorted: true}
	if err := exporter.Export(sampleEnv, opts, &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if strings.Contains(out, "export ") {
		t.Errorf("docker format should not contain 'export', got:\n%s", out)
	}
	if !strings.Contains(out, "PORT=8080") {
		t.Errorf("expected PORT=8080 in docker output, got:\n%s", out)
	}
}

func TestExport_JSONFormat_ContainsKeys(t *testing.T) {
	var buf bytes.Buffer
	opts := exporter.Options{Format: exporter.FormatJSON, Sorted: true}
	if err := exporter.Export(sampleEnv, opts, &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "\"key\"") || !strings.Contains(out, "\"value\"") {
		t.Errorf("expected JSON with key/value fields, got:\n%s", out)
	}
}

func TestExport_UnknownFormat_ReturnsError(t *testing.T) {
	var buf bytes.Buffer
	opts := exporter.Options{Format: "xml"}
	if err := exporter.Export(sampleEnv, opts, &buf); err == nil {
		t.Error("expected error for unknown format, got nil")
	}
}

func TestExport_SortedOutput_IsAlphabetical(t *testing.T) {
	var buf bytes.Buffer
	opts := exporter.Options{Format: exporter.FormatDocker, Sorted: true}
	env := map[string]string{"Z_KEY": "z", "A_KEY": "a", "M_KEY": "m"}
	if err := exporter.Export(env, opts, &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 3 || !strings.HasPrefix(lines[0], "A_KEY") {
		t.Errorf("expected sorted output starting with A_KEY, got: %v", lines)
	}
}
