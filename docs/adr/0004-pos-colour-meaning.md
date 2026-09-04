---
status: accepted
---

# One meaning per colour on /pos, chosen by risk

`/pos` used three brand colours (red, gold, green) for cart and menu controls, but each colour carried two or more meanings: red was both the quantity minus and the full Remove; gold was both the destructive Clear and the neutral side choice; green was both the quantity plus, the final Submit, and "this side is chosen." A colour that means two things teaches a novice operator nothing.

The three colours now track **risk**, not arithmetic direction or tap frequency:

- **Red** means *destructive*, and nothing else. It covers only Remove (delete a line) and Clear (empty the cart). Both already sit apart from the rest of the screen: Remove lives in the Home zone next to its line, Clear lives alone in the Guarded zone (see [ADR-0002](./0002-pos-reach-zones.md)). Colour now reinforces that spatial signal instead of diluting it.
- **Gold** (the neutral tan already in code, to become a named token) means *reversible adjustment*. It covers quantity plus, quantity minus, and side choice, both resting and chosen. A chosen side is shown by a non-colour cue (outline or checkmark; the exact glyph is [#22](https://github.com/markgoho/apple-fest-pos/issues/22)'s prototype work), not by switching to green, so gold stays one colour with one meaning regardless of state.
- **Green** means *the order is going out*, and is used in exactly one place: Submit. Nothing else on the screen is ever green. ([ADR-0005](./0005-pos-two-tap-place.md) moves that one place to Place order and gives the Review order tap gold, which applies this rule more strictly rather than bending it.)

The menu-tile price badge, currently red, moves off red to the neutral gold-tan: it sits on the most-tapped, purely additive control on the screen (tapping a tile adds an item), and a red badge there would teach the wrong association at the exact spot operators tap most.

The daylight contrast floor (7:1, decided in [Find POS touch UI guidance](https://github.com/markgoho/apple-fest-pos/issues/17)) rules out the mid-tone brand fills as button backgrounds: white text on `--apple-red` measures 6.47:1 and white text on `--leaf-green` measures 6.43:1, both short of the floor. Red and green controls move to the existing `--apple-red-dark` (10.02:1) and `--leaf-green-dark` (9.38:1) tokens; no new colours are needed. The gold-tan already in code (`#f9db95` with `--apple-red-dark` text) measures 7.44:1 and already clears the floor.

## Considered options

- **Direction: shrink / neutral / grow.** Red for anything that shrinks the order (minus, Remove, Clear), gold for the neutral side choice, green for anything that grows or confirms it (plus, Submit, and a chosen side stays green). Rejected: it keeps Remove and Clear next to the harmless minus button under the same red, so red no longer flags danger specifically, and it keeps green's third meaning ("chosen") instead of resolving it.
- **Frequency: rare / common / final.** Red reserved for Clear alone as the rarest action, gold for every common action including Remove, green for Submit alone. Rejected: it groups Remove, a destructive action, with the harmless plus/minus/side-choice taps under gold, which hides the one property (destructiveness) that most needs a distinct colour for a novice operator.

## Consequences

Red now appears in two zones (Guarded and Home) but always means the same thing, so a rotating strap and fixed zones are no longer the only defence against an accidental Clear; the colour itself now says "this removes something" everywhere it appears. Green becomes rare by construction: it is correct for an operator to see it exactly once per order, right before the order leaves. Any future control that is tempted to reuse green for anything other than Submit is a sign that control belongs in gold instead. The gold-tan value (`#f9db95`) should be promoted from an inline literal to a named custom property alongside `--apple-red`, `--leaf-green`, and their `-dark` variants, since it now carries a real meaning across every reversible control on the screen.
