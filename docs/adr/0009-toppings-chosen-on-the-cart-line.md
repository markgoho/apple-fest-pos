---
status: accepted
---

# Toppings are chosen on the cart line, and an extra sour cream packet is a menu item

A leader set the topping rule at the booth: a potato pancake carries any number of the three condiments, from none to all three, and none of them adds to the price. The sour cream a pancake comes with is one packet. A customer who wants more pays one dollar for each further packet.

The old model could not say this. `CONTEXT.md` defined a Side as **the** condiment of a pancake, one per line, and [ADR-0005](./0005-pos-two-tap-place.md) built the tile grid on it: the pancake drew four tiles, Plain plus one per Side, and the tap that chose the item also chose the condiment. One tap cannot choose a set.

The pancake now draws **one tile**, and the Sides move onto the **cart line** as three toggles. The Operator taps the pancake, then taps the toppings the customer asked for, in any order and as many as they want.

The toggles sit **inline on the line, not in a step or a dialog**. ADR-0005 already settled this shape for the checkpoint that follows: the cart lines stay live, there is no mode to leave, and touching a line is both the correction and the exit. A topping toggle is a change to the cart, so it disarms Place order exactly like a quantity tap, and the Operator reads back an order that includes the toppings.

The toggles are **gold, resting and chosen alike**, and a chosen topping is marked by a checkmark, a heavier weight and a darker edge. [ADR-0004](./0004-pos-colour-meaning.md) gives gold to every reversible adjustment and reserves green for Place order alone; a chosen topping is reversible, so it never earns a second colour.

A tap on an item that has Sides **opens a new cart line** instead of adding to the one above it. The toppings are chosen after the tap, so two taps can mean two different pancakes, and merging them would force the Operator to undo a choice they have not made yet. An item with no Sides has nothing to choose and merges as before.

**Extra Sour Cream is a menu item at one dollar**, not a packet count on the pancake line. The price of an order stays the sum of its menu items, so `insertOrder`, the per-item breakdown and the hourly revenue chart need no rule of their own, and the item joins the chart legend on its own. It also puts the "one packet is included" judgement where it belongs: the Operator applies it by tapping the tile a second time, and the software never has to know whether the free packet is counted per pancake or per order.

## Considered options

- **A packet count on the cart line, priced as `max(0, packets - 1) * 100`.** It models the rule literally. Rejected: it puts a price rule in three places (`order_store.go`, the per-item breakdown, the hourly chart), it needs a new chart series, and it forces the software to answer a question the booth has not answered — whether the free packet is per pancake or per order.
- **A step between the tile and the cart: choose the pancake, then a topping screen.** It is what the request describes literally, and it makes the topping choice unmissable. Rejected: it is the dialog ADR-0005 turned down, on the same grounds. It covers the cart at the moment the cart is being read, and it puts a mode in front of the fastest tap on the screen.
- **Keep one tile per combination.** No new interaction at all. Rejected: three condiments make eight combinations, so the pancake would need eight tiles on a grid built for four.
- **A separate `sour-cream-packet` Side with its own price.** It keeps everything on one line. Rejected: a Side that costs money breaks the sentence "a Side is a condiment, free with the pancake", and it reopens the same three pricing sites.

## Consequences

Two statements elsewhere are superseded. `CONTEXT.md`'s Side entry no longer reads "the condiment" or "chosen when the pancake goes into the order", and ADR-0005's aside that "a scout who taps Sour Cream for Applesauce has no other moment to notice" now describes a toggle on the line rather than a tile in the grid. The checkpoint it defends is unchanged and covers the toppings too, because a toggle disarms Place order.

The `/pos` grid loses three tiles and gains one, so the Potato Pancakes section holds two tiles rather than four. The API cart line carries `sides` as an array; the old `side` field is still read, because a reprint of an order placed before this decision is built from its stored request.

An order of six pancakes with the same toppings is still one line and one topping row. Six pancakes with six different topping sets is six lines, which is the honest count: the kitchen has to be told six different things.
