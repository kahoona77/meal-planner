# Modernization Plan

Goal: a simple, modern stack. Go stays the backend and keeps rendering HTML on the server; the Node toolchain goes away entirely.

## Overview

| # | Work package                                        | Status         |
|---|-----------------------------------------------------|----------------|
| 0 | Tooling & dependency refresh (mise, Go, CI)         | ✅ Done        |
| 1 | Replace Echo with `net/http`                        | ✅ Done        |
| 2 | Replace logrus with `log/slog`                      | ✅ Done        |
| 3 | Replace Lit components with plain HTML / vanilla JS | ✅ Done        |
| 4 | Tailwind 4 standalone CLI, remove Node & Vite       | 🚧 In progress |
| 5 | Migrate templates to templ                          | ⬜ Open        |
| 6 | Embed static assets, single binary                  | ⬜ Open        |
| 7 | Add htmx where it helps                             | ⬜ Open        |
| 8 | Cleanup & docs                                      | ⬜ Open        |
| 9 | Markdown descriptions instead of Trix (later)       | ⬜ Open        |

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

Output uses the [tint](https://github.com/lmittmann/tint) handler (compact, colored on a terminal, plain otherwise).

Along the way: removed duplicate error logs in handlers (returned errors are logged centrally) and fixed tag log messages that said "meal".

## 3. Replace Lit components with plain HTML / vanilla JS ✅

Prerequisite for dropping Node: no TypeScript/Lit build anymore.

- [x] `toggle-visibility` (tag list): `data-toggle` buttons + `public/toggle.js` (event delegation, no build)
- [x] `multi-select` (meal edit, wizard): checkbox chips (`_tag-select.html` partial), submitted as repeated form fields; form parsing in `web/views/meals.go` and `wizard.go` uses `formTagIds`
- [x] Remove unused `my-element.ts`, `lit.svg` and `vite.svg`
- [x] Remove `index.ts`, TypeScript and `tsconfig.json`; Vite only builds `index.css` now
- [x] Keep `image-select.js` (already vanilla); check it still works
- [x] Trix editor (vendored `trix.js`/`trix.css`): updated 1.3.1 → 2.1.19 (old files were unmodified upstream copies)
- [x] Remove `test.html` from `public/img` (not needed)

Verified: tag preselection identical to the old version (meal edit, wizard); saving meals and wizard with/without tags; tag edit toggle and Trix 2 in the browser, no console errors. Selected chips use a tinted background of the tag color (`color-mix`) so they stay readable for light and dark colors.

## 4. Tailwind 4 standalone CLI, remove Node & Vite

- [x] Add Tailwind CLI to `mise.toml` (`github:tailwindlabs/tailwindcss`, not in the mise registry)
- [x] Migrate the CSS to v4 with the official upgrade tool (`ring-opacity-*` replaced by `ring-blue-400/75` first, the tool fails on it); result in `web/styles/main.css` with `source(none)` + explicit `@source` (otherwise Tailwind scans `vendor/`)
- [x] Element styles (`body`, `a`, `button`, ...) moved into `@layer base`: unlayered CSS beats all Tailwind layers in v4, e.g. `button { background: none }` would remove `bg-primary` from buttons
- [x] Unused `secondary` color dropped, `neutral` (was an alias of v3 gray) replaced by `gray`
- [x] New asset layout: `web/static/` (images, manifest, trix, image-select, toggle, built `main.css`), CSS source in `web/styles/`
- [x] Cache busting without Vite manifest: content hash as `?v=` query parameter (`web/static.go`, recomputed when the file changes); versioned requests get `Cache-Control: immutable`, others `no-cache`
- [x] Template funcs `assetUrl`/`publicUrl` replaced by `asset`; all relative `assets/...` references go through it
- [x] Remove `vite-manifest.json` handling, `DEV_MODE`, `web/manifest.go` + tests + `test/data`
- [x] Remove `web/assets/` (`package.json`, lockfile, `node_modules`, Vite/PostCSS configs)
- [x] mise tasks: `css` (build), `dev` (Tailwind `--watch=always` + Go server); Node dropped from `mise.toml`
- [x] Dockerfile: Node stage dropped, Tailwind CLI downloaded in the Go build stage (version must match `mise.toml`)
- [ ] Visually compare all pages before/after

Verified so far: Docker build; all pages and every referenced asset return 200 in the container; cache headers; `mise run dev` rebuilds the CSS on change. Every class of the old v3 CSS exists in the new CSS, except the two renamed ones (`rounded` → `rounded-sm`, `fill-neutral-500` → `fill-gray-500`). A visual check in the browser is still open (Claude in Chrome was not available).

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

## 9. Markdown descriptions instead of Trix (later)

Decided to keep Trix for now (updated in WP 3). Later, meal descriptions (recipes: headings, ingredient lists, paragraphs) should be stored as Markdown instead of HTML. Best done after templ (WP 5).

- [ ] Edit with a plain `<textarea>`, optionally a few small toolbar buttons that insert Markdown syntax (bold, heading, list)
- [ ] Render server-side with [goldmark](https://github.com/yuin/goldmark); replaces the unchecked `htmlSafe` output of stored HTML
- [ ] goose migration converting existing HTML descriptions with [html-to-markdown](https://github.com/JohannesKaufmann/html-to-markdown) (Go migration, not SQL); check the result on a copy of the DB
- [ ] Remove `trix.js` / `trix.css`

## Open decisions

- Commit generated templ files or generate during build? (WP 5)
- Live reload in dev: templ `--proxy` or none? (WP 5)
- Which htmx interactions are worth it? (WP 7)
