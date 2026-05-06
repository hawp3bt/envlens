package differ_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/yourorg/envlens/internal/differ"
)

func TestFormat_AddedKeyAppearsInOutput(t *testing.T) {
	base := map[string]string{}
	target := map[string]string{"NEW_VAR": "hello"}

	result := differ.Diff(base, target)

	var buf bytes.Buffer
	differ.Format(&buf, result, differ.FormatOptions{Colorized: false})

	output := buf.String()
	if !strings.Contains(output, "+ NEW_VAR=hello") {
		t.Errorf("expected added key in output, got:\n%s", output)
	}
}

func TestFormat_RemovedKeyAppearsInOutput(t *testing.T) {
	base := map[string]string{"OLD_VAR": "bye"}
	target := map[string]string{}

	result := differ.Diff(base, target)

	var buf bytes.Buffer
	differ.Format(&buf, result, differ.FormatOptions{Colorized: false})

	output := buf.String()
	if !strings.Contains(output, "- OLD_VAR=bye") {
		t.Errorf("expected removed key in output, got:\n%s", output)
	}
}

func TestFormat_ChangedKeyAppearsInOutput(t *testing.T) {
	base := map[string]string{"HOST": "localhost"}
	target := map[string]string{"HOST": "prod.example.com"}

	result := differ.Diff(base, target)

	var buf bytes.Buffer
	differ.Format(&buf, result, differ.FormatOptions{Colorized: false})

	output := buf.String()
	if !strings.Contains(output, "HOST") {
		t.Errorf("expected changed key in output, got:\n%s", output)
	}
}

func TestFormat_UnchangedHiddenByDefault(t *testing.T) {
	base := map[string]string{"SAME": "value"}
	target := map[string]string{"SAME": "value"}

	result := differ.Diff(base, target)

	var buf bytes.Buffer
	differ.Format(&buf, result, differ.FormatOptions{Colorized: false, ShowUnchanged: false})

	output := buf.String()
	if strings.Contains(output, "SAME=value") {
		t.Errorf("expected unchanged key to be hidden, got:\n%s", output)
	}
}

func TestFormat_UnchangedShownWhenRequested(t *testing.T) {
	base := map[string]string{"SAME": "value"}
	target := map[string]string{"SAME": "value"}

	result := differ.Diff(base, target)

	var buf bytes.Buffer
	differ.Format(&buf, result, differ.FormatOptions{Colorized: false, ShowUnchanged: true})

	output := buf.String()
	if !strings.Contains(output, "SAME=value") {
		t.Errorf("expected unchanged key in output, got:\n%s", output)
	}
}
