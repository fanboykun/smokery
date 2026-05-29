# Session Context — Smokery Project

## What was done this session

### Frontend MVP (C.3.1-C.3.9, C.7.1-C.7.4) — ALL COMPLETE
- TailwindCSS 4 + shadcn-svelte component library
- Operations page wired to real API
- Environment & auth profile config (localStorage-backed, singleton store)
- Flow builder, Suite configurator, Plan preview with Run button
- Mermaid rendering, LayerChart trends, Vitest setup
- Pre-commit hook, breadcrumb navigation

### Backend
- SQLite adapter (pure Go, no CGO) at `internal/adapter/sqlite/`
- Adapter selection: `DB_ADAPTER=sqlite|postgres`, `STORAGE_ADAPTER=fs|minio`
- Default is sqlite+fs (zero external deps)
- Added `GET /api/projects/{project-id}/specs` and `POST /api/projects/{project-id}/specs/from-url`
- Fixed: circular $ref handling in spec parser (graceful, not fatal)
- Fixed: auth profile not resolved into SmokePlan by compiler
- Fixed: path selector wildcard matching (`/api/admin*` now works as prefix)
- Added request logging middleware (method, uri, status, latency)
- Added step-level logging in runner (method, url, status, duration per step)

### CI/CD
- Combined CI + Release workflow (no workflow_run, uses `needs`)
- Manual release workflow with bump type selector
- `ietf-tools/semver-action` for auto-versioning
- golangci-lint v2 config fixed, cache-dependency-path fixed

### Ops
- Dockerfiles in `configs/` (server + CLI)
- `v0.0.0` tag created on origin for semver baseline

## Known bugs still open
1. **Suite selector persistence** — store is now singleton, but the UI input for paths may still have binding issues with shadcn Input component (switched to plain `<input>` as workaround)
2. **Run detail Step Results** — result is base64 encoded, fixed with `atob()` decode. Verify it works.
3. **Path params not substituted** — URLs show literal `%7Bid%7D` instead of real IDs. The runner doesn't interpolate path params from the spec.

## Remaining UI/UX tasks (from audit)
- [ ] #3. Move Delete button to overflow menu with typed confirmation
- [ ] #4. Fix Plan Preview: collapse duplicate warnings into groups
- [ ] #5. Fix Plan Preview: add loading/progress state to Compile and Run buttons
- [ ] #6. Fix Run Detail: responsive layout (no horizontal overflow)
- [ ] #7. Fix Run Detail: group failures by status code, add filter
- [ ] #8. Fix Run Detail: expand Step Results height and add collapse toggle
- [ ] #9. Clean up Project Overview: move import to modal, compact config summary
- [ ] #10. Operations: add alternating rows, sticky tag navigation, auto-scroll detail panel

## Key file locations
- Frontend: `apps/web/src/`
- Backend: `apps/core/`
- Config store: `apps/web/src/lib/stores/project-config.ts` (singleton per project)
- Runner: `apps/core/internal/runner/runner.go`
- Compiler: `apps/core/internal/compiler/compiler.go`
- SQLite adapter: `apps/core/internal/adapter/sqlite/sqlite.go`
- Server main: `apps/core/cmd/server/main.go`
- CI workflow: `.github/workflows/ci.yml`
- Release workflow: `.github/workflows/release.yml`
- Phase 2 guide: `docs/ai/tasks/phase-2-task-list.md`

## Tech stack
- Go 1.26, Echo, Huma, SQLite (modernc.org/sqlite), zerolog
- SvelteKit, Svelte 5 (runes mode), TailwindCSS 4, shadcn-svelte, TanStack Query
- openapi-fetch + openapi-typescript for typed API client
