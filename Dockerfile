# syntax=docker/dockerfile:1

FROM golang:bookworm AS builder

WORKDIR /app
ENV CGO_ENABLED=0 GOTOOLCHAIN=auto

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
	go mod download

COPY . .
RUN --mount=type=cache,target=/go/pkg/mod \
	--mount=type=cache,target=/root/.cache/go-build \
	go build -trimpath -ldflags="-s -w" -o /server .

FROM debian:bookworm-slim

RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates \
	&& rm -rf /var/lib/apt/lists/*

WORKDIR /app
ENV DATABASE_PATH=/data/app.db
ENV GIN_MODE=release

COPY --from=builder /server /server

EXPOSE 8080
CMD ["/server"]
