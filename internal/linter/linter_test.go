package linter_test

import (
	"strings"
	"testing"

	"github.com/yourorg/envlens/internal/linter"
)

func TestLint_ValidKeys_NoFindings(t *testing.T) {
	env := map[string]string{
		"DATABASE_URL": "postgres://localhost/db",
		"PORT":         "8080",
		"APP_ENV":      "production",
	}
	findings := linter.Lint(env)
	if len(findings) != 0 {
		t.Errorf("expected no findings, got %d: %v", len(findings), findings)
	}
}

func TestLint_LowercaseKey_ReturnsWarning(t *testing.T) {
	env := map[string]string{"database_url": "value"}
	findings := linter.Lint(env)
	if len(findings) == 0 {
		t.Fatal("expected a finding for lowercase key")
	}
	if findings[0].Severity != linter.SeverityWarning {
		t.Errorf("expected warning severity, got %s", findings[0].Severity)
	}
}

func TestLint_MixedCaseKey_ReturnsWarning(t *testing.T) {
	env := map[string]string{"AppEnv": "value"}
	findings := linter.Lint(env)
	if len(findings) == 0 {
		t.Fatal("expected a finding for mixed-case key")
	}
}

func TestLint_KeyStartsWithDigit_ReturnsError(t *testing.T) {
	env := map[string]string{"1_BAD_KEY": "value"}
	findings := linter.Lint(env)
	var hasError bool
	for _, f := range findings {
		if f.Severity == linter.SeverityError {
			hasError = true
		}
	}
	if !hasError {
		t.Error("expected an error finding for key starting with digit")
	}
}

func TestLint_DoubleUnderscore_ReturnsWarning(t *testing.T) {
	env := map[string]string{"DB__HOST": "value"}
	findings := linter.Lint(env)
	if len(findings) == 0 {
		t.Fatal("expected a finding for double underscore")
	}
	if findings[0].Severity != linter.SeverityWarning {
		t.Errorf("expected warning, got %s", findings[0].Severity)
	}
}

func TestLint_LeadingUnderscore_ReturnsWarning(t *testing.T) {
	env := map[string]string{"_PRIVATE": "value"}
	findings := linter.Lint(env)
	if len(findings) == 0 {
		t.Fatal("expected a finding for leading underscore")
	}
}

func TestSummary_NoFindings(t *testing.T) {
	out := linter.Summary(nil)
	if out != "lint: no issues found" {
		t.Errorf("unexpected summary: %q", out)
	}
}

func TestSummary_WithFindings_ContainsKeyAndSeverity(t *testing.T) {
	findings := []linter.Finding{
		{Key: "bad_key", Message: "key should be UPPER_SNAKE_CASE", Severity: linter.SeverityWarning},
	}
	out := linter.Summary(findings)
	if !strings.Contains(out, "bad_key") {
		t.Errorf("expected key in summary, got: %s", out)
	}
	if !strings.Contains(out, string(linter.SeverityWarning)) {
		t.Errorf("expected severity in summary, got: %s", out)
	}
}
