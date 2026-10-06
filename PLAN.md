# Modernization Plan

Goal: a simple, modern stack. Go stays the backend and keeps rendering HTML on the server; the Node toolchain goes away entirely.

## Overview

| #  | Work package                                         | Status  |
|----|------------------------------------------------------|---------|
| 0  | Tooling & dependency refresh (mise, Go, CI)          | ✅ Done |
| 1  | Replace Echo with `net/http`                         | ✅ Done |
| 2  | Replace logrus with `log/slog`                       | ✅ Done |
| 3  | Replace Lit components with plain HTML / vanilla JS  | ⬜ Open |
| 4  | Tailwind 4 standalone CLI, remove Node & Vite        | ⬜ Open |
| 5  | Migrate templates to templ                           | ⬜ Open |
| 6  | Embed static assets, single binary                   | ⬜ Open |
| 7  | Add htmx where it helps                              | ⬜ Open |
| 8  | Cleanup & docs                                       | ⬜ Open |

Status: ⬜ Open · 🚧 In progress · ✅ Done

## Target stack

| Layer         | Before                                  | After                                       |
|---------------|-----------------------------------------|---------------------------------------------|
| HTTP          | Echo                                    | `net/http` (Go 1.22+ routing patterns)      |
| Templates     | `html/template` + custom renderer       | [templ](https://templ.guide)                |
| Interactivity | Lit web components, full page reloads   | htmx (vendored, no build) + small vanilla JS |
| CSS           | Tailwind 3 via PostCSS/Vite             | Tailwind 4 standalone CLI (via mise)        |
| Logging       | logrus                                  | `log/slog`                                  |
| Database      | SQLite + sqlx + goose                   | unchanged                                   |
| Delivery      | binary + `web/` directory next to it    | single binary with embedded assets          |

Required tools afterwards: Go, templ, Tailwind CLI — all pinned in `mise.toml`.

Every work package ends with: build + tests green, a smoke test of all routes against a copy of the DB, a check in the browser, one commit.

---

## 0. Tooling & dependency refresh ✅

- `mise.toml` with Go 1.27, Node 24 and tasks
- `AGENTS.md` (+ `CLAUDE.md` importing it)
- Go 1.27 and all Go dependencies updated, `vendor/` regenerated
- Dockerfile base images and GitHub Actions updated
- CI tags images: `latest` (main), `<version>` (`v*` tags), `dev` (other branches)
- Frontend dependencies updated within their major versions

## 1. Replace Echo with `net/http` ✅

Go's `ServeMux` supports methods and path parameters (`GET /meals/{id}`), so Echo is no longer needed.

- [x] Replace `core.App` / `core.Group` with a thin wrapper around `http.ServeMux` that keeps the `func(*WebContext) error` handler signature (keeps the diff in `web/views` small)
- [x] `WebContext` on top of `http.ResponseWriter` + `*http.Request`: `PathValue` instead of `Param`, `ParamAsInt`, form access, `Redirect` with `BASE_PATH` prefix, `RenderTemplate`
- [x] Central error handling (handler returns error → log + 500 page)
- [x] Mount everything under `BASE_PATH` (`http.StripPrefix` or prefixed patterns)
- [x] Static files via `http.FileServer`, `/files/{id}` with cache header
- [x] Request logging middleware; CORS middleware dropped (same-origin app)
- [x] Remove `e.Debug = true` and the `DELETE`-registers-`PUT` bug along the way
- [x] Remove Echo from `go.mod`, re-vendor

Verified: old (Echo) and new build run side by side against copies of the DB; all GET pages return byte-identical HTML, all POST flows (tag, meal with image upload, meal edit, day select, wizard, delete) lead to identical results. Only intended difference: missing records return 404 instead of a 500 with a JSON error.

## 2. Replace logrus with `log/slog` ✅

- [x] `slog` text handler on stdout, level configurable via env (e.g. `LOG_LEVEL`)
- [x] Replace all logrus calls (`core`, `web`, `web/views`, `wizard`)
- [x] Remove logrus from `go.mod`, re-vendor

Along the way: removed duplicate error logs in handlers (returned errors are logged centrally) and fixed tag log messages that said "meal".

## 3. Replace Lit components with plain HTML / vanilla JS

Prerequisite for dropping Node: no TypeScript/Lit build anymore.

- [ ] `toggle-visibility` (tag list): small vanilla JS or `<details>`; later possibly htmx (WP 7)
- [ ] `multi-select` (meal edit, wizard): checkbox chips styled with Tailwind, submitted as normal form fields; adjust form parsing in `web/views/meals.go` and `wizard.go`
- [ ] Remove unused `my-element.ts` and `vite.svg`
- [ ] Keep `image-select.js` (already vanilla); check it still works
- [ ] Trix editor (vendored `trix.js`/`trix.css`): update to the current version
- [ ] Decide: is `test.html` in `public/img` still needed?

## 4. Tailwind 4 standalone CLI, remove Node & Vite

- [ ] Add Tailwind CLI to `mise.toml`
- [ ] Migrate `index.css` to v4 (`@import "tailwindcss"`, theme colors from `tailwind.config.js` into `@theme`, `@source` for templates; check removed utilities like `ring-opacity-*`)
- [ ] New asset layout, e.g. `web/static/` (images, manifest, trix, image-select, htmx, built CSS)
- [ ] Cache busting without Vite manifest (e.g. content hash as query parameter, computed at startup)
- [ ] Remove `vite-manifest.json` handling, `DEV_MODE` dev-server logic, `web/manifest.go` + tests
- [ ] Remove `package.json`, `package-lock.json`, `node_modules`, Vite/PostCSS/TS configs
- [ ] mise tasks: `css` (build), `dev` (Tailwind `--watch` + Go server); drop Node from `mise.toml`
- [ ] Dockerfile: drop the Node stage, run the Tailwind CLI in the Go build stage
- [ ] Visually compare all pages before/after

## 5. Migrate templates to templ

- [ ] Add templ to `mise.toml` (CLI) and `go.mod` (runtime)
- [ ] Decide: commit generated `*_templ.go` files or generate in build (affects Dockerfile and CI)
- [ ] Layout component (replaces `base.html`) with header, bottom nav, asset links
- [ ] Migrate page by page: index/week, meal-of-day, meal-of-day select, meals list, meal edit, tags, wizard
- [ ] Replace template funcs (`basePath`, `assetUrl`, `fileUrl`, `formatWeekday`, `isToday`, ...) with plain Go functions
- [ ] Remove `web/renderer.go`, `web/tmpl/`, `TemplateData`
- [ ] Dev workflow: `templ generate --watch` in the `dev` task (decide on live reload via `--proxy`)

## 6. Embed static assets, single binary

- [ ] `embed.FS` for `web/static` (incl. built CSS); templ templates and migrations are already compiled in
- [ ] Binary no longer depends on the working directory
- [ ] Dockerfile: runtime image only contains the binary (+ ca-certificates)
- [ ] Update `AGENTS.md` (remove "must run from repo root")

## 7. Add htmx where it helps

- [ ] Vendor `htmx.min.js` into `web/static`, include in layout
- [ ] Candidates (decide case by case):
  - week navigation (`/offset/{n}`) without full reload
  - meal-of-day select → swap only the day card
  - tag edit inline (replaces `toggle-visibility`)
  - wizard result preview
- [ ] Handlers return fragments for htmx requests (`HX-Request` header), full pages otherwise

## 8. Cleanup & docs

- [ ] Handler tests with `httptest` for the main routes (in-memory SQLite)
- [ ] Update `AGENTS.md`, `README.md` (how to run, env vars), `mise.toml` tasks
- [ ] Check `.dockerignore`, `.gitignore`
- [ ] Optional: release tag `v1.0.0` to test the versioned image

## Open decisions

- Commit generated templ files or generate during build? (WP 5)
- Live reload in dev: templ `--proxy` or none? (WP 5)
- Which htmx interactions are worth it? (WP 7)
