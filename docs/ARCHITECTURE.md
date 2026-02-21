# Nodeveil Architecture (Phase 1)

## System Overview

```mermaid
flowchart LR
  UI[Desktop Renderer] --> MAIN[Electron Main]
  MAIN --> IPC[Engine HTTP IPC (proto-aligned)]
  IPC --> SVC[IndexService]
  SVC --> IDX[Scanner + Poll Watcher + Reconciler]
  SVC --> STO[Store Repositories]
  STO --> DB[(SQLite WAL)]
```

## Package Responsibilities

- `internal/domain`: `Root`, `Node`, `IndexEvent`, `IndexStatus`.
- `internal/ignore`: ignore matcher (glob + prefix).
- `internal/indexer`: recursive scanner + watcher abstractions.
- `internal/store`: schema/migrations/repos and all SQL.
- `internal/service`: orchestration, queues, event application, status counters.
- `internal/ipc`: API handlers for roots, status, search, version, health.

## Data ownership

Engine is authoritative for roots and node metadata. Desktop only renders engine responses.

## Contract operations

Defined in `proto/nodeveil/v1/nodeveil.proto`:

- Roots: `AddRoot`, `RemoveRoot`, `ListRoots`
- Indexing: `GetIndexStatus`
- Querying: `SearchFiles`, `GetFileById`, `ListRecentChanges`
- Health: `GetVersion`, `HealthCheck`

## Not in Phase 1

- File content hashing by default
- OCR/PDF parsing
- Cloud sync / AI augmentation
