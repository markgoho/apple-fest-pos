# Single-process server shape: facts

Research for [issue #2](https://github.com/markgoho/apple-fest-pos/issues/2), part of [issue #1](https://github.com/markgoho/apple-fest-pos/issues/1). Facts only. This file does not choose a shape.

Date: 2026-09-01.

## The fork, answered first

**Bun can run SvelteKit development with `bun:sqlite`.** `bun --bun vite dev` starts the app in `app/web` at the repo versions, and a `+server.ts` route that imports `bun:sqlite` runs correctly. `Bun.connect` (used by the printer service) is also available. The ticket asked for a plain statement if this failed. It did not fail. Bun-native Svelte is therefore **not** the only single-process option.

One extra fact goes with it: the **build** must also run under Bun. Plain `vite build` on Node fails when server code imports `bun:sqlite`.

## Versions used

Repo (installed in `node_modules`, verified 2026-09-01):

| Package | Version |
| --- | --- |
| bun | 1.4.0 |
| `@sveltejs/kit` | 2.59.1 |
| `svelte` | 5.55.5 |
| `vite` | 8.0.10 |
| `@sveltejs/adapter-static` | current (repo uses it now) |
| node (control tests only) | v24.16.0 |

A separate scratch copy (`/private/tmp/.../scratchpad/kitnode`) was used for the adapter-node tests. It installed newer packages: `@sveltejs/kit` 2.70.3, `svelte` 5.57.0, `vite` 8.2.2, `@sveltejs/adapter-node` 5.5.7. Read the adapter-node results with that drift in mind.

## Q1. Does `bun --bun vite dev` run current SvelteKit, with `bun:sqlite` and the printer service?

**Yes.** Verified in `app/web` at repo versions.

Method: start `bun --bun ../../node_modules/.bin/vite dev --port 5199`, then add a temporary route `src/routes/__bunprobe/+server.ts` that imports `Database` from `bun:sqlite`, opens `:memory:`, and reports `process.versions.bun` and `typeof Bun.connect`. The route was deleted after the test.

Result:

```
{"bun":"1.4.0","connect":"function","sqlite":7}
```

- Vite starts: `VITE v8.0.10 ready in 1033 ms`.
- `GET /pos` returns 200 with server-rendered HTML.
- `process.versions.bun` is `1.4.0`, so SvelteKit server code runs on the Bun runtime, not Node.
- `bun:sqlite` imports, creates a table, and queries it inside a SvelteKit endpoint.
- `Bun.connect` is a function, so `app/server/services/printer-service.ts` (line 42, `Bun.connect`) can run in the same process.

**Control (Node):** the same route under `node_modules/.bin/vite dev` on Node v24.16.0 returns `{"message":"Internal Error"}` with an ESM loader failure. So the `--bun` flag is what makes this work, not Vite.

**Build:** in the scratch copy, `vite build` on Node fails:

```
Error [ERR_UNSUPPORTED_ESM_URL_SCHEME]: Only URLs with a scheme in: file, data, and node are supported
by the default ESM loader. Received protocol 'bun:'
```

The same build with `bun --bun ./node_modules/.bin/vite build` succeeds. So `package.json` build scripts must use `bun --bun vite build` once server code imports `bun:sqlite`.

This contradicts a doc-only claim that `bun:` imports inside SvelteKit server code are unverified or risky. They were verified twice: in dev, and in an adapter-node production build.

**Process count in dev:** `pgrep -P <vite pid>` lists no children. `bun --bun vite dev` is a single process.

**HMR:** partly verified, with an important caveat. A headless WebSocket client connected to the Vite HMR socket (`vite-hmr` subprotocol) and received `{"type":"connected"}`. After an edit to `src/routes/pos/+page.svelte`, the Vite log printed `[vite] (ssr) page reload src/routes/pos/+page.svelte` and no client update was pushed. The **same test on Node produced exactly the same log line and the same silence.** So Bun and Node behave identically here; there is no Bun-specific HMR regression in this test. The absence of a client `update` message is most likely because no real browser had loaded the client module graph. A browser check was not done.

## Q2. Is `svelte-adapter-bun` maintained against `@sveltejs/kit` latest?

Three separate packages exist. Only one fits SvelteKit 2.x today.

### `svelte-adapter-bun` (gornostay25)

- Latest **1.0.1**, published **2025-10-22**. Last commit **2025-10-22**. 16 versions, first in 2022.
- `peerDependencies`: `@sveltejs/kit ^2.4.0`, `typescript ^5`. Runtime dependency: `rolldown ^1.0.0-beta.38`.
- Kit 2.59.1 satisfies `^2.4.0`. The adapter declares no peer on `svelte` or `vite`, so those are decided by Kit.
- **Peer conflict:** the repo `package.json` declares `typescript ^6.0.3`; the adapter asks for `^5`.
- Output in `build/`: `index.js`, `handler.js`, `env.js`, `client/`, `prerendered/`, `server/index.js`, `server/manifest.js`, `server/chunks/`.
- Start: `bun run ./build/index.js`. Bun's own SvelteKit guide points at this adapter.
- **Embeddable:** yes. `build/handler.js` exports `getHandler`, which returns `{ fetch, websocket }`. `build/index.js` spreads that into `Bun.serve()`. You can import `getHandler` and put its `fetch` inside your own `Bun.serve()`.
- Caveat: the adapter patches Kit's built `server/index.js` with a regular expression to add websocket hooks. That is fragile across Kit versions, but it affects the websocket path only, not `fetch`.
- Sources: https://registry.npmjs.org/svelte-adapter-bun , https://github.com/gornostay25/svelte-adapter-bun , https://bun.com/guides/ecosystem/sveltekit

### `svelte-adapter-bun-next` (TheOrdinaryWow)

- Latest **1.0.3**, published **2025-04-15**. Last commit **2025-04-15**. Stale for about 17 months.
- **No `peerDependencies` at all.** `engines.bun >= 1.2.5`.
- Not embeddable: its server module calls `Bun.serve` at import time and exports only the promise. No bare `fetch` export.
- Source: https://github.com/TheOrdinaryWow/svelte-adapter-bun-next

### `@sveltejs/adapter-bun` (official)

- Latest **1.0.0-next.1**, published **2026-08-21**. It is the only version. Monorepo commits to `packages/adapter-bun` continue (last **2026-08-27**).
- `peerDependencies`: **`@sveltejs/kit ^3.0.0-next.0`**. Kit dist-tags today: `latest 2.70.3`, `next 3.0.0-next.25`. So it **cannot** be used on Kit 2.59.1 without moving to the Kit 3 prerelease.
- Requires **Bun 1.4+**. Build with `bun run --bun build`; start with `bun ./build`. Optional `buildOptions.compile: true` emits one executable.
- Its handler (`export async function handler(request, bun_server)`) is a real Web `Request` handler, but it is bundled into the entrypoint, not emitted as a separate file, and no embedding entry point is documented. The docs say the generated server owns `fetch` and `routes`.
- The published docs page is still 404. Source in the repo: https://github.com/sveltejs/kit/blob/main/documentation/docs/25-build-and-deploy/45-adapter-bun.md , PR https://github.com/sveltejs/kit/pull/16695 , "require Bun 1.4" https://github.com/sveltejs/kit/pull/16880

## Q3. Does `adapter-node` output run under `bun ./build/index.js`? Can its handler go inside `Bun.serve()`?

`@sveltejs/adapter-node` latest is **5.5.7** (published 2026-06-24); peer `@sveltejs/kit ^2.4.0`.

**`bun ./build/index.js`: yes, verified empirically.** In the scratch copy (Kit 2.70.3, adapter-node 5.5.7), `PORT=5197 bun ./build/index.js` printed `Listening on http://0.0.0.0:5197`. Then:

- `GET /pos` → 200 (server-rendered page)
- `GET /__bunprobe` → `{"bun":"1.4.0","versionsBun":"1.4.0","sqlite":7}`
- `GET /_app/immutable/entry/start.*.js` → 200

So pages, endpoints, static assets, and `bun:sqlite` all work in one Bun process.

Documented caveat: Bun marks `node:http` "fully implemented", but the `fd` option of `listen()` is ignored. adapter-node uses `server.listen({ fd: 3 })` for systemd socket activation. Normal host/port start is unaffected. Source: https://bun.com/docs/runtime/nodejs-apis

**Mounting the handler: three different answers. Keep them apart.**

1. **`bun ./build/index.js`** — works, supported, one process. Above.
2. **`build/handler.js` inside `node:http`** — works, verified. A 12-line script that does `createServer((req, res) => handler(req, res, next)).listen(5195)` under `bun` served pages (200), endpoints (`bun:sqlite` probe OK), and `/_app/*` assets (200), with `bun:sqlite` used in the same file. This is the official "Custom server" path: the docs say `handler` suits Express, Connect, Polka, "or even just the built-in `http.createServer`". **This is Bun's `node:http`, not `Bun.serve`.** Source: https://svelte.dev/docs/kit/adapter-node
3. **`Bun.serve()` with adapter-node's handler — not a supported surface.** `handler` is a Connect-style `(req, res, next)` middleware and needs a Node `IncomingMessage` (through `getRequest` from `@sveltejs/kit/node`). `Bun.serve` gives a Web `Request`. There is no supported bridge.

I did test the internal route: importing `Server` and `manifest` from adapter-node's hashed chunks (`build/server/chunks/index.js-*.js`, `manifest.js-*.js`) and calling `server.respond(request, { getClientAddress })` inside `Bun.serve`. Pages and endpoints returned 200 next to my own `/api/mine` route. **But** I had to `sed`-patch the chunk files to re-export `Server` and `manifest` (they export only `S` and `m` under hashed filenames), and `/_app/immutable/*` returned **404** because `server.respond` does not serve `build/client`. Record this as "not a supported surface: it needs patched build output and your own static file serving", not as a path.

If a real `Bun.serve` entry is wanted with Kit 2, `svelte-adapter-bun` is the sanctioned way (Q2): it exports `getHandler().fetch`.

**One more fact about the shapes:** the repo uses `adapter-static` today (a client-only SPA), and `app/server/index.ts` owns `/api/*`. Once SvelteKit owns `/api/*` as `+server.ts` routes, the embedding question in Q3 mostly disappears — the count is one process either way. This is a fact about the shapes, not a recommendation.

## Q4. Does `bun-plugin-svelte` support SSR plus a client bundle, with HMR through `Bun.serve()`?

- Latest npm **0.0.6**, published **2025-03-26**. Only 6 versions ever. Peer `svelte ^5`. Source: https://registry.npmjs.org/bun-plugin-svelte , repo `oven-sh/bun`, `packages/bun-plugin-svelte`.
- Last **functional** commit 2025-03-25. Everything after is repo-wide chore work. About 18 months with no feature work.
- The word "experimental" is not in the package, but the signals are: version 0.0.x, and the README still says "not published to npm yet" in its bundler example.
- Bun's own docs do not document the plugin.

**SSR: works, although the README says it does not.** The README says server-side imports are "not yet supported". The shipped source contradicts it: `src/options.ts` has `forceSide?: "client" | "server"`, and `generateSide()` maps `target: "node" | "bun"` to `server`.

Verified: with `bunfig.toml` `preload` registering `Bun.plugin(SveltePlugin({ forceSide: "server" }))`, `render(Counter, { props: { start: 5 } })` from `svelte/server` returned `<button class="svelte-...">count is 5</button>`. It also worked inside a live `Bun.serve()` fetch handler, and next to an HTML import in the same process.

**Caveat: SSR emits no CSS.** `css: "external"` is hard-coded, and the virtual CSS import is only added when `generate != "server"`. `render()` returns an empty `head`. You must link the client bundle CSS yourself.

**Runtime plugin and bundler plugin are separate registrations. Both are needed.**

| Path | Config | Result |
| --- | --- | --- |
| Runtime `import "./X.svelte"` in a Bun process | `bunfig.toml` `preload` calling `Bun.plugin(SveltePlugin(...))` | works |
| Client bundle | `Bun.build({ target: "browser", plugins: [SveltePlugin(...)] })` | works, emits JS + CSS |
| `Bun.serve()` HTML import dev server | `bunfig.toml` `[serve.static] plugins = ["bun-plugin-svelte"]` | works, CSS auto-linked |

Three **silent** failure modes (exit code 0, wrong output, no error):

1. `bun build ./client.ts --target=browser` from the CLI with only `preload` set: the `.svelte` file is copied as an asset and the component is not in the bundle.
2. `bun build ./index.html` from the CLI ignores `[serve.static] plugins`. That key applies to `Bun.serve` / `bun ./index.html` only.
3. Dev server with `preload` but without `[serve.static] plugins`: the page bundles, but the component is not compiled and no CSS is linked.

The working production path is a `Bun.build()` script that passes the plugin explicitly. **This matters for the Pi deploy step**, which today calls a build command: a plain `bun build` CLI call does not compile Svelte.

Note: the string form `plugins = ["bun-plugin-svelte"]` loads the default export, hard-coded to `SveltePlugin({ development: true })`. bunfig gives no way to pass options.

**HMR: wired up for Svelte, not verified live.** The README says the plugin integrates with Bun's fullstack dev server "giving you HMR". `src/index.ts` passes `hmr: Boolean(args.hmr ?? process.env.NODE_ENV !== "production")` into Svelte's `compile()`, so it uses Svelte 5's own HMR. Bun's docs (`bundler/hot-reloading.mdx`) mention "React Fast Refresh or a plugin calling it for you", which covers this case. Verified statically only: the dev bundle contains HMR runtime code. A live browser hot-swap was not tested, and `hydrate()` on SSR output was not tested.

**What SvelteKit gives that this does not.** The plugin is a compiler plugin only:

- filesystem routing, `+page.svelte` / `+layout.svelte` / `+error.svelte`, nested layouts, route params, route groups
- `load` functions (universal and server), `+page.server.ts`, streaming from `load`, `depends` / `invalidate`
- form actions and `use:enhance`
- `hooks.server.ts` / `hooks.client.ts`, `handle`, `handleFetch`, `handleError`, `locals`
- client router: no `goto`, no link interception, no `data-sveltekit-preload-data`, no scroll restoration
- `$app/*` modules, `$env/*`, `$lib` alias, `$service-worker`
- adapters, prerendering, `csr` / `ssr` / `prerender` page options, CSP nonce, snapshots
- SSR-to-client hydration coordination and serialized data transfer: you hand-roll `render()` plus `hydrate()` and pass props yourself
- source maps (listed as not supported)
- **no preprocessor pipeline at all**: no `preprocess` option, nothing reads `svelte.config.js`, so no SCSS / PostCSS / Tailwind inside `<style>` blocks

Exposed options are only `forceSide`, `development`, `runes`, and a `compilerOptions` subset (`customElement`, `runes`, `modernAst`, `namespace`).

## Q5. Process counts

Today: `bun run dev` = **2 processes** (`bun --hot ./index.ts` and `bun run --cwd app/web dev`). Production = **1 process** (`bun ./index.ts` serves `/api/*` and the static build).

| Path | Dev processes | Prod processes | Notes |
| --- | --- | --- | --- |
| Today (`adapter-static` + separate API server) | 2 | 1 | The `/api` proxy in `vite.config.ts` joins them in dev |
| SvelteKit with `bun --bun vite dev`, API as `+server.ts` | **1** (verified: no child processes) | 1 | Build needs `bun --bun vite build` |
| `svelte-adapter-bun` | 1 in dev (as above) | 1 (`bun ./build/index.js`), or embedded in the repo `Bun.serve()` through `getHandler().fetch` | Only Kit-2 adapter with an embeddable Web `fetch` |
| `adapter-node` under Bun | 1 in dev (as above) | 1 (`bun ./build/index.js`), or 1 with `node:http` + `handler` | Cannot mount in `Bun.serve` |
| `bun-plugin-svelte` (no SvelteKit) | 1 (`Bun.serve` dev server) | 1 | Needs a separate `Bun.build()` step for the client bundle |

Every candidate reaches one process in production. The difference in dev is that today's split is the only 2-process case, and it is caused by the API server sitting outside the Vite server, not by a Bun limit.

## Method and reproduction

- Dev and HMR tests ran in `app/web` at repo versions. The temporary probe route was deleted; the working tree was restored.
- adapter-node tests ran in a throwaway copy at `/private/tmp/.../scratchpad/kitnode` with its own `bun install` (newer versions, listed above).
- `bun-plugin-svelte` tests ran in a throwaway copy at `/private/tmp/.../scratchpad/bunsvelte` (bun 1.4.0, svelte 5.57.0, plugin 0.0.6).
- Node control tests used `node_modules/.bin/vite` on node v24.16.0.

## Not verified

- Live browser HMR (Vite under Bun, and `bun-plugin-svelte`). Both were checked headlessly only. The Vite result was **identical on Node and Bun**.
- `hydrate()` on `bun-plugin-svelte` SSR output.
- `svelte-adapter-bun@1.0.1` was not built or run here. Its claims come from its source and README.
- systemd socket activation (`listen({ fd: 3 })`) under Bun. Bun ignores the `fd` option.
