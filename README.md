# Pismo assignment — transactions API

REST API for **accounts**, **transactions** (signed amounts by operation type), and an optional **installment EMI** extension. SQLite + Gin + GORM.

## Run locally

```bash
go run .
# or
./run
```

Server listens on `:8080` (override with `PORT`). Database file defaults to `data/app.db` (override with `DATABASE_PATH`).

## API

| Method | Path | Purpose |
|--------|------|---------|
| `POST` | `/accounts` | Create account (`document_number`) |
| `GET` | `/accounts/:accountId` | Get account |
| `POST` | `/transactions` | Create transaction |
| `POST` | `/installments/:id/pay` | Pay one EMI (extension) |

**Operation types:** `1` normal purchase, `2` purchase with installments, `3` withdrawal, `4` credit voucher. Types `1–3` are stored as **negative** amounts; type `4` as **positive**. For type `2`, send `tenure` (>1) and `start_date` (`YYYY-MM-DD`) to create an installment plan.

## Swagger UI

With the app running: [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)

Regenerate OpenAPI after handler changes:

```bash
go run github.com/swaggo/swag/cmd/swag@v1.16.6 init -g main.go -o docs --parseInternal
```

## Docker

```bash
docker compose up --build
```

## Tests

```bash
go test ./... -count=1
```

## Docs

See [docs/DOCUMENTATION.md](docs/DOCUMENTATION.md) for architecture and design notes.

## Publish to GitHub

```bash
git init
git add .
git commit -m "Initial commit: transactions API with Swagger"
git branch -M main
git remote add origin https://github.com/<your-user>/<your-repo>.git
git push -u origin main
```

Create an empty repository on GitHub first, then replace the URL above.
