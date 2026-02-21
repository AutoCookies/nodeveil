package main

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"time"

	"github.com/nodeveil/nodeveil/engine/internal/domain"
	"github.com/nodeveil/nodeveil/engine/internal/ignore"
	"github.com/nodeveil/nodeveil/engine/internal/indexer"
	"github.com/nodeveil/nodeveil/engine/internal/service"
	"github.com/nodeveil/nodeveil/engine/internal/store"
)

type Out struct {
	Bench          string  `json:"bench"`
	ScanMS         int64   `json:"scan_ms"`
	FilesIndexed   int     `json:"files_indexed"`
	EdgesInsertMS  int64   `json:"edges_insert_ms"`
	BacklinksP50MS float64 `json:"backlinks_p50_ms"`
	BacklinksP95MS float64 `json:"backlinks_p95_ms"`
	DBSizeBytes    int64   `json:"db_size_bytes"`
	RSSBytes       uint64  `json:"rss_bytes"`
}

func main() {
	tmp, _ := os.MkdirTemp("", "nodeveil-bench")
	defer func() { _ = os.RemoveAll(tmp) }()
	root := filepath.Join(tmp, "root")
	_ = os.MkdirAll(root, 0o755)
	for i := 0; i < 10000; i++ {
		d := filepath.Join(root, fmt.Sprintf("d%d", i%200))
		_ = os.MkdirAll(d, 0o755)
		_ = os.WriteFile(filepath.Join(d, fmt.Sprintf("f%d.md", i)), []byte("x"), 0o644)
	}
	st, _ := store.Open(context.Background(), filepath.Join(tmp, "bench.db"))
	idx := service.NewIndexService(st, indexer.FSScanner{Matcher: ignore.New(nil)}, indexer.PollWatcher{Interval: time.Hour})
	start := time.Now()
	_ = idx.AddRoot(context.Background(), root)
	time.Sleep(5 * time.Second)
	scanMS := time.Since(start).Milliseconds()
	nodes, _ := st.SearchFiles(context.Background(), "", "", 0, "file")

	g := service.NewGraphService(st)
	r := rand.New(rand.NewSource(7))
	eStart := time.Now()
	batch := make([]domain.Edge, 0, 1000)
	for i := 0; i < 50000 && len(nodes) > 1; i++ {
		a := nodes[r.Intn(len(nodes))].ID
		b := nodes[r.Intn(len(nodes))].ID
		if a == b {
			continue
		}
		batch = append(batch, domain.Edge{FromID: a, ToID: b, RelationType: domain.RelationRelated})
		if len(batch) == 1000 {
			_ = g.BatchCreateLinks(context.Background(), batch)
			batch = batch[:0]
		}
	}
	if len(batch) > 0 {
		_ = g.BatchCreateLinks(context.Background(), batch)
	}
	edgesMS := time.Since(eStart).Milliseconds()

	lat := make([]float64, 0, 300)
	for i := 0; i < 300 && len(nodes) > 0; i++ {
		n := nodes[r.Intn(len(nodes))].ID
		ts := time.Now()
		_, _ = g.GetBacklinks(context.Background(), n)
		lat = append(lat, float64(time.Since(ts).Microseconds())/1000.0)
	}
	sort.Float64s(lat)
	p50 := 0.0
	p95 := 0.0
	if len(lat) > 0 {
		p50 = lat[len(lat)/2]
		p95 = lat[(len(lat)*95)/100]
	}
	fi, _ := os.Stat(filepath.Join(tmp, "bench.db"))
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	_ = json.NewEncoder(os.Stdout).Encode(Out{Bench: "phase2_graph", ScanMS: scanMS, FilesIndexed: len(nodes), EdgesInsertMS: edgesMS, BacklinksP50MS: p50, BacklinksP95MS: p95, DBSizeBytes: fi.Size(), RSSBytes: m.Sys})
	idx.Stop()
}
