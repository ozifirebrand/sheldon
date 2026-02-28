# Scheduler

Appointment scheduling: book and manage 15-minute slots for a single doctor. Real-time updates via server-sent events.

## Quick start

1. **Start the backend** (terminal 1):
   ```bash
   cd backend && go mod tidy && go run ./cmd/server
   ```
2. **Start the frontend** (terminal 2):
   ```bash
   cd frontend && npm install && npm run dev
   ```
3. Open **http://localhost:5173** in your browser.

## Backend (Go gRPC + HTTP gateway)

- **Requirements:** Go 1.21+
- **Run:** From repo root, `cd backend && go mod tidy && go run ./cmd/server`
- **Options (flags override env):** `-http=:8100` (HTTP API), `-grpc=:9090` (gRPC), `-db=scheduler.db`
- **Env (optional):** `GRPC_ADDR`, `HTTP_ADDR`, `DB_PATH`
- **Layout:** `cmd/server` wires config → store, broadcast, service, gRPC server, HTTP gateway. `internal/`: config, constants, repository (interface), store (SQLite), service, broadcast, gateway (REST + SSE), server (gRPC).
- **Regenerate proto (do not edit `gen/*.pb.go`):** From `backend`, run  
  `protoc -I. --go_out=. --go_opt=module=scheduler/backend --go-grpc_out=. --go-grpc_opt=module=scheduler/backend proto/scheduler.proto`
- **API:** HTTP gateway on port 8100: `GET/POST /api/appointments`, `DELETE /api/appointments/:id`, `GET /api/appointments/stream` (SSE), `GET /api/doctors`. gRPC server on 9090 for direct clients.

## Frontend (React + TypeScript)

- **Requirements:** Node 18+
- **Run:** From repo root, `cd frontend && npm install && npm run dev`
- **URL:** http://localhost:5173 (Vite proxies `/api` to the backend at http://localhost:8100)
- **Env (optional):** `VITE_API_URL` for full API base URL when not using the proxy (e.g. production).
- **Layout:** `src/api` (API client), `src/components` (ErrorBoundary, BookingForm, AppointmentList), `src/utils` (slot helpers), `src/constants`, `src/config`.
- **Setup:** Start the backend first so the frontend can reach the API.

Both services run without errors when the backend is started before the frontend.
