package indexer

import (
	"context"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/nodeveil/nodeveil/engine/internal/domain"
	"github.com/nodeveil/nodeveil/engine/internal/ignore"
)

type Scanner interface {
	ScanRoot(context.Context, domain.Root, chan<- domain.IndexEvent) error
}

type FSScanner struct{ Matcher *ignore.Matcher }

func (s FSScanner) ScanRoot(ctx context.Context, root domain.Root, out chan<- domain.IndexEvent) error {
	return filepath.WalkDir(root.Path, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		rel, rerr := filepath.Rel(root.Path, path)
		if rerr != nil || rel == "." {
			return nil
		}
		if s.Matcher.Ignore(rel) {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		info, ierr := d.Info()
		if ierr != nil {
			return nil
		}
		kind := domain.KindFile
		if d.Type()&os.ModeSymlink != 0 {
			kind = domain.KindSymlink
		}
		if d.IsDir() {
			kind = domain.KindDir
		}
		mtime := info.ModTime().Unix()
		n := &domain.Node{RootID: root.ID, AbsPath: path, RelPath: filepath.ToSlash(rel), Kind: kind, SizeBytes: info.Size(), MtimeUnix: mtime, Mode: uint32(info.Mode()), Ext: strings.ToLower(filepath.Ext(path)), UpdatedAt: time.Now().Unix(), CreatedAt: time.Now().Unix()}
		out <- domain.IndexEvent{Type: domain.EventUpsert, RootID: root.ID, Path: path, Node: n}
		return nil
	})
}
