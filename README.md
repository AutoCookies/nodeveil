# Nodeveil

Nodeveil is a filesystem-native hard-link knowledge graph desktop app. Phase 1 delivers a production-grade local file index engine with incremental updates and crash-safe SQLite persistence.

## Getting started

```bash
cd engine && go mod download
pnpm -C desktop install
```

## Run

```bash
cd engine && go run ./cmd/nodeveil-engine
# in another shell
pnpm -C desktop build
```

## Checks

```bash
pnpm -C desktop lint
pnpm -C desktop typecheck
pnpm -C desktop test
cd engine && go test ./...
cd engine && go test -race ./...
cd engine && golangci-lint run ./...
cd engine && go run ./cmd/nodeveil-bench
```

## License summary

- Non-commercial use for official binaries
- Distributed modifications must publish full source under same terms
- No trademark grant for Nodeveil name/logo
