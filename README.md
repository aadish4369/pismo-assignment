# Pismo — assignment

REST API for accounts and transactions (Go, Gin, GORM, SQLite). JSON amounts are in rupees (stored as paisa in the DB).

| API | Description |
|-----|-------------|
| `GET /ping` | Liveness check (`{"message":"pong"}`). |
| `POST /accounts` | Create account; body `document_number`. |
| `GET /accounts/:accountId` | Account balance and active installment plans. |
| `POST /transactions` | Create transaction: `operation_type_id` **1** purchase, **2** installment (requires `tenure`; first EMI + plan), **3** withdrawal, **4** credit. |
| `POST /accounts/:accountId/installments/:planId/next` | Debit next EMI on an installment plan. |

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

The `Dockerfile` is a **multi-stage** build: compile with CGO in a Go image, then run a slim `debian:bookworm-slim` image with `libsqlite3-0`. Inside the container, `DATABASE_PATH` defaults to `/data/app.db` and `GIN_MODE=release`.

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
