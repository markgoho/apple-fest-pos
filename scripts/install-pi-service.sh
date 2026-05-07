#!/usr/bin/env bash
set -euo pipefail

APP_DIR="/home/markgoho/apple-fest-pos"
SERVICE_PATH="/etc/systemd/system/apple-fest-pos.service"

sudo cp "$APP_DIR/deploy/apple-fest-pos.service" "$SERVICE_PATH"
sudo systemctl daemon-reload
sudo systemctl enable apple-fest-pos
sudo systemctl restart apple-fest-pos
sudo systemctl --no-pager --full status apple-fest-pos
