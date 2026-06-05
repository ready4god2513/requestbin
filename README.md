# RequestBin

Inspect and debug HTTP requests in real-time. Create a bin, send requests to its unique URL, and watch them appear instantly in the browser.

## Stack

- **Backend** — Go 1.23 + [Echo](https://echo.labstack.com/) + [pgx](https://github.com/jackc/pgx)
- **Frontend** — Vue 3 (Composition API) + Vite + Vue Router
- **Database** — PostgreSQL 16
- **Realtime** — Server-Sent Events (SSE)

## Getting started

### Development (hot reload)

```bash
docker compose -f docker-compose.dev.yml up --build
```

| Service  | URL                        |
|----------|----------------------------|
| Frontend | http://localhost:5173       |
| Backend  | http://localhost:8080       |
| Postgres | localhost:5432              |

The backend uses [air](https://github.com/air-verse/air) — any `.go` file change triggers an automatic rebuild. The frontend uses Vite HMR.

> **First run:** the backend container runs `go mod tidy` on startup to sync modules. This adds a few seconds the first time.

### Production

```bash
docker compose up --build
```

The app is served at http://localhost on port 80. The frontend (nginx) proxies `/api` and `/r` to the backend.

## Usage

1. Open the app and click **Create a new bin**
2. Copy the bin URL shown in the header
3. Send HTTP requests to it:

```bash
# Any method, any path
curl -X POST http://localhost:5173/r/<bin-id>/webhook \
  -H "Content-Type: application/json" \
  -d '{"event": "user.created", "id": 42}'

curl "http://localhost:5173/r/<bin-id>/ping?foo=bar"

curl -X PUT http://localhost:5173/r/<bin-id>/items/1 \
  -H "Content-Type: application/json" \
  -d '{"name": "updated"}'
```

Requests appear in the UI instantly via SSE — no polling.

## API

| Method | Path | Description |
|--------|------|-------------|
| `POST` | `/api/bins` | Create a new bin |
| `GET` | `/api/bins/:id` | Get bin details |
| `DELETE` | `/api/bins/:id` | Delete a bin and all its requests |
| `GET` | `/api/bins/:id/requests` | List requests (newest first, max 100) |
| `DELETE` | `/api/bins/:id/requests` | Clear all requests |
| `DELETE` | `/api/bins/:id/requests/:reqID` | Delete a single request |
| `GET` | `/api/bins/:id/sse` | SSE stream of incoming requests |
| `ANY` | `/r/:id` | Capture a request |
| `ANY` | `/r/:id/*` | Capture a request with sub-path |

## Project structure

```
.
├── docker-compose.yml          # Production
├── docker-compose.dev.yml      # Development
├── backend/
│   ├── main.go                 # Server setup, routes, graceful shutdown
│   ├── handler/
│   │   ├── bin.go              # Bin CRUD handlers
│   │   ├── capture.go          # Request capture handler
│   │   ├── sse.go              # SSE streaming handler
│   │   └── hub.go              # In-memory pub/sub hub
│   ├── store/
│   │   ├── store.go            # DB pool + migrations
│   │   ├── bin.go              # Bin queries
│   │   └── request.go          # Request queries
│   ├── model/                  # Shared types (Bin, Request)
│   ├── .air.toml               # Hot reload config
│   └── Dockerfile / Dockerfile.dev
└── frontend/
    ├── src/
    │   ├── views/
    │   │   ├── HomeView.vue    # Landing page, recent bins
    │   │   └── BinView.vue     # Inspector UI
    │   └── components/
    │       ├── RequestList.vue
    │       └── RequestDetail.vue
    ├── nginx.conf              # Reverse proxy config (production)
    └── Dockerfile / Dockerfile.dev
```

## Environment variables

### Backend

| Variable | Default | Description |
|----------|---------|-------------|
| `DATABASE_URL` | `postgres://requestbin:requestbin@localhost:5432/requestbin?sslmode=disable` | Postgres connection string |
| `PORT` | `8080` | HTTP listen port |

### Frontend (dev)

| Variable | Default | Description |
|----------|---------|-------------|
| `BACKEND_URL` | `http://localhost:8080` | Vite proxy target for `/api` and `/r` |
