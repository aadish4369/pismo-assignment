# Documentation

## Architecture

- **HTTP:** Gin (`routes/router.go`) maps paths to handlers.
- **Use cases:** `services/routineService.go` orchestrates account lookup, transaction creation, and optional installment plan creation / EMI payment.
- **Domain services:** `AccountService`, `TransactionService`, `InstallmentService` encapsulate persistence and rules.
- **Persistence:** GORM repositories + SQLite (`db/db.go`, default `data/app.db`).

## Transaction rules

- Incoming `amount` in JSON is in **rupees** (float); storage uses **paisa** (int64) internally.
- `TransactionService` normalizes sign: operation types `1`, `2`, `3` → negative paisa; type `4` → positive paisa.
- Each transaction row includes `event_date` (UTC at creation).

## Installment extension (operation type `2`)

1. `POST /transactions` with `operation_type_id: 2`, positive `amount`, `tenure` > 1, and `start_date` creates:
   - A **purchase** transaction (negative amount).
   - An **installment** row with EMI split (last EMI absorbs remainder).
2. `POST /installments/:id/pay` applies one EMI: updates remaining EMIs / next due date and records a **withdrawal**-style transaction for that EMI amount.

## OpenAPI / Swagger

- Generated code lives under `docs/` (`docs.go`, `swagger.json`, `swagger.yaml`).
- UI is served at `/swagger/*` via `gin-swagger`.
- Regenerate when you change swag annotations or routes (see README).

## Configuration

| Env | Default | Meaning |
|-----|---------|---------|
| `PORT` | `8080` | HTTP port |
| `DATABASE_PATH` | `data/app.db` | SQLite file path |
