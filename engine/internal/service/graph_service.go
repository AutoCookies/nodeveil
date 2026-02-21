package service

import (
	"context"
	"os"

	"github.com/nodeveil/nodeveil/engine/internal/domain"
	"github.com/nodeveil/nodeveil/engine/internal/store"
)

type GraphService struct{ store *store.Store }

func NewGraphService(st *store.Store) *GraphService { return &GraphService{store: st} }

func (g *GraphService) CreateLink(ctx context.Context, fromID, toID, relationType, note string) error {
	return g.store.CreateLink(ctx, fromID, toID, relationType, note)
}
func (g *GraphService) RemoveLink(ctx context.Context, linkID string) error {
	return g.store.RemoveLink(ctx, linkID)
}
func (g *GraphService) BatchCreateLinks(ctx context.Context, edges []domain.Edge) error {
	return g.store.BatchCreateLinks(ctx, edges)
}
func (g *GraphService) RemoveLinkByNodes(ctx context.Context, fromID, toID, relationType string) error {
	return g.store.RemoveLinkByNodes(ctx, fromID, toID, relationType)
}
func (g *GraphService) ListLinks(ctx context.Context, nodeID string, direction domain.Direction) ([]domain.Edge, error) {
	return g.store.ListLinks(ctx, nodeID, direction)
}
func (g *GraphService) GetBacklinks(ctx context.Context, nodeID string) ([]domain.Edge, error) {
	return g.store.ListLinks(ctx, nodeID, domain.DirectionIn)
}
func (g *GraphService) GetNeighbors(ctx context.Context, nodeID string, _ int) ([]domain.Edge, error) {
	return g.store.GetNeighbors(ctx, nodeID)
}
func (g *GraphService) GetGraphStats(ctx context.Context) (domain.GraphStats, error) {
	return g.store.GraphStats(ctx)
}
func (g *GraphService) Events(ctx context.Context, limit int) ([]map[string]any, error) {
	return g.store.RecentEvents(ctx, limit)
}
func (g *GraphService) ExportGraph(ctx context.Context, outFile string) error {
	exp, err := g.store.ExportGraph(ctx)
	if err != nil {
		return err
	}
	body, err := store.EncodeExport(exp)
	if err != nil {
		return err
	}
	return os.WriteFile(outFile, body, 0o644)
}
func (g *GraphService) ImportGraph(ctx context.Context, inFile string) (domain.ImportReport, error) {
	raw, err := os.ReadFile(inFile)
	if err != nil {
		return domain.ImportReport{}, err
	}
	exp, err := store.DecodeExport(raw)
	if err != nil {
		return domain.ImportReport{}, err
	}
	return g.store.ImportGraph(ctx, exp)
}
