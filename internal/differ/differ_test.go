package differ_test

import (
	"testing"

	"github.com/yourorg/envlens/internal/differ"
)

func TestDiff_AddedKey(t *testing.T) {
	base := map[string]string{"FOO": "bar"}
	target := map[string]string{"FOO": "bar", "NEW_KEY": "value"}

	result := differ.Diff(base, target)

	if len(result.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(result.Entries))
	}
	found := findEntry(result.Entries, "NEW_KEY")
	if found == nil || found.Type != differ.Added {
		t.Errorf("expected NEW_KEY to be Added")
	}
}

func TestDiff_RemovedKey(t *testing.T) {
	base := map[string]string{"FOO": "bar", "OLD_KEY": "gone"}
	target := map[string]string{"FOO": "bar"}

	result := differ.Diff(base, target)

	found := findEntry(result.Entries, "OLD_KEY")
	if found == nil || found.Type != differ.Removed {
		t.Errorf("expected OLD_KEY to be Removed")
	}
}

func TestDiff_ChangedKey(t *testing.T) {
	base := map[string]string{"DB_HOST": "localhost"}
	target := map[string]string{"DB_HOST": "prod.db.example.com"}

	result := differ.Diff(base, target)

	found := findEntry(result.Entries, "DB_HOST")
	if found == nil || found.Type != differ.Changed {
		t.Errorf("expected DB_HOST to be Changed")
	}
	if found.OldValue != "localhost" || found.NewValue != "prod.db.example.com" {
		t.Errorf("unexpected old/new values: %q / %q", found.OldValue, found.NewValue)
	}
}

func TestDiff_UnchangedKey(t *testing.T) {
	base := map[string]string{"PORT": "8080"}
	target := map[string]string{"PORT": "8080"}

	result := differ.Diff(base, target)

	found := findEntry(result.Entries, "PORT")
	if found == nil || found.Type != differ.Unchanged {
		t.Errorf("expected PORT to be Unchanged")
	}
}

func TestResult_HasChanges(t *testing.T) {
	base := map[string]string{"A": "1"}
	target := map[string]string{"A": "2"}

	result := differ.Diff(base, target)
	if !result.HasChanges() {
		t.Error("expected HasChanges to return true")
	}
}

func TestResult_Summary(t *testing.T) {
	base := map[string]string{"A": "1", "B": "old"}
	target := map[string]string{"B": "new", "C": "added"}

	result := differ.Diff(base, target)
	summary := result.Summary()
	if summary == "" {
		t.Error("expected non-empty summary")
	}
}

func findEntry(entries []differ.Entry, key string) *differ.Entry {
	for i := range entries {
		if entries[i].Key == key {
			return &entries[i]
		}
	}
	return nil
}
