# Nodeveil

Nodeveil is a filesystem-native hard-link knowledge graph desktop app.

## Dev

```bash
cd engine && go run ./cmd/nodeveil-engine
pnpm -C desktop build
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

## Benches

```bash
cd engine && go run ./cmd/nodeveil-bench
pnpm -C desktop bench:perf
pnpm -C desktop bench:graph
```

## Phase 4 highlights

- Dedicated Graph View with focus-node neighborhood loading
- Worker-driven layout, WebGL-based render surface, pan/zoom/select/focus interactions
- Graph filters (direction/depth/relation/ext) with capped neighborhood payloads
- Graph perf bench + baseline JSON and graph e2e smoke
