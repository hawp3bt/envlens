package patcher_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/envlens/internal/patcher"
)

func baseEnv() map[string]string {
	return map[string]string{
		"APP_NAME": "myapp",
		"APP_PORT": "8080",
		"DB_PASS":  "secret",
	}
}

func TestPatch_SetNewKey(t *testing.T) {
	out, results, err := patcher.Patch(baseEnv(), []patcher.Op{{Action: "set", Key: "NEW_KEY", Value: "hello"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["NEW_KEY"] != "hello" {
		t.Errorf("expected NEW_KEY=hello, got %q", out["NEW_KEY"])
	}
	if results[0].Note != "created" {
		t.Errorf("expected note 'created', got %q", results[0].Note)
	}
}

func TestPatch_SetExistingKey_UpdatesValue(t *testing.T) {
	out, results, err := patcher.Patch(baseEnv(), []patcher.Op{{Action: "set", Key: "APP_PORT", Value: "9090"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["APP_PORT"] != "9090" {
		t.Errorf("expected APP_PORT=9090, got %q", out["APP_PORT"])
	}
	if results[0].Note != "updated" {
		t.Errorf("expected note 'updated', got %q", results[0].Note)
	}
}

func TestPatch_DeleteKey(t *testing.T) {
	out, _, err := patcher.Patch(baseEnv(), []patcher.Op{{Action: "delete", Key: "DB_PASS"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := out["DB_PASS"]; ok {
		t.Error("expected DB_PASS to be deleted")
	}
}

func TestPatch_DeleteMissingKey_NotApplied(t *testing.T) {
	_, results, err := patcher.Patch(baseEnv(), []patcher.Op{{Action: "delete", Key: "GHOST"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if results[0].Applied {
		t.Error("expected op not applied for missing key")
	}
}

func TestPatch_RenameKey(t *testing.T) {
	out, _, err := patcher.Patch(baseEnv(), []patcher.Op{{Action: "rename", Key: "APP_NAME", NewKey: "SERVICE_NAME"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := out["APP_NAME"]; ok {
		t.Error("old key APP_NAME should be gone")
	}
	if out["SERVICE_NAME"] != "myapp" {
		t.Errorf("expected SERVICE_NAME=myapp, got %q", out["SERVICE_NAME"])
	}
}

func TestPatch_UnknownAction_ReturnsError(t *testing.T) {
	_, _, err := patcher.Patch(baseEnv(), []patcher.Op{{Action: "upsert", Key: "X"}})
	if err == nil {
		t.Error("expected error for unknown action")
	}
}

func TestPatch_OriginalMapUnmutated(t *testing.T) {
	env := baseEnv()
	patcher.Patch(env, []patcher.Op{{Action: "delete", Key: "APP_NAME"}})
	if _, ok := env["APP_NAME"]; !ok {
		t.Error("original map should not be mutated")
	}
}

func TestWritePatched_CreatesFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "out.env")
	err := patcher.WritePatched(path, map[string]string{"FOO": "bar", "BAZ": "qux"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, _ := os.ReadFile(path)
	if len(data) == 0 {
		t.Error("expected non-empty output file")
	}
}
