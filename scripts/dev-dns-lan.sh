#!/usr/bin/env bash
set -euo pipefail

# Points the booth name at the Pi's LAN address, so a tablet on the home
# network gets a real HTTPS connection with the real certificate, with no
# GL.iNet needed. DuckDNS does not validate the IP it is given, so a public
# name can point at a private address; any DNS client on the same LAN then
# connects to it directly over the LAN, not the internet.
#
# This is a development-only stand-in for the GL.iNet's local DNS override
# (issue #36). Corporate DNS security tools (this Mac's included) block
# exactly this pattern as "DNS rebinding", so it will not resolve from a
# machine behind one — that is expected, not a bug.
#
# Clear it with `off` before the event: the booth is air-gapped and the
# design has no public DNS record at all.
#
# Usage: scripts/dev-dns-lan.sh on|off

REPO_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
ENV_FILE="$REPO_DIR/.env"

if [[ ! -f "$ENV_FILE" ]]; then
  echo "No $ENV_FILE. Run scripts/booth-https.sh first to set up the booth name and token." >&2
  exit 1
fi
# shellcheck disable=SC1090
source "$ENV_FILE"

case "${1:-}" in
  on)
    curl -fsS "https://www.duckdns.org/update?domains=${BOOTH_SUBDOMAIN}&token=${DUCKDNS_TOKEN}&ip=${PI_ADDRESS}"
    echo
    echo "$BOOTH_HOST now points at $PI_ADDRESS. Works for devices on the same LAN as the Pi."
    ;;
  off)
    curl -fsS "https://www.duckdns.org/update?domains=${BOOTH_SUBDOMAIN}&token=${DUCKDNS_TOKEN}&clear=true"
    echo
    echo "$BOOTH_HOST cleared. Matches production: no public DNS record at all."
    ;;
  *)
    echo "Usage: $0 on|off" >&2
    exit 1
    ;;
esac
