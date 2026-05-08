# FinScheduler API

Go backend for FinScheduler.

The API exposes personal finance resources, stores data in PostgreSQL, and runs SQL migrations automatically on startup.

## Tech Stack

- Go 1.24.
- PostgreSQL via pgx/sqlx
- OpenTelemetry.
- testcontainers-go for integration tests.

## Local Configuration

The app reads `configs/config.json` from the API project directory.

```json
{
  "env": "Local",
  "serverPort": 12345,
  "connectionString": "",
  "serviceName": "fin-scheduler-api",
  "corsSettings": {
    "allowedOrigins": ["*"],
    "allowedMethods": ["GET", "POST", "PUT", "DELETE", "OPTIONS"],
    "allowedHeaders": ["*"],
    "allowCredentials": false
  }
}
```

Viper also enables environment variables. Config keys can be overridden with uppercase names such as `SERVER_PORT`, `CONNECTION_STRING`, `SERVICE_NAME`, `CORS_ALLOWED_ORIGINS`, `CORS_ALLOWED_METHODS`, `CORS_ALLOWED_HEADERS`, and `CORS_ALLOW_CREDENTIALS`.

## Run

Start PostgreSQL with a database matching the configured connection string, then run:

```bash
go run ./cmd/finscheduler
```

The service listens on `http://localhost:8081` with the default config.

## Build

```bash
go build ./cmd/finscheduler
```

## Tests

```bash
go test ./...
```

Integration tests start PostgreSQL through Testcontainers, so Docker must be available.

## Migrations

Migration files live in `database/postgres`.

They are applied automatically during API startup through `database.RunMigrations`. Run the app from the `FinScheduler.API` directory so the `file://database/postgres` migration path resolves correctly.

## Routes

Health:

- `GET /livez`
- `GET /readyz`

Items:

- `GET /api/items`
- `POST /api/items`
- `PUT /api/items/{id}`
- `DELETE /api/items/{id}`

Tags:

- `GET /api/tags`
- `GET /api/tags/lookup`
- `POST /api/tags`
- `PUT /api/tags/{id}`

## Project Structure

```text
cmd/finscheduler/      # Application entry point
configs/               # Local config file
database/postgres/     # SQL migrations
internal/features/     # Domains, HTTP handlers, services, repositories
internal/health/       # Liveness and readiness handlers
internal/infra/        # Configuration
internal/metrics/      # Metrics setup and helpers
internal/persistence/  # Unit of work and DB factory
internal/traces/       # Tracing setup and helpers
pkg/                   # Small shared helpers
tests/                 # Integration tests
```
