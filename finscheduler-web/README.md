# FinScheduler Web

React frontend for FinScheduler.

The app provides UI pages for the dashboard, expense items, and tags. It talks to the backend through the API client in `src/api`.

## Tech Stack

- React 19.
- TypeScript.
- Vite / rolldown-vite.
- Chakra UI.
- React Router.
- lucide-react and react-icons.
- pnpm.

## Scripts

Install dependencies:

```bash
pnpm install
```

Run the dev server:

```bash
pnpm dev
```

Build:

```bash
pnpm build
```

Build the Docker image:

```bash
docker build -f finscheduler-web/Dockerfile -t finscheduler-web .
```

Run the Docker image:

```bash
docker run --rm -p 8080:80 finscheduler-web
```

Lint:

```bash
pnpm lint
```

Preview a production build:

```bash
pnpm preview
```

## Environment

The frontend uses `VITE_API_BASE_URL` to choose the backend API base URL.

Example override:

```bash
VITE_API_BASE_URL=http://localhost:8081/api pnpm dev
```

The Docker build can still override `VITE_API_BASE_URL`, and the Nginx runtime proxies `/api/*` to `API_PROXY_PASS`.

Example Docker overrides:

```bash
docker build -f finscheduler-web/Dockerfile --build-arg VITE_API_BASE_URL=/api -t finscheduler-web .
docker run --rm -p 8080:80 -e API_PROXY_PASS=http://host.docker.internal:8081/api finscheduler-web
```

## App Routes

- `/` - dashboard.
- `/items` - expense items.
- `/tags` - tags.

## API Client

The base client is `src/api/finscheduler-api-client.ts`.

Resource clients:

- `src/api/items.ts`
- `src/api/tags.ts`

## Project Structure

```text
src/api/        # API clients and shared DTO types
src/components/ # Shared UI components
src/features/   # Feature pages and local subcomponents
src/layout/     # App shell, sidebar, routes
```

## Notes

- In development and tests the backend should be running on `localhost:8081` unless `VITE_API_BASE_URL` is overridden.
- Some existing UI labels appear to have encoding issues and should be normalized to UTF-8.
