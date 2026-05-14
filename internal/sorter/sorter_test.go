package sorter_test

import (
	"testing"

	"github.com/user/envlens/internal/sorter"
)

var sampleEnv = map[string]string{
	"DB_HOST":     "localhost",
	"DB_PORT":     "5432",
	"AWS_KEY":     "abc",
	"AWS_SECRET":  "xyz",
	"PORT":        "8080",
	"APP_NAME":    "envlens",
}

func TestSort_Alpha_IsSorted(t *testing.T) {
	keys := sorter.Sort(sampleEnv, sorter.DefaultOptions())
	for i := 1; i < len(keys); i++ {
		if keys[i] < keys[i-1] {
			t.Errorf("expected ascending order, got %s before %s", keys[i-1], keys[i])
		}
	}
}

func TestSort_AlphaDesc_IsSortedDescending(t *testing.T) {
	opts := sorter.Options{Strategy: sorter.StrategyAlphaDesc}
	keys := sorter.Sort(sampleEnv, opts)
	for i := 1; i < len(keys); i++ {
		if keys[i] > keys[i-1] {
			t.Errorf("expected descending order, got %s before %s", keys[i-1], keys[i])
		}
	}
}

func TestSort_Group_ClustersByPrefix(t *testing.T) {
	opts := sorter.Options{Strategy: sorter.StrategyGroup}
	keys := sorter.Sort(sampleEnv, opts)

	// All AWS_ keys should appear before DB_ keys
	lastAWS, firstDB := -1, len(keys)
	for i, k := range keys {
		if len(k) >= 4 && k[:4] == "AWS_" {
			lastAWS = i
		}
		if len(k) >= 3 && k[:3] == "DB_" && i < firstDB {
			firstDB = i
		}
	}
	if lastAWS == -1 || firstDB == len(keys) {
		t.Fatal("expected both AWS_ and DB_ keys in output")
	}
	if lastAWS > firstDB {
		t.Errorf("expected AWS_ group before DB_ group, got lastAWS=%d firstDB=%d", lastAWS, firstDB)
	}
}

func TestSort_Priority_HonoursPriorityList(t *testing.T) {
	opts := sorter.Options{
		Strategy:     sorter.StrategyPriority,
		PriorityKeys: []string{"PORT", "APP_NAME"},
	}
	keys := sorter.Sort(sampleEnv, opts)
	if keys[0] != "PORT" {
		t.Errorf("expected PORT first, got %s", keys[0])
	}
	if keys[1] != "APP_NAME" {
		t.Errorf("expected APP_NAME second, got %s", keys[1])
	}
	// remaining keys should still be sorted alphabetically
	rest := keys[2:]
	for i := 1; i < len(rest); i++ {
		if rest[i] < rest[i-1] {
			t.Errorf("tail not sorted: %s before %s", rest[i-1], rest[i])
		}
	}
}

func TestSort_Priority_NonPriorityKeysSortedAlpha(t *testing.T) {
	opts := sorter.Options{
		Strategy:     sorter.StrategyPriority,
		PriorityKeys: []string{},
	}
	keys := sorter.Sort(sampleEnv, opts)
	for i := 1; i < len(keys); i++ {
		if keys[i] < keys[i-1] {
			t.Errorf("expected alpha fallback, got %s before %s", keys[i-1], keys[i])
		}
	}
}

func TestSort_EmptyMap_ReturnsEmptySlice(t *testing.T) {
	keys := sorter.Sort(map[string]string{}, sorter.DefaultOptions())
	if len(keys) != 0 {
		t.Errorf("expected empty slice, got %v", keys)
	}
}
