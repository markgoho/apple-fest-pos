# Apple Fest POS

The POS is **one Go binary**, cross-compiled on the Mac and copied to a Raspberry Pi. See [ADR-0002](./docs/adr/0002-one-go-binary.md).

## Go first

- The Go standard library is the default: `net/http`, `html/template`, `encoding/json`, `net.Dial`, `ListenAndServeTLS`. Add a dependency only when the standard library truly cannot do the job.
- `modernc.org/sqlite` for SQLite, because it is pure Go. Import it as `_ "modernc.org/sqlite"` and open with `sql.Open("sqlite", path)`. Do not use `mattn/go-sqlite3`, which needs cgo.
- Build with `CGO_ENABLED=0` always. This is what lets the Mac make a static `linux/arm64` binary with no cross toolchain.
- The SQLite pool holds one connection (`SetMaxOpenConns(1)`). Every read and write inside an open transaction must go through the `*sql.Tx`, never the `*sql.DB`, or the call deadlocks.
- Store timestamps as strings in the `2006-01-02T15:04:05.000Z07:00` layout in UTC. Never hand a `time.Time` to the driver: the business date is the first ten characters of that string, so a different text format breaks the sales grouping silently.
- Money is cents as `int`. Never use a float.

## Commands

```sh
go test ./...                                          # run the tests
go run ./cmd/pos                                       # run the server
go vet ./... && gofmt -l .                             # check before commit
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build ./cmd/pos   # build for the Pi
```

## Layout

- `cmd/pos/` reads the environment and starts the server.
- `internal/pos/` holds the domain: menu, database, orders, ESC/POS, printer, sales, kitchen, HTTP.
- Environment variables: `PORT`, `SQLITE_PATH`, `ORDER_NUMBER_START`, `PRINTER_ENABLED`, `PRINTER_HOST`, `PRINTER_PORT`.

## Browser code

No JavaScript framework, no bundler, no npm dependency ships in the app. The cart is hand-written JavaScript served as a static file.

## Legacy

`src/`, `tests/`, `package.json`, and `vite.config.ts` are the old SvelteKit and Bun server. They stay only as reference until the screens are rebuilt as server-rendered HTML. Do not add to them.

## Agent skills

### Issue tracker

Issues live as GitHub issues in `markgoho/apple-fest-pos`, managed with the `gh` CLI. See `docs/agents/issue-tracker.md`.

### Triage labels

The five canonical triage roles, each label string equal to its role name. See `docs/agents/triage-labels.md`.

### Domain docs

Single-context: `CONTEXT.md` and `docs/adr/` at the repo root. See `docs/agents/domain.md`.
