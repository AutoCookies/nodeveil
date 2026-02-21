package indexer

import (
	"context"
	"os"
	"path/filepath"
	"time"

	"github.com/nodeveil/nodeveil/engine/internal/domain"
)

type Watcher interface {
	Watch(context.Context, domain.Root, chan<- domain.IndexEvent) error
}

type PollWatcher struct{ Interval time.Duration }

func (w PollWatcher) Watch(ctx context.Context, root domain.Root, out chan<- domain.IndexEvent) error {
	if w.Interval <= 0 {
		w.Interval = 2 * time.Second
	}
	snap := map[string]time.Time{}
	t := time.NewTicker(w.Interval)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-t.C:
			current := map[string]time.Time{}
			_ = filepath.WalkDir(root.Path, func(path string, d os.DirEntry, err error) error {
				if err != nil || d.IsDir() {
					return nil
				}
				info, ierr := d.Info()
				if ierr != nil {
					return nil
				}
				current[path] = info.ModTime()
				if old, ok := snap[path]; !ok || !old.Equal(info.ModTime()) {
					out <- domain.IndexEvent{Type: domain.EventUpsert, RootID: root.ID, Path: path}
				}
				return nil
			})
			for p := range snap {
				if _, ok := current[p]; !ok {
					out <- domain.IndexEvent{Type: domain.EventDelete, RootID: root.ID, Path: p}
				}
			}
			snap = current
		}
	}
}
