package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/yourorg/envlens/internal/masker"
	"github.com/yourorg/envlens/internal/parser"
	"github.com/yourorg/envlens/internal/reporter"
)

func main() {
	var (
		baseFile    = flag.String("base", "", "Base .env file (required)")
		targetFile  = flag.String("target", "", "Target .env file to compare against (required)")
		format      = flag.String("format", "text", "Output format: text or json")
		showAll     = flag.Bool("show-all", false, "Show unchanged keys as well")
		maskSecrets = flag.Bool("mask", true, "Mask sensitive values")
	)
	flag.Parse()

	if *baseFile == "" || *targetFile == "" {
		fmt.Fprintln(os.Stderr, "error: --base and --target are required")
		flag.Usage()
		os.Exit(1)
	}

	baseEnv, err := parser.ParseFile(*baseFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading base file: %v\n", err)
		os.Exit(1)
	}

	targetEnv, err := parser.ParseFile(*targetFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading target file: %v\n", err)
		os.Exit(1)
	}

	m := masker.New()

	opts := reporter.DefaultOptions()
	opts.Format = *format
	opts.ShowUnchanged = *showAll
	opts.MaskSecrets = *maskSecrets
	opts.Masker = m

	if err := reporter.Generate(os.Stdout, baseEnv.ToMap(), targetEnv.ToMap(), opts); err != nil {
		fmt.Fprintf(os.Stderr, "error generating report: %v\n", err)
		os.Exit(1)
	}
}
