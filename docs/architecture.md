# Architecture — Smokery

> Hexagonal layering · Two delivery products · Pure domain · Pluggable runner

---

## 1. Two products, one core

Smokery ships as two distinct products that share the same domain code:

| Product | Delivery | Storage | Use case |
|---|---|---|---|
| **Web-backed API** | `cmd/server` (huma + echo) | PostgreSQL + MinIO | Collaborative platform with UI, persistence, live runs, reports |
| **CLI smoke runner** | `cmd/smokery` (cobra) | filesystem + memory | Standalone execution from local config, CI-friendly |

Both products import the same `internal/{model,spec,compiler,runner,assertion,report,port,app}` packages. Only the wiring (`main.go`) and the delivery layer (`internal/delivery/{http,cli}`) differ.

The CLI is **not** a stripped-down API. It is an equally-valid front to the same compiler and runner. Any feature that lives in the domain or app layer is automatically available to both products.

---

## 2. Layering (hexagonal / ports & adapters)

```text
┌─────────────────────────────────────────────────────────────┐
│  Delivery                                                   │
│   internal/delivery/http   internal/delivery/cli            │
│   apps/core/cmd/server      apps/core/cmd/smokery           │
└────────────────────────┬────────────────────────────────────┘
                         │ depends on
                         ▼
┌─────────────────────────────────────────────────────────────┐
│  Application services (use cases)                           │
│   internal/app                                              │
│     ProjectService, SpecService, RunService,                │
│     OperationService, ReportService, CommentService         │
└────────────────────────┬────────────────────────────────────┘
                         │ depends on
                         ▼
┌─────────────────────────────────────────────────────────────┐
│  Ports (interfaces)        ◀── implemented by ──┐           │
│   internal/port                                 │           │
│     ProjectRepo, SpecRepo, RunRepo, …           │           │
│     JobEnqueuer, EventBus, BlobStore            │           │
└────────────────────────┬────────────────────────┼───────────┘
                         │                        │
                         ▼                        │
┌────────────────────────────────────┐  ┌─────────┴───────────┐
│  Domain (pure Go, no infra)        │  │  Adapters           │
│   internal/model                   │  │  (infrastructure)   │
│   internal/spec                    │  │   internal/adapter/ │
│   internal/compiler                │  │     postgres/       │
│   internal/runner                  │  │     inproc/         │
│   internal/runner/hook             │  │     minio/          │
│   internal/assertion               │  │     fs/      [CLI]  │
│   internal/report                  │  │     memory/  [CLI]  │
└────────────────────────────────────┘  └─────────────────────┘
```

### Dependency direction (strict)

- `domain` imports nothing from this project except other domain packages
- `port` imports `model`
- `adapter` imports `port`, `model`, and the specific infrastructure SDK
- `app` imports `port`, `model`, `compiler`, `runner`, `assertion`, `report`
- `delivery` imports `app`, `model`, and infra SDKs only when necessary for that delivery (e.g. echo, cobra, websocket)
- `cmd/*` imports everything; `main.go` is the only place where concrete adapters meet app services

Violations are bugs.

---

## 3. Domain types vs persistence types

The domain layer (`internal/model`) uses **pure Go types only**:

```go
type Project struct {
    ID          uuid.UUID
    Name        string
    Description string
    CreatedAt   time.Time
    UpdatedAt   time.Time
}
```

The postgres adapter has its own sqlc-generated types (`pgtype.UUID`, `pgtype.Timestamptz`, etc.) and translates between them inside the repo implementation:

```go
// internal/adapter/postgres/repo.go
func (r *ProjectRepo) Get(ctx context.Context, id uuid.UUID) (*model.Project, error) {
    p, err := r.q.GetProject(ctx, pgUUID(id))
    if err != nil { return nil, mapErr(err) }
    return toModelProject(p), nil
}
```

Outside `internal/adapter/postgres/`, no code imports `pgtype`, `pgxpool`, or sqlc-generated structs. Likewise for `minio-go` outside `internal/adapter/minio/`, and `gorilla/websocket` outside `internal/delivery/http/`.

This is what makes the CLI possible: it wires a memory adapter or a file adapter, and the same app services and domain code work unchanged.

---

## 4. The compiler-runner contract

The compiler is the **bridge** between OpenAPI knowledge and runtime execution:

```text
[ OpenAPI spec + ProjectConfig ]
            │
            ▼
        Compiler            ← knows libopenapi, schemas, types
            │
            ▼
        SmokePlan           ← the AST: self-contained, executable
            │
            ▼
         Runner             ← knows HTTP, hooks, assertions
            │
            ▼
        RunResult
```

**Critical rule:** the runner does **not** import `libopenapi` or read the OpenAPI spec at runtime. Anything from the OpenAPI document that the runner needs (response schema for validation, content-type, param formats) is baked into the `SmokePlan` by the compiler.

This is what makes the runner a library rather than a service: pass it a `SmokePlan`, get back a `RunResult`. No DB, no spec parsing, no orchestration assumptions.

---

## 5. Runner as a pluggable library

The runner exposes a small, stable API:

```go
type Runner struct {
    HTTPClient     *http.Client
    PreProcessors  []PreProcessor
    PostProcessors []PostProcessor
    EventEmitter   func(Event)   // optional
}

func (r *Runner) Execute(ctx context.Context, plan *model.SmokePlan) *model.RunResult
```

### Hook interfaces

```go
type PreProcessor interface {
    BeforeRequest(rctx *RequestContext) error
}

type PostProcessor interface {
    AfterResponse(rctx *ResponseContext, step *model.StepResult) error
}
```

### Built-in hooks (`internal/runner/hook`)

| Hook | Phase | Purpose |
|---|---|---|
| `VariableInterpolator` | pre | Fill `{vars}` in path, headers, body |
| `AuthInjector` | pre | Inject bearer/basic/api-key from `plan.Auth` |
| `EnvironmentHeaders` | pre | Apply default headers from environment |
| `Capture` | post | Extract values via gjson into runtime vars |
| `Redactor` | post | Redact sensitive headers/body fields |
| `TraceExtractor` | post | Pull `X-Trace-Id` / `X-Request-Id` |
| `AssertionRunner` | post | Run all step assertions |

`runner.DefaultOptions()` composes these in the order shown. Either product can build custom options:

```go
// CLI: minimal, prints to stdout
plan := loadPlan("plan.yaml")
result := runner.New(runner.DefaultOptions()).Execute(ctx, plan)
```

```go
// API server: with event streaming
opts := runner.DefaultOptions()
opts.EventEmitter = eventBus.Emit
result := runner.New(opts).Execute(ctx, plan)
runRepo.SaveResult(ctx, runID, result)
```

Custom integrations can add hooks for AWS SigV4 signing, custom redaction lists, or additional metrics — without modifying the runner core.

---

## 6. Wiring example

### API server (`apps/core/cmd/server/main.go`)

```go
// Adapters
pool := pgxpool.New(ctx, dbURL)
projectRepo := postgres.NewProjectRepo(pool)
specRepo    := postgres.NewSpecRepo(pool)
runRepo     := postgres.NewRunRepo(pool)
eventBus    := inproc.NewEventBus()
jobs        := inproc.NewJobEnqueuer(eventBus, runRepo, runner.New(runner.DefaultOptions()))
blobs       := minio.NewBlobStore(minioClient, "smokery-artifacts")

// App services
projectSvc  := app.NewProjectService(projectRepo)
specSvc     := app.NewSpecService(specRepo, operationRepo)
runSvc      := app.NewRunService(runRepo, projectRepo, jobs)
reportSvc   := app.NewReportService(runRepo, blobs)

// Delivery
api := humaecho.New(echo.New(), huma.DefaultConfig("Smokery API", "1.0.0"))
http.RegisterProjectHandlers(api, projectSvc)
http.RegisterSpecHandlers(api, specSvc)
http.RegisterRunHandlers(api, runSvc, eventBus)
http.RegisterReportHandlers(api, reportSvc)
```

### CLI (`apps/core/cmd/smokery/main.go`)

```go
// Adapters
projectRepo := memory.NewProjectRepo()
specRepo    := memory.NewSpecRepo()
runRepo     := memory.NewRunRepo()
blobs       := fs.NewBlobStore("./smokery-output")
jobs        := inproc.NewJobEnqueuer(nil, runRepo, runner.New(runner.DefaultOptions()))

// App services (same as API)
runSvc    := app.NewRunService(runRepo, projectRepo, jobs)
reportSvc := app.NewReportService(runRepo, blobs)

// Delivery
cli.RegisterCommands(rootCmd, runSvc, reportSvc, specSvc)
rootCmd.Execute()
```

The app services and the runner are identical between products. Only the adapter choices differ.

---

## 7. Repository layout

```text
.
├── apps
│   ├── core              ← Go workspace (single go.mod, both products)
│   │   ├── cmd
│   │   │   ├── server    ← HTTP API binary
│   │   │   ├── smokery   ← CLI binary
│   │   │   └── openapi   ← internal: OpenAPI spec generator
│   │   └── internal
│   │       ├── model         ← Pure domain types
│   │       ├── spec          ← OpenAPI parser
│   │       ├── compiler      ← Config → SmokePlan
│   │       ├── runner        ← SmokePlan → RunResult
│   │       │   └── hook      ← Built-in pre/post processors
│   │       ├── assertion     ← Assertion logic
│   │       ├── report        ← RunResult → views
│   │       ├── port          ← Interfaces
│   │       ├── adapter
│   │       │   ├── postgres  ← pgx + sqlc impl
│   │       │   ├── inproc    ← in-process jobs + event bus
│   │       │   ├── minio     ← object storage
│   │       │   ├── fs        ← filesystem (CLI)
│   │       │   └── memory    ← in-memory (CLI, tests)
│   │       ├── app           ← Use case orchestration
│   │       └── delivery
│   │           ├── http      ← huma handlers (used by cmd/server)
│   │           └── cli       ← cobra commands (used by cmd/smokery)
│   └── web                   ← SvelteKit frontend
├── packages/types            ← Shared TS types (placeholder)
├── configs                   ← .air.toml, docker-compose.yml
├── docs                      ← Design + AI agent docs
├── tmp                       ← Build output (gitignored)
├── Makefile
├── AGENTS.md
└── README.md
```

The Go workspace lives at `apps/core`. The CLI binary is `smokery`. The HTTP server binary is `server`. Both are produced from one Go module.

---

## 8. What this enables

- **CLI without a database**: load a plan from YAML, run it, print or write results — no infrastructure required.
- **Testing without testcontainers**: app service tests use memory repos, runner tests use httptest servers and fake hooks.
- **Swap River later**: when persistent jobs are needed, write `internal/adapter/river/` implementing `port.JobEnqueuer`. No app or domain code changes.
- **Multiple storage backends**: S3, MinIO, GCS, filesystem — implementations of `port.BlobStore`, picked at wire-time.
- **Custom enterprise integrations**: organisations can vendor the domain packages and bring their own adapters and hooks without modifying upstream code.

---

## 9. Anti-patterns to avoid

- Importing `pgtype`, `pgxpool`, `minio-go`, or `gorilla/websocket` outside their respective adapter / delivery directories
- Returning `db.Project` from an app service or HTTP handler
- Calling `queries.X()` from inside an HTTP handler — go through an app service
- Adding business logic to a repo adapter (adapters only translate types and call the database)
- Adding HTTP-specific logic to the runner or compiler
- Using `libopenapi` inside the runner — anything spec-related must be done by the compiler
- Embedding job orchestration logic in the runner — the runner just executes one plan; orchestration is a port

---

*Last updated: 2026-05-26*
