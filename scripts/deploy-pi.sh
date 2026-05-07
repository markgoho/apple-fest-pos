#!/usr/bin/env bash
set -euo pipefail

if [ "$#" -ne 1 ]; then
  echo "Usage: $0 <pi-host>"
  exit 1
fi

PI_HOST="$1"
APP_DIR="/home/markgoho/apple-fest-pos"

ssh "$PI_HOST" "cd $APP_DIR && git pull && bun install && sudo systemctl restart apple-fest-pos && systemctl --no-pager --full status apple-fest-pos"
