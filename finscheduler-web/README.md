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

Default:

```text
http://localhost:8081/api
```

Example override:

```bash
VITE_API_BASE_URL=http://localhost:8081/api pnpm dev
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

- The backend should be running on `localhost:8081` for the default API URL.
- Some existing UI labels appear to have encoding issues and should be normalized to UTF-8.
