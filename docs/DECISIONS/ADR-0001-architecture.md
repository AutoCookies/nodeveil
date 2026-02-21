# ADR-0001: Layered Engine + Desktop Client Architecture

## Status
Accepted

## Context
Nodeveil needs strict module boundaries and predictable growth from prototype to production.

## Decision
Adopt a layered architecture:

1. Go engine as authority for graph/index metadata.
2. Electron desktop as thin client.
3. Explicit contract boundary (`proto/nodeveil/v1/nodeveil.proto`).
4. Composition root in `engine/internal/app`, pure domain in `engine/internal/graph`.

## Consequences
- Clear ownership and testability.
- UI can evolve independently from domain logic.
- Slight upfront complexity in transport contract management.
