# FinScheduler

FinScheduler - personal finance management service.

The project is split into two applications:

- [`FinScheduler.API`](./FinScheduler.API/README.md) - Go backend API.
- [`finscheduler-web`](./finscheduler-web/README.md) - React frontend.

The goal is to collect financial data, keep it structured, analyze it, build reports, and eventually automate reminders and planning around personal finances.

## Current Scope

Implemented or partially implemented:

- Expense item management.
- Tag management.

Planned ideas:

- Price history.
- Grocery planning.
- Reports and charts.
- Notifications.

## Tech Stack

Backend:

- Go 1.25.
- PostgreSQL via pgx/sqlx
- OpenTelemetry.
- testcontainers-go for integration tests.

Frontend:

- React 19.
- TypeScript.
- Vite / rolldown-vite.
- Chakra UI.
- React Router.
- lucide-react and react-icons.
- pnpm.

## Repository Structure

```text
.
+-- FinScheduler.API/       # Go API
|   +-- cmd/finscheduler/   # API entry point
|   +-- configs/            # Local JSON config
|   +-- database/postgres/  # SQL migrations
|   +-- internal/           # App features and infrastructure
|   +-- tests/              # Integration tests
+-- finscheduler-web/       # React frontend
    +-- src/
```

## Backend Setup

Run the backend:

```bash
cd FinScheduler.API
go run ./cmd/finscheduler
```

Build the backend:

```bash
cd FinScheduler.API
go build ./cmd/finscheduler
```

Run backend tests:

```bash
cd FinScheduler.API
go test ./...
```

See [`FinScheduler.API/README.md`](./FinScheduler.API/README.md) for config, routes, migrations, and test details.

## Observability

A local Grafana-based logging stack is available for the Go backend.

The repository now includes Kubernetes manifests for the current local contour. Apply everything from the repository root:

```bash
kubectl apply -k k8s/base
```

That base now includes:

- `FinScheduler.API`
- MinIO for shared object storage
- Prometheus for metrics scraping
- Mimir for long-term metrics storage
- Tempo for traces
- Pyroscope for continuous profiling
- Loki for logs
- Grafana Alloy for Kubernetes log collection
- Grafana with provisioned datasources

See [`k8s/base/storage`](./k8s/base/storage/README.md) and [`k8s/base/observability`](./k8s/base/observability/README.md) for the manifests that make up the contour.

For local test contour deployment, helper scripts live in [`scripts`](./scripts). The layers can be rolled out independently:

```bash
./scripts/deploy-storage.sh
./scripts/deploy-observability.sh
./scripts/deploy-app.sh
```

Or as one combined flow:

```bash
./scripts/deploy-test-contour.sh
```

## Frontend Setup

Install dependencies:

```bash
cd finscheduler-web
pnpm install
```

Run the frontend:

```bash
cd finscheduler-web
pnpm dev
```

Build the frontend:

```bash
cd finscheduler-web
pnpm build
```

Lint the frontend:

```bash
cd finscheduler-web
pnpm lint
```

See [`finscheduler-web/README.md`](./finscheduler-web/README.md) for scripts, environment variables, and app routes.

## Useful API Endpoints

- `GET /livez` - liveness check.
- `GET /readyz` - readiness check with database ping.
- `/api/items` - expense item operations.
- `/api/tags` - tag operations.

## Notes

- Migrations run automatically when the backend starts.
