# PhysioLink

PhysiLink is a physiotherapy platform: a Go backend serving a JSON API, a
server-rendered web portal, and an embedded marketing site, plus a Flutter
mobile app.

## Project Structure

| Path      | Contents                                                                   |
| --------- | -------------------------------------------------------------------------- |
| `backend/`| Go API (chi), PostgreSQL via pgx, and the Templ/HTMX portal.                  |
| `app/`    | Flutter mobile client (iOS/Android/web).                                     |
| `web/`    | Jaspr marketing SPA, embedded into the Go binary and served under `/site`.   |

## Prerequisites

- Go 1.26.x
- Flutter SDK (stable)
- Dart SDK (comes with Flutter)
- Podman + `podman compose` (or Docker)
- `make`

Code generation additionally needs these on `PATH`, or run `make doctor`:

```
go install github.com/a-h/templ/cmd/templ@latest
go install github.com/sqlc-dev/sqlc/cmd/sqlc@v1.31.1
go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@v2.8.0
dart pub global activate jaspr_cli 0.23.5
```

`jaspr` and `templ` must be on `PATH`; the jaspr CLI lives in
`~/.pub-cache/bin`, which is not always on `PATH`.

## Getting started

```bash
cp .env.example .env      # then set JWT_SECRET
podman compose up -d      # Postgres, Redis, Temporal, Temporal UI
make db-migrate           # apply migrations (see note below)
```

Migrations are **not** applied automatically by the container. Postgres only
runs `/docker-entrypoint-initdb.d` on first initialisation of an empty volume,
so anything added later would be silently ignored. Use `make db-migrate`, which
applies each file with the same `ON_ERROR_STOP` + single-transaction semantics
as CI.

### Running the app

```bash
make run-backend     # http://localhost:8080
make run-app         # Flutter app
```

The marketing SPA is embedded at build time, so it must be built before the
backend serves it:

```bash
make build-web       # builds web/ and copies it into the embed directory
```

Without that, `/site` returns 404 and the backend still builds - the embed
directory just contains a `placeholder.txt` sentinel. The Temporal worker is
run separately with `cd backend && go run ./cmd/worker`; it is not part of the
compose stack.

## Development

```bash
make check           # vet, staticcheck, build, tests, govulncheck
make generate        # templ, sqlc, openapi, the SPA, and app codegen
make clean           # remove generated build output
```

`make check` mirrors the CI gate. `make generate` must stay in sync with
`.github/workflows/backend-codegen.yml`, which fails CI when committed
generated code drifts.

### Layout note

The portal owns `/`; the marketing SPA is mounted at `/site` and `/` redirects
to it. The two front ends previously collided - the SPA was the router's
catch-all and shadowed the portal's own routes.

## CI

| Workflow                    | Covers                                              |
| --------------------------- | --------------------------------------------------- |
| `backend-ci`                | vet, staticcheck, govulncheck, build, unit tests, embedded SPA, Docker build |
| `backend-codegen`           | sqlc and OpenAPI generated-code drift               |
| `backend-integration`       | integration tests against a real Postgres            |
| `web-ci`                    | `dart analyze` and the Jaspr build                   |
| `app-ci`                    | codegen drift, `flutter analyze`, tests, web build  |
