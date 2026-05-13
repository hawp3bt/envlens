package merger_test

import (
	"testing"

	"github.com/yourorg/envlens/internal/merger"
)

func TestMerge_NoConflict_MergesAllKeys(t *testing.T) {
	a := map[string]string{"HOST": "localhost", "PORT": "5432"}
	b := map[string]string{"USER": "admin", "PASS": "secret"}

	res, err := merger.Merge([]map[string]string{a, b}, merger.StrategyFirst)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Env["HOST"] != "localhost" || res.Env["USER"] != "admin" {
		t.Errorf("expected all keys merged, got %v", res.Env)
	}
	if len(res.Conflicts) != 0 {
		t.Errorf("expected no conflicts, got %v", res.Conflicts)
	}
}

func TestMerge_StrategyFirst_KeepsFirstValue(t *testing.T) {
	a := map[string]string{"KEY": "first"}
	b := map[string]string{"KEY": "second"}

	res, err := merger.Merge([]map[string]string{a, b}, merger.StrategyFirst)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Env["KEY"] != "first" {
		t.Errorf("expected 'first', got %q", res.Env["KEY"])
	}
	if len(res.Conflicts) != 1 {
		t.Errorf("expected 1 conflict, got %d", len(res.Conflicts))
	}
}

func TestMerge_StrategyLast_KeepsLastValue(t *testing.T) {
	a := map[string]string{"KEY": "first"}
	b := map[string]string{"KEY": "second"}

	res, err := merger.Merge([]map[string]string{a, b}, merger.StrategyLast)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Env["KEY"] != "second" {
		t.Errorf("expected 'second', got %q", res.Env["KEY"])
	}
}

func TestMerge_StrategyError_ReturnsErrorOnConflict(t *testing.T) {
	a := map[string]string{"KEY": "alpha"}
	b := map[string]string{"KEY": "beta"}

	_, err := merger.Merge([]map[string]string{a, b}, merger.StrategyError)
	if err == nil {
		t.Fatal("expected error for conflicting key, got nil")
	}
}

func TestMerge_IdenticalValues_NoConflictRecorded(t *testing.T) {
	a := map[string]string{"KEY": "same"}
	b := map[string]string{"KEY": "same"}

	res, err := merger.Merge([]map[string]string{a, b}, merger.StrategyError)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Conflicts) != 0 {
		t.Errorf("identical values should not produce conflicts")
	}
	if res.Env["KEY"] != "same" {
		t.Errorf("expected 'same', got %q", res.Env["KEY"])
	}
}

func TestMerge_EmptySources_ReturnsEmptyEnv(t *testing.T) {
	res, err := merger.Merge([]map[string]string{}, merger.StrategyFirst)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Env) != 0 {
		t.Errorf("expected empty env, got %v", res.Env)
	}
}
