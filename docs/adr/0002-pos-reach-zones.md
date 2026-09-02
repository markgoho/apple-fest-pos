---
status: accepted
---

# Four reach zones on /pos, fixed in both orientations

The Operator holds the Pixel Tablet through a **rotating hand strap** on a case, left hand in the strap, and taps with the right index finger. The tablet is never docked and never set down during a shift. The `/pos` screen is therefore divided into four zones, and the zones do not move when the tablet turns.

- **Grip strip.** The band along the held edge, under the strap hand. It holds no control. The strap thumb rests on this glass and must never register a tap.
- **Home.** The bottom-right area, where the free index finger rests with no arm movement. It holds every frequent tap: menu tiles, quantity, and Submit. Submit sits in the far bottom-right corner.
- **Read.** The middle and upper screen. Text the Operator looks at and never touches.
- **Guarded.** The top-left, next to the grip strip and at the maximum distance from Submit. It holds Clear and nothing else.

The rule is: **frequent bottom-right, destructive top-left, held edge empty, everything else read-only.** A rotating strap keeps the grip strip on the left edge in both landscape and portrait, so one rule covers both.

## Considered options

- **A charging dock on the booth table.** This is the standing kiosk that the touch-ergonomics literature actually describes, so the 20 mm target figure would apply directly instead of being carried with a caveat. Rejected: the dock is out of stock and may not return. The strap is the posture that exists.
- **A fixed (non-rotating) hand strap.** Cheaper and available on more cases, but the held edge then lands on the left in landscape and on the top or bottom in portrait, so the layout must protect two different dead zones. Rejected: a rotating strap costs no more and removes the whole compromise. **The rotating strap is a stated hardware requirement of this design.**
- **Treating the strap thumb as a second tapping finger.** It rests on the glass already, so a shortcut there is tempting. Rejected: a resting thumb that can fire a control is how Clear gets hit by accident.
- **A customer-facing region of the screen.** Rejected: the POS handbooks surveyed in the research say to never turn the device towards the customer, and the Operator reads the total aloud instead. The whole screen serves the Operator.

## Consequences

The Operator does not handle cash. An adult stands beside them, takes the money, and never touches the tablet; the Operator then taps Submit. Submit is therefore the **last** tap of a sale, which is why it earns the largest target in the corner the finger already rests near. The `/pos` page title "Cashier POS" is wrong under this model and is corrected when the layout is built.

Because the tablet is always in the hands, the screen is never read at the steep angle of a table-top device. That removes one argument from the colour and contrast decision, which still stands on daylight glare alone.

The zones are stated as zones, not as a layout. Assigning the cart, the menu grid, the total, and the banners to them is the layout work, and it must satisfy this rule rather than restate it.
