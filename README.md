# Nodeveil

Nodeveil is a filesystem-native hard-link knowledge graph desktop app.

## Dev run

```bash
cd engine && go run ./cmd/nodeveil-engine
pnpm -C desktop build
```

## Graph CLI

```bash
cd engine && go run ./cmd/nodeveil-engine --graph-export graph.json
cd engine && go run ./cmd/nodeveil-engine --graph-import graph.json
```

## Tests

```bash
pnpm -C desktop lint
pnpm -C desktop typecheck
pnpm -C desktop test
pnpm -C desktop test:e2e
cd engine && go test ./...
cd engine && go test -race ./...
cd engine && go vet ./...
cd engine && golangci-lint run ./...
```

## Benches

```bash
cd engine && go run ./cmd/nodeveil-bench
pnpm -C desktop bench:perf
```

## Phase 3 highlights

- Three-pane UX (navigator/viewer/links)
- Link create/remove + undo for linking actions
- Backlinks/outgoing tabs and relation badges
- Engine reconnect/degraded handling and diagnostics endpoints
- E2E smoke + perf bench scripts
