package store

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/nodeveil/nodeveil/engine/internal/domain"
)

func seedNode(t *testing.T, s *Store, root string) string {
	t.Helper()
	_ = s.AddRoot(context.Background(), root)
	roots, _ := s.ListRoots(context.Background())
	n := &domain.Node{RootID: roots[0].ID, AbsPath: filepath.Join(root, "a.md"), RelPath: "a.md", Kind: domain.KindFile, Ext: ".md", SizeBytes: 1, MtimeUnix: 1, Mode: 420}
	if err := s.UpsertNode(context.Background(), n); err != nil {
		t.Fatal(err)
	}
	return n.ID
}

func TestEdgeInvariantsAndSoftDelete(t *testing.T) {
	s := newStore(t)
	r := t.TempDir()
	a := seedNode(t, s, r)
	n2 := &domain.Node{RootID: 1, AbsPath: filepath.Join(r, "b.md"), RelPath: "b.md", Kind: domain.KindFile, Ext: ".md", SizeBytes: 1, MtimeUnix: 1, Mode: 420}
	if err := s.UpsertNode(context.Background(), n2); err != nil {
		t.Fatal(err)
	}
	if err := s.CreateLink(context.Background(), a, a, "related", ""); err == nil {
		t.Fatal("expected self loop rejection")
	}
	if err := s.CreateLink(context.Background(), a, n2.ID, "related", "n"); err != nil {
		t.Fatal(err)
	}
	if err := s.CreateLink(context.Background(), a, n2.ID, "related", "n2"); err != nil {
		t.Fatal(err)
	}
	edges, _ := s.ListLinks(context.Background(), a, domain.DirectionOut)
	if len(edges) != 1 {
		t.Fatalf("expected idempotent edge, got %d", len(edges))
	}
	if err := s.RemoveLink(context.Background(), edges[0].ID); err != nil {
		t.Fatal(err)
	}
	edges, _ = s.ListLinks(context.Background(), a, domain.DirectionOut)
	if len(edges) != 0 {
		t.Fatal("expected soft-deleted edge hidden")
	}
}

func TestImportMappingRules(t *testing.T) {
	s := newStore(t)
	r := t.TempDir()
	a := seedNode(t, s, r)
	n2 := &domain.Node{RootID: 1, AbsPath: filepath.Join(r, "b.md"), RelPath: "b.md", Kind: domain.KindFile, Ext: ".md", SizeBytes: 1, MtimeUnix: 1, Mode: 420}
	_ = s.UpsertNode(context.Background(), n2)
	exp := domain.GraphExport{Version: 1, Nodes: []domain.NodeExport{{ID: a, AbsPath: filepath.Join(r, "a.md")}, {ID: n2.ID, AbsPath: filepath.Join(r, "b.md")}, {ID: "missing", AbsPath: filepath.Join(r, "missing.md")}}, Edges: []domain.Edge{{FromID: a, ToID: n2.ID, RelationType: "related"}, {FromID: "missing", ToID: n2.ID, RelationType: "related"}}}
	rep, err := s.ImportGraph(context.Background(), exp)
	if err != nil {
		t.Fatal(err)
	}
	if rep.Added != 1 || rep.Skipped != 1 {
		t.Fatalf("unexpected report %+v", rep)
	}
}

func TestMoveNodeKeepsEdges(t *testing.T) {
	s := newStore(t)
	r := t.TempDir()
	a := seedNode(t, s, r)
	n2 := &domain.Node{RootID: 1, AbsPath: filepath.Join(r, "b.md"), RelPath: "b.md", Kind: domain.KindFile, Ext: ".md", SizeBytes: 1, MtimeUnix: 1, Mode: 420}
	_ = s.UpsertNode(context.Background(), n2)
	if err := s.CreateLink(context.Background(), a, n2.ID, "related", ""); err != nil {
		t.Fatal(err)
	}
	if err := s.MoveNodePath(context.Background(), 1, filepath.Join(r, "a.md"), filepath.Join(r, "renamed.md")); err != nil {
		t.Fatal(err)
	}
	edges, _ := s.ListLinks(context.Background(), a, domain.DirectionOut)
	if len(edges) != 1 {
		t.Fatal("edge should survive move")
	}
}
