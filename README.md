# FinScheduler

FinScheduler - personal finance management service.

The project is split into two applications:

- [`FinScheduler.API`](./FinScheduler.API/README.md) - Go backend API.
- [`finscheduler-web`](./finscheduler-web/README.md) - React frontend.

The goal is to collect financial data, keep it structured, analyze it, build reports, and eventually automate reminders and planning around personal finances.

## Current Capabilities

Implemented today:

- expense item management
- tag management
- item price history persistence

Planned next steps still include broader analytics, reports, reminders, and other finance automation flows.

## Tech Stack

Backend:

- Go 1.25
- PostgreSQL via pgx/sqlx
- OpenTelemetry
- testcontainers-go for integration tests

Frontend:

- React 19
- TypeScript
- Vite via `rolldown-vite`
- Chakra UI
- React Router
- Recharts
- pnpm

## Repository Structure

```text
.
+-- FinScheduler.API/   # Go API, migrations, tests
+-- finscheduler-web/   # React app
+-- k8s/                # Kubernetes manifests for the local test contour
+-- scripts/            # Helper scripts for build/deploy/destroy flows
```

## Quick Start

Backend:

```bash
cd FinScheduler.API
go run ./cmd/finscheduler
```

Frontend:

```bash
cd finscheduler-web
pnpm install
pnpm dev
```

By default the frontend expects the backend on `http://localhost:8081` unless `VITE_API_BASE_URL` is overridden.

## Useful Commands

Backend:

```bash
cd FinScheduler.API
go build ./cmd/finscheduler
go test ./...
go test -tags=integration ./tests/integration/...
```

Frontend:

```bash
cd finscheduler-web
pnpm build
pnpm test
pnpm lint
```

## API Overview

Health:

- `GET /livez`
- `GET /readyz`
- `GET /metrics` when metrics are enabled

Resources:

- `/api/items` - item CRUD and item details, including price history
- `/api/tags` - tag CRUD
- `/api/tags/lookup` - lightweight tag lookup for selectors

## Local Kubernetes Test Contour

The repository contains a local Kubernetes contour under [`k8s/base`](./k8s/base/) with separate layers for:

- storage
- observability
- application workloads
- ingress

Helper scripts for `k3d`-based local deployment live in [`scripts`](./scripts/README.md).

Typical flow:

```bash
./scripts/create-test-cluster.sh
./scripts/deploy-storage.sh
./scripts/deploy-observability.sh
./scripts/deploy-app.sh
```

Or in one step:

```bash
./scripts/deploy-test-contour.sh
```

## Documentation

Use the app-specific READMEs for detailed setup and operational notes:

- [`FinScheduler.API/README.md`](./FinScheduler.API/README.md)
- [`finscheduler-web/README.md`](./finscheduler-web/README.md)
- [`scripts/README.md`](./scripts/README.md)
