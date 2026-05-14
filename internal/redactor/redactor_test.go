package redactor_test

import (
	"testing"

	"github.com/yourorg/envlens/internal/redactor"
)

func TestRedact_SensitiveKey_IsReplaced(t *testing.T) {
	r := redactor.New(redactor.DefaultOptions())
	env := map[string]string{"DB_PASSWORD": "s3cr3t"}
	out := r.Redact(env)
	if out["DB_PASSWORD"] != "[REDACTED]" {
		t.Errorf("expected [REDACTED], got %q", out["DB_PASSWORD"])
	}
}

func TestRedact_SafeKey_IsPreserved(t *testing.T) {
	r := redactor.New(redactor.DefaultOptions())
	env := map[string]string{"APP_ENV": "production"}
	out := r.Redact(env)
	if out["APP_ENV"] != "production" {
		t.Errorf("expected 'production', got %q", out["APP_ENV"])
	}
}

func TestRedact_CustomReplacement(t *testing.T) {
	opts := redactor.Options{Replacement: "***"}
	r := redactor.New(opts)
	env := map[string]string{"API_SECRET": "abc123"}
	out := r.Redact(env)
	if out["API_SECRET"] != "***" {
		t.Errorf("expected ***, got %q", out["API_SECRET"])
	}
}

func TestRedact_ExtraKeys_AreTreatedAsSensitive(t *testing.T) {
	opts := redactor.DefaultOptions()
	opts.ExtraKeys = []string{"INTERNAL"}
	r := redactor.New(opts)
	env := map[string]string{"INTERNAL_HOST": "10.0.0.1", "PORT": "8080"}
	out := r.Redact(env)
	if out["INTERNAL_HOST"] != "[REDACTED]" {
		t.Errorf("expected [REDACTED], got %q", out["INTERNAL_HOST"])
	}
	if out["PORT"] != "8080" {
		t.Errorf("expected 8080, got %q", out["PORT"])
	}
}

func TestRedactValue_SensitiveKey_ReturnsReplacement(t *testing.T) {
	r := redactor.New(redactor.DefaultOptions())
	result := r.RedactValue("DB_TOKEN", "mysecret")
	if result != "[REDACTED]" {
		t.Errorf("expected [REDACTED], got %q", result)
	}
}

func TestRedactValue_SafeKey_ReturnsOriginal(t *testing.T) {
	r := redactor.New(redactor.DefaultOptions())
	result := r.RedactValue("LOG_LEVEL", "debug")
	if result != "debug" {
		t.Errorf("expected debug, got %q", result)
	}
}

func TestDefaultOptions_HasSaneDefaults(t *testing.T) {
	opts := redactor.DefaultOptions()
	if opts.Replacement != "[REDACTED]" {
		t.Errorf("unexpected default replacement: %q", opts.Replacement)
	}
	if len(opts.ExtraKeys) != 0 {
		t.Errorf("expected no extra keys by default")
	}
}

func TestRedact_EmptyMap_ReturnsEmptyMap(t *testing.T) {
	r := redactor.New(redactor.DefaultOptions())
	out := r.Redact(map[string]string{})
	if len(out) != 0 {
		t.Errorf("expected empty map, got %d entries", len(out))
	}
}
