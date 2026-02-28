# Decisions made for this app

## Slots and conflicts

- Each slot is 15 minutes. You can book 1 or 2 slots in a row (15 or 30 minutes). I did it this way so we can clearly see what happens when two people try to book overlapping times.

- Conflict is if the time you want, overlaps with an existing appointment for the same doctor. This booking is rejected. Overlap means a user's start is before another user's end. Double booking is not accepted in other words.

- The app has a single doctor in a dropdown. I kept it to one to keep the scope small and focus on conflicts and real-time updates.

## Where data is stored

- I use SQLite (via `modernc.org/sqlite`) so we don’t need any external database. Everything lives in one file, `scheduler.db`, for the run. For a bigger production app I’d switch to something like Postgres, but for this, SQLite is okay.

## Concurrency 

- Creating an appointment happens inside a mutex and a single database transaction so the system allows one writer at a time per slot. We check for overlap in that same transaction before inserting.

- if two requests try to book the same slot at the same time, only one can win; the other gets a conflict error notification.

- I did not add any locking across multiple server processes. For one running server this is fine.

## Real-time updates

- The gRPC API has a stream that allows clients to call `SubscribeAppointmentUpdates` and get a stream of events. The server keeps an in-memory broadcaster and pushes to every connected stream when something is created or deleted.

- The HTTP server exposes `GET /api/appointments/stream`. The React app connects with `EventSource` and re-fetches the list on each event. So if someone books or cancels on another device, you see it right away.

## Scalability

- This is a single and simple Go process with an in-memory broadcaster and SQLite. 

- For multiple instances, I would; 
  - Use a shared, remote database, preferably an RDMS
  - Make the server stateless and use shared pub/sub technology to alert users and get notifications in real time
  - Use DB advisory locks and/or unique constraints for conflict detection when many instances are writing.


## Auth and user info

- Auth was left out to keep scope focused. Anyone can create or delete any appointment.

## Doctor dropdown

- The UI gets the doctor list from `GET /api/doctors`. For this app it returns one static doctor, Dr Smith as default.

## List API and how the UI uses it

- Filtering is not happening with`GET /api/appointments` endpoint. It keeps things simple and avoids managing timezone/range edge cases now.

- The backend always returns a JSON array (empty is `[]`, not `null`) which the frontend normalizes.

- After a successful creation of appointment, optimistic update is done internally where the frontend adds the new appointment to the list right away, then re-fetches.

## Possible conversations for scope expansion
- How to handle timezone changes and differences
- Include multiple Doctors?
- Add unit tests for the store’s conflict logic and integration tests for create/list/delete/stream.