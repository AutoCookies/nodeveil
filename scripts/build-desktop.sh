#!/usr/bin/env bash
set -euo pipefail
cd "$(dirname "$0")/.."
pnpm -C desktop install --frozen-lockfile
pnpm -C desktop build
# packaging scaffold (electron-builder recommended)
echo "Run: pnpm -C desktop exec electron-builder --linux AppImage --win nsis --mac dmg"
