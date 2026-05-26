# AI Agent Context

This project is an OpenAPI-driven smoke testing platform delivered as **two products**:

- **Web-backed API** (`cmd/server`): UI-managed platform with persistence, live runs, reports.
- **CLI smoke runner** (`cmd/smokery`): standalone executable for local config and CI.

Both products share the same domain code. The CLI is not a stripped-down API.

The platform is intended to be implemented mostly by AI agents, so the design must remain explicit, consistent, and easy to follow.

## Product identity

```text
OpenAPI-driven compiler + composable smoke-test platform + persistent collaborative reporting
                              + standalone CLI smoke runner
```

## Core pipeline

```text
OpenAPI Spec + Smoke Config
        ↓
Compiler / Composer
        ↓
Executable Smoke Plan (the AST — self-contained)
        ↓
Runner (library, no infra dependencies)
        ↓
Persistent Run Result
        ↓
Reports, Diagrams, Debug Views, Collaboration
```

The compiler is the **bridge** between OpenAPI knowledge and runtime execution. The compiler bakes everything the runner needs into the SmokePlan. The runner does not import `libopenapi` or read the OpenAPI spec at runtime.

## Main execution shapes

### Flow

Explicit, ordered, scenario-driven.

Used for:

- Login flows
- CRUD lifecycle checks
- Business journey smoke tests
- Frontend/analyst flow documentation

### Suite

Selected, generated, operation-driven.

Used for:

- Broad API coverage
- List endpoint sanity checks
- Read-only endpoint sweeps
- Generated pagination/search/filter cases

## Critical design concepts

- OpenAPI describes possibilities.
- User configuration selects intent.
- Compiler validates and composes a plan.
- Runner executes compiled plans only — pluggable, dependency-free library.
- Hooks (pre/post processors) are how the runner is extended.
- Destructive operations are blocked by default.
- Cleanup is first-class.
- List endpoint suites are first-class.
- Run results are persistent and collaborative.
- Reports serve multiple audiences.
- The API exposes its own OpenAPI 3.1 spec (via huma) consumed by the frontend.
- The same domain works for the API server and the CLI binary — no infra dependency in `internal/{model,spec,compiler,runner,assertion,report,app}`.

## Layering (hexagonal)

```text
delivery (http, cli) → app → port ← adapter (postgres, inproc, minio, fs, memory)
                              │
                              ▼
                            domain (model, spec, compiler, runner, assertion, report)
```

- Domain imports nothing infra-related.
- Adapters are the only places that import `pgtype`, `minio-go`, `gorilla/websocket`, etc.
- App services depend only on ports.
- `cmd/*/main.go` is the only place where concrete adapters meet app services.

See `docs/architecture.md` for the full architecture.

## Dev workflow

```text
make up          → start PostgreSQL + MinIO
make dev         → run API with air hot-reload
make generate    → regenerate OpenAPI spec + TS types
make test        → run Go tests
make check       → run svelte-check
make lint        → test + check
```

## Existing docs

- `docs/architecture.md`: hexagonal layering, two products, runner contract.
- `docs/design-docs.md`: product concept.
- `docs/tech-stack.md`: approved stack.
- `docs/tech-spec.html`: technical spec and MVP shape.
- `AGENTS.md`: base AI implementation rules.
- `Makefile`: all dev commands.
- `configs/`: infrastructure configs (docker-compose, air).
