package service

import (
	"context"

	"github.com/nodeveil/nodeveil/engine/internal/domain"
	"github.com/nodeveil/nodeveil/engine/internal/store"
)

type GraphQueryService struct{ store *store.Store }

func NewGraphQueryService(st *store.Store) *GraphQueryService { return &GraphQueryService{store: st} }

func (g *GraphQueryService) GetNeighborhood(ctx context.Context, nodeID string, depth int, filters domain.GraphFilters, cursor string) (domain.GraphNeighborhood, error) {
	return g.store.QueryNeighborhood(ctx, nodeID, depth, filters, cursor)
}

func (g *GraphQueryService) GetCounts(ctx context.Context, nodeID string) (domain.GraphCounts, error) {
	return g.store.QueryCounts(ctx, nodeID)
}
