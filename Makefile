.PHONY: fmt lint typecheck test ci engine-test desktop-test bench

fmt:
	gofmt -w $$(find engine -name '*.go')
	pnpm -C desktop format

lint:
	golangci-lint run ./engine/...
	pnpm -C desktop lint

typecheck:
	pnpm -C desktop typecheck

test: engine-test desktop-test

engine-test:
	cd engine && go test ./...

desktop-test:
	pnpm -C desktop test

bench:
	cd engine && go run ./cmd/nodeveil-bench

ci: fmt lint typecheck test
