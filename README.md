# Nodeveil

Nodeveil is a filesystem-native hard-link knowledge graph desktop app.

## Key features
- Local-first indexing and graph linking
- Backlinks/outgoing links and neighborhood graph view
- Runtime themes: Light, Dark, Neon, Sakura
- Graph export/import + backup

## Install / run (dev)
```bash
cd engine && go run ./cmd/nodeveil-engine
pnpm -C desktop build
```

## Building from source
```bash
./scripts/build-engine.sh
./scripts/build-desktop.sh
# Windows: scripts/build-engine.ps1 and scripts/build-desktop.ps1
```

## Packaging commands
```bash
pnpm -C desktop package:linux
pnpm -C desktop package:win
pnpm -C desktop package:mac
```

## Tests
```bash
pnpm -C desktop lint
pnpm -C desktop typecheck
pnpm -C desktop test
RUN_E2E=1 pnpm -C desktop test:e2e
cd engine && go test ./...
cd engine && go test -race ./...
cd engine && go vet ./...
cd engine && golangci-lint run ./...
```

## Benchmarks
```bash
cd engine && go run ./cmd/nodeveil-bench
pnpm -C desktop bench:perf
pnpm -C desktop bench:graph
pnpm -C desktop bench:theme
```

## License summary
- Official binaries: non-commercial use only.
- Distributed modifications must publish full source under same license.
- Nodeveil trademarks are not granted.
