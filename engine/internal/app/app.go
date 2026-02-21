package app

import (
	"context"
	"fmt"

	"github.com/nodeveil/nodeveil/engine/internal/config"
	"github.com/nodeveil/nodeveil/engine/internal/db"
	"github.com/nodeveil/nodeveil/engine/internal/ipc"
)

type App struct {
	store  *db.Store
	server *ipc.Server
}

func New(ctx context.Context, cfg config.Config) (*App, error) {
	store, err := db.Open(ctx, cfg.DBPath)
	if err != nil {
		return nil, fmt.Errorf("open store: %w", err)
	}
	server := ipc.New(cfg.HTTPAddr, store)
	return &App{store: store, server: server}, nil
}

func (a *App) Run() error { return a.server.Start() }

func (a *App) Close() error {
	if err := a.store.Close(); err != nil {
		return fmt.Errorf("close store: %w", err)
	}
	return nil
}
