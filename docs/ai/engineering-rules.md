# Engineering Rules for AI Agents

## 1. General coding rules

- Keep code explicit and readable.
- Prefer simple, testable functions.
- Avoid broad rewrites unrelated to the task.
- Preserve module boundaries.
- Add tests for compiler, runner, and assertion behavior.
- Keep API DTOs stable and documented.
- Never hide business rules in UI-only code.
- All Go binaries output to `tmp/` directory.
- Use `make` targets for all dev commands.

## 2. Backend layering rules

The backend follows hexagonal architecture. The full diagram is in `docs/architecture.md`.

### 2.1 Layer dependency direction (one-way)

```text
delivery ──► app ──► port ──┐                            ▲
                            │                            │
                            │                          implements
                            │                            │
                            ▼                            │
                          domain                       adapter
```

- `domain` (`model`, `spec`, `compiler`, `runner`, `assertion`, `report`) imports nothing from this project except other domain packages.
- `port` defines interfaces. Imports `model` only.
- `adapter` implements ports using infrastructure SDKs. Imports `port`, `model`, and the SDK.
- `app` is use-case orchestration. Imports `port`, `model`, and domain services. **Never** imports adapters.
- `delivery` is the entry layer (HTTP, CLI). Imports `app` and `model`.
- `cmd/*` is dependency-injection only. `main.go` is the only place where adapters meet app services.

### 2.2 Type translation rule

Domain types are pure Go (`uuid.UUID`, `time.Time`, primitive types). Persistence types (sqlc-generated structs with `pgtype.UUID`, `pgtype.Timestamptz`) live **only inside** `internal/adapter/postgres/`. The repo adapter is the only code that converts between them.

After every change, this command must return only files under `internal/adapter/postgres/`:

```bash
grep -rE "pgtype|pgxpool|jackc/pgx" apps/core/internal --include="*.go"
```

Same rule for:

- `minio-go` → only inside `internal/adapter/minio/`
- `gorilla/websocket` → only inside `internal/delivery/http/`
- `cobra` → only inside `internal/delivery/cli/` (when CLI is built)
- `libopenapi` → only inside `internal/spec/` and `internal/compiler/`

### 2.3 Adapter pattern

Every infrastructure interaction goes through a port. Examples:

| Concern | Port | Adapters |
|---|---|---|
| Repository (project, spec, run, …) | `port.ProjectRepo`, `port.SpecRepo`, etc. | `postgres`, `memory` |
| Job enqueueing | `port.JobEnqueuer` | `inproc`, `river` (future) |
| Event broadcasting | `port.EventBus` | `inproc`, `redis` (future) |
| Blob storage | `port.BlobStore` | `minio`, `fs`, `memory` |

App services accept ports as constructor dependencies. They never reference an adapter directly.

## 3. Runner rules

The runner is a **library**, not a service.

### 3.1 Pure runtime contract

The runner accepts only:

- A `*model.SmokePlan` (the AST emitted by the compiler)
- An optional `RunOptions` (HTTP client, hooks, event emitter)

It returns a `*model.RunResult`. It must not:

- Import `libopenapi` or any OpenAPI-spec parsing
- Import any database package
- Import any storage package
- Embed job orchestration

### 3.2 Hooks

All extension is via hooks:

```go
type PreProcessor interface {
    BeforeRequest(rctx *RequestContext) error
}

type PostProcessor interface {
    AfterResponse(rctx *ResponseContext, step *model.StepResult) error
}
```

Built-in hooks (in `internal/runner/hook`):

- `VariableInterpolator` (pre)
- `AuthInjector` (pre)
- `EnvironmentHeaders` (pre)
- `Capture` (post)
- `Redactor` (post)
- `TraceExtractor` (post)
- `AssertionRunner` (post)

`runner.DefaultOptions()` composes these. Custom hooks can be added without modifying the runner core.

### 3.3 Compiler bakes everything

Anything the runner needs at runtime — response schemas, content types, param formats — must be embedded in the SmokePlan by the compiler. The runner is OpenAPI-agnostic.

## 4. Compiler rules

Compiler output must be deterministic.

Compiler errors must include:

```text
stage
path
message
severity
related entity if available
```

Compiler must validate:

- Missing environment/base URL
- Missing auth/secret
- Invalid operation reference
- Required path/query/header params
- Required request body fields when detectable
- Destructive operation permissions
- Suite selector output
- Generated case limits
- Response shape config

The compiler **must** embed any runtime-required schema or metadata into the SmokePlan. If a step needs to validate against a JSON schema, that schema is inlined in the plan, not looked up by the runner.

## 5. Frontend rules

- Use SvelteKit 2 with Svelte 5 runes and TypeScript strict mode.
- Use `openapi-fetch` for type-safe API calls (types from `openapi-typescript`).
- Use `@tanstack/svelte-query` v6 for server state (accessor pattern: `createQuery(() => ({...}))`).
- Access query results directly (no `$` store prefix in Svelte 5).
- Zero `any` types — all types flow from the backend OpenAPI spec.
- Use `bun` as package manager and runtime.

Future additions (add when needed):

- TailwindCSS + shadcn-svelte for UI components.
- Monaco for raw JSON/YAML config editing.
- Mermaid for diagrams.
- LayerChart for run trends and latency.

## 6. Test rules

Backend minimum tests:

- Operation classification
- Compiler success/error cases
- Destructive operation blocking
- Variable interpolation hook
- JSONPath capture hook
- Status / list-shape assertions
- App service tests using memory adapters (no testcontainers needed)

Frontend minimum checks:

- `bun run check` (svelte-check) must pass with 0 errors
- Zero `any` types in handwritten code
