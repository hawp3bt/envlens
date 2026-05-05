// Package masker provides utilities for detecting and masking sensitive
// environment variable values in .env files.
//
// Sensitive keys are identified by matching against a configurable list of
// patterns (e.g. "SECRET", "PASSWORD", "TOKEN"). When a key is deemed
// sensitive, its value is replaced with a configurable mask string (default: "***").
//
// Example usage:
//
//	masker := masker.New(nil, "[REDACTED]")
//
//	env := map[string]string{
//	    "APP_ENV":     "production",
//	    "DB_PASSWORD": "supersecret",
//	}
//
//	safe := masker.MaskMap(env)
//	// safe["APP_ENV"]     == "production"
//	// safe["DB_PASSWORD"] == "[REDACTED]"
package masker
