---
type: technical-note
status: current
tags:
  - technical
  - tech-stack
related:
  - "[[Architecture]]"
  - "[[Engineering Rules]]"
---

# Tech Stack

> [!summary]
> Smokery uses a Go backend, SvelteKit frontend, SQLite-or-PostgreSQL persistence, and filesystem-or-MinIO artifact storage. The docs should describe both the approved stack and the actual local-first defaults.

## Backend

- Go 1.26
- `labstack/echo/v4`
- `danielgtaylor/huma/v2` + `humaecho`
- `spf13/cobra`
- `jackc/pgx/v5`
- `sqlc`
- `golang-migrate`
- `pb33f/libopenapi`
- `tidwall/gjson`
- `spf13/viper`
- `rs/zerolog`
- `minio/minio-go/v7`
- `gorilla/websocket`
- `modernc.org/sqlite`
- `sigs.k8s.io/yaml`
- `testify` + `testcontainers-go`

## Frontend

- SvelteKit 2 / Svelte 5
- TypeScript 6 strict mode
- `@tanstack/svelte-query` v6
- `openapi-fetch`
- `openapi-typescript`
- TailwindCSS 4
- Bits UI / shadcn-svelte style primitives
- `sveltekit-superforms`
- Mermaid
- LayerChart
- Lucide
- `@xyflow/svelte`
- Vite
- Vitest + Testing Library

## Storage And Tooling

- SQLite default for local-first development
- PostgreSQL 18 optional relational backend
- Filesystem blob store default
- MinIO / S3 optional artifact backend
- Redis remains later work only
- `bun`
- Docker Compose
- GitHub Actions
- `air`
- `Makefile`

## Current Defaults

- `DB_ADAPTER=sqlite`
- `STORAGE_ADAPTER=fs`
- `SQLITE_PATH=data/smokery.db`
- `STORAGE_PATH=data/artifacts`

## Related Notes

- [[Architecture]]
- [[Repository Layout]]
- [[Engineering Rules]]

