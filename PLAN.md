# Modernization Plan

Goal: a simple, modern stack. Go stays the backend and keeps rendering HTML on the server; the Node toolchain goes away entirely.

## Overview

| # | Work package                                        | Status         |
|---|-----------------------------------------------------|----------------|
| 0 | Tooling & dependency refresh (mise, Go, CI)         | ✅ Done        |
| 1 | Replace Echo with `net/http`                        | ✅ Done        |
| 2 | Replace logrus with `log/slog`                      | ✅ Done        |
| 3 | Replace Lit components with plain HTML / vanilla JS | ✅ Done        |
| 4 | Tailwind 4 standalone CLI, remove Node & Vite       | ✅ Done        |
| 5 | Migrate templates to templ                          | ✅ Done        |
| 6 | Embed static assets, single binary                  | ✅ Done        |
| 7 | Add htmx where it helps                             | ✅ Done (A)    |
| 8 | Cleanup & docs                                      | ✅ Done        |
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

## 4. Tailwind 4 standalone CLI, remove Node & Vite ✅

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
- [x] Visually compare all pages before/after (checked by hand)

Verified so far: Docker build; all pages and every referenced asset return 200 in the container; cache headers; `mise run dev` rebuilds the CSS on change. Every class of the old v3 CSS exists in the new CSS, except the two renamed ones (`rounded` → `rounded-sm`, `fill-neutral-500` → `fill-gray-500`).

## 5. Migrate templates to templ ✅

- [x] templ CLI via the `tool` directive in `go.mod` (`go tool templ`, version always matches the runtime library, works offline from `vendor/`); no `mise.toml` entry needed. Downside: `vendor/` grows from 12 to 28 MB (mostly `golang.org/x/tools`)
- [x] Decided: generated `*_templ.go` are **not** committed (gitignored); `mise run build`/`test`, `mise run dev` and the Dockerfile generate them
- [x] Layout component (replaces `base.html`) with header, bottom nav, asset links; page specific `<head>` tags as parameter
- [x] Migrate page by page: week, meal-of-day, meal-of-day select, meals list, meal edit, tags, wizard (one `.templ` file per area next to its handlers in `web/views`)
- [x] Replace template funcs with plain Go functions (`web.BasePath`/`web.Asset` via request context, `weekdayName`, `fileURL`, `colorStyle`)
- [x] Remove `web/renderer.go`, `web/format.go`, `web/tmpl/`, `TemplateData`, `RenderTemplate`
- [x] Decided: live reload via `templ generate --watch --proxy` (`mise run dev`, http://localhost:7331/meal-planner)

Verified: old (`html/template`) and new build side by side against copies of the DB; all pages have an identical DOM (tags, attributes, text), apart from whitespace in class attributes and a trailing `;` in inline styles. Only intended change: meal cards in the list are links (`<a href>`) instead of `<div onclick>`. All form flows work, also in the Docker image; live reload picks up `.templ` changes. Missed at first: `@source` in `web/styles/main.css` still pointed to the deleted `web/tmpl`, so utilities only used in templates were missing from the CSS (broken bottom nav); now `@source "../views/*.templ"`. After the fix, headless Brave screenshots of all pages (desktop and mobile) are pixel-identical to WP 4.

## 6. Embed static assets, single binary ✅

- [x] `embed.FS` for `web/static` (incl. built CSS, `web.StaticFS`); templ templates and migrations are already compiled in
- [x] `STATIC_DIR` env var serves the files from disk instead (set by `mise run dev`, so CSS changes show up without a restart)
- [x] Binary no longer depends on the working directory
- [x] Dockerfile: runtime image only contains the binary (+ ca-certificates)
- [x] Update `AGENTS.md` (remove "must run from repo root")
- [x] Fixed along the way: `.gitignore` rule `*_templ.go` also hid two generated files inside `vendor/` (templ CLI), so a build from a clean checkout failed; rule is now limited to `/web/**`

Verified: binary started outside the repo serves all pages and assets; build from only the files git would commit (simulated CI) works; Docker image contains only the binary; screenshots of all pages pixel-identical to WP 5; `STATIC_DIR` picks up file changes with a new hash.

## 7. Add htmx where it helps ✅ (stage A)

Decided: stage A now, B maybe later, C as a separate feature package after the modernization. Not doing: own fragment handling for the week navigation (boost covers it), wizard preview.

Stage A – basics:
- [x] htmx 4.0.0 vendored into `web/static/htmx.min.js` (htmx 4: inheritance is explicit, `hx-boost:inherited`; error responses are swapped by default)
- [x] `hx-boost:inherited="true"` on `<body>`: all links and forms load in the background and swap the body, URL and back button keep working (boosted forms push the final URL after the redirect)
- [x] All scripts (htmx, Trix, image-select, toggle) and `trix.css` are loaded in the layout on every page: a boosted swap does not load `<head>` scripts of the next page. `layout()` lost its `head` parameter
- [x] `hx-confirm` on deleting a meal (it deleted without asking before)
- [x] Wizard form got an explicit `action`: htmx resolves a missing `action` against `<base href>` instead of the current URL, so generating posted to `/meal-planner/` (405)

Verified with a scripted headless Brave (chromedp, not in the repo) against a copy of the DB: 28 checks for navigation without reload, back button, Trix/image-select after boosted navigation, saving with redirect, image upload, delete with cancel/confirm, tag edit, wizard, day select, 404 page. Screenshots pixel-identical to the previous commit.

Stage B – later, maybe: edit and add tags inline (form posts via htmx, handler returns only `tagCard` for `HX-Request`).

Stage C – separate feature package: "re-roll" a single day in the week view; pick the meal of a day directly in the week view.

## 8. Cleanup & docs ✅

- [x] Linter: `go vet` + `staticcheck` (default checks; binary via mise and downloaded in the Dockerfile, not in `go.mod`), `mise run lint`, also in the Docker build so CI fails on findings. Fixed the only finding (unused constant)
- [x] Handler tests with `httptest` for all routes (`main_test.go`: pages, 404s, base path redirect, static asset caching, create/edit/delete meal incl. image upload, tags, day select, wizard), on an in-memory SQLite
- [x] For that: routes moved into `setupRoutes`, `core.NewApp`/`App.Handler`, migrations got their own package (`migrations.Up`), shared test helper `internal/testdb` (replaces the copies in the meals and wizard tests)
- [x] Update `AGENTS.md`, `README.md` (what it is, Docker, configuration, development), `mise.toml` tasks
- [x] `.dockerignore`: also `.git` and `.idea`; `.gitignore` is fine
- [ ] Optional: release tag `v1.0.0` to test the versioned image

Verified: tests stable (20 runs), Docker build from only the files git would commit runs lint and all tests, `docker run` command from the README works with an empty data directory, goose sees the existing production schema as up to date after the move of the migrations.

Noticed while writing the wizard test (not changed, possible topics for the feature package): the wizard never picks meals without tags, and it fills the days in order starting on Saturday, so an unconstrained day can take the only meal that matches a later day's tags (that day then stays empty).

## 9. Markdown descriptions instead of Trix (later)

Decided to keep Trix for now (updated in WP 3). Later, meal descriptions (recipes: headings, ingredient lists, paragraphs) should be stored as Markdown instead of HTML. Best done after templ (WP 5).

- [ ] Edit with a plain `<textarea>`, optionally a few small toolbar buttons that insert Markdown syntax (bold, heading, list)
- [ ] Render server-side with [goldmark](https://github.com/yuin/goldmark); replaces the unchecked `htmlSafe` output of stored HTML
- [ ] goose migration converting existing HTML descriptions with [html-to-markdown](https://github.com/JohannesKaufmann/html-to-markdown) (Go migration, not SQL); check the result on a copy of the DB
- [ ] Remove `trix.js` / `trix.css`

## Bug fixes found along the way

- Wizard crashed (nil pointer in `isAlreadyPicked`) when no meal matched a day's tags: the failed day was appended as `nil`. Days without a matching meal are now skipped (test in `wizard/wizard_test.go`).
- `DeleteMeal` kept the meal's `meal_tags`, so the wizard could pick deleted meals. Now deleted in the same transaction (test in `meals/repository_test.go`).
- Migration `004_remove_orphaned_meal_tags.sql` removes `meal_tags` of meals that no longer exist (one row with `meal_id = 0` in production). Checked on a copy of the production DB: 1 row removed, wizard ran 50 times without error.

## Open decisions

- Which htmx interactions are worth it? (WP 7)
