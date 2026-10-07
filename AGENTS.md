# AGENTS.md

Guidance for AI coding agents working in this repository.

A family meal planner: server-rendered Go web app (`net/http` + [templ](https://templ.guide)) with SQLite, styled with Tailwind (standalone CLI, no Node) and a little vanilla JS.

## Conventions

- All project content (code comments, docs, commit messages, `mise.toml` descriptions) is written in **English**.

## Commands

Toolchain (Go, Tailwind CLI) is pinned in `mise.toml`, the templ CLI via the `tool` directive in `go.mod` (`go tool templ`); tasks are defined in `mise.toml`.

```sh
mise run dev            # Tailwind --watch + templ --watch with live reload proxy (:7331 -> :8080)
mise run css            # build web/static/main.css once
mise run templ          # generate *_templ.go once
mise run build          # css + templ, then CGO go build -> ./meal-planner
mise run test           # go test -mod=vendor ./...
go test -mod=vendor ./web -run TestAssets_URL   # single test
```

App URL: http://localhost:7331/meal-planner during `mise run dev` (live reload), otherwise http://localhost:8080/meal-planner

- Dependencies are **vendored** (`vendor/`): always pass `-mod=vendor`, and run `go mod vendor` after changing `go.mod`.
- `go-sqlite3` needs **CGO** (`CGO_ENABLED=1`, C compiler present).
- There is no linter configured. Tests only exist for `web/static.go`.
- Generated files are gitignored: `web/static/main.css` (`mise run css`) and `*_templ.go` (`mise run templ`). Run both after a fresh checkout before `go build`, `go vet` or `go test`; `mise run build`/`test` do it automatically.
- The binary is self-contained: templates are compiled in, `web/static` and `migrations/*.sql` are embedded (`go:embed`). `web/static/main.css` must therefore be built before `go build`.

## Configuration (env vars, `core/config.go`)

| Var          | Default               | Effect                                                                                                     |
|--------------|-----------------------|------------------------------------------------------------------------------------------------------------|
| `PORT`       | `8080`                |                                                                                                            |
| `BASE_PATH`  | `/meal-planner`       | All routes are mounted under it; static files are served under `BASE_PATH/assets`                          |
| `DB_FILE`    | `meal-planner.sqlite` | SQLite file; migrations run automatically on startup (goose)                                               |
| `LOG_LEVEL`  | `info`                | `debug`, `info`, `warn` or `error` (`log/slog`)                                                            |
| `STATIC_DIR` | (embedded)            | Serve static files from this directory instead of the embedded copy; set to `web/static` by `mise run dev` |

## Architecture

- `main.go` – route table. All routes live in one `root` group under `BASE_PATH`.
- `core/` – app bootstrap and the handler abstraction on top of `net/http` (`http.ServeMux` with Go 1.22 patterns, path params as `{id}`). Handlers have the signature `func(*core.WebContext) error` and are registered via `core.Group` (`GET`/`POST`/..., `Static`). A returned error becomes a 500 (`sql.ErrNoRows` becomes a 404). `WebContext` exposes `Db()`, `Config()`, `Param`/`ParamAsInt`, `FormValue`/`FormValues`/`FormFile`, `Render(code, component)` for templ components, `Blob`, and a `Redirect` that prefixes `BASE_PATH` automatically. `App.Use` adds middlewares.
- Domain packages `meals/`, `planner/`, `files/`, `wizard/` – each has `model.go` + `repository.go`; repositories are created per request with `NewRepository(ctx core.Context)` and use `sqlx` with raw SQL. `wizard` generates a random week plan from meals filtered by tags, using the meals and planner repositories.
- `web/views/` – one area per name: handlers in `<area>.go`, templ components in `<area>.templ` (e.g. `meals.go` + `meals.templ`). Pages wrap their content in `@layout(head)` (`layout.templ`; `head` adds page specific `<head>` tags, may be nil). Shared components: `image` (upload or placeholder), `tagSelectInput` (tag checkboxes, data from `newTagSelect` in `tagselect.go`, read back with `formTagIds`). Plain helpers (`weekdayName`, `fileURL`, `colorStyle`, `dateFormat`) live in `helpers.go`.
- `web/page.go` – `PageMiddleware` puts the base path and the asset URLs into the request context; templates call `web.BasePath(ctx)` and `web.Asset(ctx, "img/solid.svg#plus")`. Always reference static files through `web.Asset`: it adds a content hash (`?v=...`, `web/static.go`), and versioned requests are cached forever by the browser.
- Inline styles with values from the database go through `colorStyle` (only hex colors pass); meal descriptions are stored HTML from Trix and rendered with `templ.Raw`.
- `web/styles/main.css` – Tailwind 4 source (theme in `@theme`, content sources via `@source`). Element styles belong in `@layer base`, otherwise they override utilities.
- `web/static/` – embedded into the binary (`web.StaticFS`), served as-is under `BASE_PATH/assets`: built `main.css`, icons, PWA manifest, vendored Trix editor, and small vanilla JS (`image-select.js` web component, `toggle.js` for `data-toggle` buttons).
- `migrations/` – goose SQL migrations, numbered `NNN_name.sql`; add new ones there, never edit applied ones.

## Deployment

Multi-stage `Dockerfile` (builder stage downloads the Tailwind CLI, builds the CSS, generates the templ code, builds a static Go binary and runs the tests; the alpine runtime image only contains the binary). The Tailwind version in the Dockerfile must match `mise.toml`. GitHub Actions (`.github/workflows/main.yml`) builds the image on every push and pushes it to Docker Hub as `kahoona/meal-planner`: `latest` from `main`, `<version>` from `v*` tags, `dev` from other branches.
