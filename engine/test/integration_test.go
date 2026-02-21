package test

import (
	"context"
	"testing"

	"github.com/nodeveil/nodeveil/engine/internal/config"
)

func TestConfigDefaults(t *testing.T) {
	_ = context.Background()
	cfg := config.Load()
	if cfg.HTTPAddr == "" || cfg.DBPath == "" {
		t.Fatal("expected non-empty defaults")
	}
}
