# ADR-0004: Graph Edge Lifecycle and Node Move/Delete Policy

## Status
Accepted

## Decision
- Graph edges are keyed by `node_id` references (`from_id`, `to_id`), never by path.
- Node move/rename keeps the same `node_id` through `MoveNodePath` updates.
- Node delete performs soft-delete on node and cascades soft-delete for connected edges.
- Recreated files are treated as new nodes (new `node_id`) unless moved with explicit move event.

## Rationale
This preserves backlink correctness across renames while keeping delete semantics auditable and deterministic.
