# Current Work Context

## Active Work

Interactive Canvas Configuration Builder.

## User Direction

Always keep:

- design and implementation plans in `docs/`
- task tracking in `docs/tasks/`
- latest context for current work in `docs/tasks/`

## Current Implementation State

The first canvas-builder slice is implemented:

- backend model, parser, app service, and HTTP endpoint support diagram-ready operation metadata
- `ProjectConfig` has optional `canvas` metadata
- web app has `@xyflow/svelte`, canvas graph types, graph-to-config conversion, operation nodes, suite generator nodes, and a canvas-first `/projects/[id]/builder`
- plan preview/run still derive canonical `flows[]` and `suites[]` and go through compiler validation

## Latest Validation

Passed:

- `env GOCACHE=/tmp/go-build GOMODCACHE=/tmp/go-mod /usr/local/go/bin/go test ./...`
- `/home/fanboykun/.bun/bin/bun --filter @smokery/web test -- src/lib/canvas/graph-to-config.test.ts`

Still failing:

- `/home/fanboykun/.bun/bin/bun --filter @smokery/web check`

The remaining check failures are in existing shared UI primitives: calendar, sidebar, and slider components. Canvas, plan preview, and run-result route type errors have been fixed.

## Immediate Next Steps

1. Fix shared UI primitive type errors so `bun --filter @smokery/web check` can pass.
2. Regenerate OpenAPI/TypeScript types so the new canvas endpoint and `ProjectConfig.canvas` are reflected in generated API declarations.
3. Run a manual browser smoke of list users -> delete user with `data[].id`, compile preview, and verify destructive acknowledgement blocks run readiness.
