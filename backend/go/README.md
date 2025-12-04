# Physiolink — Go backend (scaffold)

This folder contains a minimal scaffold for the Go backend used during the refactor.

Quick start (local development):

1. Start Postgres and Redis with Docker Compose:

```powershell
cd backend/go
docker-compose up -d
```

2. Run the server locally:

```powershell
cd backend/go
go run ./cmd/api
```

Health endpoint: `http://localhost:8080/health`

Next steps:
- Add `sqlc` configuration and SQL query files.
- Implement services, handlers, and OpenAPI generation.
- Add migrations runner (e.g. `golang-migrate`).
