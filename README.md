# Yiqi Docs

Open-source documentation platform for multi-site Markdown docs.

## Stack

- Backend: Go, Gin, GORM, SQLite
- Frontend: Vue 3, Vite, Pinia
- Deployment: single Go binary with embedded frontend assets

## Development

Start the backend:

```bash
cd backend
JWT_SECRET=change-me SEED_ADMIN_USER=admin SEED_ADMIN_PASS=change-me go run ./cmd/server
```

Start the frontend:

```bash
cd frontend
pnpm install
pnpm dev
```

The Vite dev server proxies `/api` to `http://localhost:8080`.

## Build

```bash
./build.sh
```

This builds the frontend, copies `frontend/dist` into `backend/web/dist`, then builds the server binary.

## Runtime Config

Environment variables:

- `PORT`: HTTP port, default `8080`
- `DB_PATH`: SQLite database path, default `./data.db`
- `JWT_SECRET`: required JWT signing secret
- `SEED_ADMIN_USER`: initial super admin username, default `admin`
- `SEED_ADMIN_PASS`: initial super admin password
- `APIFOX_SHARED_DOC_TOKEN`: optional Apifox shared-doc token for the scraper. If omitted, the scraper falls back to public page scraping.

Do not commit runtime databases, logs, certificates, or secrets.
