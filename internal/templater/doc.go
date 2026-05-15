// Package templater renders text templates against a map of environment
// variables, allowing .env values to be embedded into configuration files,
// shell scripts, or any other text-based artefact.
//
// Usage:
//
//	res, err := templater.Render(tmplText, envMap, templater.DefaultOptions())
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Println(res.Output)
//
// The env helper function is available inside templates:
//
//	{{ env "MY_VAR" }}
//
// MissingKey behaviour mirrors Go's text/template option:
//   - "zero"  — substitute an empty string (default)
//   - "error" — return an error if a key is absent
package templater
