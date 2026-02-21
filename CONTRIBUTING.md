# Contributing to Nodeveil

## Development Setup

1. Install Go (see `engine/go.mod` toolchain requirement).
2. Install Node.js 20+ and pnpm 9+.
3. Install dependencies:
   - `cd engine && go mod download`
   - `pnpm -C desktop install`
4. Run checks:
   - `make lint`
   - `make typecheck`
   - `make test`

## Branch Strategy

- `main` is protected.
- Create short-lived feature branches from `main`.
- Rebase before opening a pull request.

## Commit Style

Conventional Commits are recommended:

- `feat: add graph root listing stub`
- `fix: handle sqlite migration error`
- `docs: update architecture decision`

## Pull Request Checklist

- [ ] Tests pass (`go test ./...`, `pnpm -C desktop test`)
- [ ] Lint/typecheck pass
- [ ] Documentation updated for behavior or architecture changes
- [ ] No binaries or build artifacts committed
- [ ] License and governance expectations respected

## Code Style

- Keep modules small and cohesive.
- Prefer explicit interfaces over shared mutable globals.
- Go: wrap errors with context (`fmt.Errorf("...: %w", err)`).
- TypeScript: strict mode and typed IPC contracts only.
