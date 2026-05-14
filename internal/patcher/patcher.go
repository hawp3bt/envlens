package patcher

import (
	"fmt"
	"os"
	"strings"
)

// Op represents a single patch operation.
type Op struct {
	Action string // "set", "delete", "rename"
	Key    string
	Value  string
	NewKey string // used for rename
}

// Result describes the outcome of applying a patch op.
type Result struct {
	Op      Op
	Applied bool
	Note    string
}

// Patch applies a slice of Ops to the given env map, returning results.
// The original map is not mutated; a new map is returned.
func Patch(env map[string]string, ops []Op) (map[string]string, []Result, error) {
	out := make(map[string]string, len(env))
	for k, v := range env {
		out[k] = v
	}

	results := make([]Result, 0, len(ops))

	for _, op := range ops {
		switch strings.ToLower(op.Action) {
		case "set":
			_, existed := out[op.Key]
			out[op.Key] = op.Value
			note := "created"
			if existed {
				note = "updated"
			}
			results = append(results, Result{Op: op, Applied: true, Note: note})

		case "delete":
			if _, ok := out[op.Key]; !ok {
				results = append(results, Result{Op: op, Applied: false, Note: "key not found"})
				continue
			}
			delete(out, op.Key)
			results = append(results, Result{Op: op, Applied: true, Note: "deleted"})

		case "rename":
			if op.NewKey == "" {
				return nil, nil, fmt.Errorf("rename op for %q missing new_key", op.Key)
			}
			v, ok := out[op.Key]
			if !ok {
				results = append(results, Result{Op: op, Applied: false, Note: "key not found"})
				continue
			}
			out[op.NewKey] = v
			delete(out, op.Key)
			results = append(results, Result{Op: op, Applied: true, Note: "renamed"})

		default:
			return nil, nil, fmt.Errorf("unknown patch action: %q", op.Action)
		}
	}

	return out, results, nil
}

// WritePatched serialises the patched env map to the given file path.
func WritePatched(path string, env map[string]string) error {
	var sb strings.Builder
	for k, v := range env {
		if strings.ContainsAny(v, " \t") {
			fmt.Fprintf(&sb, "%s=\"%s\"\n", k, v)
		} else {
			fmt.Fprintf(&sb, "%s=%s\n", k, v)
		}
	}
	return os.WriteFile(path, []byte(sb.String()), 0o644)
}
