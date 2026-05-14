package scanner_test

import (
	"strings"
	"testing"

	"github.com/yourorg/envlens/internal/scanner"
)

func TestScanLines_NoDuplicates_ReturnsEmpty(t *testing.T) {
	lines := []string{
		"APP_ENV=production",
		"DB_HOST=localhost",
		"DB_PORT=5432",
	}
	result := scanner.ScanLines(lines)
	if len(result.Findings) != 0 {
		t.Fatalf("expected no findings, got %d", len(result.Findings))
	}
}

func TestScanLines_DuplicateKey_ReturnsError(t *testing.T) {
	lines := []string{
		"APP_ENV=staging",
		"DB_HOST=localhost",
		"APP_ENV=production",
	}
	result := scanner.ScanLines(lines)
	if len(result.Findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(result.Findings))
	}
	f := result.Findings[0]
	if f.Key != "APP_ENV" {
		t.Errorf("expected key APP_ENV, got %q", f.Key)
	}
	if f.Severity != "error" {
		t.Errorf("expected severity error, got %q", f.Severity)
	}
	if !strings.Contains(f.Message, "line 1") || !strings.Contains(f.Message, "line 3") {
		t.Errorf("message should reference both line numbers, got: %s", f.Message)
	}
}

func TestScanLines_SkipsComments(t *testing.T) {
	lines := []string{
		"# APP_ENV=staging",
		"APP_ENV=production",
	}
	result := scanner.ScanLines(lines)
	if len(result.Findings) != 0 {
		t.Fatalf("expected no findings, got %d", len(result.Findings))
	}
}

func TestScanLines_SkipsBlankLines(t *testing.T) {
	lines := []string{
		"",
		"   ",
		"KEY=value",
	}
	result := scanner.ScanLines(lines)
	if len(result.Findings) != 0 {
		t.Fatalf("expected no findings, got %d", len(result.Findings))
	}
}

func TestScanMaps_NoCaseConflict_ReturnsEmpty(t *testing.T) {
	a := map[string]string{"APP_ENV": "prod", "DB_HOST": "localhost"}
	b := map[string]string{"APP_ENV": "staging", "DB_PORT": "5432"}
	result := scanner.ScanMaps(a, b)
	if len(result.Findings) != 0 {
		t.Fatalf("expected no findings, got %d", len(result.Findings))
	}
}

func TestScanMaps_CaseConflict_ReturnsWarning(t *testing.T) {
	a := map[string]string{"DB_HOST": "localhost"}
	b := map[string]string{"db_host": "remotehost"}
	result := scanner.ScanMaps(a, b)
	if len(result.Findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(result.Findings))
	}
	if result.Findings[0].Severity != "warning" {
		t.Errorf("expected warning severity, got %q", result.Findings[0].Severity)
	}
}

func TestResult_HasErrors_TrueWhenErrorPresent(t *testing.T) {
	r := scanner.Result{Findings: []scanner.Finding{{Key: "X", Severity: "error", Message: "dup"}}}
	if !r.HasErrors() {
		t.Error("expected HasErrors to return true")
	}
}

func TestResult_Summary_NoIssues(t *testing.T) {
	r := scanner.Result{}
	if r.Summary() != "no issues found" {
		t.Errorf("unexpected summary: %s", r.Summary())
	}
}

func TestResult_Summary_WithFindings(t *testing.T) {
	r := scanner.Result{Findings: []scanner.Finding{
		{Key: "A", Severity: "error", Message: "dup"},
		{Key: "B", Severity: "warning", Message: "case"},
	}}
	if !strings.Contains(r.Summary(), "2") {
		t.Errorf("expected summary to mention 2 issues, got: %s", r.Summary())
	}
}
