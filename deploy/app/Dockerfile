FROM golang:1.27 AS build-stage

# to run in alpine
ENV CGO_ENABLED=0

WORKDIR /app
COPY . .

# --- install dependencies ---
RUN go mod download
RUN go install github.com/go-task/task/v3/cmd/task@v3.54.0
# match with the version in internal/db/dbqueries/*.go headers
RUN go install github.com/sqlc-dev/sqlc/cmd/sqlc@v1.31.1

# version pinned as TAILWIND_VERSION in app/Taskfile.yml; these hashes come from that
# release's sha256sums.txt, update them together
ARG TAILWIND_SHA256_X64=e47695cee3f0d90623a91f3939b08bdee610eb5ad99340adff9d4145b284d6cc
ARG TAILWIND_SHA256_ARM64=ffcd0d21dbd7c605dc12743ff0c0b47debc963c86b399674084c1c8f2c96ecad
ARG TARGETARCH
RUN v="$(task --silent tailwind-version)" \
  && case "$TARGETARCH" in \
    amd64) a=x64 sum=$TAILWIND_SHA256_X64 ;; \
    arm64) a=arm64 sum=$TAILWIND_SHA256_ARM64 ;; \
    *) echo "unsupported arch $TARGETARCH"; exit 1 ;; esac \
  && echo "installing tailwindcss v$v ($a)" \
  && curl -fsSL -o /usr/local/bin/tailwindcss \
    "https://github.com/tailwindlabs/tailwindcss/releases/download/v$v/tailwindcss-linux-$a" \
  && echo "$sum  /usr/local/bin/tailwindcss" | sha256sum -c - \
  && chmod +x /usr/local/bin/tailwindcss \
  && tailwindcss --help | head -n 1

# --- build ---
RUN task build

FROM alpine:latest AS release-stage
ENV GIN_MODE=release
# no write access needed; the read-only certificates mount only has to be readable
RUN adduser -D -H app
WORKDIR /app

COPY --from=build-stage /app/web/assets web/assets
# output from `task build`
COPY --from=build-stage /app/bin/app bin/app
USER app
CMD ["bin/app"]
