#!/usr/bin/env bash
set -euo pipefail
cd "$(dirname "$0")/.."
mkdir -p build/engine
cd engine
go build -trimpath -ldflags "-s -w" -o ../build/engine/nodeveil-engine ./cmd/nodeveil-engine
