# Pismo — assignment

REST API for accounts and transactions (Go, Gin, GORM, SQLite). JSON amounts are in rupees (stored as paisa in the DB).

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/ping` | Liveness: `{"message":"pong"}`. |
| `POST` | `/accounts` | Create account. JSON body: `document_number`. Response: `account_id`, `document_number`. |
| `GET` | `/accounts/:accountId` | Load account. Response: `account_id`, `document_number`, `balance` (rupees). |
| `POST` | `/transactions` | Create transaction. JSON body: `account_id`, `operation_type_id`, `amount` (rupees). `operation_type_id`: **1** purchase, **2** installment purchase (full amount debited like 1; scheduling TBD), **3** withdrawal, **4** credit. Response includes `transaction_id`, `amount` (rupees; sign reflects debit/credit in JSON). |

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
