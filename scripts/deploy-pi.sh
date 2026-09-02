#!/usr/bin/env bash
set -euo pipefail

# Build on the Mac, push the binary and unit to the Pi over SSH, install and
# restart. The Pi does not pull from git, does not install anything, and
# does not build. The binary that lands here is the binary that runs.

if [ "$#" -ne 1 ]; then
  echo "Usage: $0 <pi-host>"
  exit 1
fi

PI_HOST="$1"
STAGE_DIR="/tmp/pos-deploy"
REPO_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

cd "$REPO_DIR"

echo "==> Check that the Pi is reachable"
ssh -o ConnectTimeout=10 "$PI_HOST" true

"$REPO_DIR/scripts/build-pos.sh"

echo "==> Push the binary, the unit file, and the install script"
ssh "$PI_HOST" "mkdir -p '$STAGE_DIR'"
scp dist-pi/pos deploy/apple-fest-pos.service scripts/install-pi-service.sh \
  "$PI_HOST:$STAGE_DIR/"
ssh "$PI_HOST" "chmod +x '$STAGE_DIR/install-pi-service.sh'"

echo "==> Install and restart (needs sudo on the Pi)"
ssh -t "$PI_HOST" "'$STAGE_DIR/install-pi-service.sh'"
