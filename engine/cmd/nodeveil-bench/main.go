package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/nodeveil/nodeveil/engine/internal/ignore"
	"github.com/nodeveil/nodeveil/engine/internal/indexer"
	"github.com/nodeveil/nodeveil/engine/internal/service"
	"github.com/nodeveil/nodeveil/engine/internal/store"
)

type Out struct {
	Bench        string `json:"bench"`
	ScanMS       int64  `json:"scan_ms"`
	FilesIndexed int    `json:"files_indexed"`
	DBSizeBytes  int64  `json:"db_size_bytes"`
	RSSBytes     uint64 `json:"rss_bytes"`
	Query1000MS  int64  `json:"query_1000_ms"`
}

func main() {
	tmp, _ := os.MkdirTemp("", "nodeveil-bench")
	defer func() { _ = os.RemoveAll(tmp) }()
	root := filepath.Join(tmp, "root")
	_ = os.MkdirAll(root, 0o755)
	for i := 0; i < 10000; i++ {
		d := filepath.Join(root, fmt.Sprintf("d%d", i%100))
		_ = os.MkdirAll(d, 0o755)
		_ = os.WriteFile(filepath.Join(d, fmt.Sprintf("f%d.md", i)), []byte("x"), 0o644)
	}
	dbPath := filepath.Join(tmp, "bench.db")
	st, _ := store.Open(context.Background(), dbPath)
	svc := service.NewIndexService(st, indexer.FSScanner{Matcher: ignore.New(nil)}, indexer.PollWatcher{Interval: time.Hour})
	start := time.Now()
	if err := svc.AddRoot(context.Background(), root); err != nil {
		panic(err)
	}
	time.Sleep(4 * time.Second)
	scanMS := time.Since(start).Milliseconds()
	count, _ := st.CountActiveFiles(context.Background())
	qstart := time.Now()
	for i := 0; i < 1000; i++ {
		_, _ = svc.SearchFiles(context.Background(), "f1", ".md", 0, "file")
	}
	qms := time.Since(qstart).Milliseconds()
	fi, _ := os.Stat(dbPath)
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	_ = json.NewEncoder(os.Stdout).Encode(Out{Bench: "phase1_index", ScanMS: scanMS, FilesIndexed: count, DBSizeBytes: fi.Size(), RSSBytes: m.Sys, Query1000MS: qms})
	svc.Stop()
}
