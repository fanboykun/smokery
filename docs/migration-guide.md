# Database Migration Guide

Smokery uses [golang-migrate](https://github.com/golang-migrate/migrate) for PostgreSQL schema management.

## Prerequisites

- PostgreSQL 16+ running (via `make up` or external)
- `DATABASE_URL` environment variable set (or use default: `postgres://smokery:smokery@localhost:5432/smokery?sslmode=disable`)

## Running Migrations

### Using Make

```bash
# Apply all pending migrations
make migrate

# Start infrastructure first if needed
make up
make migrate
```

### Using golang-migrate directly

```bash
# Install migrate CLI
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Apply all up migrations
migrate -path apps/core/db/migrations \
  -database "postgres://smokery:smokery@localhost:5432/smokery?sslmode=disable" \
  up

# Roll back the last migration
migrate -path apps/core/db/migrations \
  -database "postgres://smokery:smokery@localhost:5432/smokery?sslmode=disable" \
  down 1

# Roll back ALL migrations (destructive!)
migrate -path apps/core/db/migrations \
  -database "postgres://smokery:smokery@localhost:5432/smokery?sslmode=disable" \
  down

# Check current version
migrate -path apps/core/db/migrations \
  -database "postgres://smokery:smokery@localhost:5432/smokery?sslmode=disable" \
  version

# Force a specific version (use when dirty state)
migrate -path apps/core/db/migrations \
  -database "postgres://smokery:smokery@localhost:5432/smokery?sslmode=disable" \
  force 1
```

## Migration Files

Migrations live in `apps/core/db/migrations/` and follow the naming convention:

```
000001_init.up.sql      # Apply schema
000001_init.down.sql    # Rollback schema
```

### Current schema (000001_init)

**Up** creates:
- `projects` — top-level project container
- `specs` — imported OpenAPI specs (raw + analysis JSON)
- `environments` — target API environments (base URL, headers)
- `auth_profiles` — authentication configurations
- `operations` — classified API operations from specs
- `flows` — user-defined smoke test flows
- `suites` — generated smoke test suites
- `plans` — compiled smoke plans (snapshots)
- `runs` — execution records with lifecycle status
- `run_results` — structured JSON results per run
- `comments` — run annotations
- `artifacts` — references to stored report files

**Down** drops all tables in reverse dependency order.

## Creating New Migrations

```bash
# Create a new migration pair
migrate create -ext sql -dir apps/core/db/migrations -seq <description>

# Example
migrate create -ext sql -dir apps/core/db/migrations -seq add_tags_to_projects
```

This creates:
```
000002_add_tags_to_projects.up.sql
000002_add_tags_to_projects.down.sql
```

## Best Practices

- Always write both `up.sql` and `down.sql`
- Use `IF NOT EXISTS` / `IF EXISTS` for idempotency where possible
- Never modify an already-applied migration — create a new one
- Test rollback locally before pushing: `migrate down 1 && migrate up`
- After adding migrations, regenerate sqlc: `make generate-sqlc`

## Troubleshooting

### Dirty database state

If a migration fails halfway, the database is marked "dirty":

```bash
# Check current state
migrate version
# Output: 1 (dirty)

# Force to the last known good version
migrate force 1

# Then retry
migrate up
```

### Connection refused

Ensure PostgreSQL is running:

```bash
make up          # starts docker-compose
make migrate     # applies migrations
```

### Permission denied

The default user `smokery` needs CREATE/DROP privileges on the target database. The docker-compose setup handles this automatically.
