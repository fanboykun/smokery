---
type: task-note
status: current
tags:
  - agent
  - tasks
  - config-builder
  - persistence
related:
  - "[[Project State]]"
  - "[[MVP Task List]]"
  - "[[ADR - Config Builder Redesign]]"
  - "[[Config Persistence Blockers]]"
  - "[[Canvas Schema and Property Linking]]"
---

# Config Redesign & Persistence Task List

> [!summary]
> This note lists the concrete engineering tasks required to implement the Config Builder Redesign, configuration persistence in the SQLite/Postgres databases, and deep schema property linking in SvelteFlow. Each task is mapped back to its corresponding ADR and Technical Note.

---

## 1. Database & Repository Layer Tasks

### 1.1 Schema Extension (SQLite & Postgres)
* **Goal:** Store configuration details persistently on the server.
* **Tasks:**
  - [ ] Add `environments`, `auth_profiles`, `flows`, and `suites` tables initialization in local SQLite adapter ([sqlite.go](file:///home/fanboykun/dev/personal/smokery/apps/core/internal/adapter/sqlite/sqlite.go)).
  - [ ] Create a new Postgres/SQLite migration to add `canvas` JSONB column to the `projects` table (stores the high-level Plan canvas state).
  - [ ] Update SQLite table definitions for `projects` to include the `canvas` column.
* **Links:** `[[Config Persistence Blockers]]` (Step 1)

### 1.2 Repository Interfaces and Implementations
* **Goal:** Enable loading and saving of configuration payloads.
* **Tasks:**
  - [ ] Define `ConfigRepo` port interface in [repo.go](file:///home/fanboykun/dev/personal/smokery/apps/core/internal/port/repo.go).
  - [ ] Implement `ConfigRepo` in SQLite adapter ([sqlite.go](file:///home/fanboykun/dev/personal/smokery/apps/core/internal/adapter/sqlite/sqlite.go)).
  - [ ] Implement `ConfigRepo` in Postgres adapter (create `apps/core/internal/adapter/postgres/config_repo.go`).
  - [ ] Register the new config repository dependency injection in main server entry point ([server/main.go](file:///home/fanboykun/dev/personal/smokery/apps/core/cmd/server/main.go)).
* **Links:** `[[Config Persistence Blockers]]` (Step 2)

---

## 2. API & Application Layer Tasks

### 2.1 Configuration HTTP Endpoints
* **Goal:** Create HTTP CRUD routes for config syncing.
* **Tasks:**
  - [ ] Create HTTP request/response DTO structures in [delivery/http/types.go](file:///home/fanboykun/dev/personal/smokery/apps/core/internal/delivery/http/types.go).
  - [ ] Register `GET /api/projects/{project-id}/config` to load project environments, auth, flows, suites, and plan canvas state.
  - [ ] Register `PUT /api/projects/{project-id}/config` to save/update the config.
  - [ ] Register individual HTTP CRUD routes for plans/flows if needed, or stick to unified project-level updates.
* **Links:** `[[Config Persistence Blockers]]` (Step 3)

### 2.2 Compiler Server-Side Execution
* **Goal:** Fetch configuration from DB and compile on the server when running.
* **Tasks:**
  - [ ] Refactor `POST /api/projects/{project-id}/runs` HTTP handler in [delivery/http/run.go](file:///home/fanboykun/dev/personal/smokery/apps/core/internal/delivery/http/run.go) to fetch config from repository, compile dynamically on the backend, and create the run execution.
  - [ ] Refactor `POST /api/projects/{project-id}/plan/preview` handler in [delivery/http/plan.go](file:///home/fanboykun/dev/personal/smokery/apps/core/internal/delivery/http/plan.go) to load from DB or preview from request payload.
* **Links:** `[[ADR - Config Builder Redesign]]`, `[[Config Persistence Blockers]]` (Step 3)

---

## 3. Frontend & Canvas Layer Tasks

### 3.1 Hierarchical Canvas UI (SvelteFlow)
* **Goal:** Redesign builder view from single flat canvas to Plan Canvas + Nested Flow/Suite canvas.
* **Tasks:**
  - [ ] Add breadcrumbs navigation at top of [builder/+page.svelte](file:///home/fanboykun/dev/personal/smokery/apps/web/src/routes/projects/[id]/builder/+page.svelte) (e.g. `Project > Plan Canvas > Flow: User CRUD`).
  - [ ] Create `flowPlanNode` and `suitePlanNode` custom SvelteFlow components under `src/lib/components/canvas/`.
  - [ ] Add mode state (`currentMode`: `"plan-canvas" | "flow-canvas" | "suite-builder"`) to handle zoom and context drill-down.
  - [ ] Implement Plan Canvas sorting logic to order execution plans when connected.
* **Links:** `[[ADR - Config Builder Redesign]]` (Step 1, Step 2)

### 3.2 OpenAPI Spec Detail Panel & Property Pinning
* **Goal:** View OpenAPI spec inside builder and connect deeply nested parameters.
* **Tasks:**
  - [ ] Build a slide-out right-hand drawer in `builder/+page.svelte` showing selected operation spec, parameter tables, and schemas.
  - [ ] Add "Pin to Node" icons next to spec properties to dynamically append Handles to `OperationCanvasNode.svelte`.
  - [ ] Style Sequence Handles differently (thick, dashed lines) from Data-Link handles (thin paths).
  - [ ] Enforce connection validation in SvelteFlow (`isValidConnection` callback) based on handle prefixes.
* **Links:** `[[ADR - Config Builder Redesign]]` (Step 3, Step 4, Step 5), `[[Canvas Schema and Property Linking]]`
