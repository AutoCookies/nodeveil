package db

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/nodeveil/nodeveil/engine/internal/graph"
)

type RootRepository interface {
	ListRoots(ctx context.Context) ([]graph.Root, error)
}

type Store struct {
	path     string
	migrated bool
}

func Open(ctx context.Context, path string) (*Store, error) {
	store := &Store{path: path}
	if err := store.Migrate(ctx); err != nil {
		return nil, err
	}
	return store, nil
}

func (s *Store) Close() error {
	return nil
}

func (s *Store) IsMigrated() bool { return s.migrated }

func (s *Store) Migrate(ctx context.Context) error {
	if _, err := s.exec(ctx,
		"PRAGMA journal_mode=WAL;",
		"CREATE TABLE IF NOT EXISTS schema_migrations(version INTEGER PRIMARY KEY);",
		"CREATE TABLE IF NOT EXISTS roots(path TEXT PRIMARY KEY NOT NULL);",
		"INSERT OR IGNORE INTO schema_migrations(version) VALUES (1);"); err != nil {
		return fmt.Errorf("run migrations: %w", err)
	}
	s.migrated = true
	return nil
}

func (s *Store) ListRoots(ctx context.Context) ([]graph.Root, error) {
	out, err := s.exec(ctx, "SELECT path FROM roots ORDER BY path;")
	if err != nil {
		return nil, fmt.Errorf("query roots: %w", err)
	}
	if strings.TrimSpace(out) == "" {
		return []graph.Root{}, nil
	}
	lines := strings.Split(strings.TrimSpace(out), "\n")
	roots := make([]graph.Root, 0, len(lines))
	for _, line := range lines {
		root, err := graph.NewRoot(line)
		if err != nil {
			return nil, fmt.Errorf("validate root: %w", err)
		}
		roots = append(roots, root)
	}
	return roots, nil
}

func (s *Store) exec(ctx context.Context, stmts ...string) (string, error) {
	dbPath := s.path
	if dbPath == ":memory:" {
		tmp, err := os.CreateTemp("", "nodeveil-*.db")
		if err != nil {
			return "", fmt.Errorf("create temp db: %w", err)
		}
		dbPath = tmp.Name()
		_ = tmp.Close()
		defer func() { _ = os.Remove(dbPath) }()
	}

	joined := strings.Join(stmts, " ")
	cmd := exec.CommandContext(ctx, "sqlite3", dbPath, joined)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("sqlite3 command failed: %w (%s)", err, strings.TrimSpace(string(output)))
	}
	return string(output), nil
}
