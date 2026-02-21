package app

import (
	"context"
	"fmt"
	"time"

	"github.com/nodeveil/nodeveil/engine/internal/config"
	"github.com/nodeveil/nodeveil/engine/internal/ignore"
	"github.com/nodeveil/nodeveil/engine/internal/indexer"
	"github.com/nodeveil/nodeveil/engine/internal/ipc"
	"github.com/nodeveil/nodeveil/engine/internal/service"
	"github.com/nodeveil/nodeveil/engine/internal/store"
)

type App struct {
	store  *store.Store
	server *ipc.Server
	index  *service.IndexService
}

func New(ctx context.Context, cfg config.Config) (*App, error) {
	st, err := store.Open(ctx, cfg.DBPath)
	if err != nil {
		return nil, fmt.Errorf("open store: %w", err)
	}
	sc := indexer.FSScanner{Matcher: ignore.New(nil)}
	wa := indexer.PollWatcher{Interval: 2 * time.Second}
	idx := service.NewIndexService(st, sc, wa)
	if err := idx.Start(ctx); err != nil {
		return nil, err
	}
	return &App{store: st, index: idx, server: ipc.New(cfg.HTTPAddr, idx)}, nil
}
func (a *App) Run() error   { return a.server.Start() }
func (a *App) Close() error { a.index.Stop(); return a.store.Close() }
