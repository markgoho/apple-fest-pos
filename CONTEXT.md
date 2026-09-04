# Context

The glossary for the Apple Fest POS. Terms only, no implementation.

## Operator

A person who runs the booth: takes orders on a tablet and watches for print failures. The Operator does not handle cash. The Operator uses `/pos` and `/kitchen`. Most Operators are scouts between 11 and 17; an adult takes orders only occasionally. The Operator holds the tablet through a rotating hand strap, left hand in the strap and right index finger on the glass, and turns it to either orientation. The tablet is never docked and is not set down during a shift. Screens for the Operator are dense, work in landscape and in portrait, and assume a novice.

## Leader

An adult scout leader who watches how the day is going but does not take orders. The Leader takes the cash from the customer and reconciles the till. Day to day the Leader only reads, on a personal phone joined to the booth access point; screens for the Leader are portrait, one-thumb, and never show print status. The one exception is voiding a Placed order: that happens on a PIN-gated Leader page reached from a link in the Operator's `/pos` screen, so the Leader may use the Operator's tablet for that one action. A void corrects the sales record only; it never touches the till or the printed paper, which are still settled by hand.

An Operator and a Leader can be the same human at different moments. The distinction is the job, not the person.

## System Admin

You, the booth's technical owner. The System Admin page holds data reset and printer network setup and diagnostics, sized for troubleshooting an air-gapped event with no internet. It is PIN-gated and is never linked from `/pos`: it is your own tool, not something discovered from the tablet.

## Event day

One selling day of Apple Fest. The 2026 event is a Saturday and a Sunday, roughly 8am to 6pm each. An event day is the unit that order numbers reset on and that sales screens select.

## Start Event

The one-way switch the System Admin sets from the System Admin page once pre-event setup and testing are finished. It exists to close off the data reset tool: the wipe that clears every order can run only before Start Event, never after, so a real Event day's sales can never be erased by mistake.

## Side

The condiment served with a potato pancake: sour cream, applesauce or ketchup. A Side comes from this fixed set, and the Operator chooses it when the pancake goes into the order, not afterwards. A pancake with no Side is a normal order, not an incomplete one. A Side is not a Note: a Note is free text, and a Side is one of a known set that can be counted.

## Customer Receipt

The paper the customer takes away. It carries the order number, which is how the customer collects the food and how the booth finds the order again. One order makes exactly one Customer Receipt, and it prints on the Window Printer.

## Kitchen Ticket

The paper the cook works from. It lists what to make, and it prints on the Kitchen Printer. A Kitchen Ticket that is printed a second time carries a REPRINT mark, so the kitchen checks the order number instead of cooking the order again.

## Window Printer

The thermal printer at the booth window. It prints the Customer Receipt only.

## Kitchen Printer

The thermal printer at the cooking end of the booth. It prints the Kitchen Ticket only.

## Placed

The state of an order the Operator has confirmed and the booth has committed to. Placing assigns the order number and puts both the Customer Receipt and the Kitchen Ticket on their way. An order that is not yet Placed is a cart: it exists only on the tablet, it costs nothing to change, and it is not a sale. Placing is irreversible on paper, because nothing recalls a Kitchen Ticket that has already printed. An order Placed by mistake is settled by the Leader against the till and by hand at the kitchen, never by the Operator's tablet.

## Voided

The state of a Placed order the Leader has reversed. A Voided order drops out of the day's sales total, order count, and per-item sales breakdown, as if it were never sold. It stays in the order list, marked Voided, so the screen agrees with the Kitchen Ticket paper still at the booth about which orders happened. The void step shows the order number, so the Leader matches the paper in hand to the order on screen before voiding it.

## Sent

The state of a document that left the tablet for its printer. Sent is not proof that paper came out: the printer answers nothing about paper, so an empty roll, an open cover and a jam all count as Sent. No state in this system proves a ticket exists. Only a human eye on the tray does.
