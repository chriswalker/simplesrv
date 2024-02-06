package service

import (
	"io"
	"log/slog"
	"os"
	"path"
	"strings"
	"testing"
	"time"

	"github.com/chriswalker/simplesrv/model"
)

// The below are ratjer bare-bones tests for the service, upon
// which to build.

type testDB struct {
	items []model.Item
}

func (d *testDB) GetItems() ([]model.Item, error) {
	return d.items, nil
}

// newTestDB creates an in-memory database for testing.
func newTestDB() *testDB {
	return &testDB{
		items: []model.Item{
			{
				ID:           1,
				Title:        "Item 1",
				Description:  "Description for Item 1",
				CreationDate: time.Now(),
			},
			{
				ID:           2,
				Title:        "Item 2",
				Description:  "Description for Item 2",
				CreationDate: time.Now(),
			},
		},
	}
}

// newTestService creates a service composed of test
// components.
func newTestService() *ItemService {
	return &ItemService{
		logger: slog.New(slog.NewTextHandler(io.Discard, nil)),
		db:     newTestDB(),
	}
}

func TestNewItemService(t *testing.T) {
	testCases := map[string]struct {
		path    string
		errMsg  string
		cleanup bool
	}{
		// Trying to creater in a root-owned directory; should fail.
		"fail": {
			path:   "/etc/shouldfail.db",
			errMsg: "unable to open database file",
		},
		"success": {
			path:    path.Join(os.TempDir(), "shouldwork.db"),
			cleanup: true,
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			_, err := NewItemService(tc.path, nil)

			if tc.errMsg != "" {
				if err == nil {
					t.Errorf("expected an error, got nil")
					return
				}
				if !strings.Contains(err.Error(), tc.errMsg) {
					t.Errorf("got error '%s', want '%s'", err, tc.errMsg)
				}
				return
			}

			if tc.cleanup {
				err := os.Remove(tc.path)
				if err != nil {
					t.Fatalf("could not delete test file: %q", err)
				}
			}
		})
	}
}

func TestGetItems(t *testing.T) {
	svc := newTestService()
	items, err := svc.GetItems()
	if err != nil {
		t.Errorf("unexpected error: %q", err)
	}
	if len(items) != 2 {
		t.Errorf("got %d items, want %d", len(items), 2)
	}
}
