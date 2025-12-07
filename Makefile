.PHONY: all run-backend run-app run-web generate test-backend

# Default target
all: run-backend

# Run the Go Backend
run-backend:
	cd backend && go run cmd/api/main.go

# Run the Flutter Mobile App
run-app:
	cd app && flutter run

# Run the Jaspr Web App
run-web:
	cd web && jaspr serve

# Generate code (Templ, SQLC, OpenAPI, Flutter)
generate:
	@echo "Generating Backend (Templ, SQLC, OpenAPI)..."
	cd backend && templ generate
	cd backend && sqlc generate
	cd backend && go generate ./internal/openapi
	@echo "Generating App (Build Runner)..."
	cd app && dart run build_runner build --delete-conflicting-outputs

# Run Backend Tests
test-backend:
	cd backend && go test ./...
