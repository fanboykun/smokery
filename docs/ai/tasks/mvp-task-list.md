# MVP Task List — OpenAPI Smoke Testing Platform

This task list is designed for AI agents to implement the project incrementally.

Do not implement later phases before MVP foundations are stable.

The Go workspace is `apps/core` (one Go module, two product binaries: `cmd/server` and `cmd/smokery`).

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

---

## Phase 3 — Canonical config model ✓ (basic)

- [x] Define Go types for project config, flows, suites, plans, results.

---

## Phase 4 — Compiler MVP ✓ (basic)

- [x] Environment / auth / operation reference resolution.
- [x] Flow + suite compilation.
- [x] Default list + pagination case generation.
- [x] Destructive operation validation.

---

## Phase 5 — Runner MVP ✓ (basic, will be refactored in Phase A)

- [x] HTTP execution engine.
- [x] Variable interpolation, auth injection, JSONPath capture.
- [x] Status / list-shape / JSONPath / not_empty assertions.
- [x] Trace ID extraction, header redaction.

---

## Phase 6 — Jobs and live run streaming ✓ (basic)

- [x] Goroutine-based worker (River replacement for MVP).
- [x] Run creation API + lifecycle.
- [x] WebSocket endpoint for run events.

---

## Phase 7 — Report MVP ✓ (basic)

- [x] Backend debug report view.
- [x] CI summary report view.
- [x] Mermaid sequence diagram.
- [x] Run comments API.

---

## Phase 8 — Frontend MVP ✓ (basic)

- [x] Project list, spec import, runs list, run detail with WebSocket events, comments.
- [x] All pages use `openapi-fetch` + `@tanstack/svelte-query`. Zero `any` types.

---

## Phase 9 — Quality hardening ✓ (basic)

- [x] Backend unit tests for compiler, runner, assertion, classification.

---

## Phase A — Architecture: ports & adapters refactor

This phase aligns the code with `docs/architecture.md`.

- [ ] Rename `apps/api` → `apps/core` (Go workspace, both binaries live here).
- [ ] Add domain types in `internal/model/{project,spec,run,operation,comment,artifact}.go` using `uuid.UUID`, `time.Time` (no pgx).
- [ ] Define `internal/port/` interfaces: `ProjectRepo`, `SpecRepo`, `OperationRepo`, `RunRepo`, `CommentRepo`, `ArtifactRepo`, `JobEnqueuer`, `EventBus`, `BlobStore`.
- [ ] Move `internal/db` → `internal/adapter/postgres/db`. Write `internal/adapter/postgres/repo.go` translating sqlc types ↔ domain types.
- [ ] Move `internal/jobs` → `internal/adapter/inproc/{worker.go,eventbus.go}` implementing the corresponding ports.
- [ ] Add `internal/adapter/minio/blob.go` implementing `port.BlobStore`.
- [ ] Refactor `internal/runner/` to introduce `Options`, `PreProcessor`, `PostProcessor`. Move auth/redact/trace/capture/assertion to `internal/runner/hook/`.
- [ ] Create `internal/app/` use cases depending only on ports.
- [ ] Move huma handlers from `cmd/server/main.go` → `internal/delivery/http/{project,spec,run,report,comment,websocket}.go`.
- [ ] Refactor `cmd/server/main.go` to be dependency injection only.
- [ ] Update `cmd/openapi/main.go` to reuse delivery handler registration.
- [ ] Verify `grep -rE "pgtype|pgxpool|jackc/pgx" internal --include="*.go"` returns only `internal/adapter/postgres/`.
- [ ] Verify `grep -r "minio-go" internal --include="*.go"` returns only `internal/adapter/minio/`.
- [ ] All existing tests pass.
- [ ] Frontend regenerated types still work (UUID strings should be cleaner).

Acceptance:

- API behavior is identical.
- Domain code has zero infrastructure imports.
- App services can be tested with memory adapters (no testcontainers).
- The architecture is ready for the CLI binary.

---

## Phase B — CLI smoke runner

- [ ] Add `internal/adapter/memory/` (in-memory repos for CLI / tests).
- [ ] Add `internal/adapter/fs/` (filesystem blob store).
- [ ] Add `internal/delivery/cli/` cobra commands: `run`, `compile`, `import-spec`, `report`.
- [ ] Add `cmd/smokery/main.go` wiring memory + fs adapters.
- [ ] Add `make build-cli` and `make run-cli` targets.
- [ ] Add YAML loaders for `ProjectConfig` and `SmokePlan`.

Acceptance:

- `smokery run plan.yaml` works without any DB or HTTP server.
- Same compiler / runner / app services as the API.

---

## Phase C — Polish (parallel to Phase A/B as needed)

- [ ] JSON Schema assertion (using inlined schemas from compiler).
- [ ] Empty-result policy (`allow | warn | fail`).
- [ ] Search-from-response and enum-filter case generation.
- [ ] Cancellation API for runs.
- [ ] HTML report artifact + MinIO storage.
- [ ] Retention cleanup job.
- [ ] Frontend operation explorer / override UI.
- [ ] Frontend environment/auth config pages.
- [ ] Frontend flow builder.
- [ ] Frontend list suite configurator.
- [ ] Frontend plan preview page.
- [ ] Integration tests with PostgreSQL testcontainer.
- [ ] Redaction tests.

Acceptance:

- Full MVP user journey works end-to-end via the UI.
- CLI covers run + compile + spec-import use cases.
- CI runs all tests.

