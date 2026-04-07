FROM golang:bookworm AS builder

WORKDIR /src

ENV CGO_ENABLED=1 \
    GOTOOLCHAIN=auto

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -trimpath -ldflags="-s -w" -o /server .

FROM debian:bookworm-slim

RUN apt-get update \
	&& apt-get install -y --no-install-recommends ca-certificates libsqlite3-0 \
	&& rm -rf /var/lib/apt/lists/*

WORKDIR /data

ENV DATABASE_PATH=/data/app.db

COPY --from=builder /server /usr/local/bin/server

EXPOSE 8080

ENTRYPOINT ["/usr/local/bin/server"]
