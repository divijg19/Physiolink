# OpenAPI Code Generation (Go)

This document describes how to generate Go server types and handlers from the canonical OpenAPI spec located at `backend/go/openapi.yaml`.

Prerequisites
- Go 1.20+
- Install `oapi-codegen`:

```powershell
go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest
```

Generate

From the repository root (PowerShell):

```powershell
cd backend/go
pwsh.exe ./scripts/generate.ps1
```

This script invokes `go generate` in `backend/go/internal/openapi`, which runs the `oapi-codegen` command declared in the `generate.go` file.

Notes
- The generation will produce `openapi.gen.go` in `backend/go/internal/openapi` which contains typed models and server interfaces.
- After generation, implement the server interface in `backend/go/internal/handlers` and wire routes in `internal/server`.
