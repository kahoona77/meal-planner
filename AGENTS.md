# AGENTS.md

Guidance for AI coding agents working in this repository.

A family meal planner: server-rendered Go web app (`net/http` + `html/template`) with SQLite, styled with Tailwind (standalone CLI, no Node) and a little vanilla JS.

## Conventions

- All project content (code comments, docs, commit messages, `mise.toml` descriptions) is written in **English**.

## Commands

Toolchain (Go, Tailwind CLI) is pinned in `mise.toml`; tasks are defined there too.

```sh
mise run dev            # Tailwind --watch + Go server (:8080) in parallel
mise run css            # build web/static/main.css once
mise run build          # css, then CGO go build -> ./meal-planner
mise run test           # go test -mod=vendor ./...
go test -mod=vendor ./web -run TestAssets_URL   # single test
```

App URL: http://localhost:8080/meal-planner

- Dependencies are **vendored** (`vendor/`): always pass `-mod=vendor`, and run `go mod vendor` after changing `go.mod`.
- `go-sqlite3` needs **CGO** (`CGO_ENABLED=1`, C compiler present).
- There is no linter configured. Tests only exist for `web/static.go`.
- `web/static/main.css` is generated (gitignored); run `mise run css` before starting the server outside of `mise run dev`.
- The binary must be run from the repo root: templates (`web/tmpl`) and static files (`web/static`) are loaded from relative paths at runtime, not embedded. Only `migrations/*.sql` is embedded.

## Configuration (env vars, `core/config.go`)

| Var         | Default               | Effect                                                                            |
|-------------|-----------------------|-----------------------------------------------------------------------------------|
| `PORT`      | `8080`                |                                                                                   |
| `BASE_PATH` | `/meal-planner`       | All routes are mounted under it; static files are served under `BASE_PATH/assets` |
| `DB_FILE`   | `meal-planner.sqlite` | SQLite file; migrations run automatically on startup (goose)                      |
| `LOG_LEVEL` | `info`                | `debug`, `info`, `warn` or `error` (`log/slog`)                                   |

## Architecture

- `main.go` – route table. All routes live in one `root` group under `BASE_PATH`.
- `core/` – app bootstrap and the handler abstraction on top of `net/http` (`http.ServeMux` with Go 1.22 patterns, path params as `{id}`). Handlers have the signature `func(*core.WebContext) error` and are registered via `core.Group` (`GET`/`POST`/..., `Static`). A returned error becomes a 500 (`sql.ErrNoRows` becomes a 404). `WebContext` exposes `Db()`, `Config()`, `Param`/`ParamAsInt`, `FormValue`/`FormFile`, `RenderTemplate(code, "name.html", core.TemplateData{...})`, `Blob`, and a `Redirect` that prefixes `BASE_PATH` automatically.
- Domain packages `meals/`, `planner/`, `files/`, `wizard/` – each has `model.go` + `repository.go`; repositories are created per request with `NewRepository(ctx core.Context)` and use `sqlx` with raw SQL. `wizard` generates a random week plan from meals filtered by tags, using the meals and planner repositories.
- `web/views/` – HTTP handlers (one file per area), glue between repositories and templates.
- `web/renderer.go` – template engine. Each page template in `web/tmpl/` is parsed together with `base.html` (layout) and all `_*.html` partials; templates are **reloaded on every render**. The `funcMap` there (`basePath`, `asset`, `fileUrl`, `formatWeekday`, ...) is what templates can call. Always reference static files with `{{ asset "img/solid.svg#plus" }}`: it adds a content hash (`?v=...`, `web/static.go`), and versioned requests are cached forever by the browser.
- `web/styles/main.css` – Tailwind 4 source (theme in `@theme`, content sources via `@source`). Element styles belong in `@layer base`, otherwise they override utilities.
- `web/static/` – files served as-is under `BASE_PATH/assets`: built `main.css`, icons, PWA manifest, vendored Trix editor, and small vanilla JS (`image-select.js` web component, `toggle.js` for `data-toggle` buttons).
- Templates whose file name starts with `_` are partials and available in every page, e.g. `_tag-select.html` (tag checkboxes, data built by `newTagSelect` in `web/views/tagselect.go`, read back with `formTagIds`).
- `migrations/` – goose SQL migrations, numbered `NNN_name.sql`; add new ones there, never edit applied ones.

## Deployment

Multi-stage `Dockerfile` (builder stage downloads the Tailwind CLI, builds the CSS and a static Go binary and runs the tests; alpine runtime). The Tailwind version in the Dockerfile must match `mise.toml`. GitHub Actions (`.github/workflows/main.yml`) builds the image on every push and pushes it to Docker Hub as `kahoona/meal-planner`: `latest` from `main`, `<version>` from `v*` tags, `dev` from other branches.
