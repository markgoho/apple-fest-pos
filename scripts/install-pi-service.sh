#!/usr/bin/env bash
set -euo pipefail

# Runs on the Pi, as the deploy user, over `ssh -t` so sudo can prompt.
# Installs the binary and the systemd unit. Does not touch the certificate:
# that is scripts/booth-https.sh's job (stage "Copy the certificate to the
# Pi"). If it was skipped, this script says so and stops.

STAGE_DIR="/tmp/pos-deploy"
SERVICE_PATH="/etc/systemd/system/apple-fest-pos.service"

if ! id -u pos >/dev/null 2>&1; then
  echo "==> Create the pos system user"
  sudo useradd --system --no-create-home --shell /usr/sbin/nologin pos
fi

if [ ! -r /etc/pos/pos.crt ] || [ ! -r /etc/pos/pos.key ]; then
  echo "No certificate at /etc/pos/pos.crt and /etc/pos/pos.key." >&2
  echo "Run scripts/booth-https.sh on the Mac first (stage: Copy the" >&2
  echo "certificate to the Pi), then run this again." >&2
  exit 1
fi

echo "==> Install the binary"
sudo install -o root -g root -m 0755 "$STAGE_DIR/pos" /usr/local/bin/pos

echo "==> Install the unit file"
sudo cp "$STAGE_DIR/apple-fest-pos.service" "$SERVICE_PATH"
sudo systemctl daemon-reload
sudo systemctl enable apple-fest-pos
sudo systemctl restart apple-fest-pos
sudo systemctl --no-pager --full status apple-fest-pos
