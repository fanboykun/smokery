# Workflow 03 — Runner Execution

## Goal

Execute compiled SmokePlans safely and return structured RunResults.

The runner is a **library**, not a service. It accepts a `*model.SmokePlan`, runs it via a chain of pluggable hooks, and returns a `*model.RunResult`. The same runner powers the API server and the CLI binary.

## Contract

```go
type Runner struct {
    HTTPClient     *http.Client      // pluggable transport
    PreProcessors  []PreProcessor    // request-side hooks
    PostProcessors []PostProcessor   // response-side hooks
    EventEmitter   func(Event)       // optional, for live streaming
}

func New(opts RunOptions) *Runner
func DefaultOptions() RunOptions
func (r *Runner) Execute(ctx context.Context, plan *model.SmokePlan) *model.RunResult
```

## What the runner is allowed to import

- Standard library (`net/http`, `context`, `time`, etc.)
- `internal/model` (the AST it executes)
- `internal/runner/hook` (built-in hooks)
- `internal/assertion` (used by the AssertionRunner hook)

## What the runner must NOT import

- `libopenapi` or any OpenAPI parser
- `internal/spec`, `internal/compiler` (no compile-time knowledge at runtime)
- Any database package (`pgx`, sqlc-generated code, repositories)
- Any storage package (`minio-go`, filesystem helpers)
- Any HTTP framework (`echo`, `huma`, `gorilla/websocket`)

## Hook system

```go
type RequestContext struct {
    Plan    *model.SmokePlan
    Step    *model.PlannedStep
    Vars    map[string]any
    Request *http.Request
}

type ResponseContext struct {
    Plan     *model.SmokePlan
    Step     *model.PlannedStep
    Vars     map[string]any
    Request  *http.Request
    Response *http.Response
    Body     []byte
}

type PreProcessor interface {
    BeforeRequest(rctx *RequestContext) error
}

type PostProcessor interface {
    AfterResponse(rctx *ResponseContext, step *model.StepResult) error
}
```

### Execution order

For each step:

1. Build `*http.Request` from `PlannedStep`.
2. Run all `PreProcessor.BeforeRequest` in order.
3. `HTTPClient.Do(req)`.
4. Read response body.
5. Run all `PostProcessor.AfterResponse` in order.
6. Emit `step.result` event (if `EventEmitter` set).

### Built-in hooks (`internal/runner/hook`)

| Hook | Phase | Purpose |
|---|---|---|
| `VariableInterpolator` | pre | Fill `{vars}` in path, headers, body |
| `AuthInjector` | pre | Inject bearer/basic/api-key from `plan.Auth` |
| `EnvironmentHeaders` | pre | Apply default headers from environment |
| `Capture` | post | Extract values via gjson into `Vars` |
| `Redactor` | post | Redact sensitive headers / body fields in result |
| `TraceExtractor` | post | Pull `X-Trace-Id` / `X-Request-Id` |
| `AssertionRunner` | post | Run all step assertions, populate `step.Assertions` |

`runner.DefaultOptions()` composes these hooks in the order shown.

## Execution stages

```text
1. Start run
2. Emit run.started
3. Execute setup/auth flows if needed
4. Execute configured flows (steps in order, hook chain per step)
5. Execute configured suites (cases per operation group)
6. Attempt cleanup steps
7. Finalize summary
8. Emit run.finished
```

Persistence is **not** a runner responsibility. The caller (app service or worker) decides what to do with the returned `*model.RunResult` — save it, print it, stream it.

## Suite behavior

- Operation groups may run concurrently.
- Cases inside a group respect dependency ordering.
- Search-from-response depends on default list case.
- Pagination and enum cases may be independent.

## Run events (when `EventEmitter` is set)

```text
run.started
flow.started
flow.step.started
flow.step.result
suite.started
suite.operation.started
suite.case.started
suite.case.result
cleanup.started
cleanup.result
run.finished
run.cancelled
run.failed
```

When `EventEmitter` is nil (e.g. in CLI), no events are emitted.

## Wiring examples

### CLI (no DB, no events)

```go
plan := loadPlanFromYAML(path)
result := runner.New(runner.DefaultOptions()).Execute(ctx, plan)
report.PrintCISummary(os.Stdout, result)
```

### API server (with event streaming + persistence)

```go
opts := runner.DefaultOptions()
opts.EventEmitter = eventBus.Emit
result := runner.New(opts).Execute(ctx, plan)
runRepo.SaveResult(ctx, runID, result)
```

### Custom integration (AWS SigV4 signing)

```go
opts := runner.DefaultOptions()
opts.PreProcessors = append(opts.PreProcessors, &SigV4Signer{...})
result := runner.New(opts).Execute(ctx, plan)
```

## Acceptance criteria

- Runner can execute a compiled flow.
- Runner can execute a generated list suite.
- Runner streams live progress when `EventEmitter` is set.
- Runner extracts trace IDs/request IDs (TraceExtractor hook).
- Runner redacts sensitive values (Redactor hook).
- Runner is dependency-free: no DB, no storage, no OpenAPI parser imports.
- Same runner runs in the API server and the CLI without code changes.
