# PhysioLink

PhysioLink is a comprehensive physiotherapy management platform featuring a Go backend, a Flutter mobile app, and a Jaspr web portal.

## Project Structure

- **`backend/`**: Go API, Database (PostgreSQL), and Admin Portal (Templ/HTMX).
- **`app/`**: Flutter Mobile Application (iOS/Android).
- **`web/`**: Public Landing Page built with Jaspr (Dart for Web).

## Getting Started

### Prerequisites

- Go 1.23+
- Flutter SDK
- Docker & Docker Compose
- Make (optional, for using the Makefile)

### Running the Project

You can use the provided `Makefile` to run different parts of the application.

**Run the Backend:**
```bash
make run-backend
```
*Runs on http://localhost:8080*

**Run the Mobile App:**
```bash
make run-app
```

**Run the Web Portal:**
```bash
make run-web
```
*Runs on http://localhost:8081 (default Jaspr port)*

## Development

### Code Generation

To regenerate code for Templ, SQLC, OpenAPI, and Flutter/Riverpod:

```bash
make generate
```

### Testing

Run backend tests:

```bash
make test-backend
```
