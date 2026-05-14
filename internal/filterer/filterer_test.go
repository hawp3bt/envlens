package filterer_test

import (
	"testing"

	"github.com/yourorg/envlens/internal/filterer"
)

var sampleEnv = map[string]string{
	"DB_HOST":     "localhost",
	"DB_PASSWORD": "secret",
	"AWS_KEY":     "AKID",
	"APP_ENV":     "production",
	"PORT":        "8080",
}

func TestFilter_IncludeGlob_KeepsMatchingKeys(t *testing.T) {
	opts := filterer.Options{Mode: filterer.ModeInclude, Patterns: []string{"DB_*"}}
	out, err := filterer.Filter(sampleEnv, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := out["DB_HOST"]; !ok {
		t.Error("expected DB_HOST to be included")
	}
	if _, ok := out["PORT"]; ok {
		t.Error("expected PORT to be excluded")
	}
}

func TestFilter_ExcludeGlob_RemovesMatchingKeys(t *testing.T) {
	opts := filterer.Options{Mode: filterer.ModeExclude, Patterns: []string{"DB_*"}}
	out, err := filterer.Filter(sampleEnv, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := out["DB_HOST"]; ok {
		t.Error("expected DB_HOST to be excluded")
	}
	if _, ok := out["PORT"]; !ok {
		t.Error("expected PORT to be retained")
	}
}

func TestFilter_RegexPattern_MatchesCorrectly(t *testing.T) {
	opts := filterer.Options{Mode: filterer.ModeInclude, Patterns: []string{"^AWS"}}
	out, err := filterer.Filter(sampleEnv, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := out["AWS_KEY"]; !ok {
		t.Error("expected AWS_KEY to be included via regex")
	}
	if len(out) != 1 {
		t.Errorf("expected 1 key, got %d", len(out))
	}
}

func TestFilter_MultiplePatterns_UnionMatch(t *testing.T) {
	opts := filterer.Options{Mode: filterer.ModeInclude, Patterns: []string{"DB_*", "PORT"}}
	out, err := filterer.Filter(sampleEnv, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 3 {
		t.Errorf("expected 3 keys (DB_HOST, DB_PASSWORD, PORT), got %d", len(out))
	}
}

func TestFilter_InvalidRegex_ReturnsError(t *testing.T) {
	opts := filterer.Options{Mode: filterer.ModeInclude, Patterns: []string{"^[invalid"}}
	_, err := filterer.Filter(sampleEnv, opts)
	if err == nil {
		t.Error("expected error for invalid regex pattern")
	}
}

func TestFilter_DefaultOptions_ReturnsAll(t *testing.T) {
	opts := filterer.DefaultOptions()
	out, err := filterer.Filter(sampleEnv, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != len(sampleEnv) {
		t.Errorf("expected all %d keys, got %d", len(sampleEnv), len(out))
	}
}

func TestFilter_EmptyPatterns_IncludeMode_ReturnsEmpty(t *testing.T) {
	opts := filterer.Options{Mode: filterer.ModeInclude, Patterns: []string{}}
	out, err := filterer.Filter(sampleEnv, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 0 {
		t.Errorf("expected 0 keys with no patterns in include mode, got %d", len(out))
	}
}
