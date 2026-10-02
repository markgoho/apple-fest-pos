---
status: accepted
---

# The data-reset tool is locked out by a manual Start Event flag, not a calendar check

The System Admin page's "wipe all orders" tool ([#43](https://github.com/markgoho/apple-fest-pos/issues/43), [#46](https://github.com/markgoho/apple-fest-pos/issues/46)) deletes every row in `transactions` and clears the order-number counter, with no date scoping and no undo. It is disabled, server-side, while the System Admin has Start Event set: a new flag in the `metadata` table, set from the System Admin page.

The obvious alternative was to derive the lockout from the calendar: the event's two dates are already fixed and known, so the wipe endpoint could just refuse when today's business date falls on one of them. That was rejected because it only protects the literal event days. The System Admin wants to lock the tool out earlier, once evening setup is finished and testing is done, not wait for the calendar to turn over. A manual flag also matches how the tool is actually used: a scout-meeting practice run and pre-event dev testing already can't corrupt real sales, because sales are read per business date; the only real danger is testing on the event day itself, and only the System Admin, standing at the booth, knows when that risk starts.

## Amendment, 2 October 2026: the flag can be cleared

Start Event was first built as a one-way switch with no confirmation: one stray tap set it, and the only way back was to delete the `event_started` row from the SQLite file over SSH. That is not a repair to try at the booth, where the System Admin page is the only troubleshooting tool (ADR-0010). The System Admin page now has a Clear Start Event button, and both directions sit behind a `confirm()` dialog.

The lock still does its job. It exists to stop a wipe by mistake, not a wipe on purpose, and the System Admin is the only person behind the PIN. A wipe after Start Event now takes three deliberate actions: clear the flag, confirm, then wipe and confirm. Clearing the flag deletes no orders by itself.
