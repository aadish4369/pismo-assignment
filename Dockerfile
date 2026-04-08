FROM golang:bookworm AS builder

WORKDIR /app
RUN apt-get update && apt-get install -y --no-install-recommends gcc libc6-dev \
	&& rm -rf /var/lib/apt/lists/*

ENV CGO_ENABLED=1 GOTOOLCHAIN=auto

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -ldflags="-s -w" -o /server .

FROM debian:bookworm-slim

RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates libsqlite3-0 \
	&& rm -rf /var/lib/apt/lists/*

WORKDIR /app
ENV DATABASE_PATH=/data/app.db
ENV GIN_MODE=release

COPY --from=builder /server /server

EXPOSE 8080
CMD ["/server"]
