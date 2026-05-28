# Frontend UI/UX Plan — Smokery

## Design Philosophy

- **Spec-driven reactive UI**: The OpenAPI spec is the source of truth. Every form, dropdown, and validation derives from the imported spec.
- **Builder-first**: The smoke config builder is the core experience — interactive, reactive, showing real-time compilation feedback as users configure.
- **Dark + Green theme**: Dark background (`zinc-900`/`zinc-950`), green accents (`emerald-500`/`emerald-400`) for success/primary actions, amber for warnings, red for failures.
- **Expressive dashboard**: Dense but scannable — cards, badges, mini-charts, quick actions. Designed for power users who manage multiple projects.

---

## Tech Stack

- Svelte 5 (runes, `$state`, `$derived`, `$effect`)
- SvelteKit 2 (file-based routing, load functions)
- shadcn-svelte (component library, dark theme)
- TailwindCSS 4 (utility-first, dark mode default)
- `@tanstack/svelte-query` v6 (server state)
- `openapi-fetch` + `openapi-typescript` (type-safe API client)
- Monaco Editor (YAML/JSON editing for overrides)
- Mermaid (sequence diagrams)
- LayerChart (latency/pass-rate trends)
- lucide-svelte (icons)

---

## Theme & Visual Language

```
Background:     zinc-950 (#09090b)
Surface:        zinc-900 (#18181b)
Card:           zinc-800/50 with border zinc-700/50
Primary:        emerald-500 (#10b981)
Primary hover:  emerald-400
Text:           zinc-100
Text muted:     zinc-400
Success:        emerald-500
Warning:        amber-500
Error:          red-500
Info:           sky-500
Border:         zinc-700/50
```

Cards use subtle glass-morphism (backdrop-blur, semi-transparent backgrounds). Status badges are pill-shaped with colored backgrounds. Buttons are solid emerald for primary, ghost/outline for secondary.

---

## Route Structure

```
/                                    → Redirect to /projects
/projects                            → Project dashboard (list + stats)
/projects/[id]                       → Project overview (spec status, recent runs, quick actions)
/projects/[id]/spec                  → Spec viewer (operations tree, schemas)
/projects/[id]/operations            → Operation explorer + override editor
/projects/[id]/environments          → Environment + auth profile config
/projects/[id]/builder               → Smoke config builder (THE core page)
/projects/[id]/builder/flows/[fid]   → Flow step editor
/projects/[id]/builder/suites/[sid]  → Suite strategy configurator
/projects/[id]/plan                  → Plan preview (compiled output + errors/warnings)
/projects/[id]/runs                  → Run history (table + trend charts)
/projects/[id]/settings              → Project settings, danger zone
/runs/[runId]                        → Run detail (live WebSocket feed, step results)
/runs/[runId]/report/[view]          → Report views (debug, ci, mermaid, html)
```

---

## Page Designs

### 1. Project Dashboard (`/projects`)

```
┌─────────────────────────────────────────────────────────────┐
│  SMOKERY                                    [+ New Project] │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐        │
│  │ Total Runs  │  │ Pass Rate   │  │ Projects    │        │
│  │    142      │  │   94.3%     │  │     5       │        │
│  │  ↑12 today  │  │  ●●●●●○    │  │             │        │
│  └─────────────┘  └─────────────┘  └─────────────┘        │
│                                                             │
│  Projects                                                   │
│  ┌──────────────────────────────────────────────────────┐  │
│  │ 🟢 Payment API    v2.1.0    12 ops   Last: 2m ago   │  │
│  │    94% pass │ 3 flows │ 1 suite │ [Run] [Builder]   │  │
│  ├──────────────────────────────────────────────────────┤  │
│  │ 🟡 User Service   v1.4.2    28 ops   Last: 1h ago   │  │
│  │    87% pass │ 5 flows │ 2 suites│ [Run] [Builder]   │  │
│  ├──────────────────────────────────────────────────────┤  │
│  │ 🔴 Inventory API  v3.0.0    45 ops   Last: 3h ago   │  │
│  │    72% pass │ 2 flows │ 1 suite │ [Run] [Builder]   │  │
│  └──────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
```

**Components**: StatsCard, ProjectCard (with inline sparkline), QuickAction buttons.

---

### 2. Smoke Config Builder (`/projects/[id]/builder`)

This is the **core interactive page**. It's a split-pane layout:

```
┌────────────────────────────────────────────────────────────────────┐
│  ← Project    Payment API    [Compile] [▶ Run]    [Plan Preview]  │
├──────────────────────────────┬─────────────────────────────────────┤
│  CONFIG PANEL                │  LIVE PREVIEW                       │
│                              │                                     │
│  Environments ──────────     │  Compiled Plan                      │
│  ┌─────────────────────┐    │  ┌───────────────────────────────┐  │
│  │ 🟢 staging          │    │  │ ✓ 2 flows, 1 suite            │  │
│  │ 🔵 production       │    │  │ ✓ 14 cases generated          │  │
│  │ [+ Add]             │    │  │ ⚠ 2 warnings                  │  │
│  └─────────────────────┘    │  │                               │  │
│                              │  │ Flow: User CRUD               │  │
│  Flows ─────────────────     │  │   → POST /users (create)     │  │
│  ┌─────────────────────┐    │  │   → GET /users/{id} (read)   │  │
│  │ 📋 User CRUD        │    │  │   → PUT /users/{id} (update) │  │
│  │    4 steps, staging  │    │  │   → DELETE /users/{id} (del) │  │
│  │    [Edit] [Delete]   │    │  │                               │  │
│  ├─────────────────────┤    │  │ Suite: List Endpoints          │  │
│  │ 📋 Order Flow       │    │  │   → GET /users (default)     │  │
│  │    3 steps, staging  │    │  │   → GET /users (pagination)  │  │
│  │    [Edit] [Delete]   │    │  │   → GET /users (search)     │  │
│  └─────────────────────┘    │  │   → GET /orders (default)    │  │
│  [+ New Flow]                │  │   → GET /orders (enum:...)   │  │
│                              │  └───────────────────────────────┘  │
│  Suites ────────────────     │                                     │
│  ┌─────────────────────┐    │  Warnings                           │
│  │ 🔄 List Endpoints   │    │  ┌───────────────────────────────┐  │
│  │    auto-generated    │    │  │ ⚠ deleteUser is destructive  │  │
│  │    14 cases          │    │  │ ⚠ empty result on /archived  │  │
│  │    [Configure]       │    │  └───────────────────────────────┘  │
│  └─────────────────────┘    │                                     │
│  [+ New Suite]               │                                     │
├──────────────────────────────┴─────────────────────────────────────┤
│  Operations from spec: 28 total │ 12 list │ 8 read │ 5 create │ 3 │
└────────────────────────────────────────────────────────────────────┘
```

**Key behaviors**:
- Left panel is the config editor (environments, flows, suites)
- Right panel is a **live compilation preview** — updates reactively as user edits config
- Compilation runs client-side via the `/plan/preview` API on every meaningful change (debounced 500ms)
- Warnings and errors show inline with links to the offending config section
- Operation counts at the bottom show what the spec provides

---

### 3. Flow Step Editor (`/projects/[id]/builder/flows/[fid]`)

```
┌────────────────────────────────────────────────────────────────────┐
│  ← Builder    Flow: User CRUD    env: staging    [Save] [Delete]  │
├────────────────────────────────────────────────────────────────────┤
│                                                                    │
│  Steps                                                             │
│  ┌──────────────────────────────────────────────────────────────┐ │
│  │ 1. create-user                                                │ │
│  │    Operation: [createUser ▼]  ← dropdown from spec            │ │
│  │    POST /users                                                │ │
│  │                                                               │ │
│  │    Body:                                                      │ │
│  │    ┌─────────────────────────────────────────────────────┐   │ │
│  │    │ { "name": "Test", "email": "t@t.com" }             │   │ │
│  │    └─────────────────────────────────────────────────────┘   │ │
│  │    ↑ Schema hint: name(string,required), email(string,req)    │ │
│  │                                                               │ │
│  │    Assertions:  [+ Add]                                       │ │
│  │    ┌────────────────────────────────────────────────────┐    │ │
│  │    │ ✓ status = 201                                     │    │ │
│  │    │ ✓ jsonpath: id exists                              │    │ │
│  │    └────────────────────────────────────────────────────┘    │ │
│  │                                                               │ │
│  │    Captures:  [+ Add]                                         │ │
│  │    ┌────────────────────────────────────────────────────┐    │ │
│  │    │ user_id ← body.id                                  │    │ │
│  │    └────────────────────────────────────────────────────┘    │ │
│  └──────────────────────────────────────────────────────────────┘ │
│                                                                    │
│  ┌──────────────────────────────────────────────────────────────┐ │
│  │ 2. get-user                                                   │ │
│  │    Operation: [getUser ▼]                                     │ │
│  │    GET /users/{userId}                                        │ │
│  │    Params: userId = {{user_id}}  ← autocomplete from captures │ │
│  │    Assertions: status = 200, jsonpath: name = "Test"          │ │
│  └──────────────────────────────────────────────────────────────┘ │
│                                                                    │
│  [+ Add Step]    [+ Add Cleanup Step]                              │
│                                                                    │
│  Cleanup Steps (run on failure or completion)                      │
│  ┌──────────────────────────────────────────────────────────────┐ │
│  │ 🧹 delete-user    DELETE /users/{{user_id}}                   │ │
│  └──────────────────────────────────────────────────────────────┘ │
└────────────────────────────────────────────────────────────────────┘
```

**Spec-driven reactive features**:
- **Operation dropdown**: populated from spec operations, grouped by tag, shows method badge + path
- **Body editor**: shows schema hints (required fields, types, enums) derived from the spec's request body schema
- **Params**: auto-populated from spec path/query params with type hints
- **Assertions**: type dropdown (status, jsonpath, schema, not_empty, list_shape, pagination), path autocomplete from response schema
- **Captures**: source (body/header), path with autocomplete from response schema fields
- **Variable interpolation**: `{{var}}` syntax with autocomplete from previously captured variables in earlier steps
- **Drag-to-reorder** steps

---

### 4. Suite Configurator (`/projects/[id]/builder/suites/[sid]`)

```
┌────────────────────────────────────────────────────────────────────┐
│  ← Builder    Suite: List Endpoints    [Save] [Delete]            │
├──────────────────────────┬─────────────────────────────────────────┤
│  STRATEGY                │  GENERATED CASES (live preview)         │
│                          │                                         │
│  Environment: [staging▼] │  14 cases from 6 operations:           │
│                          │                                         │
│  Selector ───────────    │  ┌─────────────────────────────────┐   │
│  Tags: [users] [orders]  │  │ GET /users                      │   │
│  Classifications:        │  │   • default_list  ✓             │   │
│    [✓] list              │  │   • pagination (page, limit) ✓  │   │
│    [✓] read              │  │   • search (q) ✓                │   │
│    [ ] create            │  │   • filter: status=active ✓     │   │
│    [ ] update            │  │   • filter: status=inactive ✓   │   │
│    [ ] delete            │  ├─────────────────────────────────┤   │
│  Exclude: [deleteUser]   │  │ GET /orders                     │   │
│                          │  │   • default_list  ✓             │   │
│  Strategy ───────────    │  │   • pagination (page, limit) ✓  │   │
│  [✓] Default list        │  │   • filter: status=pending ✓   │   │
│  [✓] Pagination          │  │   • filter: status=shipped ✓   │   │
│  [✓] Search from resp    │  │   • filter: status=delivered ✓ │   │
│  [✓] Enum filters        │  └─────────────────────────────────┘   │
│  Max cases/op: [5]       │                                         │
│                          │  Excluded (destructive):                 │
│  Empty result: [warn ▼]  │  ┌─────────────────────────────────┐   │
│                          │  │ ⚠ POST /users (createUser)      │   │
│                          │  │ ⚠ DELETE /users/{id}             │   │
│                          │  └─────────────────────────────────┘   │
└──────────────────────────┴─────────────────────────────────────────┘
```

**Spec-driven reactive features**:
- **Selector**: tags come from spec tags, classifications from our analysis
- **Matched operations**: updates live as selector changes
- **Generated cases**: shows exactly what the compiler will produce — pagination params from spec's query hints, search params, enum values
- **Excluded operations**: shows why (destructive flag) with option to override

---

### 5. Run Detail (`/runs/[runId]`)

```
┌────────────────────────────────────────────────────────────────────┐
│  Run #abc123    🟢 passed    Duration: 1.2s    [Cancel] [Re-run]  │
├────────────────────────────────────────────────────────────────────┤
│  [Summary] [Steps] [Diagram] [Report]                              │
├────────────────────────────────────────────────────────────────────┤
│                                                                    │
│  Flow: User CRUD ─── 🟢 passed (450ms)                            │
│  ┌──────────────────────────────────────────────────────────────┐ │
│  │ ✓ create-user   POST /users         201   120ms   trace:abc │ │
│  │ ✓ get-user      GET /users/42       200    45ms   trace:def │ │
│  │ ✓ update-user   PUT /users/42       200    89ms   trace:ghi │ │
│  │ ✓ delete-user   DELETE /users/42    204    32ms   trace:jkl │ │
│  └──────────────────────────────────────────────────────────────┘ │
│                                                                    │
│  Suite: List Endpoints ─── 🟢 passed (750ms)                      │
│  ┌──────────────────────────────────────────────────────────────┐ │
│  │ ✓ listUsers_default       GET /users       200   65ms       │ │
│  │ ✓ listUsers_pagination    GET /users       200   72ms       │ │
│  │ ✓ listUsers_search        GET /users       200   58ms       │ │
│  │ ✓ listOrders_default      GET /orders      200   81ms       │ │
│  │ ⚠ listOrders_filter_...   GET /orders      200   44ms  warn│ │
│  └──────────────────────────────────────────────────────────────┘ │
│                                                                    │
│  Comments                                                          │
│  ┌──────────────────────────────────────────────────────────────┐ │
│  │ @alice: Looks good, the empty result on archived is expected │ │
│  │ [Add comment...]                                             │ │
│  └──────────────────────────────────────────────────────────────┘ │
└────────────────────────────────────────────────────────────────────┘
```

**Live features**:
- WebSocket connection for in-progress runs (events: flow.step.started, suite.case.result)
- Expandable rows showing request/response details, assertions, captures
- Click trace ID to copy or link to observability tool
- Mermaid tab shows sequence diagram
- Report tab links to HTML artifact

---

### 6. Operation Explorer (`/projects/[id]/operations`)

```
┌────────────────────────────────────────────────────────────────────┐
│  Operations    28 total    [Filter: ___] [Group by: tag ▼]        │
├────────────────────────────────────────────────────────────────────┤
│                                                                    │
│  users (12 operations)                                             │
│  ┌──────────────────────────────────────────────────────────────┐ │
│  │ GET    /users           listUsers       list        safe     │ │
│  │ GET    /users/{id}      getUser         read        safe     │ │
│  │ POST   /users           createUser      create      ⚠ destr │ │
│  │ PUT    /users/{id}      updateUser      update      ⚠ destr │ │
│  │ DELETE /users/{id}      deleteUser      delete      ⚠ destr │ │
│  └──────────────────────────────────────────────────────────────┘ │
│                                                                    │
│  Selected: listUsers                                               │
│  ┌──────────────────────────────────────────────────────────────┐ │
│  │ Classification: [list ▼]    Destructive: [ ]                 │ │
│  │                                                               │ │
│  │ Query Hints (from spec):                                      │ │
│  │   Pagination: page, limit                                     │ │
│  │   Search: q                                                   │ │
│  │   Enum filters: status (active, inactive, archived)           │ │
│  │                                                               │ │
│  │ Response schema: { type: "object", properties: { data: ... }} │ │
│  │                                                               │ │
│  │ Overrides (JSON):                                             │ │
│  │ ┌────────────────────────────────────────────────────────┐   │ │
│  │ │ { "custom_headers": { "X-Test": "true" } }            │   │ │
│  │ └────────────────────────────────────────────────────────┘   │ │
│  │                                              [Save Override]  │ │
│  └──────────────────────────────────────────────────────────────┘ │
└────────────────────────────────────────────────────────────────────┘
```

---

## Reactive Data Flow

```
┌─────────────────────────────────────────────────────────────┐
│                    BROWSER STATE                             │
│                                                             │
│  specAnalysis ─────────────────────────────────────────┐    │
│  (from GET /api/projects/{id}/spec)                    │    │
│    • operations[]                                      │    │
│    • schemas                                           │    │
│    • queryHints per operation                          │    │
│                                                        │    │
│  projectConfig ────────────────────────────────────┐   │    │
│  ($state, user-editable)                           │   │    │
│    • environments[]                                │   │    │
│    • flows[]                                       │   │    │
│    • suites[]                                      │   │    │
│                                                    │   │    │
│  compiledPlan ─────────────────────────────────────┤   │    │
│  ($derived, from POST /plan/preview)               │   │    │
│    • plan (or null)                                │   │    │
│    • errors[]                                      │   │    │
│    • warnings[]                                    │   │    │
│                                                    │   │    │
│  On every config change (debounced 500ms):         │   │    │
│    POST /api/projects/{id}/plan/preview            │   │    │
│      body: projectConfig                           │   │    │
│      → updates compiledPlan                        │   │    │
│      → UI shows errors/warnings/generated cases    │   │    │
└─────────────────────────────────────────────────────────────┘
```

**Key reactive patterns (Svelte 5 runes)**:

```svelte
<script>
  let config = $state(initialConfig);        // user-editable config
  let specOps = $derived(spec.operations);   // from loaded spec

  // Debounced compilation preview
  let compiledPlan = $state(null);
  $effect(() => {
    const timeout = setTimeout(async () => {
      compiledPlan = await api.POST('/projects/{id}/plan/preview', { body: config });
    }, 500);
    return () => clearTimeout(timeout);
  });

  // Derived: available operations for dropdowns
  let operationOptions = $derived(
    specOps.map(op => ({ value: op.operation_id, label: `${op.method} ${op.path}`, ...op }))
  );

  // Derived: captured variables available for interpolation
  let availableVars = $derived(
    config.flows.flatMap(f => f.steps.flatMap(s => s.captures?.map(c => c.name) ?? []))
  );
</script>
```

---

## Component Architecture

```
src/lib/components/
├── ui/                    ← shadcn-svelte primitives (button, card, badge, etc.)
├── layout/
│   ├── Sidebar.svelte
│   ├── TopBar.svelte
│   └── SplitPane.svelte
├── project/
│   ├── ProjectCard.svelte
│   ├── StatsCard.svelte
│   └── ProjectNav.svelte
├── builder/
│   ├── ConfigPanel.svelte        ← left pane: env/flow/suite list
│   ├── LivePreview.svelte        ← right pane: compiled plan view
│   ├── FlowEditor.svelte         ← step list with drag-reorder
│   ├── StepEditor.svelte         ← single step form
│   ├── SuiteConfigurator.svelte  ← strategy + selector form
│   ├── OperationPicker.svelte    ← spec-driven dropdown
│   ├── AssertionEditor.svelte    ← assertion type + value
│   ├── CaptureEditor.svelte      ← capture source + path
│   ├── VariableInput.svelte      ← input with {{var}} autocomplete
│   └── SchemaHint.svelte         ← shows expected schema inline
├── operations/
│   ├── OperationTable.svelte
│   ├── OperationDetail.svelte
│   └── OverrideEditor.svelte     ← Monaco for JSON overrides
├── runs/
│   ├── RunTable.svelte
│   ├── RunDetail.svelte
│   ├── StepResultRow.svelte
│   ├── LiveRunFeed.svelte        ← WebSocket event stream
│   └── RunTrends.svelte          ← LayerChart sparklines
├── reports/
│   ├── MermaidDiagram.svelte
│   ├── DebugView.svelte
│   └── CISummary.svelte
└── environments/
    ├── EnvironmentForm.svelte
    └── AuthProfileForm.svelte
```

---

## Implementation Order

1. **C.3.7** — TailwindCSS 4 + shadcn-svelte setup, dark+green theme tokens
2. **C.3.1** — Operation explorer (reads spec, shows operations, override editor)
3. **C.3.2** — Environment config pages
4. **C.3.3** — Auth profile config
5. **C.3.4** — Flow builder (step editor with spec-driven dropdowns)
6. **C.3.5** — Suite configurator (strategy form + live case preview)
7. **C.3.6** — Plan preview page (full compilation output)
8. **C.3.8** — Mermaid rendering on run detail
9. **C.3.9** — LayerChart trends on runs list

Each step builds on the previous — the operation explorer establishes the spec-data patterns, the builder uses those patterns for dropdowns and hints.

---

## API Endpoints Used by Frontend

| Page | Endpoint | Purpose |
|------|----------|---------|
| Dashboard | `GET /api/projects` | List projects |
| Spec | `GET /api/projects/{id}/specs` | Get spec analysis |
| Operations | `GET /api/operations?spec_id=X` | List operations |
| Operations | `PUT /api/operations/{id}/classification` | Override classification |
| Operations | `PUT /api/operations/{id}/overrides` | Set JSON overrides |
| Builder | `POST /api/projects/{id}/plan/preview` | Live compilation |
| Runs | `POST /api/projects/{id}/runs` | Start a run |
| Runs | `GET /api/projects/{id}/runs` | List runs |
| Runs | `GET /api/runs/{id}` | Get run status |
| Runs | `GET /api/runs/{id}/result` | Get run result |
| Runs | `POST /api/runs/{id}/cancel` | Cancel run |
| Reports | `GET /api/runs/{id}/report/debug` | Debug view |
| Reports | `GET /api/runs/{id}/report/ci` | CI summary |
| Reports | `GET /api/runs/{id}/report/mermaid` | Mermaid diagram |
| Comments | `GET /api/runs/{id}/comments` | List comments |
| Comments | `POST /api/runs/{id}/comments` | Add comment |
| WebSocket | `WS /ws/runs/{id}` | Live run events |
| Health | `GET /health` | Status check |
