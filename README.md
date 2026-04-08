# Pismo — assignment

REST API for accounts and transactions (Go, Gin, GORM, SQLite).

## What’s included

- Create / fetch accounts (`document_number`)
- Post transactions by operation type (amounts in rupees; stored in paisa)
- Optional installment flow: type `2` with `tenure` debits the first EMI and creates a plan; `POST /accounts/:accountId/installments/:planId/next` records further EMIs

Operation types: `1` purchase, `2` installment purchase, `3` withdrawal, `4` credit.

## Run

```bash
go run .
```

Defaults: port `8080` (`PORT`), DB `data/app.db` (`DATABASE_PATH`).

## API docs (Swagger)

<http://localhost:8080/swagger/index.html>

Regenerate OpenAPI after changing handlers:

```bash
go run github.com/swaggo/swag/cmd/swag@v1.16.6 init -g main.go -o docs --parseInternal
```

## Tests

```bash
go test ./... -count=1
```

## Docker

```bash
docker compose up --build
```
