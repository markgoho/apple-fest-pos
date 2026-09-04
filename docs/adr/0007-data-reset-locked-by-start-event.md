---
status: accepted
---

# The data-reset tool is locked out by a manual Start Event flag, not a calendar check

The System Admin page's "wipe all orders" tool ([#43](https://github.com/markgoho/apple-fest-pos/issues/43), [#46](https://github.com/markgoho/apple-fest-pos/issues/46)) deletes every row in `transactions` and clears the order-number counter, with no date scoping and no undo. It is disabled, server-side, once the System Admin has set Start Event: a new flag in the `metadata` table, toggled once from the System Admin page.

The obvious alternative was to derive the lockout from the calendar: the event's two dates are already fixed and known, so the wipe endpoint could just refuse when today's business date falls on one of them. That was rejected because it only protects the literal event days. The System Admin wants to lock the tool out earlier, once evening setup is finished and testing is done, not wait for the calendar to turn over. A manual flag also matches how the tool is actually used: a scout-meeting practice run and pre-event dev testing already can't corrupt real sales, because sales are read per business date; the only real danger is testing on the event day itself, and only the System Admin, standing at the booth, knows when that risk starts.
