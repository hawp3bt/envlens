package reporter_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/yourorg/envlens/internal/reporter"
)

func TestGenerate_TextFormat_ShowsAddedKey(t *testing.T) {
	base := map[string]string{"PORT": "8080"}
	target := map[string]string{"PORT": "8080", "NEW_KEY": "value"}
	var buf bytes.Buffer
	opts := reporter.DefaultOptions()
	opts.Out = &buf
	opts.MaskSecrets = false
	if err := reporter.Generate(base, target, opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "NEW_KEY") {
		t.Errorf("expected NEW_KEY in output, got:\n%s", buf.String())
	}
}

func TestGenerate_TextFormat_MasksSecrets(t *testing.T) {
	base := map[string]string{}
	target := map[string]string{"API_SECRET": "super-secret-value"}
	var buf bytes.Buffer
	opts := reporter.DefaultOptions()
	opts.Out = &buf
	opts.MaskSecrets = true
	_ = reporter.Generate(base, target, opts)
	if strings.Contains(buf.String(), "super-secret-value") {
		t.Error("secret value should be masked in output")
	}
}

func TestGenerate_JSONFormat_ContainsKey(t *testing.T) {
	base := map[string]string{}
	target := map[string]string{"HOST": "localhost"}
	var buf bytes.Buffer
	opts := reporter.DefaultOptions()
	opts.Format = "json"
	opts.Out = &buf
	opts.MaskSecrets = false
	_ = reporter.Generate(base, target, opts)
	if !strings.Contains(buf.String(), "HOST") {
		t.Errorf("expected HOST in JSON output, got:\n%s", buf.String())
	}
}

func TestGenerate_JSONFormat_HidesUnchangedByDefault(t *testing.T) {
	base := map[string]string{"PORT": "8080"}
	target := map[string]string{"PORT": "8080"}
	var buf bytes.Buffer
	opts := reporter.DefaultOptions()
	opts.Format = "json"
	opts.Out = &buf
	opts.MaskSecrets = false
	_ = reporter.Generate(base, target, opts)
	var out map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	diff := out["diff"].([]interface{})
	if len(diff) != 0 {
		t.Errorf("expected empty diff for identical envs, got %d entries", len(diff))
	}
}

func TestDefaultOptions_HasSaneDefaults(t *testing.T) {
	opts := reporter.DefaultOptions()
	if opts.Format != "text" {
		t.Errorf("expected text format, got %s", opts.Format)
	}
	if !opts.MaskSecrets {
		t.Error("MaskSecrets should be true by default")
	}
	if opts.ShowSame {
		t.Error("ShowSame should be false by default")
	}
}

func TestGenerate_WithAudit_ShowsIssues(t *testing.T) {
	base := map[string]string{}
	target := map[string]string{"DB_PASS": "changeme"}
	var buf bytes.Buffer
	opts := reporter.DefaultOptions()
	opts.Out = &buf
	opts.MaskSecrets = false
	opts.Audit = true
	_ = reporter.Generate(base, target, opts)
	if !strings.Contains(buf.String(), "Audit Issues") {
		t.Errorf("expected audit section in output, got:\n%s", buf.String())
	}
}
