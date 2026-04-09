# Pismo — assignment

REST API for accounts and transactions (Go, Gin, GORM, SQLite). JSON amounts are in rupees (stored as paisa in the DB).

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/ping` | Health check. |
| `POST` | `/accounts` | Create account (`document_number`). |
| `GET` | `/accounts/:accountId` | Get account and balance (rupees). |
| `POST` | `/transactions` | Post transaction (`account_id`, `operation_type_id`, `amount`). Types: **1** purchase, **2** installment, **3** withdrawal, **4** credit. |

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

```bash
docker compose up --build
```

```bash
docker compose down
```

Swagger in Docker: `http://localhost:8080/swagger/index.html`
