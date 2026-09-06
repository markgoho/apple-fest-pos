---
status: accepted
---

# Toppings are chosen in a dialog on the way into the cart, and an extra sour cream packet is a menu item

A leader set the topping rule at the booth: a potato pancake carries any number of the three condiments, from none to all three, and none of them adds to the price. The sour cream a pancake comes with is one packet. A customer who wants more pays one dollar for each further packet.

The old model could not say this. `CONTEXT.md` defined a Side as **the** condiment of a pancake, one per line, and [ADR-0005](./0005-pos-two-tap-place.md) built the tile grid on it: the pancake drew four tiles, Plain plus one per Side, and the tap that chose the item also chose the condiment. One tap cannot choose a set.

The pancake now draws **one tile**, and tapping it opens a **dialog** that asks for the toppings before the pancake reaches the cart. The dialog offers **Plain** and the three toppings. Plain and the toppings are mutually exclusive: choosing Plain clears the toppings, and choosing a topping clears Plain. **Add to cart** stays inert until one of them is chosen, because a pancake nobody was asked about and a plain pancake look the same in the cart and only one of them is an order.

The choice is a **dialog and not a row on the cart line**. The first build put the toggles on the line and nowhere else, which meant the only place to choose a topping was a list that grows all day: by the tenth order the Operator is hunting down the cart for the line that just appeared, while a customer waits. Choosing the toppings on the way in puts the choice where the attention already is, at the tile that was just tapped.

The **cart line keeps the toggles too**, so a topping asked for late is one tap on the line and not a delete and a redo. Editing a line re-keys it, and a line that becomes the twin of another folds into it.

This does not reopen [ADR-0005](./0005-pos-two-tap-place.md). That decision turned a dialog down for the **Place order checkpoint**, because a dialog there would cover the cart lines at the exact moment they are being read back. This dialog opens at add time and is gone before the read-back; the lines stay live and editable through the checkpoint exactly as ADR-0005 requires, and a topping toggle on a line is a cart change, so it disarms Place order like a quantity tap.

The toggles are **gold, resting and chosen alike** ([ADR-0004](./0004-pos-colour-meaning.md) gives gold to every reversible adjustment and reserves green for Place order alone). A chosen topping is a **solid fill**: the dark gold ink behind cream text, 10.04:1, well over ADR-0004's 7:1 daylight floor. An earlier build marked it with a checkmark in the corner of a pale button; the mark was small, it read as decoration rather than state, and it was one more thing to hunt for. A filled block is legible from across the booth and needs no glyph. The gold-tan fill was measured first and rejected: `--gold-tan` behind `--gold-tan-ink` is 4.78:1, short of the floor.

Because the toppings are known before the line exists, a line is keyed by its item and its topping set, so **six taps on the pancake tile with the same toppings make one line reading `6`**, and the kitchen ticket says `6 POTATO PANCAKE` rather than printing the same thing six times.

**Extra Sour Cream draws a minor tile**: as tall as the pancake beside it and a fifth as wide, but muted — parchment rather than cream, small text, no shadow. The demotion is done with colour and type, not with size. It is an add-on to a pancake, not something anyone comes to the booth to buy, and a tile the size of the pancake's said otherwise. It keeps a full tap target, because a small tile must not be a small tap.

**Extra Sour Cream is a menu item at one dollar**, not a packet count on the pancake line. The price of an order stays the sum of its menu items, so `insertOrder`, the per-item breakdown and the hourly revenue chart need no rule of their own, and the item joins the chart legend on its own. It also puts the "one packet is included" judgement where it belongs: the Operator applies it by tapping the tile a second time, and the software never has to know whether the free packet is counted per pancake or per order.

## Considered options

- **A packet count on the cart line, priced as `max(0, packets - 1) * 100`.** It models the rule literally. Rejected: it puts a price rule in three places (`order_store.go`, the per-item breakdown, the hourly chart), it needs a new chart series, and it forces the software to answer a question the booth has not answered — whether the free packet is per pancake or per order.
- **Toggles on the cart line and nowhere else.** No dialog, no mode, and the fastest possible tile tap. Built first and rejected in use: it makes the toppings reachable only through a list that grows all day, so the choice gets further from the tap that caused it with every order. It also forces every tile tap to open a new line, since the toppings are not yet known, which puts the same pancake on the ticket several times over.
- **Keep one tile per combination.** No new interaction at all. Rejected: three condiments make eight combinations, so the pancake would need eight tiles on a grid built for four.
- **A separate `sour-cream-packet` Side with its own price.** It keeps everything on one line. Rejected: a Side that costs money breaks the sentence "a Side is a condiment, free with the pancake", and it reopens the same three pricing sites.

## Consequences

Two statements elsewhere are superseded. `CONTEXT.md`'s Side entry no longer reads "the condiment" or "chosen when the pancake goes into the order", and ADR-0005's aside that "a scout who taps Sour Cream for Applesauce has no other moment to notice" now describes a toggle on the line rather than a tile in the grid. The checkpoint it defends is unchanged and covers the toppings too, because a toggle disarms Place order.

The `/pos` grid loses three tiles and gains one, so the Potato Pancakes section holds two tiles rather than four. The API cart line carries `sides` as an array; the old `side` field is still read, because a reprint of an order placed before this decision is built from its stored request.

Six pancakes with six different topping sets are six lines, which is the honest count: the kitchen has to be told six different things. Six with the same toppings are one line, whether the Operator taps the tile six times or uses the quantity stepper.

Every pancake now costs one extra tap, on top of the one ADR-0005 already added: tile, topping, Add to cart. That is the price of asking a question that has four right answers.

Landscape splits the screen three quarters menu to one quarter cart. A topping row on a cart line needs the width, and the cart is what the Operator reads back. In portrait the cart runs the full width of the screen, so a line reads across — name, toppings, quantity — rather than down. Stacked, three topping toggles make a line tall enough that the strip shows a single order line, and the Operator has to scroll to read back an order they should see at a glance. Portrait also gives the menu only the height its tiles need and the cart everything left over, rather than the other way round, and a toastie tile is about half the height it takes in landscape. Every row the menu gives up is a line of the order the Operator can read without scrolling, and the menu is five tiles that do not need the room.

The kitchen ticket prints **one topping per line**, not a joined list. It prints at double width, so an 80mm roll holds about 24 characters and nothing in the builder wraps them; a joined list of two toppings already overflows, and a break mid-word is the one thing the cook must not read. The customer receipt prints at normal width and keeps the joined form.

The booth confirmed the rule this ADR left open: the included packet is **per pancake**, not pooled across an order, so the pricing here stands. What stayed missing was legible only to whoever placed the order — a customer who watched a second pancake go in plain, then saw a dollar added for "Extra Sour Cream," had no way to see that the charge was a real second packet and not a double charge for the first. The cart line for Extra Sour Cream now carries a **Note**, `MenuItem.Note`, drawn under the name the same way a **Plain** tag is: "Each pancake includes 1 free packet." `Note` is a plain string on the menu item, not a pricing input, so it changes nothing about `insertOrder`, the per-item breakdown, or the hourly chart.
