package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/yourorg/envlens/internal/masker"
	"github.com/yourorg/envlens/internal/parser"
	"github.com/yourorg/envlens/internal/snapshot"
)

// snapshotFlags holds parsed CLI flags for the snapshot subcommand.
type snapshotFlags struct {
	envFile     string
	environment string
	snapshotDir string
	listEnv     string
}

// runSnapshot executes the snapshot subcommand logic.
func runSnapshot(args []string) error {
	fs := flag.NewFlagSet("snapshot", flag.ContinueOnError)
	var f snapshotFlags

	fs.StringVar(&f.envFile, "file", "", "path to .env file to snapshot")
	fs.StringVar(&f.environment, "env", "", "environment label (e.g. production)")
	fs.StringVar(&f.snapshotDir, "dir", ".envlens/snapshots", "directory to store snapshots")
	fs.StringVar(&f.listEnv, "list", "", "list snapshots for given environment label")

	if err := fs.Parse(args); err != nil {
		return err
	}

	if f.listEnv != "" {
		return listSnapshots(f.snapshotDir, f.listEnv)
	}

	if f.envFile == "" || f.environment == "" {
		return fmt.Errorf("snapshot: --file and --env are required")
	}

	entries, err := parser.ParseFile(f.envFile)
	if err != nil {
		return fmt.Errorf("snapshot: parse: %w", err)
	}

	m := masker.New()
	masked := m.MaskMap(entries.ToMap())

	path, err := snapshot.Save(f.snapshotDir, f.environment, masked)
	if err != nil {
		return fmt.Errorf("snapshot: save: %w", err)
	}

	fmt.Fprintf(os.Stdout, "snapshot saved: %s\n", path)
	return nil
}

func listSnapshots(dir, env string) error {
	files, err := snapshot.List(dir, env)
	if err != nil {
		return err
	}
	if len(files) == 0 {
		fmt.Fprintf(os.Stdout, "no snapshots found for environment %q\n", env)
		return nil
	}
	for _, f := range files {
		fmt.Fprintln(os.Stdout, f)
	}
	return nil
}
