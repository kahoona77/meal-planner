# AGENTS.md

Guidance for AI coding agents working in this repository.

A family meal planner: server-rendered Go web app (`net/http` + `html/template`) with SQLite, plus a small Vite/Lit/Tailwind frontend bundle.

## Conventions

- All project content (code comments, docs, commit messages, `mise.toml` descriptions) is written in **English**.

## Commands

Toolchain (Go, Node) is pinned in `mise.toml`; tasks are defined there too.

```sh
mise run dev            # Vite dev server (:5173) + Go server in DEV_MODE (:8080) in parallel
mise run build          # npm build in web/assets, then CGO go build -> ./meal-planner
mise run test           # go test -mod=vendor ./...
go test -mod=vendor ./web -run TestManifest_File   # single test
```

App URL: http://localhost:8080/meal-planner

- Dependencies are **vendored** (`vendor/`): always pass `-mod=vendor`, and run `go mod vendor` after changing `go.mod`.
- `go-sqlite3` needs **CGO** (`CGO_ENABLED=1`, C compiler present).
- There is no linter configured. Tests only exist for `web/manifest.go` (fixtures in `test/data/`).
- The binary must be run from the repo root: templates (`web/tmpl`), assets (`web/assets/dist`) and the Vite manifest are loaded from relative paths at runtime, not embedded. Only `migrations/*.sql` is embedded.

## Configuration (env vars, `core/config.go`)

| Var         | Default               | Effect                                                                                                                                                      |
|-------------|-----------------------|-------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `DEV_MODE`  | `false`               | `true`: asset URLs point to the Vite dev server `http://localhost:5173`, no manifest needed. `false`: requires a built `web/assets/dist/vite-manifest.json` |
| `PORT`      | `8080`                |                                                                                                                                                             |
| `BASE_PATH` | `/meal-planner`       | All routes are mounted under it. Vite's `base` (`/meal-planner/assets`) is hardcoded in `vite.config.ts` and must match                                     |
| `DB_FILE`   | `meal-planner.sqlite` | SQLite file; migrations run automatically on startup (goose)                                                                                                |

## Architecture

- `main.go` – route table. All routes live in one `root` group under `BASE_PATH`.
- `core/` – app bootstrap and the handler abstraction on top of `net/http` (`http.ServeMux` with Go 1.22 patterns, path params as `{id}`). Handlers have the signature `func(*core.WebContext) error` and are registered via `core.Group` (`GET`/`POST`/..., `Static`). A returned error becomes a 500 (`sql.ErrNoRows` becomes a 404). `WebContext` exposes `Db()`, `Config()`, `Param`/`ParamAsInt`, `FormValue`/`FormFile`, `RenderTemplate(code, "name.html", core.TemplateData{...})`, `Blob`, and a `Redirect` that prefixes `BASE_PATH` automatically.
- Domain packages `meals/`, `planner/`, `files/`, `wizard/` – each has `model.go` + `repository.go`; repositories are created per request with `NewRepository(ctx core.Context)` and use `sqlx` with raw SQL. `wizard` generates a random week plan from meals filtered by tags, using the meals and planner repositories.
- `web/views/` – HTTP handlers (one file per area), glue between repositories and templates.
- `web/renderer.go` – template engine. Each page template in `web/tmpl/` is parsed together with `base.html` (layout) and all `_*.html` partials; templates are **reloaded on every render**. The `funcMap` there (`basePath`, `assetUrl`, `publicUrl`, `fileUrl`, `formatWeekday`, `json`, ...) is what templates can call. `assetUrl` resolves hashed file names through the Vite manifest in production.
- `web/assets/` – Vite project: Lit web components and Tailwind CSS. Entry points are `src/index.ts` and `src/index.css` (set in `vite.config.ts`).
- `migrations/` – goose SQL migrations, numbered `NNN_name.sql`; add new ones there, never edit applied ones.

## Deployment

Multi-stage `Dockerfile` (Node builds the frontend, Go builds a static binary and runs the tests, alpine runtime). GitHub Actions (`.github/workflows/main.yml`) builds the image on every push and pushes it to Docker Hub as `kahoona/meal-planner:latest`.
