# Technical decisions

## Conflict rule

- **15-minute slots:** Each slot is 15 minutes. A user can book 1 or 2 consecutive slots (15 or 30 min). This use case was chosen to show how conflicts are resolved when a slot choice overlaps with another user's.
- **Conflict:** If any of the requested slots overlap with an existing appointment for the same doctor, the booking is rejected. Overlap is defined as: existing_start < our_end AND our_start < existing_end (in slot-time).
- **Single doctor:** We stick to one doctor, hence the dropdown has one option and this is in order to avoid extra product scope.

## Data persistence

- **SQLite** via `modernc.org/sqlite` was chosen for simplicity and zero external dependencies; one file `scheduler.db` per run. We would switch to Postgres for production scale projects but SQLite gets the job done for now.

## Auth and user details

- **No auth.** Anyone can create or delete any appointment. Chosen to keep scope minimal.
- **User details at submit:** Patient name and contact are collected in the booking form and stored on the appointment; no login or user accounts.

## Single-doctor dropdown

- The UI shows a doctor dropdown populated from `GET /api/doctors`. For this assessment it returns one static doctor. The backend defaults to that doctor if `doctor_id` is omitted.

## List API and appointments UI

- **All appointments:** `GET /api/appointments` with no `from`/`to` (or both 0) returns every appointment in the store (backend uses range 0 to 1<<62). The UI does not send a time range and shows all appointments, past and future. This avoids timezone/range bugs and keeps behaviour simple.