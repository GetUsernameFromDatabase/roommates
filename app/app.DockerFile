FROM golang:1.27 AS build-stage

# make sure this matches with the version from `tailwindcss --help`
# easy copy paste `tailwindcss --help | grep -oE 'v[0-9]+\.[0-9]+\.[0-9]+'`
ARG TAILWIND_VERSION=4.3.3
# to run in alpine
ENV CGO_ENABLED=0

WORKDIR /app
COPY . .

# install dependencies
RUN go mod download

RUN go install github.com/go-task/task/v3/cmd/task@latest
# make sure this matches with the version from
# `tailwindcss --help | head -n 1`
RUN curl -fsSL -o /usr/local/bin/tailwindcss \
  "https://github.com/tailwindlabs/tailwindcss/releases/download/v${TAILWIND_VERSION}/tailwindcss-linux-x64" \
  && chmod +x /usr/local/bin/tailwindcss

# build
RUN task build

FROM alpine:latest AS release-stage
WORKDIR /app

COPY --from=build-stage /app/assets assets
COPY --from=build-stage /app/certificates certificates
# output from `task build`
COPY --from=build-stage /app/bin/app bin/app
CMD ["bin/app"]