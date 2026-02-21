# Nodeveil Architecture (Phase 3)

```mermaid
flowchart LR
  APP[Renderer App Shell] --> ROOTS[roots feature]
  APP --> SEARCH[search feature]
  APP --> NAV[navigator feature]
  APP --> VIEWER[viewer feature]
  APP --> LINKS[links feature]
  APP --> DIAG[diagnostics feature]
  ROOTS --> IPC[typed ipc client]
  SEARCH --> IPC
  VIEWER --> IPC
  LINKS --> IPC
  DIAG --> IPC
  IPC --> ENGINE[engine ipc server]
  ENGINE --> INDEX[IndexService]
  ENGINE --> GRAPH[GraphService]
  INDEX --> STORE
  GRAPH --> STORE
  STORE[(SQLite WAL)]
```

## Frontend boundaries

- `features/*` are isolated by concern (no feature-to-feature imports).
- each feature has explicit state transitions (`state.ts`, `actions.ts`, `selectors.ts`, `controller.ts` as needed).
- `shared/ipc` owns retries/timeouts/cancellation.

## Backend boundaries

- `internal/service/index_service.go` for indexing lifecycle.
- `internal/service/graph_service.go` for link semantics.
- `internal/store/*` contains all SQL.

## Stability goals

- resilient reconnect behavior when engine is down.
- paged/virtualized rendering in navigator/search paths.
- diagnostics via `events_log` read model.
