package service

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/nodeveil/nodeveil/engine/internal/domain"
	"github.com/nodeveil/nodeveil/engine/internal/store"
)

func TestGraphCreateAndBacklinks(t *testing.T) {
	st, _ := store.Open(context.Background(), filepath.Join(t.TempDir(), "g.db"))
	_ = st.AddRoot(context.Background(), t.TempDir())
	roots, _ := st.ListRoots(context.Background())
	a := &domain.Node{RootID: roots[0].ID, AbsPath: "/tmp/a", RelPath: "a", Kind: domain.KindFile, SizeBytes: 1, MtimeUnix: 1, Mode: 420}
	b := &domain.Node{RootID: roots[0].ID, AbsPath: "/tmp/b", RelPath: "b", Kind: domain.KindFile, SizeBytes: 1, MtimeUnix: 1, Mode: 420}
	_ = st.UpsertNode(context.Background(), a)
	_ = st.UpsertNode(context.Background(), b)
	g := NewGraphService(st)
	if err := g.CreateLink(context.Background(), a.ID, b.ID, "related", ""); err != nil {
		t.Fatal(err)
	}
	out, _ := g.ListLinks(context.Background(), a.ID, domain.DirectionOut)
	in, _ := g.GetBacklinks(context.Background(), b.ID)
	if len(out) != 1 || len(in) != 1 {
		t.Fatal("expected in/out edge")
	}
}
