# Pismo — assignment

REST API for accounts and transactions (Go, Gin, GORM, SQLite). JSON amounts are in rupees (stored as paisa in the DB).

| API | Description |
|-----|-------------|
| `GET /ping` | Liveness check (`{"message":"pong"}`). |
| `POST /accounts` | Create account; body `document_number`. |
| `GET /accounts/:accountId` | Account balance (rupees). |
| `POST /transactions` | Create transaction: `operation_type_id` **1** purchase, **2** installment purchase (full amount debited like 1; scheduling TBD), **3** withdrawal, **4** credit. |

## Run

```bash
go run .
```

Defaults: port `8080` (`PORT`), database `data/app.db` (`DATABASE_PATH`).

## API docs (Swagger)

<http://localhost:8080/swagger/index.html>

## Tests

```bash
go test ./... -count=1
```

## Docker

The `Dockerfile` is a **multi-stage** build with **`CGO_ENABLED=0`**: SQLite uses [glebarez/sqlite](https://github.com/glebarez/sqlite) (pure Go, no `gcc` / `libsqlite3`). BuildKit **cache mounts** speed up repeated `go build` after `go.mod` changes. Runtime image is `debian:bookworm-slim` with `ca-certificates` only. `DATABASE_PATH` defaults to `/data/app.db`, `GIN_MODE=release`.

**Compose** (foreground, API on `http://localhost:8080`). Use **`--no-log-prefix`** so you only see your app’s log lines (no `api-1 |` prefix from Compose):

```bash
docker compose up --build
```

Detached:

```bash
docker compose up -d --build
docker compose logs -f --no-log-prefix api
```

Stop:

```bash
docker compose down
```

**Without Compose:**

```bash
docker build -t pismo-api .
docker run --rm -p 8080:8080 pismo-api
```

Persist SQLite on the host (`DATABASE_PATH` in the image is `/data/app.db`):

```bash
docker run --rm -p 8080:8080 -v "$(pwd)/data:/data" pismo-api
```

Swagger in Docker: `http://localhost:8080/swagger/index.html`
