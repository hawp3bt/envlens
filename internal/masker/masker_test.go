package masker

import (
	"testing"
)

func TestIsSensitive_MatchesKnownPatterns(t *testing.T) {
	m := New(nil, "")

	sensitiveKeys := []string{
		"DB_PASSWORD",
		"API_KEY",
		"AUTH_TOKEN",
		"STRIPE_SECRET",
		"PRIVATE_KEY_PATH",
		"DATABASE_URL",
	}

	for _, key := range sensitiveKeys {
		if !m.IsSensitive(key) {
			t.Errorf("expected key %q to be sensitive", key)
		}
	}
}

func TestIsSensitive_SafeKeys(t *testing.T) {
	m := New(nil, "")

	safeKeys := []string{
		"APP_ENV",
		"PORT",
		"LOG_LEVEL",
		"REGION",
	}

	for _, key := range safeKeys {
		if m.IsSensitive(key) {
			t.Errorf("expected key %q to NOT be sensitive", key)
		}
	}
}

func TestMaskValue_MasksSensitive(t *testing.T) {
	m := New(nil, "****")

	masked := m.MaskValue("DB_PASSWORD", "supersecret")
	if masked != "****" {
		t.Errorf("expected masked value, got %q", masked)
	}
}

func TestMaskValue_PreservesSafe(t *testing.T) {
	m := New(nil, "****")

	val := m.MaskValue("APP_ENV", "production")
	if val != "production" {
		t.Errorf("expected original value, got %q", val)
	}
}

func TestMaskMap_MasksCorrectly(t *testing.T) {
	m := New(nil, "[REDACTED]")

	env := map[string]string{
		"APP_ENV":     "production",
		"DB_PASSWORD": "s3cr3t",
		"PORT":        "8080",
		"API_KEY":     "key-abc-123",
	}

	result := m.MaskMap(env)

	if result["APP_ENV"] != "production" {
		t.Errorf("APP_ENV should not be masked")
	}
	if result["PORT"] != "8080" {
		t.Errorf("PORT should not be masked")
	}
	if result["DB_PASSWORD"] != "[REDACTED]" {
		t.Errorf("DB_PASSWORD should be masked")
	}
	if result["API_KEY"] != "[REDACTED]" {
		t.Errorf("API_KEY should be masked")
	}
}

func TestNew_CustomPatterns(t *testing.T) {
	m := New([]string{"CUSTOM_FIELD"}, "XXX")

	if !m.IsSensitive("MY_CUSTOM_FIELD") {
		t.Error("expected custom pattern to match")
	}
	if m.IsSensitive("API_KEY") {
		t.Error("default patterns should not apply when custom patterns provided")
	}
}
