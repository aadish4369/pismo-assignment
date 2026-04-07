# Documentation

## Architecture

- **HTTP:** Gin (`routes/router.go`) maps paths to handlers.
- **Use cases:** `services/routineService.go` orchestrates accounts, transactions, and optional `InstallmentPlan` rows.
- **Domain services:** `AccountService`, `TransactionService` encapsulate persistence and rules.
- **Persistence:** GORM repositories + SQLite (`db/db.go`, default `data/app.db`).

## Transaction rules

- Incoming `amount` in JSON is in **rupees** (float); storage uses **paisa** (int64) internally.
- `TransactionService` normalizes sign: operation types `1`, `2`, `3` → negative paisa; type `4` → positive paisa.
- Debits fail if the resulting balance would be negative (`insufficient balance`).

## Installment plan (optional, operation type `2`)

1. **`POST /transactions`** with `operation_type_id: 2` and positive `amount`:
   - Without `tenure`: one **purchase** debit only (same storage as type `1`).
   - With `tenure` > 1: same full debit **and** an **`InstallmentPlan`** linked by `transaction_id` (unique): `TotalPaisa`, `Tenure`, `EMIPaisa`, `LastEMIPaisa`, `PaidEMIs` = 0, `NextDueDate`.
2. **`POST /accounts/:accountId/installments/:planId/pay`** (account must own the plan):
   - If `PaidEMIs < Tenure`, posts a **credit voucher (type 4)** for this EMI amount (adds to balance).
   - Increments `PaidEMIs` and moves `NextDueDate` forward by one month.
3. **`GET /accounts/:accountId`** includes **`active_installment_plans`**: plans where `PaidEMIs < Tenure`.

## OpenAPI / Swagger

- Generated code lives under `docs/` (`docs.go`, `swagger.json`, `swagger.yaml`).
- UI is served at `/swagger/*` via `gin-swagger`.
- Regenerate when you change swag annotations or routes (see README).

## Configuration

| Env | Default | Meaning |
|-----|---------|---------|
| `PORT` | `8080` | HTTP port |
| `DATABASE_PATH` | `data/app.db` | SQLite file path |
