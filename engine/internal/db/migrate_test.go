package db

import (
	"context"
	"path/filepath"
	"testing"
)

func TestMigrateMarksStore(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "nodeveil.db")
	store, err := Open(context.Background(), dbPath)
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	defer func() { _ = store.Close() }()

	if !store.IsMigrated() {
		t.Fatal("expected migrated store")
	}

	out, err := store.exec(context.Background(), "SELECT name FROM sqlite_master WHERE type='table' AND name='roots';")
	if err != nil {
		t.Fatalf("query sqlite_master: %v", err)
	}
	if out == "" {
		t.Fatal("expected roots table")
	}
}
