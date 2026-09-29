.PHONY: all check lint vet vuln build test test-backend test-integration db-migrate doctor \
        generate build-web run-backend run-app run-web clean

# Default target: verify everything (matches the CI gate).
all: check

# Report which toolchain binaries are missing. Several are needed by
# `make generate` and are not installed by any package manager by default.
doctor:
	@fail=0; \
	for tool in go dart flutter podman git; do \
		if command -v $$tool >/dev/null 2>&1; then \
			printf '  ok      %s\n' "$$tool"; \
		else \
			printf '  MISSING %s\n' "$$tool"; fail=1; \
		fi; \
	done; \
	for tool in templ sqlc oapi-codegen staticcheck gotestsum govulncheck; do \
		if command -v $$tool >/dev/null 2>&1; then \
			printf '  ok      %s\n' "$$tool"; \
		else \
			printf '  MISSING %s (needed by make generate / make check)\n' "$$tool"; fail=1; \
		fi; \
	done; \
	if [ -x "$$HOME/.pub-cache/bin/jaspr" ]; then \
		printf '  ok      jaspr (add ~/.pub-cache/bin to PATH)\n'; \
	else \
		printf '  MISSING jaspr (dart pub global activate jaspr_cli 0.23.5)\n'; fail=1; \
	fi; \
	exit $$fail

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

# Apply every migration in order with the same strict semantics CI uses.
# Each file runs in its own transaction and aborts on the first error, so a
# broken migration is loud instead of half-applied.
db-migrate:
	@for f in backend/migrations/*.sql; do \
		echo "applying $$f"; \
		psql -v ON_ERROR_STOP=1 --single-transaction "$$DATABASE_URL" -f "$$f" || exit 1; \
	done

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
	       backend/internal/server/jaspr/packages \
	       backend/internal/server/jaspr/.dart_tool \
	       backend/coverage.out backend/unit-tests.xml
