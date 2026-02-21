package service

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"time"

	"github.com/nodeveil/nodeveil/engine/internal/domain"
	"github.com/nodeveil/nodeveil/engine/internal/indexer"
	"github.com/nodeveil/nodeveil/engine/internal/store"
)

type IndexService struct {
	store   *store.Store
	scanner indexer.Scanner
	watcher indexer.Watcher
	status  domain.IndexStatus
	events  chan domain.IndexEvent
	ctx     context.Context
	cancel  context.CancelFunc
}

func NewIndexService(st *store.Store, sc indexer.Scanner, w indexer.Watcher) *IndexService {
	ctx, cancel := context.WithCancel(context.Background())
	s := &IndexService{store: st, scanner: sc, watcher: w, events: make(chan domain.IndexEvent, 4096), ctx: ctx, cancel: cancel}
	go s.consume()
	return s
}

func (s *IndexService) consume() {
	for {
		select {
		case <-s.ctx.Done():
			return
		case ev := <-s.events:
			s.status.QueueDepth = len(s.events)
			s.applyEvent(ev)
		}
	}
}

func (s *IndexService) applyEvent(ev domain.IndexEvent) {
	ctx := context.Background()
	switch ev.Type {
	case domain.EventDelete:
		if err := s.store.SoftDeleteByPath(ctx, ev.RootID, ev.Path); err != nil {
			atomic.AddInt64(&s.status.Errors, 1)
			return
		}
		atomic.AddInt64(&s.status.DBUpserts, 1)
	case domain.EventUpsert:
		n := ev.Node
		if n == nil {
			fi, err := os.Lstat(ev.Path)
			if err != nil {
				return
			}
			kind := domain.KindFile
			if fi.IsDir() {
				kind = domain.KindDir
			}
			rel := filepath.Base(ev.Path)
			n = &domain.Node{RootID: ev.RootID, AbsPath: ev.Path, RelPath: rel, Kind: kind, Ext: filepath.Ext(ev.Path), SizeBytes: fi.Size(), MtimeUnix: fi.ModTime().Unix(), Mode: uint32(fi.Mode()), UpdatedAt: time.Now().Unix(), CreatedAt: time.Now().Unix()}
		}
		if err := s.store.UpsertNode(ctx, n); err != nil {
			atomic.AddInt64(&s.status.Errors, 1)
			return
		}
		atomic.AddInt64(&s.status.DBUpserts, 1)
	}
}

func (s *IndexService) Start(ctx context.Context) error {
	roots, err := s.store.ListRoots(ctx)
	if err != nil {
		return err
	}
	for _, r := range roots {
		rc := r
		go func(rr domain.Root) { _ = s.fastInitialScan(rr); _ = s.watcher.Watch(s.ctx, rr, s.events) }(rc)
		go s.reconcileLoop(rc)
	}
	return nil
}

func (s *IndexService) Stop() { s.cancel() }

func (s *IndexService) AddRoot(ctx context.Context, p string) error {
	info, err := os.Stat(p)
	if err != nil || !info.IsDir() {
		return os.ErrInvalid
	}
	if err := s.store.AddRoot(ctx, p); err != nil {
		return err
	}
	roots, _ := s.store.ListRoots(ctx)
	for _, r := range roots {
		if r.Path == filepath.Clean(p) {
			go func(rr domain.Root) {
				_ = s.fastInitialScan(rr)
				_ = s.watcher.Watch(s.ctx, rr, s.events)
			}(r)
			go s.reconcileLoop(r)
		}
	}
	return nil
}
func (s *IndexService) RemoveRoot(ctx context.Context, p string) error {
	return s.store.RemoveRoot(ctx, p)
}
func (s *IndexService) ListRoots(ctx context.Context) ([]domain.Root, error) {
	return s.store.ListRoots(ctx)
}
func (s *IndexService) SearchFiles(ctx context.Context, q, ext string, rootID int64, kind string) ([]domain.Node, error) {
	return s.store.SearchFiles(ctx, q, ext, rootID, kind)
}
func (s *IndexService) GetFileByID(ctx context.Context, id string) (*domain.Node, error) {
	return s.store.GetFileByID(ctx, id)
}
func (s *IndexService) ListRecentChanges(ctx context.Context, since int64) ([]domain.Node, error) {
	return s.store.RecentChanges(ctx, since)
}
func (s *IndexService) GetIndexStatus() domain.IndexStatus {
	s.status.QueueDepth = len(s.events)
	return s.status
}

func (s *IndexService) reconcileLoop(root domain.Root) {
	t := time.NewTicker(3 * time.Second)
	defer t.Stop()
	for {
		select {
		case <-s.ctx.Done():
			return
		case <-t.C:
			_ = s.scanner.ScanRoot(s.ctx, root, s.events)
		}
	}
}

func (s *IndexService) fastInitialScan(root domain.Root) error {
	nodes := make([]*domain.Node, 0, 2048)
	err := filepath.WalkDir(root.Path, func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}
		rel, rerr := filepath.Rel(root.Path, path)
		if rerr != nil {
			return nil
		}
		if strings.HasPrefix(rel, ".git") || strings.Contains(rel, "node_modules") {
			return nil
		}
		fi, ierr := d.Info()
		if ierr != nil {
			return nil
		}
		nodes = append(nodes, &domain.Node{RootID: root.ID, AbsPath: path, RelPath: filepath.ToSlash(rel), Kind: domain.KindFile, Ext: strings.ToLower(filepath.Ext(path)), SizeBytes: fi.Size(), MtimeUnix: fi.ModTime().Unix(), Mode: uint32(fi.Mode()), UpdatedAt: time.Now().Unix(), CreatedAt: time.Now().Unix()})
		return nil
	})
	if err != nil {
		return err
	}
	if err := s.store.BulkUpsert(context.Background(), nodes); err != nil {
		return err
	}
	atomic.AddInt64(&s.status.DBUpserts, int64(len(nodes)))
	atomic.AddInt64(&s.status.FilesIndexed, int64(len(nodes)))
	return nil
}
