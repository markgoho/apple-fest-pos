# ITPP047 printer capabilities

## Source and confidence tiers

No Munbyn ITPP047 manual PDF exists in the repo root (checked with `ls -la *.pdf`; only the ITPP068 manual is there). This file uses three web sources instead. Each claim below is tagged with one of these tiers.

- **Spec sheet** — the official Munbyn product page, `https://pos.munbyn.com/munbyn-itpp047-series-thermal-receipt-printer/`. Fetched directly. Treated as first-party and reliable.
- **Help Center** — the official Munbyn support site, `support.munbyn.com`. A direct fetch of this domain returns HTTP 403 (bot-blocked). Content came through a text-extraction proxy (`r.jina.ai`) instead. This is still first-party Munbyn text, but the extraction is not a verbatim manual scan — one table came back visibly garbled (noted below) and had to be dropped.
- **Hardware Manual FAQ** — a 30-page "ITPP047 Thermal Printer Hardware Manual Version 1," read through a mirrored FAQ page (manualshelf.com) that exposes the manual's own text layer under Q&A headings. Third-party mirror, but the quoted lines are the manual's original text, with page numbers.
- **Derived** — reasoning from a confirmed fact, not a direct quote. Marked as such every time.

Nothing here is marked "confirmed by manual" the way the ITPP068 file does, because no manual PDF was available to read directly.

Sources used:
- Spec sheet: `https://pos.munbyn.com/munbyn-itpp047-series-thermal-receipt-printer/`
- Help Center, HEX command article: `https://support.munbyn.com/hc/en-us/articles/11503298932243-ITPP047-esc-pos-HEX-command`
- Help Center, Linux/Raspberry Pi setup: `https://support.munbyn.com/hc/en-us/articles/4602141661843-Set-Up-The-ITPP047-for-First-Use-on-Linux-Raspberry-Pi`
- Help Center, Epson mode: `https://support.munbyn.com/hc/en-us/articles/30531813601683-ITPP047-turns-on-Epson-mode`
- Hardware Manual FAQ mirror: `https://www.manualshelf.com/manual/munbyn/itpp047/frequently-asked-questions-english.html`

## Summary

The ITPP047 is an 80 mm thermal printer with a 72 mm printable width, same nominal geometry as the ITPP068. Two things make it different from the ITPP068 file's findings: the printable width is a **printer setting**, not a fixed property of the paper (a HEX command changes it), and Munbyn's own spec sheet confirms barcode and QR printing, buffer size, and print speed — none of which the ITPP068 manual stated. No source for either printer documents the real-time status queries (`DLE EOT n`, `GS r n`) that issue #20 and #34 care about.

## 1. Paper width and characters per line

Spec sheet, confirmed: printable width is 72 mm, off 80 mm stock. The paper width is adjustable in hardware between 48 mm and 80 mm (48/52/56/64/68/76/80 mm), and print resolution is stated as "576dots/line or 512dots/line."

This is the one finding this printer has that the ITPP068 file does not: the dot count per line is not fixed. The Help Center HEX command article documents a command to set it (`1F 1B 1F E1 13 14 n`, "Set the print content width"), with `n` selecting a width from 48 mm to 80 mm. The full table of `n` values came back garbled through the text-extraction proxy (duplicate and skipped rows), so this file does not reproduce it. The actionable fact survives: the width is a stored printer setting, readable off the self-test page, not a constant. Issue #34 must print a self-test page and read the active width before trusting any character count.

Two dot counts are confirmed by the spec sheet: 576 dots/line and 512 dots/line. Characters per line at each, derived from the standard ESC/POS Font A cell width (12 dots) and Font B (9 dots) — not stated by any Munbyn source, same assumption the ITPP068 file used:

| Dots/line | Font A (12 dots/char) | Font A double-width | Font B (9 dots/char) |
|---|---|---|---|
| 576 (72 mm, likely default) | 48 | 24 | 64 |
| 512 | 42 | 21 | 56 |

If the ITPP047 is at its 72 mm / 576-dot setting, its character counts match the ITPP068's derived counts exactly. If it has been set narrower, they do not. This is worth one line item on the self-test-page check in #34.

## 2. Font, bold, underline, alignment

Not documented by any source found. The Help Center HEX command article was fetched and explicitly does not list commands for font size, bold, underline, or alignment — it covers initialisation, paper cut, cash drawer, logo, and WiFi/network settings only. The spec sheet states the command set as "ESC/POS emulation supported" but gives no byte-level list.

Derived, not confirmed: since the spec sheet advertises ESC/POS emulation, the standard Epson bytes the current code already trusts for the ITPP068 (`ESC @` init, `GS ! n` character size, `ESC E` bold, `ESC -` underline, `ESC a` alignment) are likely to work here too. This is not proven by any Munbyn document for this model and must be checked on real hardware, the same way `escpos.go`'s existing bytes were.

## 3. Barcode or QR code printing

Spec sheet, confirmed: the printer prints one-dimensional barcodes — UPC-A, UPC-E, JAN13/EAN13, JAN8/EAN8, CODE39, ITF, CODABAR, CODE93, CODE128 — and two-dimensional QR codes. This is a direct, named capability claim from Munbyn's own product page, stronger evidence than the ITPP068 file had (which only inferred barcode capability indirectly, from a self-test page that happened to include a QR code).

Not confirmed: the exact command bytes. No source gives the HEX sequences for barcode or QR printing. Derived: since the spec sheet claims ESC/POS emulation, the standard commands are `GS k` (1D bar codes) and `GS ( k` (2D symbols, including QR) — the usual Epson-compatible bytes — but this is an assumption, not a citation.

## 4. Cut behaviour

Three partial-cut byte sequences turned up across sources. All three are partial cuts. No source for the ITPP047 documents a full cut.

- Help Center HEX command article: `ESC d 5` (`1B 64 05`, feed 5 lines) followed by `GS V 1` (`1D 56 01`, partial cut). This is Munbyn's own documented recipe for this specific model.
- Hardware Manual FAQ, p.10 of 30: `1D,56,42,00` — `GS V 66 0`, function-B partial cut with zero extra feed. Quoted directly: "cutter commands in the software--1D,56,42,00 Such as Loyverse."
- The current code (`internal/pos/escpos.go`): `GS V 0x42 0x08` — function-B partial cut with an 8-dot feed. Same command family as both sources above, different feed amount.

The Hardware Manual FAQ also confirms, for the Mac driver, "in the Printer features, choose Cut options, select Partial cut and print" (p.9) — the driver-level option is explicitly a partial cut, not full.

Practical read: the code's existing cut command is in the right family and should work unchanged. Munbyn's own recipe feeds 5 lines before cutting; `encodeLines` in the current code already appends three blank lines before the cut bytes run. Issue #34's real-device test should confirm the cut clears the last printed line with the feed the code already sends.

## 5. Print speed and buffer limits

Spec sheet, confirmed: print speed is 230 mm/s. Input buffer is 256 KB.

Both figures are more specific than anything the ITPP068 file found (that manual gave "300mm/s" from a self-test printout and said nothing about buffer size). Neither source says how the 256 KB buffer behaves under a burst of tickets, or whether tickets queue, block, or drop when the buffer fills. If order bursts become a concern, the same recommendation as the ITPP068 file applies: a timed real-device test, sending several tickets back to back and watching for dropped or corrupted output.

## 6. Real-time status queries (DLE EOT n, GS r n)

Not documented by any source found. The Help Center HEX command article was fetched directly for this question and does not mention `DLE EOT`, `GS r`, or any real-time status query. It lists only initialisation, cut, cash drawer, logo, and network-configuration commands.

One relevant, actionable fact turned up that is not about the status query itself but about a precondition for it: the Help Center has a separate article, "ITPP047 turns on Epson mode," describing a printer setting called "USB EPSON Mode," checkable on the printer's self-test page ("Yes" means Epson mode is on). The article does not explain what Epson mode changes about command handling, and does not say whether it affects status queries specifically. But since `DLE EOT n` and `GS r n` are Epson-standard commands, and this printer has an explicit Epson-compatibility toggle, issue #34 should print a self-test page and confirm "USB EPSON Mode: Yes" before testing status queries — a "no" answer would be a simpler explanation for a non-response than the printer lacking the feature outright.

The Help Center's Linux/Raspberry Pi setup article was also checked directly for raw-socket evidence, the same way the ITPP068 manual gave a `socket://hostname:9100` URI. It gives none: the ITPP047 Linux guide covers USB connection and the CUPS web admin page (`http://localhost:631/admin`) only, no network/socket/port-9100 detail. This does not mean the printer lacks raw-socket TCP printing — Ethernet is listed as an interface on the spec sheet, and the current code already targets it that way — it means the port and protocol are not confirmed by any source found and should be checked in #34 alongside everything else.

Same code-level blocker as the ITPP068 finding: `sendEscPos` in `internal/pos/printer.go` writes the payload and closes the socket without reading a reply (`defer connection.Close()`, no `Read` call). Even if the printer answers a status query on the same TCP connection — which is how Epson-compatible printers generally behave — the current code cannot see it. This file does not change `printer.go`.

## Proposal: item name wrap and Side placement

Menu data (`internal/pos/menu.go`), counted directly: the longest item name is "The Harvest Toastie" at 19 characters. The longest Side names are "Sour cream" and "Applesauce," both 10 characters ("Ketchup" is 7).

This proposal uses the 576-dot / default-width character counts from Q1 (48 normal, 24 double-width). Those counts are derived, not confirmed, and depend on the printer's width setting actually being at 72 mm — see the #34 self-test check in Q1. If the printer is ever found set to 512 dots (42 normal, 21 double-width), the kitchen line below no longer fits; that is called out below.

**Kitchen ticket.** Prints in double width, budget about 24 characters per line at the default setting. Current line format is `"%d  %s"` (quantity, two spaces, upper-case name). For the longest item: `"1  THE HARVEST TOASTIE"` is 1 + 2 + 19 = 22 characters — fits in 24, with 2 characters to spare. Note: at the narrower 512-dot / 21-character setting, this line would not fit (22 > 21); this is a second reason #34 must confirm the active width setting before this ships. Propose: keep the item name on its own line, same as the ITPP068 file's proposal. Print the Side directly below it, same double width, indented two spaces:

```
1  POTATO PANCAKE
  SOUR CREAM
```

**Customer receipt.** Prints at normal width, budget about 48 characters per line (42 at the narrower setting). Current line format is `"%d x %s"` followed by a price line. The longest line, `"1 x The Harvest Toastie"`, is 23 characters — well inside budget at either width setting, so wrapping is not a practical concern here regardless of which width the printer is set to. Propose the same layout as the ITPP068 file: add the Side as its own indented line between the item name and the price, at normal size, matching the exact casing menu.go already uses ("Sour cream," not "Sour Cream"):

```
1 x Potato Pancake
  Sour cream
  $10.00
```

Both placements reuse the existing one-line-per-element layout in `encodeLines`; this needs a new line in each build function's input, not a new ESC/POS command.

## Comparison to issue #33 (ITPP068)

Same in ways that matter: both are 80 mm paper with a 72 mm printable width. Both derive to 576 dots at that width (ITPP068: derived from an assumed 203 dpi print head; ITPP047: confirmed directly by the spec sheet). Both use the same cut command family — `GS V 66 n`, function-B partial cut — and neither printer's documentation gives a full-cut byte sequence. Neither printer's documentation answers `DLE EOT n` or `GS r n`; both are silent on real-time status queries entirely.

Different in ways that matter:

- **Width is configurable on the ITPP047, fixed on the ITPP068.** The ITPP047 has a documented command to change its printable width (48–80 mm); nothing in the ITPP068 manual suggests its width is adjustable. This means the ITPP047's character-per-line count is a setting that must be verified on the specific unit in the booth, not a printer-model constant.
- **Barcode/QR is a named, confirmed capability on the ITPP047; only inferred on the ITPP068.** The ITPP047 spec sheet names nine 1D symbologies plus QR. The ITPP068 finding only inferred QR capability from a self-test printout that happened to contain one.
- **Speed and buffer differ.** ITPP047: 230 mm/s, 256 KB input buffer (both confirmed by spec sheet). ITPP068: 300 mm/s (confirmed by self-test printout), buffer size not documented anywhere.
- **Network path is less confirmed for the ITPP047.** The ITPP068 manual's Linux install guide shows an explicit `socket://hostname:9100` raw-socket URI. No equivalent was found for the ITPP047; its Linux setup guide covers only USB and the CUPS web admin page. The current code (`internal/pos/printer.go`) already dials a raw TCP socket for printing, so this is not necessarily a blocker, but it is one more thing #34 must confirm on the real ITPP047 rather than read off a document.
- **The ITPP047 has an "Epson mode" toggle; nothing equivalent is documented for the ITPP068.** This could explain a printer that does not answer standard Epson status queries even though the command family is otherwise correct.

Verdict for #34 and #35: the two printers are close enough to share one ticket format and one status-query test procedure, provided the ITPP047 is confirmed to be at its 72 mm / 576-dot width setting and its "USB EPSON Mode" is on. Both of those are unverified assumptions, not facts, and both are checkable from the ITPP047's own self-test page in the same pass as the status-query test #34 already needs to run. If either check comes back different from assumed, the ITPP047's character-per-line budget (Q1) and possibly its response to status queries (Q6) will diverge from the ITPP068, and #35's ticket-building code would need a per-printer width, not a shared constant.
