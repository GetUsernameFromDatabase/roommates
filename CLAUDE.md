# Roommates

Go web app: Gin + templ + htmx + Tailwind v4 / shadcn-templ + sqlc + PostgreSQL. Layout is described in README.md.

## Commands

Run from `app/` (the Go module root):

- `task dev` — hot reload (tailwind watch + shadcn-templ bundle watch + `templ generate --watch` running `scripts/dev-run.sh` under dlv on :2345)
- `task build` — full production build → `bin/app`
- `go build ./...` / `go vet ./...` — quick check
- `sh scripts/lint.sh` — staticcheck + golangci-lint
- `sh scripts/format.sh` — `templ fmt` + `gofmt`

Code generation (regenerate after editing the source, never edit the output):

- `sqlc generate` — `internal/db/queries.sql` + `internal/db/migrations/` → `internal/db/dbqueries/`
- `go tool templ generate` — `*.templ` → `*_templ.go`
- `go generate ./tools/...` — `internal/locales/{et,en}/*.yaml` → `locales-keys.go`, `locales-languages.go` (et.yaml is the baseline)
- `go tool shadcn-templ bundle` — component scripts → `web/assets/public/js/`

From the repo root: `docker compose up --build` runs db + app (app on https://127.0.0.1:8000). Env comes from `.env` (see `.env.example`). The server loads `./certificates/server.{pem,key}` relative to the working directory, so run it from `app/`.

## Conventions

- Commits follow [Conventional Commits 1.0.0](https://www.conventionalcommits.org/en/v1.0.0/) (`feat:`, `fix:`, `refactor:`, `chore:`, `docs:` …).
- Work stays on the `fresh-start` branch for now.
- Keep CRLF line endings on files that already use them (e.g. `docker-compose.yml`, `sqlc.yml`, `.vscode/launch.json`, Dockerfiles); `sed -i` strips them.

## graphify

This project has a knowledge graph at graphify-out/ with god nodes, community structure, and cross-file relationships.

Rules:
- When graphify-out/graph.json exists, run `graphify query "<question>"` before your first grep, search or source-file read for a codebase question, and put this rule in every subagent prompt that explores code. Read raw files only for lines the graph or an earlier read already located. Use `graphify path "<A>" "<B>"` for relationships and `graphify explain "<concept>"` for focused concepts. These return a scoped subgraph, usually much smaller than GRAPH_REPORT.md or raw grep output.
- If graphify-out/wiki/index.md exists, use it for broad navigation instead of raw source browsing.
- Read graphify-out/GRAPH_REPORT.md only for broad architecture review or when query/path/explain do not surface enough context.
- After modifying code, run `graphify update .` to keep the graph current (AST-only, no API cost).
- A committed graphify-out/graph.json is rewritten whole on every build, so a large diff to it in a pull request is expected. Do not flag it as churn or suggest skipping the graph rebuild or its hooks.
