// Package templater renders .env values into text templates,
// substituting {{ .KEY }} placeholders with environment values.
package templater

import (
	"bytes"
	"fmt"
	"text/template"
)

// Options controls templater behaviour.
type Options struct {
	// MissingKey controls what happens when a placeholder key is absent.
	// Accepted values: "zero" (empty string), "error" (return error).
	MissingKey string
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{
		MissingKey: "zero",
	}
}

// Result holds the rendered output and any keys that were referenced.
type Result struct {
	Output       string
	ReferencedKeys []string
}

// Render executes tmplText as a Go text/template, substituting values
// from env. Keys referenced in the template are recorded in Result.
func Render(tmplText string, env map[string]string, opts Options) (Result, error) {
	referenced := []string{}
	tracked := make(map[string]string)

	// Wrap env so we can record which keys are accessed.
	proxy := &trackingMap{env: env, seen: map[string]bool{}}

	missingKey := "zero"
	if opts.MissingKey == "error" {
		missingKey = "error"
	}

	tmpl, err := template.New("env").
		Option("missingkey=" + missingKey).
		Funcs(template.FuncMap{
			"env": func(key string) string {
				proxy.seen[key] = true
				return env[key]
			},
		}).
		Parse(tmplText)
	if err != nil {
		return Result{}, fmt.Errorf("templater: parse error: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, tracked); err != nil {
		return Result{}, fmt.Errorf("templater: execute error: %w", err)
	}

	for k := range proxy.seen {
		referenced = append(referenced, k)
	}

	return Result{Output: buf.String(), ReferencedKeys: referenced}, nil
}

// trackingMap records which keys are accessed via the env() func.
type trackingMap struct {
	env  map[string]string
	seen map[string]bool
}
