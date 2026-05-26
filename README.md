# OpenAPI Smoke Testing Platform

A collaborative OpenAPI-driven smoke testing platform delivered as **two products**:

- **Web-backed API** — UI-managed platform with persistence, live runs, and reports.
- **CLI smoke runner** (`smokery`) — standalone executable for local config and CI.

Both products share the same compiler, runner, and domain code. The CLI is not a stripped-down API.

The platform imports an OpenAPI spec, analyzes the API surface, helps teams configure smoke test flows and generated suites, compiles those configurations into executable smoke plans, runs them through a Go runner, and persists rich reports for backend developers, frontend developers, analysts, QA, and DevOps.

## Core idea

```text
OpenAPI Spec + Smoke Configuration
        ↓
Compiler / Composer
        ↓
Executable Smoke Plan
        ↓
Runner
        ↓
Persistent Run Result
        ↓
Reports, Diagrams, Debug Views, Collaboration
```

## Main goals

- Import and analyze OpenAPI specs.
- Classify API operations and allow user overrides.
- Compose smoke tests visually through UI.
- Support explicit scenario flows.
- Support generated list endpoint suites.
- Compile configuration into deterministic executable plans.
- Execute plans safely with destructive-operation protection.
- Persist structured run results.
- Generate audience-specific reports and diagrams.
- Correlate smoke failures with server logs through request IDs and trace IDs.

## Tech stack

- Backend: Go 1.26, Echo, Huma (OpenAPI 3.1), pgx, sqlc, libopenapi, gjson.
- Frontend: SvelteKit 2, Svelte 5, TypeScript 6, openapi-fetch, openapi-typescript, TanStack Query.
- Storage: PostgreSQL 18, MinIO/S3.
- Dev tooling: Docker Compose, GitHub Actions, bun workspaces, air, Makefile.

See:

- `docs/architecture.md` for the hexagonal architecture and two-product layering.
- `docs/design-docs.md` for the product and concept design.
- `docs/tech-stack.md` for selected libraries and rationale.
- `docs/tech-spec.html` for the technical implementation specification.
- `AGENTS.md` for AI agent implementation rules.
- `docs/ai/` for split AI-agent workflows and task lists.
- `Makefile` for all dev commands.
- `configs/` for infrastructure configs.

## Expected repository layout

```text
.
├── Makefile
├── AGENTS.md
├── README.md
├── configs
│   ├── .air.toml
│   └── docker-compose.yml
├── apps
│   ├── core              ← Go workspace (one go.mod, both products)
│   │   ├── cmd
│   │   │   ├── server    ← HTTP API binary
│   │   │   ├── smokery   ← CLI binary
│   │   │   └── openapi   ← spec generator
│   │   └── internal
│   │       ├── model, spec, compiler, runner, assertion, report  ← domain
│   │       ├── port             ← interfaces
│   │       ├── adapter          ← infra (postgres, inproc, minio, fs, memory)
│   │       ├── app              ← use cases
│   │       └── delivery         ← http, cli
│   └── web               ← SvelteKit frontend
├── packages
│   └── types
└── docs
    ├── architecture.md
    ├── tech-spec.html
    ├── tech-stack.md
    ├── design-docs.md
    └── ai
```

## Development principle

This project is expected to be implemented primarily by AI agents. Agents must follow the design documents, preserve the compiler-first architecture, and avoid inventing unrelated product behavior without updating the design docs first.

## Development

```bash
# Install dependencies
make install

# Start infrastructure (PostgreSQL + MinIO)
make up

# Run database migrations
make migrate

# Generate OpenAPI spec + TypeScript types
make generate

# Build API server (outputs to tmp/server)
make build

# Run API with hot-reload (air)
make dev

# Run tests and checks
make lint
```
