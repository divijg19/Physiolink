# Physiolink — Go backend

Production backend implemented in Go. The prior Node.js code has been removed.

## Stack
- HTTP: Chi v5
- DB: Postgres (pgx v5), migrations in `migrations/`
- Jobs: Temporal (optional in dev)
- Views: Templ (SSR for public pages)
- Contracts: OpenAPI (`backend/openapi.yaml`)
- SQL codegen: sqlc (`sqlc.yaml`, queries in `internal/db/queries`)

## Quick start

1) Infra
```powershell
cd backend
docker compose up -d
```

2) API server
```powershell
cd backend
go run .\cmd\api
```

3) Worker (optional)
```powershell
cd backend
go run .\cmd\worker
```

Health: http://localhost:8080/health

## OpenAPI
Spec lives at `backend/openapi.yaml` and matches mobile clients (e.g., `_id` fields).

## sqlc — Type-safe Database Access

This project uses [sqlc](https://sqlc.dev/) to generate type-safe Go code from SQL queries.

**Adding new queries:**

1. Write your SQL query in `internal/db/queries/` (`.sql` files organized by domain)
2. Annotate with sqlc directives:
   ```sql
   -- name: GetUserByID :one
   SELECT id, email, role FROM users WHERE id = $1;
   ```
3. Generate Go code:
   ```powershell
   sqlc generate
   ```
4. Use the generated method in your service:
   ```go
   user, err := s.db.Queries.GetUserByID(ctx, userID)
   ```

**Query result types:**
- `:one` — single row (returns struct or error)
- `:many` — multiple rows (returns slice)
- `:exec` — no result (INSERT/UPDATE/DELETE)

**Transactions:**
For transactional operations, use `WithTx`:
```go
tx, _ := s.db.SQL.BeginTx(ctx, nil)
qtx := s.db.Queries.WithTx(tx)
// use qtx for queries within transaction
tx.Commit()
```

See existing queries in `internal/db/queries/` for examples.

## Testing
Run tests:
```powershell
cd backend
go test ./...
```

Current coverage focuses on middleware and basic handlers. More tests welcome.
