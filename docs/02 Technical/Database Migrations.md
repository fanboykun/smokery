---
type: technical-note
status: current
tags:
  - technical
  - database
  - migrations
related:
  - "[[Tech Stack]]"
  - "[[Repository Layout]]"
  - "[[Project State]]"
---

# Database Migrations

> [!summary]
> Smokery keeps SQL migrations under `apps/core/db/migrations` and supports two runtime modes: SQLite by default for local-first work, and PostgreSQL when the server is configured to use the postgres adapter.

## Current Migration Reality

- The checked-in migration set is PostgreSQL-oriented and managed with `golang-migrate`.
- `make migrate` runs the migrate CLI against `DATABASE_URL`.
- Local default development can still use SQLite without requiring PostgreSQL to be up first.

## Commands

```bash
make up
make migrate
make generate-sqlc
```

## Current Caveat

> [!warning]
> Do not describe PostgreSQL migrations as if they are required for every local run. The current server defaults to SQLite unless `DB_ADAPTER=postgres` is set.

## Related Notes

- [[Tech Stack]]
- [[Repository Layout]]
- [[Project State]]

