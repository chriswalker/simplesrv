package db

import (
	"os"
	"path"
	"testing"
)

// newDB is a helper function that creates the new test database, and
// provides teardown for it between tests.
func newDB(t *testing.T) *DB {
	file := path.Join(os.TempDir(), "test_items.db")
	db, err := New(file)
	if err != nil {
		t.Error(err)
	}
	t.Cleanup(func() {
		os.Remove(file)
	})

	return db
}

func TestGetItems(t *testing.T) {
	db := newDB(t)

	got, err := db.GetItems()
	if err != nil {
		t.Errorf("unepxected error loading items: %s", err)
	}

	if len(got) != 4 {
		t.Errorf("got %d items from db, want %d", len(got), 4)
	}
}
