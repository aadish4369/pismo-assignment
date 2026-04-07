# Pismo assignment — transactions API

REST API for **accounts** and **transactions** (signed amounts by operation type). SQLite + Gin + GORM.

**Core behavior:** types `1`–`3` are stored as negative amounts, type `4` as positive. Amounts are sent in rupees; storage uses **paisa**.

**Optional installment metadata (type `2`):** omit **`tenure`** for a lump debit only (same as type `1`, different label). With **`tenure`** > 1, the purchase is still one **full debit**, plus an **`InstallmentPlan`** row for EMI split/progress. **`POST /accounts/:accountId/installments/:planId/pay`** posts a **credit voucher (type 4)** for one EMI. **`GET /accounts/:id`** lists incomplete plans in `active_installment_plans`.

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
| `GET` | `/accounts/:accountId` | Get account, balance, active installment plans |
| `POST` | `/transactions` | Create transaction (`tenure` optional for type `2`) |
| `POST` | `/accounts/:accountId/installments/:planId/pay` | Pay one EMI (credit voucher) |

**Operation types:** `1` normal purchase, `2` purchase with installments, `3` withdrawal, `4` credit voucher.

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
