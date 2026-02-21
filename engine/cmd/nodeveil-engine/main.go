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
	"github.com/nodeveil/nodeveil/engine/internal/version"
)

func main() {
	showVersion := flag.Bool("version", false, "print build info")
	flag.Parse()

	if *showVersion {
		if err := json.NewEncoder(os.Stdout).Encode(version.Get()); err != nil {
			fmt.Fprintf(os.Stderr, "encode version: %v\n", err)
			os.Exit(1)
		}
		return
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cfg := config.Load()
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
