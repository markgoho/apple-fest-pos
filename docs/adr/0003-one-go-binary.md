---
status: accepted
---

# One Go binary replaces SvelteKit and Bun

The Apple Fest POS now ships as one Go binary, cross-compiled on the Mac for the Raspberry Pi. It serves the whole POS with the Go standard library: `net/http` for routing, `html/template` for the screens, `net.Dial` for the ESC/POS printers, and `ListenAndServeTLS` for HTTPS. SQLite goes through `modernc.org/sqlite`, which is pure Go, so `CGO_ENABLED=0` gives a statically linked `linux/arm64` executable that needs no toolchain, no JavaScript runtime, and no `node_modules` on the Pi.

This supersedes [ADR-0001](./0001-sveltekit-3-adapter-bun.md), which put SvelteKit 3 and `@sveltejs/adapter-bun` in charge of the server.

## Why the earlier decision was reversed

ADR-0001 accepted prerelease risk deliberately, and the risk came due. The whole event depended on `@sveltejs/kit@3.0.0-next.25` and `@sveltejs/adapter-bun@1.0.0-next.1`, an adapter with one published version and a documentation page that returns 404.

Two facts then made the trade one-sided:

- **HTTPS is a requirement, not a nicety.** Issue [#6](https://github.com/markgoho/apple-fest-pos/issues/6) proved on a real tablet that plain HTTP refuses the screen wake lock, so the POS screen sleeps mid-shift. `adapter-bun` cannot pass a `tls` option through to `Bun.serve()` at all. Go terminates TLS in one line.
- **The framework was carrying almost nothing.** The server was 625 lines of TypeScript across seven files: a SQLite schema, an order service, an ESC/POS byte builder, a TCP write, and three JSON routes. None of it used framework leverage that Go's standard library does not also give.

## Consequences

- The deploy is one `scp` of one file, plus a `systemctl restart`. Issue [#28](https://github.com/markgoho/apple-fest-pos/issues/28) confirmed a `CGO_ENABLED=0 GOOS=linux GOARCH=arm64` build runs on the Pi (`aarch64`, Debian 13 trixie) and that `modernc.org/sqlite v1.57.0` opens a database and writes a row there.
- The cart stays client-side, but as hand-written JavaScript served as a static file. There is no bundler and no build step for the browser code.
- The `src/` SvelteKit tree stays in the repository until the three screens are rebuilt as server-rendered HTML. It is legacy, not the running server.
- The Pi stays the single point of failure by design. Nothing in the Go server hides a dead Pi behind a queue or a cache.
