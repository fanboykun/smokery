# Smokery - Unified Progress & Roadmap

## Current State

**MVP (Phases 0-9, A-C)**: ✅ Complete  
**Phase 2 FE Tier 1 (Reporting UI)**: ✅ Complete (mock data)  
**Phase 2 FE Tier 2 (Failure Classification UI)**: ✅ Complete (mock data)  
**Phase 2 FE Tier 3 (Spec Evolution UI)**: ✅ Complete (mock data)  
**Phase 2 FE Tier 4 (Analytics UI)**: ✅ Complete (mock data)  
**Phase 2 FE Tier 5 (Governance UI)**: ✅ Complete (mock data)  
**Backend Phase 2 endpoints**: ✅ All tiers done (reports, classification, spec evolution, analytics, governance)

---

## FE ↔ BE Dependency Matrix

The FE has built UI against mock data. The BE needs to implement these endpoints to make them real.

### Tier 1: Reporting — FE ✅ Done, BE ✅ Done

| FE Page | BE Endpoint | Status |
|---------|-------------|--------|
| `/runs/[runId]/report/contract` | `GET /api/runs/{id}/report/contract` | ✅ |
| `/runs/[runId]/report/analyst` | `GET /api/runs/{id}/report/analyst` | ✅ |
| `/runs/[runId]/report/qa` | `GET /api/runs/{id}/report/qa` | ✅ |
| `/runs/[runId]/report/correlation` | `GET /api/runs/{id}/correlations` | BE: TODO |

**BE implementation**: Add `internal/report/contract.go`, `analyst.go`, `qa.go`. Transform `RunResult` into view-specific structs. Wire to delivery layer.

### Tier 2: Failure Classification — FE ✅ Done, BE ✅ Done

| FE Component/Page | BE Endpoint | Status |
|-------------------|-------------|--------|
| `FailureClassifier.svelte` | `PUT /api/runs/{id}/failure-classification` | ✅ |
| `AssigneeSelector.svelte` | `GET /api/team/members` | BE: TODO (use project members) |
| `FailureTimeline.svelte` | `GET /api/runs/{id}/failure-classification` | ✅ |
| Run detail classification card | `GET /api/projects/{project-id}/failure-classifications` | ✅ |

**BE implementation**: Add `failure_classifications` table (run_id, classification, assignee, note, author, created_at). CRUD API.

### Tier 3: Spec Evolution — FE ✅ Done, BE ✅ Done

| FE Page | BE Endpoint |
|---------|-------------|
| `/projects/[id]/spec/versions` | `GET /api/projects/{project-id}/specs` (existing) |
| `/projects/[id]/spec/diff` | `GET /api/specs/{from-id}/diff/{to-id}` |
| `/projects/[id]/impact` | `GET /api/projects/{project-id}/impact/spec/{spec-id}` |

### Tier 4: Analytics — FE ✅ Done, BE ✅ Done

| FE Page | BE Endpoint |
|---------|-------------|
| `/projects/[id]/analytics/latency` | `GET /api/projects/{project-id}/analytics/latency?range=7d` |
| `/projects/[id]/analytics/flaky` | `GET /api/projects/{project-id}/analytics/flaky-operations` |
| `/projects/[id]/analytics/trends` | `GET /api/projects/{project-id}/analytics/health-trends` |

### Tier 5: Governance — FE ✅ Done, BE ✅ Done

| FE Page | BE Endpoint |
|---------|-------------|
| `/projects/[id]/settings/members` | `GET/POST /api/projects/{project-id}/members` |
| `/projects/[id]/settings/audit` | `GET /api/projects/{project-id}/audit-log` |
| `/projects/[id]/settings/webhooks` | `GET/POST /api/projects/{project-id}/webhooks` |
| `/projects/[id]/settings/notifications` | `GET/POST /api/projects/{project-id}/notifications` |

---

## BE Priority 1: Critical Enhancements — ✅ Already Implemented

These were already in the codebase:

| # | Enhancement | Status |
|---|-------------|--------|
| 1 | `GET /api/projects` — per-project stats (ProjectWithStats) | ✅ |
| 2 | `POST /api/projects/{project-id}/plan/preview` — structured error codes | ✅ |
| 3 | `GET /api/runs/{id}` — dynamic `websocket_url` from Host headers | ✅ |
| 4 | `GET /api/runs/{id}/result` — full RunResult with steps/assertions | ✅ |

---

## Implementation Order

### FE Phase 2: ✅ COMPLETE (all tiers built with mock data)

All UI pages are built and type-safe. They use `mockGet*` functions from `mock-phase2.ts`.  
When BE implements an endpoint, swap the mock call for a real `api.GET(...)` call.

### BE Phase 2 Tier 1-2: ✅ COMPLETE

New endpoints implemented:
- `GET /api/runs/{id}/report/contract` — contract compliance view
- `GET /api/runs/{id}/report/analyst` — root cause analysis
- `GET /api/runs/{id}/report/qa` — QA summary
- `PUT /api/runs/{id}/failure-classification` — classify failures
- `GET /api/runs/{id}/failure-classification` — get classification
- `GET /api/projects/{project-id}/failure-classifications` — list all

### BE Phase 2 Tier 3-5: ✅ COMPLETE

- `GET /api/specs/{from-id}/diff/{to-id}` — spec diff
- `GET /api/projects/{project-id}/impact/spec/{spec-id}` — impact analysis
- `GET /api/projects/{project-id}/analytics/latency` — latency trends
- `GET /api/projects/{project-id}/analytics/flaky-operations` — flaky detection
- `GET /api/projects/{project-id}/analytics/health-trends` — health over time
- `GET/POST /api/projects/{project-id}/members` — team members
- `GET /api/projects/{project-id}/audit-log` — audit trail
- `GET/POST /api/projects/{project-id}/webhooks` — webhook config
- `GET/POST /api/projects/{project-id}/notifications` — alert rules

### Next Steps:

1. **FE integration** — swap `mockGet*()` calls for real `api.GET(...)` calls
2. **Regenerate OpenAPI types** — run `bun run generate` to update `v1.d.ts`
3. **Add DELETE endpoints** for webhooks/notifications/members
4. **Persist governance data** — move from in-memory to SQLite tables

---

## Type Contracts

All FE ↔ BE type contracts are defined in:
- **FE types**: `apps/web/src/lib/api/phase2.ts`
- **FE mocks**: `apps/web/src/lib/api/mock-phase2.ts`
- **API schema**: `apps/web/src/lib/api/v1.d.ts` (generated from BE OpenAPI spec)

When BE implements an endpoint, FE replaces the `mockGet*` call with a real `api.GET(...)` call.
