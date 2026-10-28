# Research on how to better set it up

Project structure:

- <https://templ.guide/project-structure/project-structure>
- <https://github.com/golang-standards/project-layout>

Front-End:

- <https://github.com/larsartmann/templ-components>
- <https://dev.to/peterintech/go-templ-htmx-tailwind-css-v4-setup-in-2026-1ghj>
- <https://github.com/axadrn/shadcn-templ>

We are using [templ](https://templ.guide/) and [templui](https://templui.io/docs/how-to-use#cli-workflow)

install [tailwindcss](https://github.com/tailwindlabs/tailwindcss/releases/tag/v4.3.3)

Live-Reload (nice for dev):

- <https://templ.guide/developer-tools/live-reload/> -- will try this before air
- <https://github.com/air-verse/air>

DB should be good, use the same as before

## /scripts

Scripts to perform various build, install, analysis, etc operations.

These scripts keep the root level Makefile small and simple (e.g., <https://github.com/hashicorp/terraform/blob/main/Makefile>).

See the /scripts directory for examples.

```sh
#!/bin/bash
# for generating code
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest # DB
go install github.com/a-h/templ/cmd/templ@latest # templ
go install github.com/axadrn/shadcn-templ/v2/cmd/shadcn-templ@latest # templui components
# go install github.com/swaggo/swag/cmd/swag@latest # swagger
```

## authentication

Thinking session bassed
`sessions/postgres`

<https://gin-gonic.com/en/docs/middleware/session-management/#cookie-based-sessions>
