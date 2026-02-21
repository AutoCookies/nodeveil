# Nodeveil

Nodeveil is a filesystem-native hard-link knowledge graph desktop application. Phase 0 establishes foundation quality gates, architecture boundaries, and governance.

## Positioning

- **Engine-first architecture**: Go engine owns graph/index metadata.
- **Desktop as client**: Electron renderer is a presentation layer.
- **Contract-driven IPC**: Protobuf contract defines `GetVersion`, `HealthCheck`, `ListRoots`.

## Repository Layout

- `engine/` — Go services, domain, db, IPC skeleton, benchmark CLI.
- `desktop/` — Electron + TypeScript app shell with strict lint/typecheck/tests.
- `proto/` — contract definitions.
- `docs/` — architecture docs and ADRs.

## Getting Started

### Prerequisites

- Go 1.23+
- Node.js 20+
- pnpm 9+

### Setup

```bash
cd engine && go mod download
pnpm -C desktop install
```

### Development workflow

```bash
# run quality checks
make lint
make typecheck
make test

# engine version output
cd engine && go run ./cmd/nodeveil-engine --version

# benchmark harness
cd engine && go run ./cmd/nodeveil-bench
```

## License summary

Nodeveil uses a source-available non-commercial copyleft license:

- Official binaries are for non-commercial use only.
- If you distribute modified versions, you must publish full corresponding source under the same license.
- No trademark rights are granted for the Nodeveil name/logo.

See `LICENSE` for full terms.
