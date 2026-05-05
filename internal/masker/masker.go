package masker

import "strings"

// DefaultSensitiveKeys contains common patterns for sensitive environment variable keys.
var DefaultSensitiveKeys = []string{
	"SECRET",
	"PASSWORD",
	"PASSWD",
	"TOKEN",
	"API_KEY",
	"APIKEY",
	"PRIVATE_KEY",
	"AUTH",
	"CREDENTIAL",
	"DSN",
	"DATABASE_URL",
	"DB_URL",
}

// Masker handles masking of sensitive environment variable values.
type Masker struct {
	sensitivePatterns []string
	maskValue         string
}

// New creates a new Masker with the given sensitive key patterns.
// If patterns is nil, DefaultSensitiveKeys is used.
func New(patterns []string, maskValue string) *Masker {
	if patterns == nil {
		patterns = DefaultSensitiveKeys
	}
	if maskValue == "" {
		maskValue = "****"
	}
	return &Masker{
		sensitivePatterns: patterns,
		maskValue:         maskValue,
	}
}

// IsSensitive returns true if the key matches any sensitive pattern.
func (m *Masker) IsSensitive(key string) bool {
	upper := strings.ToUpper(key)
	for _, pattern := range m.sensitivePatterns {
		if strings.Contains(upper, strings.ToUpper(pattern)) {
			return true
		}
	}
	return false
}

// MaskValue returns the masked value if the key is sensitive, otherwise returns the original value.
func (m *Masker) MaskValue(key, value string) string {
	if m.IsSensitive(key) {
		return m.maskValue
	}
	return value
}

// MaskMap returns a copy of the map with sensitive values masked.
func (m *Masker) MaskMap(env map[string]string) map[string]string {
	masked := make(map[string]string, len(env))
	for k, v := range env {
		masked[k] = m.MaskValue(k, v)
	}
	return masked
}
