# Roommates

Go web app: Gin + templ + htmx + Tailwind v4 / shadcn-templ + sqlc + PostgreSQL. Layout is described in README.md.

## Commands

Run from `app/` (the Go module root). Anything that touches the database (`go test`, `task dev`, the server) needs the compose `db` service on 127.0.0.1:5433, and the server also needs `rustfs` (S3 API on 127.0.0.1:9000): run `docker compose up -d --wait db rustfs` first (`task dev` does it itself). DB tests use an in-memory store, not RustFS.

- `task dev` — hot reload (tailwind watch + shadcn-templ bundle watch + `templ generate --watch`; `restart-watch` signals `scripts/dev-run.sh` to restart the server under dlv on :2345 when a Go or embedded `web/assets` file changes. Don't give templ a `--cmd`: on Windows one failed TASKKILL stops it running for good)
- `task build` — full production build → `bin/app`
- `go build ./...` / `go vet ./...` — quick check
- `sh scripts/lint.sh` — staticcheck + golangci-lint
- `sh scripts/format.sh` — `templ fmt` + `gofmt`
- `go test ./...` — DB tests are skipped unless `TEST_POSTGRES_DB` names a throwaway database (the migration test drops all tables); never point it at the dev database

Code generation (regenerate after editing the source, never edit the output):

- `sqlc generate` — `internal/db/queries/*.sql` + `internal/db/migrations/` → `internal/db/dbqueries/`
- `go tool templ generate` — `*.templ` → `*_templ.go`
- `go generate ./tools/...` — `internal/locales/{et,en}/*.yaml` → `locales-keys.go`, `locales-languages.go` (et.yaml is the baseline)
- `go tool shadcn-templ bundle` — component scripts → `web/assets/public/js/`
- `task vendor` (`go run ./tools/vendor-assets`) — downloads htmx and the fonts pinned in `web/assets/vendor.json` (hashes in `vendor.lock.json`; the files themselves are never committed); `--check` lists newer versions (exit 1 if any), `--update [pkg…]` moves to the latest

From the repo root: `docker compose up --build` runs db + app (app on https://127.0.0.1:8000). Env comes from `.env` (see `.env.example`). The server loads `./certificates/server.{pem,key}` relative to the working directory, so run it from `app/`.

## Conventions

- Commits follow [Conventional Commits 1.0.0](https://www.conventionalcommits.org/en/v1.0.0/) (`feat:`, `fix:`, `refactor:`, `chore:`, `docs:` …).
- Work stays on the `fresh-start` branch for now.

## graphify

This project has a knowledge graph at graphify-out/ with god nodes, community structure, and cross-file relationships.

Rules:
- When graphify-out/graph.json exists, run `graphify query "<question>"` before your first grep, search or source-file read for a codebase question, and put this rule in every subagent prompt that explores code. Read raw files only for lines the graph or an earlier read already located. Use `graphify path "<A>" "<B>"` for relationships and `graphify explain "<concept>"` for focused concepts. These return a scoped subgraph, usually much smaller than GRAPH_REPORT.md or raw grep output.
- If graphify-out/wiki/index.md exists, use it for broad navigation instead of raw source browsing.
- Read graphify-out/GRAPH_REPORT.md only for broad architecture review or when query/path/explain do not surface enough context.
- After modifying code, run `graphify update .` to keep the graph current (AST-only, no API cost).
- A committed graphify-out/graph.json is rewritten whole on every build, so a large diff to it in a pull request is expected. Do not flag it as churn or suggest skipping the graph rebuild or its hooks.
