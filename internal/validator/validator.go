// Package validator checks .env values against expected types or patterns,
// such as URLs, integers, booleans, and non-empty requirements.
package validator

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

// Rule defines a validation rule for a specific key.
type Rule struct {
	Key      string
	Type     string // "bool", "int", "url", "nonempty"
}

// Finding represents a single validation failure.
type Finding struct {
	Key     string
	Value   string
	Rule    string
	Message string
}

func (f Finding) String() string {
	return fmt.Sprintf("[%s] %s=%q — %s", f.Rule, f.Key, f.Value, f.Message)
}

// Validate checks the provided env map against the given rules and returns
// any findings for keys that fail validation.
func Validate(env map[string]string, rules []Rule) []Finding {
	var findings []Finding
	for _, rule := range rules {
		val, exists := env[rule.Key]
		if !exists {
			findings = append(findings, Finding{
				Key:     rule.Key,
				Value:   "",
				Rule:    rule.Type,
				Message: "key not found in env",
			})
			continue
		}
		if f := checkRule(rule.Key, val, rule.Type); f != nil {
			findings = append(findings, *f)
		}
	}
	return findings
}

func checkRule(key, val, ruleType string) *Finding {
	switch strings.ToLower(ruleType) {
	case "nonempty":
		if strings.TrimSpace(val) == "" {
			return &Finding{Key: key, Value: val, Rule: ruleType, Message: "value must not be empty"}
		}
	case "bool":
		norm := strings.ToLower(strings.TrimSpace(val))
		if norm != "true" && norm != "false" && norm != "1" && norm != "0" {
			return &Finding{Key: key, Value: val, Rule: ruleType, Message: "value must be a boolean (true/false/1/0)"}
		}
	case "int":
		if _, err := strconv.Atoi(strings.TrimSpace(val)); err != nil {
			return &Finding{Key: key, Value: val, Rule: ruleType, Message: "value must be an integer"}
		}
	case "url":
		u, err := url.ParseRequestURI(strings.TrimSpace(val))
		if err != nil || u.Scheme == "" || u.Host == "" {
			return &Finding{Key: key, Value: val, Rule: ruleType, Message: "value must be a valid URL"}
		}
	}
	return nil
}
