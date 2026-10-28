# Assets

Installation sites:

- [HTMX](https://htmx.org/docs/#installing)

## Third-party files

htmx and the fonts are not committed. `task vendor` downloads them from
jsdelivr into `public/lib/` and `public/fonts/`. `tw-animate.css` and
`shadcn-tailwind.css` are written by `shadcn-templ add`/`init` instead.

- `vendor.json` pins each npm package and maps its files to paths here
- `vendor.lock.json` holds their sha256, a download that does not match fails
- `go run ./tools/vendor-assets --check` lists packages with a newer version
- `go run ./tools/vendor-assets --update [pkg...]` moves them to the latest,
  editing a version in `vendor.json` and running `task vendor` works too
