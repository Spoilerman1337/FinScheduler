# FinScheduler API

Go backend for FinScheduler.

The API exposes personal finance resources, stores data in PostgreSQL, and runs SQL migrations automatically on startup.

## Tech Stack

- Go 1.25.
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
  "observability": {
    "serviceName": "fin-scheduler-api",
    "metrics": {
      "enabled": true,
      "exportEndpoint": "/metrics"
    },
    "traces": {
      "enabled": false,
      "exportEndpoint": "http://localhost:4318",
      "rootTraceSamplingRatio": 1
    },
    "profiling": {
      "enabled": false,
      "pushURL": "http://localhost:4040"
    }
  },
  "corsSettings": {
    "allowedOrigins": ["*"],
    "allowedMethods": ["GET", "POST", "PUT", "DELETE", "OPTIONS"],
    "allowedHeaders": ["*"],
    "allowCredentials": false
  }
}
```

Viper also enables environment variables. Config keys can be overridden with uppercase names such as `SERVER_PORT`, `CONNECTION_STRING`, `OBSERVABILITY_SERVICE_NAME`, `METRICS_ENABLED`, `METRICS_EXPORT_ENDPOINT`, `TRACES_ENABLED`, `TRACES_EXPORT_ENDPOINT`, `TRACES_ROOT_TRACE_SAMPLING_RATIO`, `PROFILING_ENABLED`, `PROFILING_PUSH_URL`, `CORS_ALLOWED_ORIGINS`, `CORS_ALLOWED_METHODS`, `CORS_ALLOWED_HEADERS`, and `CORS_ALLOW_CREDENTIALS`.

## Run

Start PostgreSQL with a database matching the configured connection string, then run:

```bash
go run ./cmd/finscheduler
```

The service listens on `http://localhost:8081` with the default config.

When metrics are enabled, Prometheus-compatible metrics are exposed on `http://localhost:8081/metrics`.

## Build

```bash
go build ./cmd/finscheduler
```

## Tests

```bash
go test ./...
```

Integration tests start PostgreSQL through Testcontainers, so Docker must be available.

## Observability

The backend now supports:

- Prometheus metrics on `/metrics`
- OTLP/HTTP trace export to Tempo
- Pyroscope profiling
- JSON logs with `trace_id` and `span_id`

For Kubernetes deployments, the repository now includes Grafana, Prometheus, Tempo, and Pyroscope manifests under `k8s/base/observability`.

Apply them together with the existing base manifests from the repository root:

```bash
kubectl apply -k k8s/base
```

Grafana is available through ingress on `grafana.finscheduler.local` or via port-forward:

```bash
kubectl port-forward svc/grafana 3000:3000
```

See [`../k8s/base/observability/README.md`](../k8s/base/observability/README.md) for details.

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
