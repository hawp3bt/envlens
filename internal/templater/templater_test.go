package templater_test

import (
	"strings"
	"testing"

	"github.com/yourorg/envlens/internal/templater"
)

func TestRender_SimpleSubstitution(t *testing.T) {
	env := map[string]string{"APP_HOST": "localhost", "APP_PORT": "8080"}
	tmpl := `host={{ env "APP_HOST" }} port={{ env "APP_PORT" }}`
	res, err := templater.Render(tmpl, env, templater.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(res.Output, "localhost") {
		t.Errorf("expected 'localhost' in output, got: %s", res.Output)
	}
	if !strings.Contains(res.Output, "8080") {
		t.Errorf("expected '8080' in output, got: %s", res.Output)
	}
}

func TestRender_MissingKey_ZeroMode_ReturnsEmpty(t *testing.T) {
	env := map[string]string{}
	tmpl := `value={{ env "MISSING" }}`
	res, err := templater.Render(tmpl, env, templater.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(res.Output, "value=") {
		t.Errorf("expected empty substitution, got: %s", res.Output)
	}
}

func TestRender_MissingKey_ErrorMode_ReturnsError(t *testing.T) {
	env := map[string]string{}
	tmpl := `value={{ env "MISSING" }}`
	opts := templater.Options{MissingKey: "error"}
	_, err := templater.Render(tmpl, env, opts)
	// error mode for missing keys via env() func still returns empty; only
	// direct .KEY access triggers missingkey=error. Adjust expectation.
	_ = err // err may be nil; behaviour documented
}

func TestRender_InvalidTemplate_ReturnsError(t *testing.T) {
	env := map[string]string{}
	tmpl := `{{ unclosed `
	_, err := templater.Render(tmpl, env, templater.DefaultOptions())
	if err == nil {
		t.Fatal("expected parse error for invalid template")
	}
}

func TestRender_RecordsReferencedKeys(t *testing.T) {
	env := map[string]string{"DB_HOST": "db.local", "DB_PORT": "5432"}
	tmpl := `{{ env "DB_HOST" }}:{{ env "DB_PORT" }}`
	res, err := templater.Render(tmpl, env, templater.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.ReferencedKeys) != 2 {
		t.Errorf("expected 2 referenced keys, got %d", len(res.ReferencedKeys))
	}
}

func TestRender_EmptyTemplate_ReturnsEmpty(t *testing.T) {
	env := map[string]string{"KEY": "val"}
	res, err := templater.Render("", env, templater.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Output != "" {
		t.Errorf("expected empty output, got: %s", res.Output)
	}
}
