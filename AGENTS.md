# AGENTS.md — OpenAPI Smoke Testing Platform

This file is the primary instruction document for AI agents implementing this repository.

Agents must read this file before making changes. When deeper context is needed, read the referenced documents in `docs/` and `docs/ai/`.

---

## 1. Mission

Build a collaborative OpenAPI-driven smoke testing platform delivered as **two products** that share one core:

- **Web-backed API** (`cmd/server`): UI-managed platform with persistence, live runs, reports, and collaboration.
- **CLI smoke runner** (`cmd/smokery`): standalone executable that runs compiled plans from local config, suitable for CI and local development.

Both products are first-class. Any feature that lives in the domain or application layer must work for both — never assume a database, an HTTP server, or a UI.

Core pipeline:

```text
OpenAPI Spec + User-Composed Smoke Configuration
        ↓
Compiler / Composer
        ↓
Executable Smoke Plan (the AST)
        ↓
Runner (library, not service)
        ↓
Persistent Run Result
        ↓
Reports, Diagrams, Debug Views, Collaboration, Insights
```

The platform must help teams:

- Import and analyze OpenAPI specs.
- Classify and override API operations.
- Compose smoke tests through UI or YAML.
- Run explicit flows and generated suites.
- Safely execute plans against environments.
- Persist structured run results.
- Generate useful reports for backend, frontend, analysts, QA, DevOps, and CI.
- Trace failures using request IDs, trace IDs, and observability links.

For the full architecture, see `docs/architecture.md`.

---

## 2. Source-of-truth documents

Read these documents before implementing major features:

1. `docs/design-docs.md`
   - Product concept and end-to-end design.
   - This is the primary product source of truth.

2. `docs/tech-stack.md`
   - Approved technologies and libraries.
   - Do not replace core stack choices without explicit approval.

3. `docs/tech-spec.html`
   - Technical specification, architecture, database sketch, routes, compiler, runner, UI, and MVP scope.

4. `docs/ai/context.md`
   - Condensed agent context.

5. `docs/ai/engineering-rules.md`
   - Coding, architecture, and implementation rules.

6. `docs/ai/product-rules.md`
   - Product behavior rules that must not be violated.

7. `docs/ai/workflows/*.md`
   - Feature-specific workflows.

8. `docs/ai/tasks/*.md`
   - Implementation task lists.

If these documents conflict, follow this priority:

```text
AGENTS.md
  > docs/design-docs.md
  > docs/tech-spec.html
  > docs/tech-stack.md
  > docs/ai/workflows
  > docs/ai/tasks
```

When a conflict is found, do not silently choose a behavior. Note the conflict in the implementation summary and prefer the higher-priority document.

---

## 3. Non-negotiable architecture rules

### 3.1 Compiler-first execution

The UI must not directly execute arbitrary form state.

The required execution path is:

```text
UI / CLI input
  → canonical project configuration
  → compiler validation
  → executable SmokePlan (the AST)
  → runner execution
  → structured RunResult
```

The runner must execute compiled plans only. The compiler is responsible for **baking everything the runner needs at runtime** into the SmokePlan — response schemas, content types, param formats. The runner does not import `libopenapi` and does not read the OpenAPI spec at runtime.

### 3.2 Control plane and execution plane separation

Maintain a clear boundary:

Control plane (lives at the application + delivery layers):

- Project management
- Spec import and analysis
- Operation registry
- Environment/auth config
- Flow and suite builders
- Plan preview
- Run orchestration
- Reports and collaboration

Execution plane (lives in `internal/runner` as a pure library):

- HTTP request execution
- Pluggable pre/post processors (auth injection, variable interpolation, capture, redaction, trace extraction, assertions)
- Structured result generation

The runner is a library. It accepts a `*model.SmokePlan` and an optional `RunOptions` (HTTP client, hooks, event emitter) and returns a `*model.RunResult`. It does not import database packages, HTTP frameworks, or the OpenAPI parser. This is what allows the same runner to power both the API server and the CLI binary.

See `docs/architecture.md` § 5 for the hook system and `docs/ai/workflows/03-runner.md` for the full runner contract.

### 3.3 Flow and Suite are both first-class

The system has two execution shapes:

```text
Flow
  Explicit, ordered, scenario-driven.

Suite
  Selected, generated, operation-driven.
```

Do not reduce suites to hand-authored flows. Do not reduce flows to endpoint lists.

### 3.4 List endpoint suites are first-class

List endpoints must support generated smoke coverage:

- Default list request.
- Response shape detection.
- Empty result policy: `allow | warn | fail`.
- Pagination sanity.
- Search-from-response.
- Enum filter cases.
- Operation selector rules.
- Operation-level overrides.
- Concurrent operation grouping.
- Per-operation and per-case reporting.

### 3.5 Destructive operation safety

Destructive operations must be protected by default.

Treat these as unsafe unless explicitly configured otherwise:

- `POST`
- `PUT`
- `PATCH`
- `DELETE`
- Business action endpoints such as approve, submit, recalculate, send, publish, pay, cancel.

The compiler must block destructive execution unless the flow/suite/project explicitly allows it.

### 3.6 Persistent run results

Smoke output must not be only terminal logs.

Persist structured results for:

- Run summary
- Flow results
- Step results
- Suite results
- Generated case results
- Request/response metadata
- Assertions
- Captures
- Trace context
- Comments
- Failure classification
- Artifacts

### 3.7 Reports must serve multiple audiences

Reports should support different views:

- Backend Debug
- Frontend Contract
- Analyst Flow
- QA Evidence
- Ops / Environment Health
- CI Summary

Do not design report data only for backend developers.

---

## 4. Approved tech stack

Use the selected stack unless the user explicitly changes it.

Backend:

- Go 1.26
- `danielgtaylor/huma/v2` + `humaecho` (OpenAPI 3.1 spec generation)
- `labstack/echo/v4`
- `jackc/pgx/v5`
- `sqlc`
- `golang-migrate`
- `pb33f/libopenapi`
- `tidwall/gjson`
- `spf13/viper`
- `rs/zerolog`
- `google/uuid`
- `gorilla/websocket`
- `minio/minio-go/v7`
- `testify`
- `testcontainers-go`

Frontend:

- SvelteKit 2.x / Svelte 5
- TypeScript 6.x strict mode
- `openapi-fetch` (type-safe API client)
- `openapi-typescript` (type generation from OpenAPI spec)
- `@tanstack/svelte-query` v6

Future frontend additions (not yet installed):

- TailwindCSS 4.x + shadcn-svelte
- `monaco-editor`
- `mermaid`
- `layerchart`
- `lucide-svelte`

Storage:

- PostgreSQL 18
- MinIO/S3
- Redis only in phase 2 for pub/sub/session scaling

Tooling:

- `bun` (package manager, runtime)
- `air` (Go hot-reload)
- `Makefile` (all dev commands)
- Docker Compose (infra)
- GitHub Actions (CI)

Do not introduce ORM magic. Use `sqlc` for database queries.

---

## 5. Repository expectations

Expected repository shape:

```text
.
├── Makefile
├── AGENTS.md
├── README.md
├── package.json
├── bun.lock
├── configs
│   ├── .air.toml
│   └── docker-compose.yml
├── apps
│   ├── core                   ← Go workspace (single go.mod, both products)
│   │   ├── go.mod
│   │   ├── sqlc.yaml
│   │   ├── cmd
│   │   │   ├── server         ← HTTP API binary
│   │   │   ├── smokery        ← CLI binary
│   │   │   └── openapi        ← internal: OpenAPI spec generator
│   │   └── internal
│   │       ├── model          ← Pure domain types (uuid.UUID, time.Time)
│   │       ├── spec           ← OpenAPI parser
│   │       ├── compiler       ← Config → SmokePlan
│   │       ├── runner         ← SmokePlan → RunResult (pluggable hooks)
│   │       │   └── hook       ← Built-in pre/post processors
│   │       ├── assertion      ← Assertion logic
│   │       ├── report         ← RunResult → views
│   │       ├── port           ← Interfaces (repos, jobs, events, blob)
│   │       ├── adapter
│   │       │   ├── postgres   ← pgx + sqlc impl (only place importing pgtype)
│   │       │   ├── inproc     ← in-process jobs + event bus
│   │       │   ├── minio      ← object storage
│   │       │   ├── fs         ← filesystem (CLI)
│   │       │   └── memory     ← in-memory (CLI, tests)
│   │       ├── app            ← Use case orchestration
│   │       └── delivery
│   │           ├── http       ← huma handlers + websocket (used by cmd/server)
│   │           └── cli        ← cobra commands (used by cmd/smokery)
│   └── web                    ← SvelteKit frontend
├── packages
│   └── types
├── docs
│   ├── architecture.md
│   ├── tech-stack.md
│   ├── design-docs.md
│   ├── tech-spec.html
│   └── ai
└── tmp                        ← Build output (gitignored)
```

The Go workspace is `apps/core`. It produces multiple binaries (`server`, `smokery`, `openapi`) from one Go module. Create directories only when they are needed by an implementation task.

---

## 6. Implementation behavior for AI agents

Before making code changes:

1. Identify the feature area.
2. Read the relevant docs.
3. Preserve existing architecture boundaries.
4. Make the smallest coherent change.
5. Add or update tests where practical.
6. Summarize changed files and reasoning.

Do not:

- Rewrite large unrelated areas.
- Replace the selected stack.
- Add new infrastructure without justification.
- Add background jobs outside River for MVP.
- Store real secrets in config or fixtures.
- Execute raw UI config directly.
- Treat OpenAPI as always correct.
- Skip compiler validation.
- Flatten rich run results into plain logs only.

---

## 7. Backend implementation rules

### 7.1 Layering

The backend follows hexagonal layering. See `docs/architecture.md` for the full diagram.

```text
delivery (http, cli) → app (use cases) → port (interfaces) ← adapter (infra)
                                              │
                                              ▼
                                         domain (model, spec, compiler, runner, assertion, report)
```

Strict rules:

- The `domain` layer imports nothing from this project except other domain packages.
- The `app` layer imports `port`, `model`, and domain services. It never imports adapters.
- The `delivery` layer imports `app` and `model`. It is the only layer that may import HTTP frameworks (`echo`, `gorilla/websocket` for `delivery/http`) or CLI frameworks (`cobra` for `delivery/cli`).
- The `cmd/*` packages are the **only** place where concrete adapters meet app services. `main.go` is dependency-injection.

### 7.2 Internal package boundaries

```text
internal/model        Pure domain types (uuid.UUID, time.Time)
internal/spec         OpenAPI parsing and operation classification
internal/compiler     Project config → executable SmokePlan
internal/runner       SmokePlan → RunResult (pluggable library)
internal/runner/hook  Built-in pre/post processors
internal/assertion    Assertion and validation helpers
internal/report       RunResult → report views/artifacts
internal/port         Interfaces for storage, jobs, events, blob
internal/adapter/*    Concrete infrastructure implementations
internal/app          Use case orchestration (depends only on ports)
internal/delivery/*   HTTP / CLI handlers (depend only on app services)
```

### 7.3 Type translation rule

Domain types are pure Go (`uuid.UUID`, `time.Time`, `[]byte`). Persistence types (sqlc-generated `db.Project` with `pgtype.UUID`, etc.) live **only** inside `internal/adapter/postgres/`. The repo adapter is the only code that converts between them.

After refactor, this command must return only files under `internal/adapter/postgres/`:

```bash
grep -r "pgtype\|pgxpool\|jackc/pgx" apps/core/internal --include="*.go"
```

Same principle for `minio-go` (only inside `internal/adapter/minio/`) and `gorilla/websocket` (only inside `internal/delivery/http/`).

### 7.4 Database access

- Use PostgreSQL 18.
- Use pgx and sqlc-generated queries inside the postgres adapter only.
- Use migrations through golang-migrate.
- Prefer explicit SQL over hidden ORM behavior.
- Store flexible run/config payloads in JSONB where appropriate.
- Store large artifacts via `port.BlobStore`, never in PostgreSQL.

### 7.5 Storage (blobs)

- Define a single `port.BlobStore` interface.
- Provide implementations under `internal/adapter/{minio,fs,memory}`.
- App services and use cases depend on `port.BlobStore`, never on `minio-go`.

### 7.6 Jobs

- Define `port.JobEnqueuer` and `port.EventBus` interfaces.
- The MVP implementation lives in `internal/adapter/inproc` (goroutines + channels).
- River can be added later as `internal/adapter/river/` without changes to app or domain code.
- A run is represented by a database row + an enqueued job + result row.

### 7.7 Logging

Use zerolog structured logging. Never log secrets or unredacted auth headers.

### 7.8 Validation

- Use huma's struct-tag validation for HTTP request inputs (handled in `delivery/http`).
- Use explicit compiler validation for smoke configs.

Compiler errors must include:

- Stage
- Field path
- Message
- Severity: error or warning
- Related operation/flow/suite if any

### 7.9 Runner contract

The runner accepts a `*model.SmokePlan` and an optional `RunOptions`, and returns a `*model.RunResult`. It must not:

- Import `libopenapi` or read the OpenAPI spec
- Import any database package
- Import `minio-go` or any storage package
- Embed orchestration logic (job queueing, persistence)

All extension is done through `PreProcessor` / `PostProcessor` hooks. See `docs/ai/workflows/03-runner.md`.

---

## 8. Frontend implementation rules

### 8.1 Product UX

The UI should guide users through:

```text
Project
  → Spec import
  → Operation review
  → Environment/auth setup
  → Suite/flow composition
  → Plan preview
  → Run
  → Report/collaboration
```

### 8.2 UI must reflect compiler state

Plan preview and validation errors are first-class UI states.

Do not let users start a run from invalid configuration.

### 8.3 SvelteKit patterns

- Use TypeScript strict mode.
- Use SvelteKit routes and load functions appropriately.
- Use TanStack Query for client-side server state.
- Use Superforms + Zod for forms.
- Use Monaco for JSON/YAML override editors.
- Use Mermaid for diagrams.
- Use LayerChart for trends/latency.
- Keep components small and feature-focused.

### 8.4 Required UI areas

Implement around these pages:

```text
/projects
/projects/[id]/spec
/projects/[id]/operations
/projects/[id]/environments
/projects/[id]/flows/[fid]
/projects/[id]/suites/[sid]
/projects/[id]/plan
/projects/[id]/runs
/runs/[runId]
/runs/[runId]/report/[view]
/projects/[id]/settings
```

---

## 9. Security and data handling rules

- Never store real secrets in exported config.
- Redact auth headers and sensitive JSON paths.
- Support retention policies for request/response bodies.
- Default to storing full bodies only for failed cases.
- Always redact known sensitive fields.
- Do not expose runner tokens or internal IDs unnecessarily.
- Treat persisted run results as potentially sensitive.

Default redaction targets:

```text
Authorization
Cookie
Set-Cookie
password
token
accessToken
refreshToken
apiKey
secret
```

---

## 10. Testing expectations

Backend:

- Unit test compiler behavior.
- Unit test operation classification.
- Unit test variable interpolation.
- Unit test JSONPath capture.
- Unit test assertions.
- Integration test persistence with testcontainers PostgreSQL.
- Integration test runner against local test HTTP server.

Frontend:

- Type-check all changes.
- Test critical form schemas.
- Test route data assumptions where practical.
- Keep UI components deterministic and easy to inspect.

Minimum quality gate:

```text
make test
make check
make build
```

Available make targets:

```text
make build          Build API server to tmp/server
make dev            Run API with air hot-reload (watches apps/api)
make generate       Generate OpenAPI spec + TypeScript types
make generate-sqlc  Regenerate sqlc queries
make test           Run Go tests
make check          Run svelte-check
make lint           Run test + check
make up             Start Docker Compose (PostgreSQL + MinIO)
make down           Stop Docker Compose
make migrate        Run database migrations
make install        Install all dependencies
make clean          Remove tmp/ and .svelte-kit
```

All Go binaries must be output to the `tmp/` directory.

Infrastructure configs live in `configs/`:

```text
configs/
├── .air.toml          Air hot-reload config
└── docker-compose.yml PostgreSQL + MinIO
```

---

## 11. MVP boundaries

MVP should include:

- Project + spec import + analysis.
- Operation registry + overrides.
- Environment + auth profile config.
- Flow builder.
- Generated list suite.
- Compiler + plan preview.
- River-backed runner job.
- WebSocket live run feed.
- Persistent run result + history.
- Backend debug report.
- CI summary report.
- Request ID / trace extraction.
- Basic Mermaid sequence diagram.
- Run comments.
- JSON + HTML report artifacts.
- Redaction + retention policy.

Phase 2 features should not block MVP unless required by architecture.

---

## 12. Completion summary format

When finishing an implementation task, provide:

```text
Implemented:
- ...

Changed files:
- path/to/file: reason

Validation:
- command/result

Notes / follow-ups:
- ...
```

Be honest about skipped tests or incomplete pieces.
