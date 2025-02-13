# syntax=docker/dockerfile:1.3
ARG GO_VERSION

FROM golang:${GO_VERSION}-alpine AS base
RUN apk add --no-cache gcc musl-dev
WORKDIR /src

FROM golangci/golangci-lint:v1.37-alpine AS golangci-lint

FROM base AS lint
RUN --mount=type=bind,target=. \
  --mount=type=cache,target=/root/.cache/go-build \
  --mount=type=cache,target=/root/.cache/golangci-lint \
  --mount=from=golangci-lint,source=/usr/bin/golangci-lint,target=/usr/bin/golangci-lint \
  golangci-lint run --timeout 10m0s ./...
