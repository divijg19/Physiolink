.PHONY: all run-backend run-app run-web build-web generate test-backend

# Default target
all: run-backend

# Run the Go Backend
run-backend:
	cd backend && go run ./cmd/api

# Run the Flutter Mobile App
run-app:
	cd app && flutter run

# Run the Jaspr Web App (dev server)
run-web:
	cd web && jaspr serve

# Build the Jaspr marketing site and copy to backend for embedding.
# `packages/` (Dart SDK sources) is excluded to avoid bloating the binary.
build-web:
	cd web && jaspr build
	cp -r web/build/jaspr/* backend/internal/server/jaspr/
	rm -rf backend/internal/server/jaspr/packages

# Generate code (Templ, SQLC, OpenAPI, Flutter, Web)
generate:
	@echo "Generating Backend (Templ, SQLC, OpenAPI)..."
	cd backend && templ generate
	cd backend && sqlc generate
	cd backend && go generate ./internal/openapi
	@echo "Generating Web (Jaspr)..."
	$(MAKE) build-web
	@echo "Generating App (Build Runner)..."
	cd app && dart run build_runner build --delete-conflicting-outputs

# Run Backend Tests
test-backend:
	cd backend && go test -race ./...
