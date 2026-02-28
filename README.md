# Scheduler

Appointment scheduling: book and manage 15-minute slots for a single doctor. Real-time updates via server-sent events.

## Quick start

1. **Start the backend** (terminal 1):
   ```bash
   cd backend && go mod tidy && go run .
   ```
2. **Start the frontend** (terminal 2):
   ```bash
   cd frontend && npm install && npm run dev
   ```
3. Open **http://localhost:5173** in your browser.

## Backend (Go gRPC + HTTP gateway)

- **Requirements:** Go 1.21+
- **Run:** From repo root: `cd backend && go mod tidy && go run .`
- **Flags (override env):** `-http=:8100`, `-grpc=:9090`, `-db=scheduler.db`
- **Env (optional):** `HTTP_ADDR`, `GRPC_ADDR`, `DB_PATH` (defaults: `:8100`, `:9090`, `scheduler.db`)
- **Layout:** `main.go` at backend root loads config, creates store (SQLite), broadcaster, service, gRPC server, and HTTP gateway. `internal/`: `config`, `constants`, `repository` (interface), `store`, `service`, `broadcast`, `gateway` (REST + SSE), `server` (gRPC). Generated code: `gen/schedulerv1/*.pb.go` from `proto/scheduler.proto`.
- **Regenerate proto (do not edit `gen/schedulerv1/*.pb.go` by hand):** From `backend`:
  ```bash
  protoc -I. --go_out=. --go_opt=module=scheduler/backend --go-grpc_out=. --go-grpc_opt=module=scheduler/backend proto/scheduler.proto
  ```
- **API:** HTTP on port 8100: `GET/POST /api/appointments`, `DELETE /api/appointments/:id`, `GET /api/appointments/stream` (SSE), `GET /api/doctors`. gRPC on 9090 for direct clients.

## Frontend (React + TypeScript)

- **Requirements:** Node 18+
- **Run:** From repo root: `cd frontend && npm install && npm run dev`
- **URL:** http://localhost:5173 (Vite proxies `/api` to the backend at http://localhost:8100)
- **Env (optional):** `VITE_API_URL` — full API base URL when not using the dev proxy (e.g. production).
- **Layout:** `src/api/index.ts` (API client), `src/components/` (BookingForm, AppointmentList, ErrorBoundary), `src/utils/slotUtils.ts`, `src/constants.ts`, `src/config.ts`, `src/App.tsx`, `src/main.tsx`.
- **Start the backend first so the frontend can reach the API.**
