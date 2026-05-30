# Interactive Canvas Configuration Builder Task Tracker

## Current Status

Implementation reached the first usable canvas-builder slice on 2026-05-30.

The builder now uses a graph authoring model and derives the existing canonical `ProjectConfig` before plan preview/run. Compiler and runner semantics remain unchanged: the canvas is metadata and authoring state, while `flows[]` and `suites[]` are still the execution inputs.

## Completed Changes

- Added backend canvas graph model types on `ProjectConfig.canvas`.
- Extended spec operation metadata with request schema, response schema, and parameters for diagram ports.
- Added a canvas operation metadata use case and HTTP endpoint:
  - `GET /api/specs/{spec-id}/operations/canvas`
- Installed `@xyflow/svelte` in the web app.
- Added frontend canvas graph and operation metadata types.
- Added graph-to-config adapter that converts:
  - operation nodes to flow steps
  - sequence edges to flow order
  - data-link edges to captures and target bindings
  - suite generator nodes to generated suite config
- Reworked `/projects/[id]/builder` into the primary Xyflow canvas page.
- Added operation and suite generator canvas node components.
- Kept plan preview/run behind compiler diagnostics.
- Updated existing plan/run pages to match current generated API response shapes.

## Implementation Checklist

- [x] Backend: add canvas graph model types.
- [x] Backend: extend `spec.OperationInfo` with request schema and parameters.
- [x] Backend: extract OpenAPI request body schema.
- [x] Backend: extract path/query/header/cookie parameters.
- [x] Backend: add canvas operation metadata endpoint.
- [x] Frontend: add `@xyflow/svelte`.
- [x] Frontend: add canvas graph types.
- [x] Frontend: add graph-to-config adapter.
- [x] Frontend: replace `/projects/[id]/builder` with canvas.
- [x] Frontend: keep plan preview panel and disable run on errors.
- [x] Validation: run focused Go tests.
- [x] Validation: run focused frontend tests.
- [x] Validation: run full Go tests.
- [ ] Validation: full `bun --filter @smokery/web check` passes.

## Validation Results

Passed:

- `env GOCACHE=/tmp/go-build GOMODCACHE=/tmp/go-mod /usr/local/go/bin/go test ./...`
- `/home/fanboykun/.bun/bin/bun --filter @smokery/web test -- src/lib/canvas/graph-to-config.test.ts`

Not passing yet:

- `/home/fanboykun/.bun/bin/bun --filter @smokery/web check`

Current `svelte-check` blockers are existing shared UI primitive type errors in:

- `apps/web/src/lib/components/ui/calendar/calendar.svelte`
- `apps/web/src/lib/components/ui/sidebar/sidebar-menu-button.svelte`
- `apps/web/src/lib/components/ui/sidebar/sidebar-menu-item.svelte`
- `apps/web/src/lib/components/ui/sidebar/sidebar-menu.svelte`
- `apps/web/src/lib/components/ui/sidebar/sidebar-rail.svelte`
- `apps/web/src/lib/components/ui/slider/slider.svelte`

The previous canvas/plan/run route type errors are cleared.

## Latest Context

The user explicitly requested that all plans, descriptions, designs, implementation steps, task tracking, and latest context for this work live under `docs/` and `docs/tasks/`.

This tracker is the latest implementation context for the interactive canvas work. The next useful slice is to fix the shared UI primitive type errors, regenerate OpenAPI/TypeScript types so the new canvas endpoint and `ProjectConfig.canvas` are reflected in generated web API types, then do a manual browser smoke of list users -> delete user with `data[].id`.
