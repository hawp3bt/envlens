package parser

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// EnvEntry represents a single key-value pair from a .env file.
type EnvEntry struct {
	Key     string
	Value   string
	Comment string
	LineNum int
}

// EnvFile holds all parsed entries from a .env file.
type EnvFile struct {
	Path    string
	Entries []EnvEntry
}

// ParseFile reads and parses a .env file, returning an EnvFile.
func ParseFile(path string) (*EnvFile, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %q: %w", path, err)
	}
	defer f.Close()

	envFile := &EnvFile{Path: path}
	scanner := bufio.NewScanner(f)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines
		if line == "" {
			continue
		}

		// Extract inline comment
		comment := ""
		if idx := strings.Index(line, " #"); idx != -1 {
			comment = strings.TrimSpace(line[idx+2:])
			line = strings.TrimSpace(line[:idx])
		}

		// Skip full-line comments
		if strings.HasPrefix(line, "#") {
			continue
		}

		// Parse KEY=VALUE
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.Trim(strings.TrimSpace(parts[1]), `"`)

		envFile.Entries = append(envFile.Entries, EnvEntry{
			Key:     key,
			Value:   value,
			Comment: comment,
			LineNum: lineNum,
		})
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file %q: %w", path, err)
	}

	return envFile, nil
}

// ToMap converts EnvFile entries into a key-value map.
func (e *EnvFile) ToMap() map[string]string {
	m := make(map[string]string, len(e.Entries))
	for _, entry := range e.Entries {
		m[entry.Key] = entry.Value
	}
	return m
}
