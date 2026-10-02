---
status: accepted
---

# A worker's free meal is a Placed order the Leader comps afterwards

Each adult or scout who works the booth gets one free meal. The Operator places that order the normal way, and the Leader then marks it Comped on the Leader page, with a Comp button beside Void on each order. Comped is a third outcome for a Placed order, distinct from Voided: the order and its items still count as served, and only the money drops out.

The Leader page is the place for it because a comp is a statement about the till: no cash came in for this order. `CONTEXT.md` gives the cash and the till to the Leader and says the Operator does not handle cash. The path already exists for void (More, Leader, PIN, Orders), and the comp reuses all of it, including the confirm that shows the order number.

The money figures leave a Comped order out everywhere: the day's revenue, the event total, the per-item revenue column, and the revenue-by-hour chart. The count figures keep it: the order count and the per-item quantity. This keeps every revenue figure equal to the cash in the till, and every quantity equal to what the kitchen made. The stored order keeps its original total, so the order list can show the price that was waived.

## Considered options

- **A Comp button in the Place sequence on `/pos`, with a PIN dialog.** It would let the Customer Receipt print $0. Rejected: [ADR-0005](./0005-pos-two-tap-place.md) keeps dialogs and money decisions out of the Place sequence, and [ADR-0006](./0006-admin-pin-no-session.md)'s PIN gate is a form-first page, so this needs a second PIN pattern inside the cart script. The receipt showing the full price is accepted, the same way a void never touches the paper.
- **A $0 "Worker Meal" menu item.** Rejected: a scout can ring it up alone, and the per-item quantities lose what the kitchen made.
- **Void the worker's order.** Rejected: Voided means "as if never sold", so the food drops out of the item quantities although the kitchen made it.
- **Enforce one meal per worker.** Rejected: the booth has no list of workers, and the Leader who types the PIN is the control.

## Consequences

The `status` column now holds two kinds of fact: how the printing went, and what the Leader decided. A Voided or Comped order therefore no longer shows its print result in `status` (the two print-status columns still hold it), and a reprint must leave a Voided or Comped status alone. Before this ADR a reprint overwrote Voided with the print result and put the order back into the sales total; that is fixed with it.
