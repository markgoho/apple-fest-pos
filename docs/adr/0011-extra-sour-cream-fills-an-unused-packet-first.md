---
status: reverted
---

# Tapping Extra Sour Cream fills an unused packet on a pancake in the cart before it charges for a new one

**Reverted.** The Operator found the auto-merge broke in ordinary use once a converted line was edited afterward: the search in this ADR happily merges different customers' pancakes onto one shared `quantity: N` cart line, but the per-line Side toggle and the quantity stepper both act on that whole line at once, not one unit of it. Toggle Sour Cream off a merged line and every pancake on it loses the Side together; step its quantity down and the wrong customer's order can vanish. A cart edited after the fact stopped reliably matching what anyone had ordered, which is worse than the dollar this ADR set out to save. `pos.js` is back to always charging for Extra Sour Cream, the behaviour [ADR-0009](./0009-toppings-chosen-on-the-cart-line.md) originally chose; its line about the software not needing to know per-pancake versus per-order stands again. The Note on the cart line (ADR-0009's Consequences) still explains the charge. This document stays for the reasoning already spent and the option it rules out.

The Operator hit the free-rider problem [ADR-0009](./0009-toppings-chosen-on-the-cart-line.md) declined to solve. An order held one plain pancake and one Sour Cream pancake; the Sour Cream customer asked for a second packet, and the Operator's only tool was the Extra Sour Cream tile, which always charges a dollar. The plain pancake's own included packet sat unused in the same order the whole time.

ADR-0009 rejected answering this in software: "the software never has to know whether the free packet is counted per pancake or per order," and it explicitly rejected a packet-count formula because that formula would have to live in three pricing sites (`insertOrder`, the per-item breakdown, the hourly chart). The booth has since answered the question the formula begged: the packet is **per pancake**, one each, never pooled into a shared count across the order. That answer does not require a pricing formula. It requires the tile to look at the cart before deciding what to add to it.

Tapping **Extra Sour Cream** now searches the cart for a pancake line that has not chosen the Sour Cream Side — plain, or carrying only Applesauce or Ketchup — and gives that line the Side for free, the same edit the toggle on the line already makes. Only when every pancake in the cart already carries Sour Cream does the tap fall back to the old behaviour: a new one-dollar line. A line of more than one pancake splits: one unit peels off carrying the new Side, the rest stay as they were, so tapping the tile once never hands out more than one packet.

This is a client-side search over `CartLine.sides` and `CartLine.chosen`, not a pricing rule. `insertOrder`, the per-item breakdown, and the hourly chart are unchanged: they still price whatever mix of menu items and Sides the cart holds when Place order is tapped, the same sum they always computed. What changed is which lines the cart holds by the time it gets there.

## Considered options

- **Leave it charging every time**, and rely on the Note ADR-0009's Consequences added ("Each pancake includes 1 free packet") to let the Operator explain the charge. Rejected: a Note explains a charge; it does not stop the booth from taking a customer's dollar for something already sitting free in the same order. The Operator reported this in use as the actual problem, not a communication gap.
- **A packet-count formula**, ADR-0009's first rejected option, now with the pooling question answered. Still rejected for the reason ADR-0009 gave: it puts a price rule in three places for a fact that is fully decided by what is already in the cart before Place order is ever tapped.
- **Ask the Operator which pancake the extra packet is for**, with a picker on the tap. Rejected: the answer is never ambiguous when only one pancake in the cart lacks the Side, which is the common case this fixes; a picker adds a tap and a dialog to solve a case the search already resolves silently, and it still needs the search underneath to know when to skip itself.

## Consequences

ADR-0009's line "the software never has to know whether the free packet is counted per pancake or per order" is superseded: the software now knows, because the booth answered it, and the Extra Sour Cream tap acts on the answer. The rest of ADR-0009 stands — Sides are still free, chosen on a dialog and a per-line toggle, and Extra Sour Cream is still a menu item at one dollar for the case this ADR's search cannot satisfy.

The Note on the Extra Sour Cream cart line (ADR-0009's Consequences) now reads correctly less often, because the tap resolves for free whenever it can. It still matters for the case that does charge: a cart with no pancake at all, or one where every pancake already carries Sour Cream.

A pancake's Side toggle and the Extra Sour Cream tile can now produce the same edit to the same line. An Operator who taps the toggle instead of the tile, or the tile instead of the toggle, reaches the same cart either way, so there is no wrong tile to have tapped for a pancake that still has an unused packet.
