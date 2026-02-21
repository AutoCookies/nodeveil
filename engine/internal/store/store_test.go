package store

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/nodeveil/nodeveil/engine/internal/domain"
)

func newStore(t *testing.T) *Store {
	t.Helper()
	p := filepath.Join(t.TempDir(), "test.db")
	s, err := Open(context.Background(), p)
	if err != nil {
		t.Fatal(err)
	}
	return s
}

func TestPathNormalize(t *testing.T) {
	if NormalizePath("a/../b") != "b" {
		t.Fatal("normalize failed")
	}
}

func TestRootCRUDAndNodeUpsertIdempotent(t *testing.T) {
	s := newStore(t)
	ctx := context.Background()
	rd := t.TempDir()
	if err := s.AddRoot(ctx, rd); err != nil {
		t.Fatal(err)
	}
	roots, _ := s.ListRoots(ctx)
	if len(roots) != 1 {
		t.Fatal("expected root")
	}
	n := &domain.Node{RootID: roots[0].ID, AbsPath: filepath.Join(rd, "a.txt"), RelPath: "a.txt", Kind: domain.KindFile, Ext: ".txt", SizeBytes: 1, MtimeUnix: 1, Mode: 420}
	if err := s.UpsertNode(ctx, n); err != nil {
		t.Fatal(err)
	}
	n.SizeBytes = 2
	if err := s.UpsertNode(ctx, n); err != nil {
		t.Fatal(err)
	}
	res, _ := s.SearchFiles(ctx, "a.txt", "", roots[0].ID, "")
	if len(res) != 1 || res[0].SizeBytes != 2 {
		t.Fatal("idempotent upsert failed")
	}
	_ = os.Remove(n.AbsPath)
}
