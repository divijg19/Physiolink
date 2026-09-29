.PHONY: all check lint vet vuln build test test-backend test-integration \
        generate build-web run-backend run-app run-web clean

# Default target: verify everything (matches the CI gate).
all: check

# Full backend gate, mirroring .github/workflows/backend-ci.yml.
check: lint build test vuln

lint:
	cd backend && go vet ./...
	cd backend && staticcheck ./...

build:
	cd backend && go build ./...

# Unit tests only. Integration tests need a live Postgres; see test-integration.
test-backend:
	cd backend && go test -race $$(go list ./... | grep -v '/integration$')

test: test-backend

# Requires DATABASE_URL to point at a Postgres with the migrations applied.
test-integration:
	cd backend && go test -race -v ./integration

vuln:
	cd backend && govulncheck ./...

# Build the Jaspr marketing site and copy it into the backend for embedding.
# packages/ (Dart SDK sources) and .dart_tool/ are excluded: go:embed skips
# dotfiles but not packages/, which would bloat the binary.
build-web:
	cd web && jaspr build
	cp -r web/build/jaspr/. backend/internal/server/jaspr/
	rm -rf backend/internal/server/jaspr/packages backend/internal/server/jaspr/.dart_tool

# Generate code. Commands must stay in sync with backend-codegen.yml, which
# fails CI when the committed output drifts.
generate:
	@echo "Generating Backend (Templ, SQLC, OpenAPI)..."
	cd backend && templ generate
	cd backend && sqlc generate
	cd backend && go generate ./internal/openapi/...
	@echo "Generating Web (Jaspr)..."
	$(MAKE) build-web
	@echo "Generating App (Build Runner)..."
	cd app && dart run build_runner build --delete-conflicting-outputs

# Run the Go Backend
run-backend:
	cd backend && go run ./cmd/api

# Run the Flutter Mobile App
run-app:
	cd app && flutter run

# Run the Jaspr Web App (dev server)
run-web:
	cd web && jaspr serve

clean:
	rm -rf backend/internal/server/jaspr/main.dart.js \
	       backend/internal/server/jaspr/styles.css \
	       backend/internal/server/jaspr/favicon.ico \
	       backend/internal/server/jaspr/index.html \
	       backend/internal/server/jaspr/icons \
	       backend/internal/server/jaspr/images \
	       backend/coverage.out backend/unit-tests.xml
