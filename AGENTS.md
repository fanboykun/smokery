# AGENTS.md — Smokery

This file is the primary instruction document for AI agents implementing this repository.

Agents must read this file before making changes. When deeper context is needed, start from `docs/00 Index.md` and `docs/03 Agent/Project State.md`.

---

## 1. Mission

Build **Smokery**, a collaborative OpenAPI-driven smoke testing platform delivered as two products that share one core:

- **API server** (`cmd/server`): UI-managed platform with persistence, live runs, reports, analytics, and governance surfaces.
- **CLI smoke runner** (`cmd/smokery`): standalone executable for local config and CI.

Both products are first-class. Any feature that lives in the domain or application layer must work for both unless it is explicitly delivery-specific.

Core pipeline:

```text
OpenAPI Spec + User-Composed Smoke Configuration
        ↓
Compiler / Composer
        ↓
Executable SmokePlan
        ↓
Runner
        ↓
Persistent RunResult
        ↓
Reports, Diagrams, Debug Views, Collaboration, Insights
```

---

## 2. Hard Documentation Rules

### 2.1 Obsidian Markdown Is Mandatory For `docs/`

All Markdown written under `docs/` **must** be Obsidian-compatible Markdown.

Required for every `docs/**/*.md` file:

- YAML properties at the top.
- Obsidian wikilinks for internal documentation references, for example `[[Architecture]]`.
- Callouts for summaries, warnings, todos, and decisions, for example `> [!summary]`.
- Task lists for implementation work, for example `- [ ]` and `- [x]`.
- Stable note names written for humans, not only file paths.
- A link from `docs/00 Index.md` or a relevant parent note.

Required note pattern for every new or rewritten `docs/**/*.md` file:

```markdown
---
type: <index|product-note|technical-note|agent-note|decision|workflow-note|task-note>
status: <current|planned|reference|legacy-reference|archive>
tags:
  - ...
related:
  - "[[Related Note]]"
---

# Note Title

> [!summary]
> One short summary of the note's purpose and how it fits the vault.
```

When editing an existing `docs/**/*.md` file:

- Preserve or improve the Obsidian pattern. Do not remove frontmatter, summary callouts, or vault links.
- Prefer updating the existing note in place over creating parallel files with slightly different names.
- If the note is legacy or compatibility-only, mark it explicitly in `status` and link back to the canonical vault note.
- If you add a new note, also add its incoming link from `docs/00 Index.md` or the appropriate parent note in `01 Product`, `02 Technical`, `03 Agent`, or `04 Decisions`.

Avoid in `docs/**/*.md`:

- Raw HTML unless the file is explicitly an HTML artifact.
- Unlinked orphan notes.
- Broken wikilinks to notes that do not exist.
- Recreating flat duplicate note trees like the old `docs/ai` layout.

### 2.2 `docs/` Is An Obsidian Vault

The docs directory must follow this structure:

```text
docs/
  00 Index.md
  01 Product/
  02 Technical/
  03 Agent/
  04 Decisions/
```

Use the folders as semantic areas:

- `01 Product/` — product vision, behavior rules, user lifecycle.
- `02 Technical/` — architecture, stack, compiler, runner, repository layout.
- `03 Agent/` — latest context, engineering rules, workflows, task lists.
- `04 Decisions/` — ADR-style decisions.

Any future normalization, migration, or doc rewrite under `docs/` must move the content closer to this vault pattern, never away from it.

### 2.3 Project State Must Always Live In Docs

`docs/03 Agent/Project State.md` is the canonical latest-context file for AI agents.

Agents must update it whenever a change affects:

- current project phase
- completed or remaining work
- architecture constraints
- documentation structure
- implementation context
- major design decisions
- next recommended tasks

Do not rely only on chat history, commit messages, or memory for latest context. Put the state in the docs.

---

## 3. Source Of Truth

Read these before major work:

1. `docs/03 Agent/Project State.md`
2. `docs/00 Index.md`
3. `docs/01 Product/Product Vision.md`
4. `docs/01 Product/Product Rules.md`
5. `docs/02 Technical/Architecture.md`
6. `docs/02 Technical/Compiler Pipeline.md`
7. `docs/02 Technical/Runner Execution.md`
8. `docs/03 Agent/Engineering Rules.md`
9. `docs/03 Agent/MVP Task List.md`
10. `docs/tech-spec.html` when detailed legacy technical-spec context is needed

Priority order when documents conflict:

```text
AGENTS.md
  > docs/03 Agent/Project State.md
  > docs/01 Product/Product Rules.md
  > docs/01 Product/Product Vision.md
  > docs/02 Technical/Architecture.md
  > docs/02 Technical/Compiler Pipeline.md
  > docs/02 Technical/Runner Execution.md
  > docs/03 Agent/Engineering Rules.md
  > docs/03 Agent/* Task List.md
```

When a conflict is found, do not silently choose behavior. Prefer the higher-priority document and mention the conflict in the implementation summary.

---

## 4. Non-Negotiable Architecture Rules

### 4.1 Compiler-First Execution

The UI must not directly execute arbitrary form state.

Required path:

```text
UI / CLI input
  → canonical project configuration
  → compiler validation
  → executable SmokePlan
  → runner execution
  → structured RunResult
```

The runner executes compiled plans only. The compiler bakes everything the runner needs at runtime into the SmokePlan.

### 4.2 Control Plane And Execution Plane Separation

Control plane:

- project management
- spec import and analysis
- operation registry
- environment/auth config
- flow and suite builders
- plan preview
- run orchestration
- reports and collaboration

Execution plane:

- HTTP request execution
- pluggable pre/post processors
- auth injection
- variable interpolation
- capture
- redaction
- trace extraction
- assertions
- structured result generation

### 4.3 Flow And Suite Are Both First-Class

```text
Flow
  Explicit, ordered, scenario-driven.

Suite
  Selected, generated, operation-driven.
```

Do not reduce suites to hand-authored flows. Do not reduce flows to endpoint lists.

### 4.4 Destructive Operation Safety

Treat these as unsafe unless explicitly configured otherwise:

- `POST`
- `PUT`
- `PATCH`
- `DELETE`
- business action endpoints such as approve, submit, recalculate, send, publish, pay, cancel

The compiler must block destructive execution unless the flow, suite, or project explicitly allows it.

### 4.5 Persistent Run Results

Smoke output must not be only terminal logs. Persist structured results for summaries, flows, steps, suites, generated cases, request/response metadata, assertions, captures, traces, comments, failure classification, and artifacts.

---

## 5. Approved Tech Stack

Backend:

- Go 1.26
- `danielgtaylor/huma/v2` + `humaecho`
- `labstack/echo/v4`
- `jackc/pgx/v5`
- `sqlc`
- `golang-migrate`
- `pb33f/libopenapi`
- `tidwall/gjson`
- `spf13/viper`
- `rs/zerolog`
- `google/uuid`
- `gorilla/websocket`
- `minio/minio-go/v7`
- `modernc.org/sqlite`
- `testify`
- `testcontainers-go`

Frontend:

- SvelteKit 2.x / Svelte 5
- TypeScript 6.x strict mode
- `openapi-fetch`
- `openapi-typescript`
- `@tanstack/svelte-query` v6
- TailwindCSS 4
- Mermaid
- LayerChart
- `sveltekit-superforms`

Storage:

- SQLite default for local-first use
- PostgreSQL 18 optional
- Filesystem artifacts default
- MinIO/S3 optional
- Redis only in later horizontal-scaling work

Tooling:

- `bun`
- `air`
- `Makefile`
- Docker Compose
- GitHub Actions

Do not introduce ORM magic. Use `sqlc` for PostgreSQL queries.
