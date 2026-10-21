meal-planner
============

A family meal planner: collect your meals with photo, tags and recipe, plan the week by hand
or let the wizard fill it based on tags per weekday.

Server-rendered Go app ([templ](https://templ.guide), [htmx](https://htmx.org), Tailwind CSS)
with a SQLite database, shipped as a single binary.

## Run with Docker

```sh
docker run -p 8080:8080 -v ./data:/data -e DB_FILE=/data/meal-planner.sqlite kahoona/meal-planner
```

Then open http://localhost:8080/meal-planner. The database is created and migrated on startup.

Image tags: `latest` (main branch), `<version>` (release tags like `v1.2.3`), `dev` (other branches).

## Configuration

| Variable     | Default               | Description                                    |
|--------------|-----------------------|------------------------------------------------|
| `PORT`       | `8080`                | HTTP port                                      |
| `BASE_PATH`  | `/meal-planner`       | URL path the app is served under               |
| `DB_FILE`    | `meal-planner.sqlite` | SQLite database file                           |
| `LOG_LEVEL`  | `info`                | `debug`, `info`, `warn` or `error`             |
| `STATIC_DIR` | (embedded)            | Serve static files from disk (development)     |

## Development

All tools are managed with [mise](https://mise.jdx.dev); `mise install` installs Go, the
Tailwind CLI and staticcheck, templ comes via `go tool`.

```sh
mise run dev     # live reload on http://localhost:7331/meal-planner
mise run test    # tests
mise run lint    # go vet + staticcheck
mise run build   # ./meal-planner
```

More details on the architecture and conventions are in [AGENTS.md](AGENTS.md).
