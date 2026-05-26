# Tech Stack — OpenAPI Smoke Testing Platform

> Go backend · SvelteKit frontend · PostgreSQL · MinIO

---

## Backend (Go 1.26)

| Concern | Package | Notes |
|---|---|---|
| HTTP framework | `labstack/echo/v4` | Middleware-first, fast, WebSocket support |
| OpenAPI generation | `danielgtaylor/huma/v2` + `humaecho` | Auto-generates OpenAPI 3.1 spec from typed operations |
| CLI framework | `spf13/cobra` | For the `smokery` CLI binary |
| DB driver | `jackc/pgx/v5` | Native PG driver (only inside `internal/adapter/postgres`) |
| Query generation | `sqlc` | Type-safe SQL, avoids ORM magic |
| Migrations | `golang-migrate` | File-based, CI-friendly |
| OpenAPI parsing | `pb33f/libopenapi` | Used by compiler, **not** by runner |
| JSONPath extraction | `tidwall/gjson` | Used by Capture hook |
| Config | `spf13/viper` | Env + file config |
| Logging | `rs/zerolog` | Structured JSON, zero-alloc |
| Object storage | `minio/minio-go/v7` | Only inside `internal/adapter/minio` |
| WebSocket | `gorilla/websocket` | Only inside `internal/delivery/http` |
| Testing | `testify` + `testcontainers-go` | Integration tests with real PG containers |
| Hot reload | `air-verse/air` | Dev-time file watching and rebuild |

### Hook architecture

The runner is a pluggable library. All extension is via `PreProcessor` / `PostProcessor` interfaces. Built-in hooks (auth injection, variable interpolation, capture, redaction, trace extraction, assertions) live in `internal/runner/hook`. See `docs/architecture.md` § 5 and `docs/ai/workflows/03-runner.md`.

---

## Frontend (SvelteKit 2 / Svelte 5)

| Concern | Package | Notes |
|---|---|---|
| Framework | `SvelteKit 2.x` / `Svelte 5` | Runes, file-based routing |
| Language | TypeScript 6.x | Strict mode throughout |
| Data fetching | `@tanstack/svelte-query` v6 | Accessor pattern, cache, mutations |
| API client | `openapi-fetch` | Type-safe fetch from OpenAPI schema |
| Type generation | `openapi-typescript` | Generates `.d.ts` from OpenAPI spec |
| Build | `vite` 8.x | Fast bundler |

Future additions (not yet implemented):

- TailwindCSS 4.x + shadcn-svelte for UI
- Monaco for JSON/YAML editing
- Mermaid for diagrams
- LayerChart for trends/latency

---

## Data Storage

| Store | Version | Use |
|---|---|---|
| PostgreSQL | 18 | All relational data, JSONB for run payloads |
| MinIO / S3 | — | Full run result artifacts, HTML reports |
| Redis | *(phase 2)* | Session cache, pub/sub at scale |

---

## Infrastructure & Tooling

| Concern | Tool | Notes |
|---|---|---|
| Package manager | `bun` | Runtime, package manager, script runner |
| Containerisation | Docker + Docker Compose | Dev environment via `make up` |
| CI/CD | GitHub Actions | Go build + bun check |
| Hot reload | `air` | Watches Go source, rebuilds to `tmp/` |
| Build system | `Makefile` | All dev commands via make targets |
| Monorepo | bun workspaces | `apps/core`, `apps/web`, `packages/types` |
| Type sharing | `openapi-typescript` | Generate TS types from API's OpenAPI spec |
| Configs | `configs/` dir | docker-compose.yml, .air.toml |

---

## Go Module Versions (go.mod snapshot)

```go
require (
    github.com/danielgtaylor/huma/v2       v2.38.0
    github.com/labstack/echo/v4            v4.15.2
    github.com/jackc/pgx/v5               v5.9.2
    github.com/pb33f/libopenapi           v0.36.6
    github.com/tidwall/gjson              v1.19.0
    github.com/google/uuid                v1.6.0
    github.com/gorilla/websocket          v1.5.3
    github.com/rs/zerolog                 v1.35.1
    github.com/spf13/viper                v1.21.0
)
```

---

## Frontend Package Versions (package.json snapshot)

```json
{
  "dependencies": {
    "@tanstack/svelte-query": "^6.1.33",
    "openapi-fetch": "^0.17.0"
  },
  "devDependencies": {
    "@sveltejs/kit": "^2.57.0",
    "openapi-typescript": "^7.13.0",
    "svelte": "^5.55.2",
    "typescript": "^6.0.2",
    "vite": "^8.0.7"
  }
}
```

---

*Last updated: 2026-05-26*
