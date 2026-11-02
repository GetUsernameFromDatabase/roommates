<!-- omit in toc -->
# Roommates

- [Project structure](#project-structure)
- [Quick start](#quick-start)
- [Setup](#setup)
- [Common commands](#common-commands)

## Project structure

```txt
roommates/
├─ docker-compose.yml   db, rustfs, mailpit and app services
├─ deploy/              PostgreSQL, RustFS, Mailpit and app images
└─ app/                 Go module (run go/task commands from here)
   ├─ cmd/server/       entry point (main.go) and routes
   ├─ internal/         application code
   │  ├─ config/        env, routes, cookie/header keys
   │  ├─ db/            migrations, queries/, sqlc output (dbqueries/)
   │  ├─ handlers/      HTTP layer
   │  ├─ services/      business logic used by handlers
   │  ├─ storage/       file contents in an S3 API (RustFS)
   │  ├─ mailer/        SMTP email (Mailpit in development)
   │  ├─ push/          web push notifications
   │  ├─ events/ images/ middleware/ models/ validators/ sessions/ logger/ utils/
   │  └─ locales/       translations (en, et) and generated keys
   ├─ web/
   │  ├─ components/    shared templ components (shadcn-templ ui)
   │  ├─ pages/         full pages (home, chat, expenses, files, notes, shopping, …)
   │  ├─ renderer/      gin <-> templ renderer
   │  └─ assets/        css, favicon, pwa, public/ (embedded static files)
   ├─ tools/            dev and build time programs, see tools/README.md
   ├─ scripts/          dev, format and lint scripts
   └─ certificates/     local mkcert TLS cert/key (not committed)
```

## Quick start

After [Setup](#setup). Run from `app/` unless noted.

1. From the repo root: `cp .env.example .env`, then replace the `REPLACE_WITH_*` and `` `openssl …` `` values
2. Create the TLS certificate: [app/certificates/README.md](./app/certificates/README.md)
3. `task vendor` downloads htmx and the fonts (not committed, `task dev` does not fetch them), see [app/web/assets/README.md](./app/web/assets/README.md)
4. `docker compose up -d mailpit` gives an inbox at <http://127.0.0.1:8025>. Without it, comment out `SMTP_HOST` in `.env` so emails go to the log instead of failing
5. `task dev` starts db + rustfs and serves the app with hot reload on <https://localhost:8090>. The database is migrated on start
6. `task seed` adds example users and a household

Whole stack in Docker instead: `docker compose up --build` from the repo root, app on <https://127.0.0.1:8000>.

## Setup

Install golang
> <https://go.dev/doc/install>

Install [Docker](https://docs.docker.com/get-started/get-docker/) with compose (PostgreSQL, RustFS and Mailpit run in it)

Install [mkcert](https://github.com/FiloSottile/mkcert#installation) for the local certificate

Install binaries with go

- `go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest`
  - [sqlc](https://docs.sqlc.dev/en/latest/overview/install.html)
  - generate go code to interact with database
- `go install github.com/go-task/task/v3/cmd/task@latest`
  - [taskfile](https://taskfile.dev/docs/installation#build-from-source)
  - similar to makefile
- `go install github.com/go-delve/delve/cmd/dlv@latest`
  - [delve](https://github.com/go-delve/delve/tree/master/Documentation/installation)
  - `task dev` runs the server under it (debugger on 127.0.0.1:2345)
- `go install honnef.co/go/tools/cmd/staticcheck@latest`
  - [staticcheck](https://staticcheck.dev/docs/getting-started/), used by `scripts/lint.sh`

[templ](https://templ.guide) and [shadcn-templ](https://shadcn-templ.com) need no install: they are pinned as `tool` in [go.mod](./app/go.mod) and run with `go tool templ` / `go tool shadcn-templ`

Install [tailwindcss](https://github.com/tailwindlabs/tailwindcss/releases) at the version pinned as `TAILWIND_VERSION` in [Taskfile.yml](./app/Taskfile.yml) so that it is usable with `tailwindcss` command (`task tailwind` and `task build` refuse other versions)
> used in [Taskfile.yml](./app/Taskfile.yml)
> required by <https://shadcn-templ.com/docs/installation#existing-project>

Install [golangci-lint](https://golangci-lint.run/docs/welcome/install/local/)

## Common commands

From `app/`:

- `task build` builds for production → `bin/app`
- `sh scripts/lint.sh` runs staticcheck + golangci-lint
- `sh scripts/format.sh` runs `templ fmt` + `gofmt`
- `go test ./...` skips DB tests unless `TEST_POSTGRES_DB` names a throwaway database (the migration test drops all tables), never the dev database

Code generation (regenerate after editing the source, never edit the output):

- `sqlc generate`: `internal/db/queries/*.sql` → `internal/db/dbqueries/`
- `go tool templ generate`: `*.templ` → `*_templ.go`
- `go generate ./tools/...`: `internal/locales/` yaml → locale keys
