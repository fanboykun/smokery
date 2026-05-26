# MVP Task List — OpenAPI Smoke Testing Platform

This task list is designed for AI agents to implement the project incrementally.

The Go workspace is `apps/core` (one Go module, two product binaries: `cmd/server` and `cmd/smokery`).

---

## Status snapshot

| Phase | Description | Status |
|---|---|---|
| 0 | Repository foundation | ✓ |
| 1 | Backend project + DB foundation | ✓ |
| 2 | OpenAPI ingestion (basic) | ✓ |
| 3 | Canonical config model | ✓ |
| 4 | Compiler MVP (basic) | ✓ |
| 5 | Runner MVP (basic) | ✓ |
| 6 | Jobs + live run streaming (basic) | ✓ |
| 7 | Report MVP (basic) | ✓ |
| 8 | Frontend MVP (basic) | ✓ |
| 9 | Quality hardening (basic) | ✓ |
| **A** | **Architecture: ports & adapters refactor** | **✓** |
| **B** | **CLI smoke runner (`smokery`)** | **✓** |
| **C** | **Polish & MVP completion** | **In progress** |

---

## Phase 0 — Repository foundation ✓

- [x] Create monorepo structure.
- [x] Add `apps/core` Go module (one module for both binaries).
- [x] Add `apps/web` SvelteKit app.
- [x] Add `packages/types` placeholder.
- [x] Add Docker Compose with PostgreSQL and MinIO (under `configs/`).
- [x] Add root README and docs references.
- [x] Add basic GitHub Actions skeleton.
- [x] Add `Makefile` with all dev commands.
- [x] Add air hot-reload (`configs/.air.toml`).

---

## Phase 1 — Backend project and database foundation ✓

- [x] Add Go config loading with Viper.
- [x] Add Echo server with health route.
- [x] Add huma + humaecho for typed operations + OpenAPI 3.1 spec generation.
- [x] Add zerolog middleware.
- [x] Add pgx connection pool.
- [x] Add golang-migrate setup.
- [x] Add sqlc config.
- [x] Create initial migrations for projects, specs, environments, auth profiles, fixtures, flows, suites, plans, runs, results, comments, artifacts.
- [x] Add basic CRUD queries for projects.

---

## Phase 2 — OpenAPI ingestion ✓ (basic)

- [x] Implement `internal/spec` parser using libopenapi.
- [x] Add spec import API.
- [x] Store raw spec and analysis JSON.
- [x] Extract operation list.
- [x] Classify operations.
- [x] Add operation registry API.
- [x] Add operation override API.

Deferred to Phase C:

- [ ] Detect pagination/search/enum query hints from OpenAPI parameter schemas.

---

## Phase 3 — Canonical config model ✓

- [x] Define Go types for project config, flows, suites, plans, results.
- [x] Add domain entity types in `internal/model/entity.go` (Phase A).

---

## Phase 4 — Compiler MVP ✓ (basic)

- [x] Environment / auth / operation reference resolution.
- [x] Flow + suite compilation.
- [x] Default list + pagination case generation.
- [x] Destructive operation validation.

Deferred to Phase C (see C.1):

- [ ] Search-from-response case generation.
- [ ] Enum-filter case generation.
- [ ] Plan preview API (no DB writes, returns compile result).
- [ ] Inline JSON schemas into SmokePlan for runtime validation.

---

## Phase 5 — Runner MVP ✓ (basic, refactored in Phase A)

- [x] HTTP execution engine.
- [x] Variable interpolation, auth injection (via hooks).
- [x] JSONPath capture (via hook).
- [x] Status / list-shape / JSONPath / not_empty assertions.
- [x] Trace ID extraction (via hook).
- [x] Header redaction (via hook).
- [x] Pluggable hook architecture (Phase A).
- [x] Cleanup step execution.

Deferred to Phase C (see C.1, C.2):

- [ ] JSON Schema assertion (using inlined schemas).
- [ ] Empty-result policy (`allow | warn | fail`).
- [ ] Pagination sanity assertion.

---

## Phase 6 — Jobs and live run streaming ✓ (basic)

- [x] Goroutine-based worker (replaces River for MVP).
- [x] Run creation API + lifecycle.
- [x] WebSocket endpoint for run events.
- [x] EventBus port + inproc adapter (Phase A).

Deferred to Phase C (see C.2):

- [ ] Run cancellation API + worker support.
- [ ] Per-step / per-case events emitted from runner (currently only run_started / run_finished).

---

## Phase 7 — Report MVP ✓ (basic)

- [x] Backend debug report view.
- [x] CI summary report view.
- [x] Mermaid sequence diagram.
- [x] Run comments API.

Deferred to Phase C (see C.2):

- [ ] HTML report artifact generation.
- [ ] Persist report artifacts via `port.BlobStore`.
- [ ] JSON artifact upload to MinIO.

---

## Phase 8 — Frontend MVP ✓ (basic)

- [x] Project list, spec import, runs list, run detail with WebSocket events, comments.
- [x] All pages use `openapi-fetch` + `@tanstack/svelte-query`. Zero `any` types.

Deferred to Phase C (see C.3): operation explorer, environment/auth config, flow builder, list suite configurator, plan preview page.

---

## Phase 9 — Quality hardening ✓ (basic)

- [x] Backend unit tests for compiler, runner, assertion, classification.
- [x] Test HTTP server for runner integration (httptest).
- [x] Frontend typecheck in CI.

Deferred to Phase C (see C.4):

- [ ] Integration tests with PostgreSQL testcontainer.
- [ ] App service tests using memory adapters.
- [ ] Redaction tests.
- [ ] Retention cleanup job + tests.
- [ ] CLI smoke tests.

---

## Phase A — Architecture: ports & adapters refactor ✓

This phase aligned the code with `docs/architecture.md`.

- [x] Rename `apps/api` → `apps/core` (Go workspace, both binaries live here).
- [x] Add domain types in `internal/model/entity.go` using `uuid.UUID`, `time.Time` (no pgx).
- [x] Define `internal/port/` interfaces: `ProjectRepo`, `SpecRepo`, `OperationRepo`, `RunRepo`, `CommentRepo`, `ArtifactRepo`, `JobEnqueuer`, `EventBus`, `BlobStore`.
- [x] Move `internal/db` → `internal/adapter/postgres/db`. Wrote per-entity repo files translating sqlc types ↔ domain types.
- [x] Move `internal/jobs` → `internal/adapter/inproc/{worker.go,eventbus.go}` implementing the corresponding ports.
- [x] Add `internal/adapter/minio/blob.go` implementing `port.BlobStore`.
- [x] Refactor `internal/runner/` to introduce `Options`, `PreProcessor`, `PostProcessor`. Moved auth/redact/trace/capture/assertion to `internal/runner/hook/`.
- [x] Create `internal/app/` use cases depending only on ports.
- [x] Move huma handlers from `cmd/server/main.go` → `internal/delivery/http/`.
- [x] Refactor `cmd/server/main.go` to be dependency injection only.
- [x] Update `cmd/openapi/main.go` to reuse delivery handler registration via no-op port impls.
- [x] Verified `pgtype/pgxpool/jackc/pgx` only in `internal/adapter/postgres/`.
- [x] Verified `minio-go` only in `internal/adapter/minio/`.
- [x] Verified `gorilla/websocket` only in `internal/delivery/http/`.
- [x] All existing tests pass.
- [x] Frontend regenerated types still work.

---

## Phase B — CLI smoke runner ✓

- [x] Add `internal/adapter/memory/` (in-memory repos for CLI / tests).
- [x] Add `internal/adapter/fs/` (filesystem blob store).
- [x] Add `internal/cli/loader/` (YAML/JSON loaders via `sigs.k8s.io/yaml` to respect json tags).
- [x] Add `internal/delivery/cli/` cobra commands: `run`, `compile`, `import-spec`, `report`.
- [x] Add `cmd/smokery/main.go` wiring memory adapters + runner.
- [x] Add `make build-cli`, `make run-cli`, `make install-cli` targets.
- [x] End-to-end test against httpbin.org passes.

---

## Phase C — Polish & MVP completion

Tasks are grouped by area for parallel execution. Each section is independent of the others except where noted.

### C.1 — Compiler completeness

- [ ] **C.1.1** Detect query parameter hints from OpenAPI (pagination keys: `page`, `limit`, `offset`, `cursor`; search keys: `q`, `query`, `search`; enum filters from schema).
- [ ] **C.1.2** Inline relevant JSON schemas into the `SmokePlan` so the runner can validate without reading the spec.
- [ ] **C.1.3** Implement search-from-response case generation: compile a list call → capture → search call.
- [ ] **C.1.4** Implement enum-filter case generation: emit one case per enum value.
- [ ] **C.1.5** Implement empty-result policy field on suite strategy and pass it through to runner-side assertion.
- [ ] **C.1.6** Add plan preview API (`POST /api/projects/{id}/plan/preview` — compiles config without persisting). Wired into existing app service.
- [ ] **C.1.7** Compiler tests for each new case generator.

### C.2 — Runtime / orchestration

- [ ] **C.2.1** JSON Schema assertion (`type: schema`) using `santhosh-tekuri/jsonschema/v6`. Reads inlined schema from the step.
- [ ] **C.2.2** Pagination sanity assertion (`type: pagination` — confirms paginated list returns expected shape).
- [ ] **C.2.3** Empty-result policy enforcement at runtime (`allow | warn | fail`).
- [ ] **C.2.4** Per-step / per-case event emission from runner (`flow.step.started`, `suite.case.result`, etc.).
- [ ] **C.2.5** Run cancellation API: `POST /api/runs/{id}/cancel`. Worker context cancellation, status `cancelled`.
- [ ] **C.2.6** HTML report artifact generation (`internal/report/html.go`).
- [ ] **C.2.7** Persist artifacts via `port.BlobStore` (server uses minio adapter, CLI uses fs adapter).
- [ ] **C.2.8** Health check that probes DB + MinIO connectivity.
- [ ] **C.2.9** Retention cleanup job (delete old run results / artifacts past TTL).

### C.3 — Frontend builders

- [ ] **C.3.1** Operation explorer / override UI (`/projects/[id]/operations`).
- [ ] **C.3.2** Environment config pages (`/projects/[id]/environments`).
- [ ] **C.3.3** Auth profile config pages (under environments).
- [ ] **C.3.4** Flow builder UI (`/projects/[id]/flows/[fid]`).
- [ ] **C.3.5** List suite configurator UI (`/projects/[id]/suites/[sid]`).
- [ ] **C.3.6** Plan preview page showing compiler errors/warnings + generated cases (`/projects/[id]/plan`).
- [ ] **C.3.7** TailwindCSS + shadcn-svelte component library setup.
- [ ] **C.3.8** Mermaid rendering on the run detail page (currently shows raw text).
- [ ] **C.3.9** LayerChart-based latency / pass-rate trends on the runs list page.

### C.4 — Quality & reliability

- [ ] **C.4.1** App service tests using memory adapters (no testcontainers needed).
- [ ] **C.4.2** Integration tests with PostgreSQL testcontainer for the postgres adapter.
- [ ] **C.4.3** CLI command tests (`run`, `compile`, `import-spec`, `report`).
- [ ] **C.4.4** Redaction hook tests (sensitive headers, sensitive JSON paths).
- [ ] **C.4.5** Hook unit tests for each built-in hook (auth, interpolate, capture, etc.).
- [ ] **C.4.6** Consistent error model across delivery: `huma.ErrorModel` is fine, but ensure all errors include `path`, `code`, `message`.
- [ ] **C.4.7** Add `golangci-lint` config and run in CI.
- [ ] **C.4.8** Add minimum code coverage threshold to CI.

### C.5 — Documentation & examples

- [ ] **C.5.1** Add example project configs to `/examples/configs/`.
- [ ] **C.5.2** Add example smoke plans to `/examples/plans/`.
- [ ] **C.5.3** CLI README at `apps/core/cmd/smokery/README.md` with usage examples.
- [ ] **C.5.4** Architecture diagram (PNG/SVG) in `docs/architecture.md`.
- [ ] **C.5.5** Update `docs/tech-spec.html` to reference current architecture (currently outdated).
- [ ] **C.5.6** Add migration guide for the postgres schema (golang-migrate up/down examples).

### C.6 — Operational / deployment

- [ ] **C.6.1** Production Dockerfile for `cmd/server`.
- [ ] **C.6.2** Production Dockerfile for `cmd/smokery` (small static binary image).
- [ ] **C.6.3** GitHub Actions: build + push container images on tag.
- [ ] **C.6.4** GitHub Actions: build + publish CLI binaries for linux / macos / windows.
- [ ] **C.6.5** Sample Kubernetes manifests under `deploy/k8s/` (server + postgres + minio).
- [ ] **C.6.6** `.env.example` documenting all configurable environment variables.

### C.7 — Frontend developer experience

- [ ] **C.7.1** `make generate` should restart `air` after spec regeneration (or document the workflow).
- [ ] **C.7.2** Pre-commit hook to regenerate FE types when API changes.
- [ ] **C.7.3** Vitest setup for the SvelteKit app.
- [ ] **C.7.4** Add a few component tests for critical UI paths.

---

## MVP Completion Definition

The MVP is "done" when:

- [ ] Phase C.1 (compiler completeness) is finished.
- [ ] Phase C.2.1–C.2.7 (runtime + reports) is finished.
- [ ] Phase C.3.1–C.3.6 (core frontend builders) is finished.
- [ ] Phase C.4.1, C.4.2, C.4.4 (essential tests) are finished.
- [ ] Phase C.5.1, C.5.2 (examples) are finished.

C.6 and C.7 are nice-to-have polish that can ship post-MVP.

---

## After MVP

See `docs/ai/tasks/phase-2-task-list.md` for post-MVP roadmap (drift detection, RBAC, scheduled runs, AI-assisted features, etc.).
