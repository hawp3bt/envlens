package parser

import (
	"os"
	"testing"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	tmp, err := os.CreateTemp("", "*.env")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	if _, err := tmp.WriteString(content); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	tmp.Close()
	t.Cleanup(func() { os.Remove(tmp.Name()) })
	return tmp.Name()
}

func TestParseFile_BasicEntries(t *testing.T) {
	path := writeTempEnv(t, "APP_ENV=production\nDB_HOST=localhost\nDB_PORT=5432\n")

	envFile, err := ParseFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(envFile.Entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(envFile.Entries))
	}

	if envFile.Entries[0].Key != "APP_ENV" || envFile.Entries[0].Value != "production" {
		t.Errorf("unexpected first entry: %+v", envFile.Entries[0])
	}
}

func TestParseFile_SkipsComments(t *testing.T) {
	path := writeTempEnv(t, "# This is a comment\nKEY=value\n")

	envFile, err := ParseFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(envFile.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(envFile.Entries))
	}
}

func TestParseFile_QuotedValues(t *testing.T) {
	path := writeTempEnv(t, `SECRET="my_secret_value"` + "\n")

	envFile, err := ParseFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if envFile.Entries[0].Value != "my_secret_value" {
		t.Errorf("expected unquoted value, got %q", envFile.Entries[0].Value)
	}
}

func TestParseFile_ToMap(t *testing.T) {
	path := writeTempEnv(t, "FOO=bar\nBAZ=qux\n")

	envFile, err := ParseFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	m := envFile.ToMap()
	if m["FOO"] != "bar" || m["BAZ"] != "qux" {
		t.Errorf("unexpected map contents: %v", m)
	}
}

func TestParseFile_NotFound(t *testing.T) {
	_, err := ParseFile("/nonexistent/.env")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}
