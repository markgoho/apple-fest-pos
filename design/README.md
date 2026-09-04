# Design canvas

`pos.pen` is a [pen.dev](https://pen.dev) canvas holding every `/pos` layout explored for [#16](https://github.com/markgoho/apple-fest-pos/issues/16): the baselines imported from the running server, the tile-grid variants, a Crazy 8s row of nine layouts, four colour schemes, and the states that shipped as [#41](https://github.com/markgoho/apple-fest-pos/issues/41) and [#42](https://github.com/markgoho/apple-fest-pos/issues/42). It is a record of what was rejected and why, which the ADRs summarise but do not show.

It ships nothing. No build step reads it, and [ADR-0003](../docs/adr/0003-one-go-binary.md) still holds: the POS is one Go binary.

## Open it

```sh
open -a Pen design/pos.pen
```

## Put the live screen on the canvas

Every frame here came from the real page, not from drawing. The pen.dev CLI (`npx @pen.dev/cli`, or `./node_modules/.bin/pen` after `bun install`) drives it:

1. The CLI's headless mode has **no browser**. Use `pen interactive --app desktop`, with Pen.app running **and** `pos.pen` open in it. `--in` alone is not enough.
2. The browser pane is portrait and the kiosk `Start shift` gate hides the register, so do not point it at the server. Make a static snapshot instead: `curl` `/pos`, rewrite `/static/...` to relative paths, drop the `kiosk.js` script tag and the two kiosk buttons, and append CSS forcing `:root{--app-height:800px}`, `html,body{width:1280px;height:800px}` and the landscape `.body { grid-template-columns: 1fr 15.5rem }`. Load it with a `file://` URL.
3. Seed a cart by appending a `<script>` that clicks `.tile` buttons on load. Same trick reaches any JavaScript state: click `#more` for the sheet, click `#submit` for the armed button.
4. `execute` verbs are `Insert` / `Update` / `Copy` / `Replace` / `Move` / `Delete` / `Get` / `Export`. The `U(...)` shorthand in `pen --help` does not exist; `read_skill({path:"execute.md"})` has the real API.

## Note

The file is about 2 MB of pretty-printed JSON, and a save rewrites all of it. Expect the repo to grow if it is saved often.
