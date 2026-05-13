package interpolator

import (
	"os"
	"testing"
)

func TestInterpolate_BraceStyle_ResolvesReference(t *testing.T) {
	env := map[string]string{
		"HOST": "localhost",
		"URL":  "http://${HOST}:8080",
	}
	out, err := Interpolate(env, DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := out["URL"]; got != "http://localhost:8080" {
		t.Errorf("expected http://localhost:8080, got %q", got)
	}
}

func TestInterpolate_DollarStyle_ResolvesReference(t *testing.T) {
	env := map[string]string{
		"NAME": "world",
		"MSG":  "hello $NAME",
	}
	out, err := Interpolate(env, DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := out["MSG"]; got != "hello world" {
		t.Errorf("expected 'hello world', got %q", got)
	}
}

func TestInterpolate_UnresolvedRef_EmptyStringByDefault(t *testing.T) {
	env := map[string]string{
		"VAL": "prefix_${MISSING}_suffix",
	}
	opts := Options{FallbackToOS: false, ErrorOnMissing: false}
	out, err := Interpolate(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := out["VAL"]; got != "prefix__suffix" {
		t.Errorf("expected 'prefix__suffix', got %q", got)
	}
}

func TestInterpolate_ErrorOnMissing_ReturnsError(t *testing.T) {
	env := map[string]string{
		"VAL": "${DOES_NOT_EXIST}",
	}
	opts := Options{FallbackToOS: false, ErrorOnMissing: true}
	_, err := Interpolate(env, opts)
	if err == nil {
		t.Fatal("expected error for missing variable, got nil")
	}
}

func TestInterpolate_FallbackToOS_ResolvesFromEnviron(t *testing.T) {
	os.Setenv("_TEST_ENVLENS_VAR", "fromOS")
	defer os.Unsetenv("_TEST_ENVLENS_VAR")

	env := map[string]string{
		"RESULT": "value=${_TEST_ENVLENS_VAR}",
	}
	opts := Options{FallbackToOS: true, ErrorOnMissing: false}
	out, err := Interpolate(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := out["RESULT"]; got != "value=fromOS" {
		t.Errorf("expected 'value=fromOS', got %q", got)
	}
}

func TestInterpolate_NoReferences_ValueUnchanged(t *testing.T) {
	env := map[string]string{
		"PLAIN": "no-refs-here",
	}
	out, err := Interpolate(env, DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := out["PLAIN"]; got != "no-refs-here" {
		t.Errorf("expected 'no-refs-here', got %q", got)
	}
}

func TestDefaultOptions_FallbackEnabled(t *testing.T) {
	opts := DefaultOptions()
	if !opts.FallbackToOS {
		t.Error("expected FallbackToOS to be true by default")
	}
	if opts.ErrorOnMissing {
		t.Error("expected ErrorOnMissing to be false by default")
	}
}
