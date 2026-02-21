# Nodeveil

Nodeveil is a filesystem-native hard-link knowledge graph desktop app.

## Getting started

```bash
cd engine && go mod download
pnpm -C desktop install
```

## Run engine + desktop

```bash
cd engine && go run ./cmd/nodeveil-engine
# another terminal
pnpm -C desktop build
```

## Graph export/import CLI

```bash
cd engine && go run ./cmd/nodeveil-engine --graph-export graph.json
cd engine && go run ./cmd/nodeveil-engine --graph-import graph.json
```

## Quality gates

```bash
pnpm -C desktop lint
pnpm -C desktop typecheck
pnpm -C desktop test
cd engine && go test ./...
cd engine && go test -race ./...
cd engine && go vet ./...
cd engine && golangci-lint run ./...
cd engine && go run ./cmd/nodeveil-bench
```

## Example usage

1. Add two roots in desktop.
2. Search and copy two node IDs.
3. Create link from A to B in Links panel.
4. View outgoing/incoming backlinks.

## License summary

- Non-commercial use for official binaries
- Distributed modifications must publish full source under same terms
- No trademark grant for Nodeveil name/logo
