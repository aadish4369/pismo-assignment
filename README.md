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

```bash
docker compose up --build
```
