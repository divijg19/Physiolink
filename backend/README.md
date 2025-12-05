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

## sqlc
Generate DB access code:
```powershell
cd backend
sqlc generate
```

## Testing
Run tests:
```powershell
cd backend
go test ./...
```

Current coverage focuses on middleware and basic handlers. More tests welcome.
