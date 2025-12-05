#!/usr/bin/env pwsh
Set-StrictMode -Version Latest

Write-Host "Ensure oapi-codegen is installed: go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest"

Push-Location (Join-Path $PSScriptRoot "..\internal\openapi")
try {
    Write-Host "Running go generate for OpenAPI codegen..."
    go generate
} finally {
    Pop-Location
}
