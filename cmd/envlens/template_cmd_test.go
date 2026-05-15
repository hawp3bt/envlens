package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeTempTemplate(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.tmpl")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatal(err)
	}
	f.Close()
	return f.Name()
}

func TestTemplateCmd_RendersValue(t *testing.T) {
	env := writeTempEnv(t, "APP_NAME=envlens\nAPP_ENV=production\n")
	tmpl := writeTempTemplate(t, `name={{ env "APP_NAME" }} env={{ env "APP_ENV" }}`)

	out := filepath.Join(t.TempDir(), "out.txt")
	err := runTemplate([]string{"--file", env, "--tmpl", tmpl, "--out", out})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, _ := os.ReadFile(out)
	if !strings.Contains(string(data), "envlens") {
		t.Errorf("expected 'envlens' in output, got: %s", data)
	}
	if !strings.Contains(string(data), "production") {
		t.Errorf("expected 'production' in output, got: %s", data)
	}
}

func TestTemplateCmd_MissingFileFlag_ReturnsError(t *testing.T) {
	tmpl := writeTempTemplate(t, `hello`)
	err := runTemplate([]string{"--tmpl", tmpl})
	if err == nil || !strings.Contains(err.Error(), "--file") {
		t.Errorf("expected --file error, got: %v", err)
	}
}

func TestTemplateCmd_MissingTmplFlag_ReturnsError(t *testing.T) {
	env := writeTempEnv(t, "KEY=val\n")
	err := runTemplate([]string{"--file", env})
	if err == nil || !strings.Contains(err.Error(), "--tmpl") {
		t.Errorf("expected --tmpl error, got: %v", err)
	}
}

func TestTemplateCmd_NonExistentEnvFile_ReturnsError(t *testing.T) {
	tmpl := writeTempTemplate(t, `hello`)
	err := runTemplate([]string{"--file", "/no/such/file.env", "--tmpl", tmpl})
	if err == nil {
		t.Error("expected error for missing env file")
	}
}

func TestTemplateCmd_MissingKey_ZeroMode_WritesEmpty(t *testing.T) {
	env := writeTempEnv(t, "PRESENT=yes\n")
	tmpl := writeTempTemplate(t, `{{ env "ABSENT" }}`)
	out := filepath.Join(t.TempDir(), "out.txt")

	err := runTemplate([]string{"--file", env, "--tmpl", tmpl, "--out", out, "--missing", "zero"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, _ := os.ReadFile(out)
	if strings.TrimSpace(string(data)) != "" {
		t.Errorf("expected empty output for missing key, got: %s", data)
	}
}
