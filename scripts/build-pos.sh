#!/usr/bin/env bash
set -euo pipefail

# Cross-compile the POS for the Pi. Static assets and templates are
# //go:embed'd into cmd/pos, so this one file is the whole deploy payload.

REPO_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
OUT="$REPO_DIR/dist-pi/pos"

cd "$REPO_DIR"
mkdir -p "$(dirname "$OUT")"

echo "==> go vet"
go vet ./...

echo "==> go test"
go test ./...

echo "==> build linux/arm64 (CGO disabled)"
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o "$OUT" ./cmd/pos

echo "==> built $OUT"
file "$OUT" 2>/dev/null || true
