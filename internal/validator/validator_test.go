package validator

import (
	"strings"
	"testing"
)

func TestValidate_NonemptyRule_FailsOnEmpty(t *testing.T) {
	env := map[string]string{"APP_NAME": ""}
	rules := []Rule{{Key: "APP_NAME", Type: "nonempty"}}
	findings := Validate(env, rules)
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(findings))
	}
	if !strings.Contains(findings[0].Message, "empty") {
		t.Errorf("unexpected message: %s", findings[0].Message)
	}
}

func TestValidate_NonemptyRule_PassesOnValue(t *testing.T) {
	env := map[string]string{"APP_NAME": "envlens"}
	rules := []Rule{{Key: "APP_NAME", Type: "nonempty"}}
	if findings := Validate(env, rules); len(findings) != 0 {
		t.Errorf("expected no findings, got %v", findings)
	}
}

func TestValidate_BoolRule_FailsOnInvalid(t *testing.T) {
	env := map[string]string{"DEBUG": "yes"}
	rules := []Rule{{Key: "DEBUG", Type: "bool"}}
	findings := Validate(env, rules)
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(findings))
	}
}

func TestValidate_BoolRule_PassesOnTrue(t *testing.T) {
	for _, v := range []string{"true", "false", "1", "0"} {
		env := map[string]string{"DEBUG": v}
		rules := []Rule{{Key: "DEBUG", Type: "bool"}}
		if findings := Validate(env, rules); len(findings) != 0 {
			t.Errorf("value %q should pass bool rule, got %v", v, findings)
		}
	}
}

func TestValidate_IntRule_FailsOnString(t *testing.T) {
	env := map[string]string{"PORT": "abc"}
	rules := []Rule{{Key: "PORT", Type: "int"}}
	findings := Validate(env, rules)
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(findings))
	}
}

func TestValidate_IntRule_PassesOnNumber(t *testing.T) {
	env := map[string]string{"PORT": "8080"}
	rules := []Rule{{Key: "PORT", Type: "int"}}
	if findings := Validate(env, rules); len(findings) != 0 {
		t.Errorf("expected no findings, got %v", findings)
	}
}

func TestValidate_URLRule_FailsOnBadURL(t *testing.T) {
	env := map[string]string{"API_URL": "not-a-url"}
	rules := []Rule{{Key: "API_URL", Type: "url"}}
	findings := Validate(env, rules)
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(findings))
	}
}

func TestValidate_URLRule_PassesOnValidURL(t *testing.T) {
	env := map[string]string{"API_URL": "https://api.example.com/v1"}
	rules := []Rule{{Key: "API_URL", Type: "url"}}
	if findings := Validate(env, rules); len(findings) != 0 {
		t.Errorf("expected no findings, got %v", findings)
	}
}

func TestValidate_MissingKey_ReturnsError(t *testing.T) {
	env := map[string]string{}
	rules := []Rule{{Key: "REQUIRED_KEY", Type: "nonempty"}}
	findings := Validate(env, rules)
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding for missing key, got %d", len(findings))
	}
	if !strings.Contains(findings[0].Message, "not found") {
		t.Errorf("unexpected message: %s", findings[0].Message)
	}
}
