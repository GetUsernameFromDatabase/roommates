# Tools

go programs run at dev and build time, never shipped (those live in `cmd/`)

- `generate-locales-keys`, `generate-locales-languages`: generate go code from the locale yaml files (`go generate ./tools/...`)
- `vendor-assets`: downloads htmx and the fonts pinned in `web/assets/vendor.json` (`task vendor`)

currently I want to keep each file as self-contained as possible
so I will keep duplicate code with these files
