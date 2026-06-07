---
type: technical-note
status: current
tags:
  - technical
  - architecture
  - hexagonal
related:
  - "[[Product Vision]]"
  - "[[Repository Layout]]"
  - "[[Compiler Pipeline]]"
  - "[[Runner Execution]]"
  - "[[ADR - Hexagonal Architecture]]"
---

# Architecture

> [!summary]
> Smokery uses hexagonal architecture so the API server and CLI smoke runner share one domain core while swapping delivery and infrastructure adapters. The current runtime defaults are SQLite plus filesystem storage, with PostgreSQL and MinIO available as optional adapters.

## Two Products, One Core

| Product | Delivery | Storage | Use case |
|---|---|---|---|
| API server | `cmd/server` | default `sqlite + fs`, optional `postgres + minio` | Local-first collaborative platform with UI, persistence, reports, analytics, governance, and live runs |
| CLI smoke runner | `cmd/smokery` | filesystem + memory | Standalone execution from local config, CI-friendly |

## Current Runtime Composition

- `cmd/server` wires repositories from either `internal/adapter/sqlite` or `internal/adapter/postgres`.
- Blob artifacts come from `internal/adapter/fs` by default, with `internal/adapter/minio` as the alternative.
- Run execution is orchestrated through the in-process worker and event bus in `internal/adapter/inproc`.
- Production builds can embed the frontend into the server binary through `-tags embed_frontend`.
- The runner remains a pure library under `internal/runner`.
- Additional app services now back failure classification, spec evolution, analytics, and governance endpoints.

## Dependency Direction

```text
delivery → app → port ← adapter
              │
              ▼
            domain
```

Strict rules:

- `domain` imports nothing from this project except other domain packages.
- `port` imports `model` only.
- `adapter` imports `port`, `model`, and infrastructure SDKs.
- `app` imports `port`, `model`, and domain services. It never imports adapters.
- `delivery` imports `app` and `model`.
- `cmd/*/main.go` is dependency injection only.

> [!warning]
> Any dependency inversion violation is a bug.

## Boundary Rules

- `pgtype`, `pgxpool`, and sqlc-generated persistence types stay in `internal/adapter/postgres/`.
- `modernc.org/sqlite` stays in `internal/adapter/sqlite/`.
- `minio-go` stays in `internal/adapter/minio/`.
- `gorilla/websocket` stays in `internal/delivery/http/`.
- `cobra` stays in `internal/delivery/cli/`.
- `libopenapi` stays in `internal/spec/` and `internal/compiler/`.

## Current Public Surface

- HTTP API: projects, specs, operations, plan preview, runs, reports, comments, artifacts, websocket run events, failure classification, spec evolution, analytics, governance, and health checks.
- CLI: `import-spec`, `compile`, `run`, and `report`.
- Frontend: project list/detail, spec review, versions, diff, operations, environments, builder, flows, suites, plan preview, runs, run detail, report pages, analytics pages, impact page, and settings pages.

> [!note]
> Some frontend report and governance pages may still be partially mock-backed or not fully wired to the newer backend endpoints. Current-state docs should say that explicitly instead of collapsing planned and shipped behavior together.

## Related Notes

- [[Repository Layout]]
- [[Compiler Pipeline]]
- [[Runner Execution]]
- [[Tech Stack]]

