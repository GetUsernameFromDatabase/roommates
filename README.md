# Roommates

## Setup

Install golang
> <https://go.dev/doc/install>

Install binaries with go

- `go install github.com/a-h/templ/cmd/templ@latest`
  - [templ](https://templ.guide/quick-start/installation#go-install-global)
  - used for server-side rendering
- `go install github.com/templui/templui/cmd/templui@latest`
  - [shadcn-templ](https://shadcn-templ.com/docs/installation#existing-project)
  - ui components for templ
- `go install github.com/axadrn/shadcn-templ/v2/cmd/shadcn-templ@latest`
  - [sqlc](https://docs.sqlc.dev/en/latest/overview/install.html)
  - generate go code to interact with database
- `go install github.com/go-task/task/v3/cmd/task@latest`
  - [taskfile](https://taskfile.dev/docs/installation#build-from-source)
  - similar to makefile

Install [tailwindcss](https://github.com/tailwindlabs/tailwindcss/releases/latest) so that it is usable with `tailwindcss` command
> used in [Taskfile.yml](./app/configs/Taskfile.yml)
> required by <https://templui.io/docs/how-to-use#tailwind-css-v41>

Install [golangci-lint](https://golangci-lint.run/docs/welcome/install/local/)

```sh
staticcheck -checks "all,-ST1000" ./...
golangci-lint run
```
