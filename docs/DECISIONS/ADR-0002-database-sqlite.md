# ADR-0002: SQLite Storage + WAL + SQL Migration Bootstrap

## Status
Accepted

## Context
Phase 0 needs a deterministic SQLite-backed persistence path before adopting a compiled Go SQLite driver.

## Decision
Use SQLite as the data format and run SQL migrations through the `sqlite3` command-line client from `internal/db`.

- Set WAL mode (`PRAGMA journal_mode=WAL`).
- Apply bootstrap migrations from code in deterministic order.
- Keep the repository boundary in `internal/db` so driver internals can be swapped later without domain changes.

## Consequences
- Maintains SQLite semantics in Phase 0 with no third-party Go module downloads.
- Requires `sqlite3` to be available in the runtime environment.
- Planned follow-up: switch to a pure-Go embedded driver for tighter portability.
