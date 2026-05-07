package auditor_test

import (
	"strings"
	"testing"

	"github.com/yourorg/envlens/internal/auditor"
)

func TestAudit_EmptyValue(t *testing.T) {
	env := map[string]string{"DB_HOST": ""}
	report := auditor.Audit(env)
	if len(report.Issues) == 0 {
		t.Fatal("expected at least one issue for empty value")
	}
	if report.Issues[0].Severity != "warn" {
		t.Errorf("expected warn severity, got %s", report.Issues[0].Severity)
	}
}

func TestAudit_PlaceholderValue(t *testing.T) {
	env := map[string]string{"API_KEY": "changeme"}
	report := auditor.Audit(env)
	if !report.HasErrors() {
		t.Fatal("expected error for placeholder value")
	}
}

func TestAudit_UnquotedSpaces(t *testing.T) {
	env := map[string]string{"APP_NAME": "my app"}
	report := auditor.Audit(env)
	if len(report.Issues) == 0 {
		t.Fatal("expected issue for unquoted spaces")
	}
	if report.Issues[0].Severity != "warn" {
		t.Errorf("expected warn, got %s", report.Issues[0].Severity)
	}
}

func TestAudit_CleanEnv(t *testing.T) {
	env := map[string]string{
		"DB_HOST": "localhost",
		"PORT":    "5432",
	}
	report := auditor.Audit(env)
	if len(report.Issues) != 0 {
		t.Errorf("expected no issues, got %d", len(report.Issues))
	}
}

func TestReport_Summary_NoIssues(t *testing.T) {
	report := &auditor.Report{}
	if report.Summary() != "No issues found." {
		t.Errorf("unexpected summary: %s", report.Summary())
	}
}

func TestReport_Summary_WithIssues(t *testing.T) {
	report := &auditor.Report{
		Issues: []auditor.Issue{
			{Key: "FOO", Severity: "warn", Message: "empty value"},
		},
	}
	summary := report.Summary()
	if !strings.Contains(summary, "FOO") {
		t.Errorf("summary should contain key name, got: %s", summary)
	}
	if !strings.Contains(summary, "WARN") {
		t.Errorf("summary should contain severity, got: %s", summary)
	}
}

func TestAudit_TodoPlaceholder(t *testing.T) {
	env := map[string]string{"SECRET": "TODO"}
	report := auditor.Audit(env)
	if !report.HasErrors() {
		t.Fatal("expected error for TODO placeholder")
	}
}
