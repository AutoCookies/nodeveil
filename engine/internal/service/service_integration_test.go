package service

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/nodeveil/nodeveil/engine/internal/ignore"
	"github.com/nodeveil/nodeveil/engine/internal/indexer"
	"github.com/nodeveil/nodeveil/engine/internal/store"
)

func TestScanDeleteRenameFlow(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "a.txt")
	if err := os.WriteFile(f, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	st, _ := store.Open(context.Background(), filepath.Join(t.TempDir(), "n.db"))
	svc := NewIndexService(st, indexer.FSScanner{Matcher: ignore.New(nil)}, indexer.PollWatcher{Interval: 200 * time.Millisecond})
	if err := svc.AddRoot(context.Background(), dir); err != nil {
		t.Fatal(err)
	}
	time.Sleep(500 * time.Millisecond)
	res, _ := svc.SearchFiles(context.Background(), "a.txt", "", 0, "")
	if len(res) == 0 {
		t.Fatal("scan missing file")
	}
	if err := os.Rename(f, filepath.Join(dir, "b.txt")); err != nil {
		t.Fatal(err)
	}
	time.Sleep(700 * time.Millisecond)
	res2, _ := svc.SearchFiles(context.Background(), "b.txt", "", 0, "")
	if len(res2) == 0 {
		t.Fatal("rename not reflected")
	}
	if err := os.Remove(filepath.Join(dir, "b.txt")); err != nil {
		t.Fatal(err)
	}
	time.Sleep(700 * time.Millisecond)
	res3, _ := svc.SearchFiles(context.Background(), "b.txt", "", 0, "")
	if len(res3) != 0 {
		t.Fatal("delete not reflected")
	}
	svc.Stop()
}
