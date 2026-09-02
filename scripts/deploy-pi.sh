#!/usr/bin/env bash
set -euo pipefail

# Build at home over wifi, then push the build to the Pi over SSH.
# The Pi does not pull from git, install dependencies, or build. The build
# that is on the Pi when you leave the house is the build that runs.

if [ "$#" -ne 1 ]; then
  echo "Usage: $0 <pi-host>"
  exit 1
fi

PI_HOST="$1"
APP_DIR="/home/markgoho/apple-fest-pos"
REPO_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

cd "$REPO_DIR"

echo "==> Check that the Pi is reachable"
ssh -o ConnectTimeout=10 "$PI_HOST" true

echo "==> Compare the Bun versions"
LOCAL_BUN="$(bun --version)"
REMOTE_BUN="$(ssh "$PI_HOST" 'export PATH=$HOME/.bun/bin:$PATH && bun --version')"
if [ "$LOCAL_BUN" != "$REMOTE_BUN" ]; then
  echo "Bun version mismatch: local $LOCAL_BUN, Pi $REMOTE_BUN." >&2
  echo "The build is Bun bundler output. Run 'bun upgrade' on the Pi, or" >&2
  echo "install the same version here, then deploy again." >&2
  exit 1
fi
echo "Bun $LOCAL_BUN on both."

echo "==> Build"
bun run build

echo "==> Push the build, the unit file, and the install script"
# COPYFILE_DISABLE keeps macOS tar from adding ._* AppleDouble files.
COPYFILE_DISABLE=1 tar -czf - build deploy scripts/install-pi-service.sh \
  | ssh "$PI_HOST" "set -euo pipefail \
    && mkdir -p '$APP_DIR' \
    && rm -rf '$APP_DIR/build' \
    && tar -xzf - -C '$APP_DIR' \
    && chmod +x '$APP_DIR/scripts/install-pi-service.sh'"

echo "==> Install the unit file and restart the service"
ssh "$PI_HOST" "'$APP_DIR/scripts/install-pi-service.sh'"
