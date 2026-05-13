package profiler_test

import (
	"strings"
	"testing"

	"github.com/user/envlens/internal/profiler"
)

func TestAnalyze_CountsTotalKeys(t *testing.T) {
	env := map[string]string{"FOO": "bar", "BAZ": "qux"}
	p := profiler.Analyze(env)
	if p.TotalKeys != 2 {
		t.Errorf("expected 2 total keys, got %d", p.TotalKeys)
	}
}

func TestAnalyze_CountsSecretKeys(t *testing.T) {
	env := map[string]string{
		"DATABASE_PASSWORD": "secret",
		"API_KEY":           "abc123",
		"APP_NAME":          "myapp",
	}
	p := profiler.Analyze(env)
	if p.SecretKeys < 2 {
		t.Errorf("expected at least 2 secret keys, got %d", p.SecretKeys)
	}
}

func TestAnalyze_CountsEmptyValues(t *testing.T) {
	env := map[string]string{"FOO": "", "BAR": "   ", "BAZ": "value"}
	p := profiler.Analyze(env)
	if p.EmptyValues != 2 {
		t.Errorf("expected 2 empty values, got %d", p.EmptyValues)
	}
}

func TestAnalyze_UniqueValues(t *testing.T) {
	env := map[string]string{"A": "same", "B": "same", "C": "different"}
	p := profiler.Analyze(env)
	if p.UniqueValues != 2 {
		t.Errorf("expected 2 unique values, got %d", p.UniqueValues)
	}
}

func TestAnalyze_SecretRatio(t *testing.T) {
	env := map[string]string{
		"API_KEY":  "secret",
		"APP_NAME": "myapp",
	}
	p := profiler.Analyze(env)
	if p.SecretRatio <= 0 || p.SecretRatio > 1 {
		t.Errorf("unexpected secret ratio: %f", p.SecretRatio)
	}
}

func TestAnalyze_TopKeys_AtMostFive(t *testing.T) {
	env := map[string]string{
		"ALPHA": "1", "BETA": "2", "GAMMA": "3",
		"DELTA": "4", "EPSILON": "5", "ZETA": "6",
	}
	p := profiler.Analyze(env)
	if len(p.TopKeys) > 5 {
		t.Errorf("expected at most 5 top keys, got %d", len(p.TopKeys))
	}
}

func TestAnalyze_EmptyMap(t *testing.T) {
	p := profiler.Analyze(map[string]string{})
	if p.TotalKeys != 0 || p.SecretRatio != 0 {
		t.Errorf("expected zero profile for empty map")
	}
}

func TestProfile_String_ContainsSummary(t *testing.T) {
	env := map[string]string{"SECRET_KEY": "abc", "APP": "test"}
	p := profiler.Analyze(env)
	s := p.String()
	if !strings.Contains(s, "Total keys") {
		t.Errorf("expected 'Total keys' in profile string, got: %s", s)
	}
}
