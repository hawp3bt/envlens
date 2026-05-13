package transformer_test

import (
	"testing"

	"github.com/user/envlens/internal/transformer"
)

func TestTransform_NoOp_ReturnsSameKeys(t *testing.T) {
	env := map[string]string{"FOO": "bar", "BAZ": "qux"}
	out, err := transformer.Transform(env, transformer.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["FOO"] != "bar" || out["BAZ"] != "qux" {
		t.Errorf("expected unchanged output, got %v", out)
	}
}

func TestTransform_Prefix_AppendsToKeys(t *testing.T) {
	env := map[string]string{"DB_HOST": "localhost"}
	opts := transformer.DefaultOptions()
	opts.Prefix = "APP_"
	out, err := transformer.Transform(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := out["APP_DB_HOST"]; !ok {
		t.Errorf("expected key APP_DB_HOST, got %v", out)
	}
}

func TestTransform_Suffix_AppendsToKeys(t *testing.T) {
	env := map[string]string{"PORT": "8080"}
	opts := transformer.DefaultOptions()
	opts.Suffix = "_VAR"
	out, err := transformer.Transform(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := out["PORT_VAR"]; !ok {
		t.Errorf("expected key PORT_VAR, got %v", out)
	}
}

func TestTransform_UppercaseKeys(t *testing.T) {
	env := map[string]string{"db_host": "localhost", "api_key": "secret"}
	opts := transformer.DefaultOptions()
	opts.UppercaseKeys = true
	out, err := transformer.Transform(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["DB_HOST"] != "localhost" || out["API_KEY"] != "secret" {
		t.Errorf("expected uppercase keys, got %v", out)
	}
}

func TestTransform_LowercaseKeys(t *testing.T) {
	env := map[string]string{"DB_HOST": "localhost"}
	opts := transformer.DefaultOptions()
	opts.LowercaseKeys = true
	out, err := transformer.Transform(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["db_host"] != "localhost" {
		t.Errorf("expected lowercase key db_host, got %v", out)
	}
}

func TestTransform_RenameMap_RenamesKey(t *testing.T) {
	env := map[string]string{"OLD_KEY": "value"}
	opts := transformer.DefaultOptions()
	opts.RenameMap = map[string]string{"OLD_KEY": "NEW_KEY"}
	out, err := transformer.Transform(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["NEW_KEY"] != "value" {
		t.Errorf("expected NEW_KEY=value, got %v", out)
	}
	if _, ok := out["OLD_KEY"]; ok {
		t.Errorf("expected OLD_KEY to be removed")
	}
}

func TestTransform_PreservesValues(t *testing.T) {
	env := map[string]string{"SECRET": "super-secret-value"}
	opts := transformer.DefaultOptions()
	opts.Prefix = "TEST_"
	out, _ := transformer.Transform(env, opts)
	if out["TEST_SECRET"] != "super-secret-value" {
		t.Errorf("expected value to be preserved, got %v", out["TEST_SECRET"])
	}
}
