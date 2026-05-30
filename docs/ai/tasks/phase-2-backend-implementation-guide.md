# Phase 2 Implementation Guide

This document is the implementation guide for post-MVP features. Each section describes the feature area, implementation approach, affected layers, and task breakdown.

Prerequisites: All Phase C (MVP) tasks must be complete. The system has a working compiler, runner, SQLite/Postgres adapters, fs/minio blob storage, CLI, and frontend with shadcn-svelte.

---

## 1. Reporting and Collaboration

### Goal
Extend the report system beyond Backend Debug and CI Summary to serve frontend developers, analysts, QA, and ops teams. Add collaboration features for failure triage.

### Architecture

Reports are generated from `model.RunResult` by the `internal/report` package. Each new view is a function that transforms a `RunResult` into a view-specific struct. The delivery layer exposes them as `GET /api/runs/{id}/report/{view}`.

```text
internal/report/
├── debug.go        (existing)
├── ci.go           (existing)
├── mermaid.go      (existing)
├── contract.go     ← new
├── analyst.go      ← new
├── qa.go           ← new
└── pdf.go          ← new (uses go-pdf or chromedp)
```

Collaboration features (failure classification, assignees, labels) extend the `comments` table or add new tables.

### Tasks

#### 1.1 Frontend Contract Report View
- **What**: Shows request/response schemas, field coverage, type mismatches, and breaking changes per endpoint.
- **Backend**: Add `report.ContractView` struct with per-endpoint schema conformance. Add `GET /api/runs/{id}/report/contract`.
- **Frontend**: New page at `/runs/[runId]/report/contract` showing a table of endpoints with schema pass/fail status.
- **Files**: `internal/report/contract.go`, `internal/delivery/http/report.go`, `apps/web/src/routes/runs/[runId]/report/contract/+page.svelte`

#### 1.2 Analyst Flow Report View
- **What**: Shows business flow outcomes — which user journeys passed/failed, with timing and data captured at each step.
- **Backend**: Add `report.AnalystView` struct grouping results by flow with captured business values.
- **Frontend**: Page at `/runs/[runId]/report/analyst` with flow diagrams and captured data tables.
- **Files**: `internal/report/analyst.go`, frontend route

#### 1.3 QA Evidence Report View
- **What**: Provides exportable test evidence — assertions, screenshots of responses, timestamps, environment info.
- **Backend**: Add `report.QAView` struct with assertion details, request/response pairs, and environment metadata.
- **Frontend**: Page at `/runs/[runId]/report/qa` with expandable assertion details and export button.
- **Files**: `internal/report/qa.go`, frontend route

#### 1.4 Report Template Builder UI
- **What**: Let users configure which sections appear in reports, custom headers, and branding.
- **Model**: Add `report_templates` table (id, project_id, name, config JSONB, created_at).
- **Backend**: CRUD API for report templates. Report generation accepts an optional template ID.
- **Frontend**: Page at `/projects/[id]/settings/reports` with drag-and-drop section ordering.
- **Files**: Migration, `internal/adapter/*/report_template_repo.go`, `internal/app/report_template.go`, delivery + frontend

#### 1.5 PDF Report Export
- **What**: Generate downloadable PDF reports from any report view.
- **Approach**: Use `chromedp` to render the HTML report artifact to PDF, or use `go-pdf/fpdf` for a pure-Go approach.
- **Backend**: Add `POST /api/runs/{id}/report/{view}/pdf` that generates and stores a PDF artifact.
- **Dependency**: `github.com/nicholasgasior/gofpdf` or `github.com/nicholasgasior/chromedp` (choose based on fidelity needs).
- **Files**: `internal/report/pdf.go`, delivery endpoint

#### 1.6 Failure Classification UI
- **What**: Let users classify failures (bug, flaky, environment, expected, config error).
- **Model**: Add `failure_classifications` table (id, run_id, step_name, classification, note, author, created_at).
- **Backend**: CRUD API. Classification is per-step within a run.
- **Frontend**: On run detail page, add a dropdown per failed step to classify the failure.
- **Files**: Migration, adapter, app service, delivery, frontend component

#### 1.7 Failure Assignee Support
- **What**: Assign failures to team members for follow-up.
- **Model**: Add `assignee` column to `failure_classifications` table.
- **Frontend**: Assignee picker on failure classification card.

#### 1.8 Labels and Status Transitions
- **What**: Add labels (tags) to runs and support status transitions (open → investigating → resolved).
- **Model**: Add `run_labels` table and `resolution_status` column to runs.
- **Backend**: API to add/remove labels, update resolution status.
- **Frontend**: Label badges on run list, status dropdown on run detail.

---

## 2. Spec Evolution

### Goal
Detect and surface API changes when a new spec version is imported. Help teams understand the impact on existing flows and suites.

### Architecture

Spec diffing lives in `internal/spec/diff.go`. It compares two parsed specs and produces a structured diff. The impact analysis lives in `internal/compiler/impact.go` — it takes a diff and the current project config and reports which flows/suites are affected.

```text
internal/spec/
├── parser.go       (existing)
├── classify.go     (existing)
└── diff.go         ← new

internal/compiler/
├── compiler.go     (existing)
└── impact.go       ← new
```

### Tasks

#### 2.1 Spec Diff on Re-import
- **What**: When importing a new spec for a project that already has one, compute and store the diff.
- **Backend**: `spec.Diff(old, new) DiffResult` — compares operations, schemas, parameters. Store diff as JSON in a new `spec_diffs` table.
- **API**: `GET /api/specs/{id}/diff` returns the diff against the previous version.
- **Files**: `internal/spec/diff.go`, migration, adapter, delivery

#### 2.2 Operation Added/Removed/Changed Detection
- **What**: Classify each operation as added, removed, unchanged, or changed (with change details).
- **Backend**: Part of `spec.Diff`. Each operation gets a status and a list of changes (path params, query params, request body, response schema).
- **Frontend**: Show diff badges on the operations page after re-import.

#### 2.3 Impact Report for Flows and Suites
- **What**: Given a spec diff, report which flows and suites reference affected operations.
- **Backend**: `compiler.Impact(diff, config) ImpactReport` — maps changed operations to flows/suites that use them.
- **API**: `GET /api/projects/{id}/impact` (computes from latest diff + current config).
- **Frontend**: Impact page at `/projects/[id]/impact` showing affected flows/suites with severity.

#### 2.4 Contract Drift Trend View
- **What**: Track spec changes over time and show a trend of breaking vs non-breaking changes.
- **Backend**: Aggregate spec diffs over time into a trend dataset.
- **Frontend**: LayerChart on the project overview showing drift over time.

---

## 3. Run History and Analytics

### Goal
Provide comparative analysis across runs to detect regressions, flakiness, and performance trends.

### Architecture

Analytics queries live in the adapter layer (SQL queries for postgres/sqlite). The app layer exposes them through an `AnalyticsService`. The frontend uses TanStack Query to fetch and LayerChart to render.

```text
internal/app/analytics.go       ← new service
internal/adapter/sqlite/analytics.go  ← new queries
internal/adapter/postgres/analytics.go
```

### Tasks

#### 3.1 Run Comparison: Current vs Previous
- **What**: Side-by-side comparison of two runs showing new failures, resolved failures, and regressions.
- **Backend**: `GET /api/runs/{id}/compare/{otherId}` — returns a diff of step results.
- **Frontend**: Comparison page at `/runs/[runId]/compare/[otherId]`.

#### 3.2 Run Comparison: Current vs Last Successful
- **What**: Shortcut to compare against the most recent passing run for the same project.
- **Backend**: Query to find last successful run, then reuse comparison logic.
- **API**: `GET /api/runs/{id}/compare/last-success`

#### 3.3 Latency Trend Charts
- **What**: Show p50/p95/p99 latency per endpoint over time.
- **Backend**: Aggregate step durations from stored run results. Return time-series data.
- **API**: `GET /api/projects/{id}/analytics/latency?operation_id=X&window=7d`
- **Frontend**: LayerChart line chart on the operations page or a dedicated analytics page.

#### 3.4 Endpoint Flakiness Detection
- **What**: Identify endpoints that intermittently fail across runs.
- **Backend**: Compute flakiness score = (runs with mixed pass/fail for same step) / total runs.
- **API**: `GET /api/projects/{id}/analytics/flaky`
- **Frontend**: Badge on operations page, dedicated flaky endpoints list.

#### 3.5 Environment Comparison
- **What**: Compare run results across different environments (staging vs production).
- **Backend**: Query runs by environment, group results, compute deltas.
- **Frontend**: Side-by-side environment health dashboard.

---

## 4. Runner and Infrastructure

### Goal
Support distributed execution, scheduled runs, and team notifications.

### Architecture

The runner remains a library. Distributed execution is handled by the job system. Phase 2 replaces the in-process worker with a proper job queue when scaling is needed.

```text
internal/adapter/
├── inproc/         (existing — goroutine worker)
├── redis/          ← new (pub/sub for scaled events)
└── river/          ← new (optional, for persistent job queue)

internal/delivery/http/
└── runner_registration.go  ← new
```

### Tasks

#### 4.1 Private/Self-hosted Runner Mode
- **What**: Allow runners to execute outside the server process, pulling jobs from a queue.
- **Approach**: Add a runner agent binary (`cmd/runner-agent`) that polls for jobs via HTTP or connects via WebSocket.
- **Protocol**: Server enqueues run → agent polls `GET /api/runner/jobs` → agent executes → agent POSTs result.
- **Files**: `cmd/runner-agent/main.go`, `internal/delivery/http/runner_jobs.go`

#### 4.2 Runner Registration Tokens
- **What**: Runners authenticate with a registration token to claim jobs.
- **Model**: `runner_tokens` table (id, project_id, token_hash, name, created_at, last_seen_at).
- **Backend**: Token generation API, middleware to validate runner auth.
- **Files**: Migration, adapter, `internal/app/runner_auth.go`, middleware

#### 4.3 Scheduled Smoke Runs
- **What**: Configure cron-like schedules for automatic smoke runs.
- **Model**: `schedules` table (id, project_id, cron_expr, config JSONB, enabled, last_run_at).
- **Backend**: Scheduler goroutine that checks due schedules and enqueues runs.
- **Dependency**: `github.com/robfig/cron/v3`
- **Files**: Migration, adapter, `internal/app/scheduler.go`, `internal/adapter/inproc/scheduler.go`

#### 4.4 Redis Pub/Sub for Scaled Live Events
- **What**: Replace in-process EventBus with Redis pub/sub so multiple server instances can broadcast WebSocket events.
- **Approach**: Implement `port.EventBus` backed by Redis pub/sub.
- **Config**: `EVENT_BUS=inproc|redis`, `REDIS_URL=redis://localhost:6379`
- **Dependency**: `github.com/redis/go-redis/v9`
- **Files**: `internal/adapter/redis/eventbus.go`, config update, main.go wiring

#### 4.5 Slack/Teams Notifications
- **What**: Send notifications on run completion (pass/fail) to Slack or Teams webhooks.
- **Model**: `notification_channels` table (id, project_id, type, webhook_url, config JSONB).
- **Backend**: Post-run hook that sends a formatted message to configured webhooks.
- **Files**: `internal/app/notify.go`, `internal/adapter/webhook/slack.go`, `internal/adapter/webhook/teams.go`

---

## 5. Observability & Log Correlation

### Goal
**The killer feature**: When a smoke test fails, instantly jump to the exact server logs, traces, and spans that handled that request. No manual searching, no timestamp guessing, no grepping through log aggregators. One click from a failed assertion to the server-side root cause.

### Why This Matters

Smoke tests tell you *what* failed. Server logs tell you *why*. The gap between these two is where debugging time is wasted. Smokery bridges this gap by:

1. **Injecting correlation IDs** into every smoke request (so the server tags its logs).
2. **Extracting correlation IDs** from every smoke response (so Smokery knows which logs to find).
3. **Generating deep links** to the exact log query / trace view for each failed step.
4. **Displaying a unified timeline** showing smoke step + server logs side by side.

### Architecture

```text
┌─────────────────────────────────────────────────────────────────┐
│ Smoke Runner                                                     │
│                                                                   │
│  PreProcessor: OTelPropagation                                   │
│    → Injects traceparent, X-Request-Id, X-Correlation-Id         │
│                                                                   │
│  PostProcessor: TraceExtractor (existing)                        │
│    → Captures trace_id, request_id, span_id from response        │
│                                                                   │
│  PostProcessor: DeepLinkGenerator                                │
│    → Builds URLs: Loki query, Tempo trace, Kibana search         │
│                                                                   │
│  PostProcessor: LogFetcher (optional)                            │
│    → Fetches relevant log lines via Loki/ES API                  │
│                                                                   │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│ StepResult                                                       │
│                                                                   │
│  response.trace_id      = "abc123..."                            │
│  response.request_id    = "req-456..."                           │
│  correlation.loki_url   = "https://grafana.../explore?..."       │
│  correlation.tempo_url  = "https://tempo.../trace/abc123"        │
│  correlation.kibana_url = "https://kibana.../app/discover?..."   │
│  correlation.logs[]     = [{timestamp, level, message}, ...]     │
│                                                                   │
└─────────────────────────────────────────────────────────────────┘
```

```text
internal/runner/hook/
├── trace.go        (existing — extracts trace/request IDs from response)
├── otel.go         ← new (injects traceparent + correlation headers)
├── deeplink.go     ← new (generates observability URLs per step)
└── logfetch.go     ← new (optional: fetches log snippets from Loki/ES)

internal/model/types.go
└── StepResult.Correlation  ← new field

internal/report/
└── correlation.go  ← correlation-focused report view

internal/delivery/http/
└── correlation.go  ← API for fetching logs on-demand
```

### Model Extension

```go
// Added to StepResult
type Correlation struct {
    TraceID    string            `json:"trace_id,omitempty"`
    SpanID     string            `json:"span_id,omitempty"`
    RequestID  string            `json:"request_id,omitempty"`
    Links      map[string]string `json:"links,omitempty"`      // "loki" → URL, "tempo" → URL, etc.
    LogSnippet []LogLine         `json:"log_snippet,omitempty"` // fetched server logs
}

type LogLine struct {
    Timestamp string `json:"timestamp"`
    Level     string `json:"level"`
    Message   string `json:"message"`
    Service   string `json:"service,omitempty"`
}
```

### Configuration

```env
# --- Correlation ID injection ---
# Header name to inject a unique correlation ID per request
CORRELATION_HEADER=X-Correlation-Id

# --- Trace propagation ---
# Inject W3C traceparent header (enables distributed tracing)
OTEL_PROPAGATION=true

# --- Deep link templates ---
# Use ${trace_id}, ${request_id}, ${correlation_id}, ${timestamp} as placeholders
LOKI_URL_TEMPLATE=https://grafana.example.com/explore?left={"queries":[{"expr":"{request_id=\"${request_id}\"}"}]}
TEMPO_URL_TEMPLATE=https://grafana.example.com/explore?left={"queries":[{"queryType":"traceql","query":"${trace_id}"}]}
JAEGER_URL_TEMPLATE=https://jaeger.example.com/trace/${trace_id}
KIBANA_URL_TEMPLATE=https://kibana.example.com/app/discover#/?_a=(query:(query:'request_id:${request_id}'))

# --- Log fetching (optional, for inline log display) ---
LOG_FETCH_ENABLED=false
LOG_FETCH_PROVIDER=loki|elasticsearch
LOKI_API_URL=http://loki:3100
LOKI_QUERY={app="target-service"} |= "${request_id}"
ES_URL=http://elasticsearch:9200
ES_INDEX=app-logs-*
ES_QUERY_FIELD=request_id
```

### Tasks

#### 5.0 Correlation Model Extension (do first)
- **What**: Add `Correlation` struct to `model.StepResult`. This is the foundation for all correlation features.
- **Backend**: Add `Correlation` field to `StepResult`. Update the runner to populate it from trace extractor output.
- **Migration**: None (stored in run_results JSONB).
- **Files**: `internal/model/types.go`, `internal/runner/execute.go`

#### 5.1 Correlation ID Injection (Pre-processor)
- **What**: Inject a unique correlation ID into every smoke request so the target server can tag its logs.
- **Approach**: New `hook.CorrelationInjector` pre-processor. Generates a UUID per request, injects it as a configurable header (default `X-Correlation-Id`). Also injects `X-Request-Id` if not already present.
- **Why**: Even if the target server doesn't support OpenTelemetry, it likely logs request headers. This gives you a searchable ID.
- **Files**: `internal/runner/hook/correlation.go`

#### 5.2 OpenTelemetry Traceparent Propagation
- **What**: Inject W3C `traceparent` header so smoke requests appear in the target service's distributed traces.
- **Approach**: Generate a valid trace ID + span ID per request. Inject `traceparent: 00-{trace_id}-{span_id}-01`. Store the generated trace_id in the step result.
- **Why**: If the target uses OTel, this means the smoke request's entire server-side processing is captured as a trace you can view.
- **Dependency**: None needed — trace ID is just 32 hex chars, span ID is 16 hex chars. Pure Go generation.
- **Files**: `internal/runner/hook/otel.go`

#### 5.3 Deep Link Generation (Post-processor)
- **What**: After each step, generate clickable URLs to Loki, Tempo, Jaeger, Kibana, or any log viewer.
- **Approach**: `hook.DeepLinkGenerator` post-processor. Takes URL templates from config, substitutes `${trace_id}`, `${request_id}`, `${correlation_id}`, `${timestamp}`. Attaches URLs to `StepResult.Correlation.Links`.
- **Config**: Multiple templates supported (one per observability tool).
- **Files**: `internal/runner/hook/deeplink.go`

#### 5.4 Inline Log Fetching (Optional Post-processor)
- **What**: After a step fails, automatically fetch the relevant server log lines and attach them to the result.
- **Approach**: `hook.LogFetcher` post-processor. Queries Loki or Elasticsearch API using the correlation/request ID. Fetches logs within a time window around the request. Attaches as `StepResult.Correlation.LogSnippet`.
- **Providers**: Loki (LogQL query), Elasticsearch (term query on request_id field).
- **Config**: Only runs on failed steps by default (configurable). Time window: request timestamp ± 5s.
- **Files**: `internal/runner/hook/logfetch.go`, `internal/runner/hook/logfetch_loki.go`, `internal/runner/hook/logfetch_es.go`

#### 5.5 Correlation Report View
- **What**: A dedicated report view showing the correlation timeline — smoke step on the left, server logs on the right.
- **Backend**: `report.CorrelationView` — groups steps with their correlation data, links, and log snippets.
- **API**: `GET /api/runs/{id}/report/correlation`
- **Frontend**: Page at `/runs/[runId]/report/correlation` with a split-pane timeline view.
- **Files**: `internal/report/correlation.go`, delivery, frontend route

#### 5.6 On-Demand Log Fetch API
- **What**: Let the frontend fetch logs for a specific step on demand (for when inline fetching is disabled).
- **API**: `POST /api/runs/{id}/steps/{stepName}/logs` — fetches logs using the step's correlation IDs.
- **Frontend**: "Fetch Logs" button on failed steps that loads server logs inline.
- **Files**: `internal/delivery/http/correlation.go`, `internal/app/correlation.go`

#### 5.7 Correlation in Debug Report
- **What**: Enhance the existing debug report to include correlation links and log snippets for failed steps.
- **Backend**: Update `report.DebugView` to include correlation data.
- **Frontend**: Show "View Logs" / "View Trace" links on the run detail page for each failed step.

#### 5.8 Prometheus Metrics Export
- **What**: Export smoke test metrics (pass rate, latency, run count) as Prometheus metrics.
- **Approach**: Add `/metrics` endpoint using `prometheus/client_golang`.
- **Metrics**: `smokery_runs_total`, `smokery_run_duration_seconds`, `smokery_step_pass_total`, `smokery_step_fail_total`, `smokery_correlation_logs_fetched_total`.
- **Dependency**: `github.com/prometheus/client_golang`
- **Files**: `internal/delivery/http/metrics.go`

### Frontend UX for Correlation

The run detail page should show correlation prominently:

```text
┌─────────────────────────────────────────────────────────┐
│ Step: POST /users (createUser)                    FAILED │
├─────────────────────────────────────────────────────────┤
│ Assertion: status expected 201, got 500                  │
│                                                          │
│ ┌─ Correlation ────────────────────────────────────────┐ │
│ │ Request ID:  req-a1b2c3d4                            │ │
│ │ Trace ID:    00-abcdef1234567890abcdef1234567890     │ │
│ │                                                      │ │
│ │ [🔗 View Logs in Grafana]  [🔗 View Trace in Tempo] │ │
│ │                                                      │ │
│ │ ── Server Logs (3 lines) ──────────────────────────  │ │
│ │ 12:03:01.234 ERROR users/create.go:45                │ │
│ │   "duplicate key violation: email already exists"     │ │
│ │ 12:03:01.235 ERROR middleware/error.go:12            │ │
│ │   "returning 500 to client"                          │ │
│ └──────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────┘
```

### CLI Support

The CLI should also benefit from correlation:

```bash
$ smokery run -c config.yaml -s spec.json --correlation

Step: POST /users (createUser) .............. FAILED
  Assertion: status expected 201, got 500
  Request-ID: req-a1b2c3d4
  Trace-ID:   abcdef1234567890abcdef1234567890
  Logs URL:   https://grafana.example.com/explore?...

  Server logs:
    12:03:01 ERROR duplicate key violation: email already exists
    12:03:01 ERROR returning 500 to client
```

---

## 6. Security and Governance

### Goal
Add access control, approval workflows, and audit trails for team environments.

### Architecture

Auth is handled at the delivery layer via middleware. RBAC is stored in the database. The app layer checks permissions before executing operations.

```text
internal/app/auth.go            ← permission checks
internal/delivery/http/
├── middleware_auth.go          ← JWT/session validation
└── middleware_rbac.go          ← role enforcement
```

### Tasks

#### 6.1 Role-Based Access Control
- **What**: Define roles (admin, editor, viewer) with different permissions.
- **Model**: `users` table, `project_members` table (user_id, project_id, role).
- **Backend**: Auth middleware that extracts user from JWT/session, RBAC middleware that checks role for the target project.
- **Approach**: Start with simple JWT tokens. Can integrate with OAuth2/OIDC later.
- **Files**: Migrations, `internal/app/auth.go`, middleware, delivery updates

#### 6.2 Project-Level Permissions
- **What**: Each project has its own member list with roles.
- **Backend**: `port.MemberRepo` interface, permission check in app services.
- **Frontend**: Project settings page with member management.

#### 6.3 Approval Workflow for Destructive Flows
- **What**: Flows containing destructive operations require approval before execution.
- **Model**: `approvals` table (id, run_id, flow_id, status, approver, created_at).
- **Backend**: Compiler marks flows needing approval. Run creation checks approval status.
- **Frontend**: Approval request UI, approval inbox for admins.

#### 6.4 Audit Trail
- **What**: Log all significant actions (spec import, run start, classification change, config update).
- **Model**: `audit_log` table (id, project_id, user_id, action, resource_type, resource_id, metadata JSONB, created_at).
- **Backend**: Audit logger called from app services.
- **Frontend**: Audit log page at `/projects/[id]/settings/audit`.

---

## 7. AI-Assisted Features

### Goal
Use LLMs to reduce manual work in classification, flow design, failure analysis, and reporting.

### Architecture

AI features are implemented as optional services that call an LLM API (OpenAI, Anthropic, or local). They live in a dedicated package and are wired only when configured.

```text
internal/ai/
├── client.go       ← LLM API client (configurable provider)
├── classify.go     ← operation classification suggestions
├── suggest.go      ← flow suggestions from spec
├── explain.go      ← failure explanations
├── summarize.go    ← report summaries
└── drift.go        ← config drift recommendations
```

**Config**: `AI_ENABLED=true`, `AI_PROVIDER=openai|anthropic`, `AI_API_KEY=...`, `AI_MODEL=gpt-4o`

### Tasks

#### 7.1 AI-Assisted Operation Classification Review
- **What**: After spec import, suggest classifications for operations that the heuristic classifier is uncertain about.
- **Approach**: Send operation metadata (method, path, summary, parameters) to LLM with a classification prompt.
- **UI**: Show AI suggestions as badges on the operations page with accept/reject buttons.
- **Files**: `internal/ai/classify.go`, delivery endpoint, frontend component

#### 7.2 AI-Assisted Flow Suggestion from OpenAPI
- **What**: Given an OpenAPI spec, suggest logical smoke test flows (e.g., CRUD sequences).
- **Approach**: Send operation list + schemas to LLM, ask for flow definitions in the project config format.
- **UI**: "Suggest Flows" button on the flow builder that populates suggested flows.
- **Files**: `internal/ai/suggest.go`, delivery endpoint, frontend integration

#### 7.3 AI-Assisted Failure Explanation
- **What**: When a step fails, provide a human-readable explanation of what likely went wrong.
- **Approach**: Send the step result (request, response, assertion, error) to LLM for analysis.
- **UI**: "Explain" button on failed steps that shows the AI explanation inline.
- **Files**: `internal/ai/explain.go`, delivery endpoint, frontend component

#### 7.4 AI-Assisted Report Summary
- **What**: Generate a natural-language summary of a run result for non-technical stakeholders.
- **Approach**: Send the run result summary to LLM, ask for a 2-3 paragraph executive summary.
- **UI**: "Generate Summary" button on the report page.
- **Files**: `internal/ai/summarize.go`, delivery endpoint, frontend integration

#### 7.5 AI-Assisted Config Drift Recommendation
- **What**: When spec drift is detected, suggest config changes (new flows, updated assertions, removed operations).
- **Approach**: Send the spec diff + current config to LLM, ask for recommended config updates.
- **UI**: Drift recommendation panel on the impact page.
- **Files**: `internal/ai/drift.go`, delivery endpoint, frontend integration

---

## Implementation Order

Recommended execution order based on dependencies and value:

```text
Phase 2A — Log Correlation (highest value, differentiator):
  5.0  Correlation model extension
  5.1  Correlation ID injection
  5.2  OpenTelemetry traceparent propagation
  5.3  Deep link generation
  5.7  Correlation in debug report (frontend links)
  5.4  Inline log fetching (Loki/ES)
  5.5  Correlation report view
  5.6  On-demand log fetch API

Phase 2B — Analytics & Spec Evolution:
  3.3  Latency trend charts
  3.4  Endpoint flakiness detection
  2.1  Spec diff on re-import
  2.2  Operation added/removed/changed detection
  1.1  Frontend Contract report view
  1.6  Failure classification UI

Phase 2C — Collaboration:
  1.2  Analyst Flow report view
  1.3  QA Evidence report view
  1.7  Failure assignee support
  1.8  Labels and status transitions
  2.3  Impact report for flows and suites
  3.1  Run comparison: current vs previous
  4.5  Slack/Teams notifications

Phase 2D — Infrastructure scaling:
  4.3  Scheduled smoke runs
  4.4  Redis pub/sub for scaled live events
  5.8  Prometheus metrics export

Phase 2E — Security and governance:
  6.1  Role-based access control
  6.2  Project-level permissions
  6.3  Approval workflow for destructive flows
  6.4  Audit trail

Phase 2F — Advanced features:
  4.1  Private/self-hosted runner mode
  4.2  Runner registration tokens
  1.4  Report template builder UI
  1.5  PDF report export
  2.4  Contract drift trend view
  3.2  Run comparison: current vs last successful
  3.5  Environment comparison

Phase 2G — AI-assisted (requires LLM API access):
  7.1  AI-assisted operation classification
  7.2  AI-assisted flow suggestion
  7.3  AI-assisted failure explanation
  7.4  AI-assisted report summary
  7.5  AI-assisted config drift recommendation
```

---

## New Dependencies (Phase 2)

| Feature | Package | Purpose |
|---------|---------|---------|
| Scheduled runs | `github.com/robfig/cron/v3` | Cron expression parsing |
| Redis events | `github.com/redis/go-redis/v9` | Pub/sub for scaled WebSocket |
| PDF export | `github.com/nicholasgasior/gofpdf` or `chromedp` | PDF generation |
| OTel propagation | `go.opentelemetry.io/otel` | Trace context generation |
| Metrics | `github.com/prometheus/client_golang` | Prometheus metrics |
| AI features | `github.com/sashabaranov/go-openai` | LLM API client |

---

## Database Migrations (Phase 2)

New tables needed:

```sql
-- 1.4 Report templates
CREATE TABLE report_templates (...);

-- 1.6-1.7 Failure classification
CREATE TABLE failure_classifications (...);

-- 1.8 Run labels
CREATE TABLE run_labels (...);
ALTER TABLE runs ADD COLUMN resolution_status TEXT;

-- 2.1 Spec diffs
CREATE TABLE spec_diffs (...);

-- 4.2 Runner tokens
CREATE TABLE runner_tokens (...);

-- 4.3 Schedules
CREATE TABLE schedules (...);

-- 4.5 Notifications
CREATE TABLE notification_channels (...);

-- 6.1-6.2 Auth
CREATE TABLE users (...);
CREATE TABLE project_members (...);

-- 6.3 Approvals
CREATE TABLE approvals (...);

-- 6.4 Audit
CREATE TABLE audit_log (...);
```

Each migration must work for both SQLite and PostgreSQL. Use compatible SQL (TEXT for UUIDs in SQLite, UUID in PostgreSQL). The SQLite adapter auto-applies schema on open; PostgreSQL uses golang-migrate.
