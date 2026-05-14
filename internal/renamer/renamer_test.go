package renamer

import (
	"testing"
)

func TestRename_ExactMatch_RenamesKey(t *testing.T) {
	env := map[string]string{"OLD_KEY": "value"}
	opts := Options{Rules: []Rule{{From: "OLD_KEY", To: "NEW_KEY"}}}
	out, results, err := Rename(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["NEW_KEY"] != "value" {
		t.Errorf("expected NEW_KEY=value, got %v", out)
	}
	if _, exists := out["OLD_KEY"]; exists {
		t.Error("OLD_KEY should not exist after rename")
	}
	if !results[0].Applied {
		t.Error("expected Applied=true")
	}
}

func TestRename_PrefixMatch_RenamesPrefix(t *testing.T) {
	env := map[string]string{"APP_HOST": "localhost", "APP_PORT": "8080"}
	opts := Options{Rules: []Rule{{From: "APP_", To: "SVC_", IsPrefix: true}}}
	out, _, err := Rename(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["SVC_HOST"] != "localhost" {
		t.Errorf("expected SVC_HOST=localhost, got %v", out)
	}
	if out["SVC_PORT"] != "8080" {
		t.Errorf("expected SVC_PORT=8080, got %v", out)
	}
}

func TestRename_NoMatchingRule_KeyUnchanged(t *testing.T) {
	env := map[string]string{"UNRELATED": "val"}
	opts := Options{Rules: []Rule{{From: "OTHER", To: "THING"}}}
	out, results, err := Rename(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["UNRELATED"] != "val" {
		t.Errorf("expected UNRELATED unchanged, got %v", out)
	}
	if results[0].Applied {
		t.Error("expected Applied=false for unmatched key")
	}
}

func TestRename_Collision_ReturnsError(t *testing.T) {
	env := map[string]string{"OLD": "a", "NEW": "b"}
	opts := Options{Rules: []Rule{{From: "OLD", To: "NEW"}}}
	_, _, err := Rename(env, opts)
	if err == nil {
		t.Fatal("expected collision error, got nil")
	}
}

func TestRename_IgnoreCase_MatchesLowerInput(t *testing.T) {
	env := map[string]string{"old_key": "v"}
	opts := Options{
		Rules:      []Rule{{From: "OLD_KEY", To: "NEW_KEY"}},
		IgnoreCase: true,
	}
	out, _, err := Rename(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["NEW_KEY"] != "v" {
		t.Errorf("expected NEW_KEY=v, got %v", out)
	}
}

func TestRename_EmptyEnv_ReturnsEmpty(t *testing.T) {
	out, results, err := Rename(map[string]string{}, Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 0 || len(results) != 0 {
		t.Error("expected empty output for empty input")
	}
}
