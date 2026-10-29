FROM golang:1.27 AS build-stage

# to run in alpine
ENV CGO_ENABLED=0

WORKDIR /app
COPY . .

# --- install dependencies ---
RUN go mod download
RUN go install github.com/go-task/task/v3/cmd/task@latest
# match with the version in internal/db/dbqueries/*.go headers
RUN go install github.com/sqlc-dev/sqlc/cmd/sqlc@v1.31.1

# version pinned as TAILWIND_VERSION in app/Taskfile.yml
RUN v="$(task --silent tailwind-version)" \
  && echo "installing tailwindcss v$v" \
  && curl -fsSL -o /usr/local/bin/tailwindcss \
    "https://github.com/tailwindlabs/tailwindcss/releases/download/v$v/tailwindcss-linux-x64" \
  && chmod +x /usr/local/bin/tailwindcss \
  && tailwindcss --help | head -n 1

# --- build ---
RUN task build

FROM alpine:latest AS release-stage
WORKDIR /app

COPY --from=build-stage /app/web/assets web/assets
# output from `task build`
COPY --from=build-stage /app/bin/app bin/app
CMD ["bin/app"]
