package test

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/nodeveil/nodeveil/engine/internal/ignore"
	"github.com/nodeveil/nodeveil/engine/internal/indexer"
	"github.com/nodeveil/nodeveil/engine/internal/service"
	"github.com/nodeveil/nodeveil/engine/internal/store"
)

func TestRandomOpsConverge(t *testing.T) {
	dir := t.TempDir()
	st, _ := store.Open(context.Background(), filepath.Join(t.TempDir(), "s.db"))
	svc := service.NewIndexService(st, indexer.FSScanner{Matcher: ignore.New(nil)}, indexer.PollWatcher{Interval: 100 * time.Millisecond})
	if err := svc.AddRoot(context.Background(), dir); err != nil { t.Fatal(err) }
	r := rand.New(rand.NewSource(1))
	for i := 0; i < 50; i++ {
		p := filepath.Join(dir, fmt.Sprintf("f%d.txt", r.Intn(20)))
		_ = os.WriteFile(p, []byte("x"), 0o644)
		if r.Intn(3) == 0 { _ = os.Remove(p) }
	}
	time.Sleep(2 * time.Second)
	rows, _ := svc.SearchFiles(context.Background(), "", ".txt", 0, "file")
	set := map[string]struct{}{}
	for _, n := range rows { set[n.AbsPath] = struct{}{} }
	_ = filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err == nil && !d.IsDir() {
			if _, ok := set[path]; !ok { t.Fatalf("missing file in db snapshot: %s", path) }
		}
		return nil
	})
	svc.Stop()
}
