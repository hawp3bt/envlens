// Package merger implements multi-source .env file merging for envlens.
//
// It accepts an ordered slice of parsed environment maps and combines them
// into a single map using one of three conflict resolution strategies:
//
//   - StrategyFirst  — the first file wins on key conflicts.
//   - StrategyLast   — the last file wins on key conflicts.
//   - StrategyError  — returns an error when conflicting values are detected.
//
// Conflicts where all sources agree on the same value are silently resolved
// regardless of the chosen strategy.
//
// Typical usage:
//
//	res, err := merger.Merge([]map[string]string{base, override}, merger.StrategyLast)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Println(res.Env)
//	for _, c := range res.Conflicts {
//		fmt.Printf("conflict on %s: %v\n", c.Key, c.Values)
//	}
package merger
