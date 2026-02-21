package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/nodeveil/nodeveil/engine/internal/app"
	"github.com/nodeveil/nodeveil/engine/internal/config"
	"github.com/nodeveil/nodeveil/engine/internal/ignore"
	"github.com/nodeveil/nodeveil/engine/internal/indexer"
	"github.com/nodeveil/nodeveil/engine/internal/service"
	"github.com/nodeveil/nodeveil/engine/internal/store"
	"github.com/nodeveil/nodeveil/engine/internal/version"
)

func main() {
	showVersion := flag.Bool("version", false, "print build info")
	graphExport := flag.String("graph-export", "", "export graph to file")
	graphImport := flag.String("graph-import", "", "import graph from file")
	flag.Parse()

	if *showVersion {
		_ = json.NewEncoder(os.Stdout).Encode(version.Get())
		return
	}
	cfg := config.Load()
	if *graphExport != "" || *graphImport != "" {
		st, err := store.Open(context.Background(), cfg.DBPath)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		gs := service.NewGraphService(st)
		if *graphExport != "" {
			if err := gs.ExportGraph(context.Background(), *graphExport); err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			fmt.Println("exported", *graphExport)
			return
		}
		rep, err := gs.ImportGraph(context.Background(), *graphImport)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		_ = json.NewEncoder(os.Stdout).Encode(rep)
		return
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	_ = ignore.New(nil)
	_ = indexer.PollWatcher{}
	application, err := app.New(ctx, cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "create app: %v\n", err)
		os.Exit(1)
	}
	defer func() { _ = application.Close() }()
	if err := application.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "run app: %v\n", err)
		os.Exit(1)
	}
}
