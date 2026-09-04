---
status: accepted
---

# Placing an order on /pos takes two taps, and green marks only the second

[ADR-0002](./0002-pos-reach-zones.md) made Submit the last tap of a sale and gave it the largest target in the corner the finger already rests near. It weighed reach and it did not weigh who is doing the reaching: `CONTEXT.md` says most Operators are scouts between 11 and 17 and that their screens assume a novice. Against that, one tap is the whole distance between a novice's slip and a Kitchen Ticket on the griddle, and Placed is irreversible on paper.

The rail's primary control now has three states:

- **Review order**, gold. The cart has items. This tap commits nothing.
- **Check the order…**, faded and inert, for about half a second.
- **Place order**, green. This tap assigns the order number and sends both documents.

Any change to the cart — a tile, a quantity — returns the control to **Review order**. An Operator therefore cannot Place an order that differs from the one they just read back.

The checkpoint exists to catch the **wrong item**, not the wrong total. The total is on the rail the whole time and the Operator reads it aloud anyway; a scout who taps Sour Cream for Applesauce has no other moment to notice.

The guard is **time, not distance**. Both taps stay in the Home zone, because moving the second one into the Read zone would put a per-order tap outside comfortable reach on a strap-held tablet, against ADR-0002's rule. Two taps in one place are defeated by a fumbled double tap, so the inert half second absorbs it.

There is **no dialog and no Back**. A dialog would cover the cart lines, which are the thing being checked. With the lines left live there is no mode to leave, so there is nothing for a Back control to do; touching a quantity is both the correction and the exit.

The **Guarded zone holds one More button** instead of Clear alone. Clear order and the Admin link live inside it, so nothing destructive is on screen during normal selling. Voiding a placed order is not on the Operator's screen at all: it cannot recall paper, and reconciling the till is the Leader's job. It waits behind the password in [#40](https://github.com/markgoho/apple-fest-pos/issues/40).

## Considered options

- **Keep one tap.** Fastest, and it matches ADR-0002 as written. Rejected: it leaves a novice Operator with no checkpoint between a mis-tap and a printed Kitchen Ticket, and the cost of the mistake is paid in food and in queue time, not in the half second saved.
- **A confirm dialog over the screen.** The familiar pattern, and it makes the commit boundary unmissable. Rejected: it hides the cart lines at the exact moment they need reading. A checkpoint that covers the thing being checked is theatre.
- **Put the second tap somewhere else on the screen.** Genuinely immune to a double tap rather than merely resistant. Rejected: it moves a control used on every single order into the Read zone, which ADR-0002 reserves for text the Operator never touches.
- **Press and hold to place.** One control, no second target, and no double tap can defeat it. Rejected: a hold is slow under queue pressure, and a novice does not discover a gesture that nothing on the screen describes.
- **Void on the Operator's screen.** Rejected: see above. It is a Leader action, and it belongs behind the password.

## Consequences

Two statements in ADR-0002 are superseded. "Submit is therefore the last tap of a sale" now reads Place order, and "[Guarded] holds Clear and nothing else" now reads More. The reach rule itself is unchanged: frequent bottom-right, destructive top-left, held edge empty.

[ADR-0004](./0004-pos-colour-meaning.md)'s green moves from Submit to Place order, and the arming state takes gold. The rule that green means "the order is going out" and appears in exactly one place is not bent by this: it is applied more strictly, because under one tap the green button was armed from the moment the cart had items, and now it is green only when the next tap actually sends the order. Gold is correct for Review order for the same reason it is correct for the quantity steppers — that tap is reversible.

Every sale costs one extra tap and about half a second. At a booth that is a real price, paid deliberately, on the judgement that a novice Operator's mis-tap costs more.

The Admin row in the More sheet promises a password that `/admin` does not have yet. Until #40 lands, that link hands the sales screen to anyone holding the tablet.
