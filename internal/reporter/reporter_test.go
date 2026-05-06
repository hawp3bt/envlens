package reporter_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/envlens/internal/reporter"
)

func TestGenerate_TextFormat_ShowsAddedKey(t *testing.T) {
	base := map[string]string{"APP_NAME": "myapp"}
	target := map[string]string{"APP_NAME": "myapp", "NEW_KEY": "value"}

	var buf bytes.Buffer
	opts := reporter.DefaultOptions()
	opts.Output = &buf
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
	target := map[string]string{"DB_PASSWORD": "supersecret"}

	var buf bytes.Buffer
	opts := reporter.DefaultOptions()
	opts.Output = &buf
	opts.MaskSecrets = true

	if err := reporter.Generate(base, target, opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if strings.Contains(buf.String(), "supersecret") {
		t.Errorf("expected secret to be masked, got:\n%s", buf.String())
	}
}

func TestGenerate_JSONFormat_ContainsKey(t *testing.T) {
	base := map[string]string{"HOST": "localhost"}
	target := map[string]string{"HOST": "prod.example.com"}

	var buf bytes.Buffer
	opts := reporter.DefaultOptions()
	opts.Format = reporter.FormatJSON
	opts.Output = &buf
	opts.MaskSecrets = false

	if err := reporter.Generate(base, target, opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, `"HOST"`) {
		t.Errorf("expected HOST key in JSON output, got:\n%s", out)
	}
	if !strings.Contains(out, `"changed"`) {
		t.Errorf("expected changed status in JSON output, got:\n%s", out)
	}
}

func TestGenerate_JSONFormat_HidesUnchangedByDefault(t *testing.T) {
	base := map[string]string{"STABLE": "same", "HOST": "a"}
	target := map[string]string{"STABLE": "same", "HOST": "b"}

	var buf bytes.Buffer
	opts := reporter.DefaultOptions()
	opts.Format = reporter.FormatJSON
	opts.Output = &buf
	opts.MaskSecrets = false
	opts.ShowUnchanged = false

	if err := reporter.Generate(base, target, opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if strings.Contains(buf.String(), `"unchanged"`) {
		t.Errorf("expected unchanged entries to be hidden, got:\n%s", buf.String())
	}
}

func TestDefaultOptions_HasSaneDefaults(t *testing.T) {
	opts := reporter.DefaultOptions()
	if opts.Format != reporter.FormatText {
		t.Errorf("expected default format to be text, got %s", opts.Format)
	}
	if !opts.MaskSecrets {
		t.Error("expected MaskSecrets to be true by default")
	}
}
