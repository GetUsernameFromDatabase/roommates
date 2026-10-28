<!-- omit in toc -->
# Roommates

- [Project structure](#project-structure)
- [Setup](#setup)

## Project structure

```txt
roommates/
├─ docker-compose.yml   db + app services
├─ deploy/postgres/     PostgreSQL image and config
└─ app/                 Go module (run go/task commands from here)
   ├─ cmd/server/       entry point (main.go) and routes
   ├─ internal/         application code
   │  ├─ config/        env, routes, cookie/header keys
   │  ├─ db/            migrations, queries/, sqlc output (dbqueries/)
   │  ├─ handlers/      HTTP layer
   │  ├─ services/      business logic used by handlers
   │  ├─ middleware/ models/ validators/ sessions/ logger/ utils/
   │  └─ locales/       translations (en, et) and generated keys
   ├─ web/
   │  ├─ components/    shared templ components (shadcn-templ ui)
   │  ├─ pages/         full pages (login, register)
   │  ├─ renderer/      gin <-> templ renderer
   │  └─ assets/        css, favicon, public/ (embedded static files)
   ├─ tools/            go:generate code generators
   ├─ scripts/          dev, format and lint scripts
   └─ certificates/     local TLS cert/key (not committed)
```

## Setup

Install golang
> <https://go.dev/doc/install>

Install binaries with go

- `go install github.com/a-h/templ/cmd/templ@latest`
  - [templ](https://templ.guide/quick-start/installation#go-install-global)
  - used for server-side rendering
- `go install github.com/axadrn/shadcn-templ/v2/cmd/shadcn-templ@latest`
  - [shadcn-templ](https://shadcn-templ.com/docs/installation#existing-project)
  - ui components for templ
- `go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest`
  - [sqlc](https://docs.sqlc.dev/en/latest/overview/install.html)
  - generate go code to interact with database
- `go install github.com/go-task/task/v3/cmd/task@latest`
  - [taskfile](https://taskfile.dev/docs/installation#build-from-source)
  - similar to makefile

Install [tailwindcss](https://github.com/tailwindlabs/tailwindcss/releases/latest) so that it is usable with `tailwindcss` command
> used in [Taskfile.yml](./app/Taskfile.yml)
> required by <https://templui.io/docs/how-to-use#tailwind-css-v41>

Install [golangci-lint](https://golangci-lint.run/docs/welcome/install/local/)

```sh
staticcheck -checks "all,-ST1000" ./...
golangci-lint run
```
