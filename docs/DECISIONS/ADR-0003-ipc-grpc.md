# ADR-0003: gRPC Contract-First IPC (HTTP skeleton in Phase 0)

## Status
Accepted

## Context
Engine/Desktop communication must use explicit contracts and remain evolvable.

## Decision
Define canonical IPC contract in protobuf (`proto/nodeveil/v1/nodeveil.proto`) with gRPC service semantics.

For Phase 0, implement a lightweight HTTP JSON skeleton in engine that mirrors the same RPC operations, allowing incremental transition to full gRPC transport without changing operation semantics.

## Consequences
- Contract-first governance starts immediately.
- Desktop can consume stable operation names and payload shapes.
- Full gRPC transport wiring deferred to a later phase.
