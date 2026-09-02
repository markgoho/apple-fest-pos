---
status: superseded by ADR-0002
---

# SvelteKit 3 with the official adapter-bun

Superseded by [ADR-0002](./0002-one-go-binary.md): the POS now ships as one Go binary.

The Apple Fest POS ran two processes in development: a `Bun.serve()` API in `app/server/` and a Vite dev server for the Svelte client. To collapse these into one process on one URL, SvelteKit now owns the server. The app moves to `@sveltejs/kit@next` (3.x) with the official `@sveltejs/adapter-bun`, which serves through a real `Bun.serve()`. Development is `bun --bun vite dev`, the build is `bun run --bun build`, and the server starts with `bun ./build`.

## Considered options

- **Bun-native Svelte** (`Bun.serve()` with HTML imports and `bun-plugin-svelte`, no Vite or SvelteKit). Rejected: research proved that `bun --bun vite dev` already runs SvelteKit, `bun:sqlite`, and `Bun.connect` in a single process, so hand-rolled routing and server rendering would buy nothing.
- **`@sveltejs/adapter-node` on SvelteKit 2**, started with `bun ./build/index.js`. This was verified end to end on the repo's exact versions and carries no prerelease risk, but it runs `node:http` on Bun rather than `Bun.serve()`. Rejected in favour of the Bun-native path.
- **`svelte-adapter-bun`** (community, gornostay25). A real `Bun.serve()` and embeddable, but it declares a `typescript ^5` peer against this repo's `^6`, and it patches SvelteKit's built server with a regular expression to add websocket hooks. Rejected as too fragile.

## Consequences

Both `@sveltejs/kit@3` and `@sveltejs/adapter-bun` are prereleases. Kit 3 is published on the `next` tag only, `adapter-bun` has one published version (`1.0.0-next.1`), and its documentation page returns 404. This risk was accepted deliberately, before an event with a fixed date. If the upgrade breaks, **the rule is to fix forward on Kit 3**; there is no fallback to `adapter-node`.

Two downstream effects follow from SvelteKit owning the server:

- `app/server/index.ts` and `app/server/http/cors.ts` are deleted. The adapter owns the entry point, and one origin serves everything, so no CORS layer is needed.
- The client API wrapper `app/web/src/lib/services/api.ts` is deleted. Page reads call the services in process from a `load` function. Only the polling screens (admin sales, kitchen display) and the order submit keep an HTTP endpoint.

Because a compiled single binary would be architecture-locked to the Raspberry Pi's arm64, `buildOptions.compile` stays off and the Pi deploy keeps copying a `build/` tree over SSH.
