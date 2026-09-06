---
status: accepted
---

# The System Admin page discovers printers on the network and assigns them at runtime

[Give both printers addresses on the booth network, and switch printing on](https://github.com/markgoho/apple-fest-pos/issues/39) plans static IPs and DHCP reservations on the GL.iNet, set once into `deploy/apple-fest-pos.service`'s env vars. That plan still stands and stays the stable path: known addresses, no scanning, no admin action needed on a normal event day.

It has one failure mode with no recovery: the booth network ([#52](https://github.com/markgoho/apple-fest-pos/issues/52)) is air-gapped, no internet, 3-4 October 2026. If a printer arrives with a factory address instead of its reserved one, a DHCP reservation doesn't take, or a spare unit swaps in mid-event, fixing it meant editing the systemd unit's env vars and restarting the service — over SSH, from a Mac, on a network Claude Code (or any tool that needs the internet) cannot reach. On event day, at the booth, that fix does not exist. The System Admin page is the only troubleshooting surface that works there at all (CONTEXT.md, System Admin), so it is where this has to live.

## Decision

`DiscoverPrinters` (`internal/pos/discovery.go`) scans the server's own local IPv4 subnets for hosts answering on port 9100 — the raw ESC/POS port every ITPP047(P) listens on — and verifies each hit with the same `DLE EOT` status query `checkPrinter` already uses (ADR-0008), so an unrelated open port doesn't get mistaken for a printer. The System Admin page's new "Discover printers on network" button runs it and lists what it finds; a form lets the admin type or click-fill a Window and a Kitchen host/port and save. "Save" persists the assignment to the `metadata` table (the same key/value table `business_date`/`next_order_number` already use) and makes it the live config immediately — no restart.

`PrinterConfig.Enabled` is no longer read from the page: it's derived as `WindowHost != "" || KitchenHost != ""` whenever the admin saves. Requiring a separate switch alongside the host fields was a second way to leave printing looking "not configured" after an admin had just set an address. Clearing both hosts is now the only way to turn printing off from the page. `PRINTER_ENABLED` in the deploy env only matters before the first save; the deployed unit currently pins it `false` regardless, which used to mean printing stayed off even after `WINDOW_PRINTER_HOST`/`KITCHEN_PRINTER_HOST` were set — this removes that trap.

On startup, `OrderService.LoadPrinterConfig` overlays any saved `metadata` row onto the env-derived config `main.go` built; a row that's never been written falls back to the env var untouched, so a fresh deploy behaves exactly as before. `OrderService.Printer` is read through a new mutex-guarded `printerConfig()`/`setPrinterConfig()` pair, because the page can now change it while a request handler is mid-print.

This is an addition to ADR-0008's scope, not a reversal: 0008 decided the System Admin page may read status bytes from the two configured printers. This decision is about how those two printers get configured in the first place, which 0008's "stays scoped to the two manual troubleshooting actions" note didn't cover. [Decide the scope of the System Admin printer tooling](https://github.com/markgoho/apple-fest-pos/issues/50) closed with that narrower scope; this reopens the question only for address assignment, not for the status-reading design 0008 settled.

## Alternative considered

Editing the systemd unit and restarting, i.e., what #39 already does, kept as the only path. Rejected as the sole mechanism: it depends on SSH and a Mac with the repo, neither of which exists at an air-gapped booth. The static-IP/DHCP-reservation plan stays as the default; this page is the fallback when that plan doesn't hold on the day.
