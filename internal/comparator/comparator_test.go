package comparator_test

import (
	"testing"

	"github.com/yourorg/envlens/internal/comparator"
)

func TestCompare_AllEnvsHaveKey_NoDivergence(t *testing.T) {
	envs := map[string]map[string]string{
		"dev":  {"APP_ENV": "development"},
		"prod": {"APP_ENV": "development"},
	}
	res := comparator.Compare(envs)
	if len(res.Statuses) != 1 {
		t.Fatalf("expected 1 key status, got %d", len(res.Statuses))
	}
	ks := res.Statuses[0]
	if ks.Diverges {
		t.Error("expected no divergence for identical values")
	}
}

func TestCompare_DifferentValues_Diverges(t *testing.T) {
	envs := map[string]map[string]string{
		"dev":  {"DB_HOST": "localhost"},
		"prod": {"DB_HOST": "db.prod.example.com"},
	}
	res := comparator.Compare(envs)
	if !res.Statuses[0].Diverges {
		t.Error("expected divergence for differing values")
	}
}

func TestCompare_MissingKeyInOneEnv(t *testing.T) {
	envs := map[string]map[string]string{
		"dev":  {"FEATURE_FLAG": "true"},
		"prod": {},
	}
	res := comparator.Compare(envs)
	ks := res.Statuses[0]
	if ks.Status["dev"] != comparator.Present {
		t.Errorf("expected dev to be present, got %s", ks.Status["dev"])
	}
	if ks.Status["prod"] != comparator.Missing {
		t.Errorf("expected prod to be missing, got %s", ks.Status["prod"])
	}
}

func TestCompare_KeysSortedAlphabetically(t *testing.T) {
	envs := map[string]map[string]string{
		"dev": {"ZEBRA": "1", "ALPHA": "2", "MIDDLE": "3"},
	}
	res := comparator.Compare(envs)
	expected := []string{"ALPHA", "MIDDLE", "ZEBRA"}
	for i, ks := range res.Statuses {
		if ks.Key != expected[i] {
			t.Errorf("position %d: expected %s, got %s", i, expected[i], ks.Key)
		}
	}
}

func TestCompare_EnvLabelsSortedAlphabetically(t *testing.T) {
	envs := map[string]map[string]string{
		"staging": {"X": "1"},
		"dev":     {"X": "1"},
		"prod":    {"X": "1"},
	}
	res := comparator.Compare(envs)
	expected := []string{"dev", "prod", "staging"}
	for i, lbl := range res.Envs {
		if lbl != expected[i] {
			t.Errorf("position %d: expected %s, got %s", i, expected[i], lbl)
		}
	}
}

func TestCompare_EmptyInputs_ReturnsEmptyResult(t *testing.T) {
	res := comparator.Compare(map[string]map[string]string{})
	if len(res.Statuses) != 0 {
		t.Errorf("expected 0 statuses for empty input, got %d", len(res.Statuses))
	}
}
