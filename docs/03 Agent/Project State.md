---
type: project-state
status: current
last_updated: 2026-06-07
tags:
  - agent
  - project-state
  - context
related:
  - "[[Agent Context]]"
  - "[[MVP Task List]]"
  - "[[Phase 2 Task List]]"
---

# Project State

> [!summary]
> This note is the mandatory latest-context file for AI agents. The current repo already has a substantial working implementation, but the older flat docs and `docs/ai` tree no longer describe the current state accurately enough.

## Current State Snapshot

Smokery is an OpenAPI-driven smoke testing platform delivered as two products:

- API server and UI-managed platform
- standalone `smokery` CLI runner

The repository now includes:

- a real Huma + Echo API server
- a real standalone CLI
- SQLite and filesystem defaults for local-first usage
- optional PostgreSQL and MinIO adapters
- a SvelteKit frontend with project, spec, operations, builder, plan, runs, comments, analytics, spec evolution, impact, reports, and governance routes
- backend services for failure classification, spec evolution, analytics, and governance

The major remaining MVP work is finishing end-to-end product behavior, especially persistence and frontend integration, rather than creating missing pages from scratch.

## Current Main Phase

```text
Phase C.FE — End-to-end MVP completion and doc consolidation
```

## Remaining Areas

- Persist environments, auth profiles, flows, suites, and canvas state outside browser-local storage.
- Replace or complete mock-backed frontend integrations where backend endpoints already exist.
- Harden builder-to-preview-to-run flows across frontend, API, and CLI.
- Verify artifact, classification, analytics, and governance flows against the current adapters.
- Continue UI polish and validation.

## Documentation State

The `docs/` directory is organized as an Obsidian-style vault:

```text
docs/
  00 Index.md
  01 Product/
  02 Technical/
  03 Agent/
  04 Decisions/
```

Deprecated documentation sets are being replaced:

- `docs/ai/` is no longer the canonical documentation system.
- flat root duplicate notes like `docs/architecture.md` and `docs/design-docs.md` are deprecated in favor of the vault notes.

## Hard Context Rules

- [[Project State]] is the latest project context source.
- Agents must update this note after any change that affects implementation status, design direction, open tasks, docs layout, or architecture constraints.
- Agents must link any new documentation note from [[Smokery Documentation Index]] or from a relevant parent note.
- Agents must prefer Obsidian wikilinks for internal docs references.
- Every `docs/**/*.md` write must preserve the Obsidian note pattern: frontmatter, note title, summary callout, and related vault links.
- Doc normalization is one-way: agents may improve legacy docs toward the vault pattern, but must not reintroduce plain standalone markdown outside that pattern.

## Recent Context Log

### 2026-06-13 — herdr/kiro workflow, Config Redesign ADR, Technical Notes, and Tasks documented
- Documented the exact workflow for spawning `kiro-cli chat` inside a new pane in `herdr` to ensure instant setup in future agent sessions.
- Created `ADR - Config Builder Redesign.md` in `docs/04 Decisions/` detailing the hierarchical canvas workflow, OpenAPI spec drawer, and connection validation choices.
- Created `Config Persistence Blockers.md` in `docs/03 Agent/` detailing database, repository, HTTP endpoint, and frontend changes needed to persist configurations.
- Created `Canvas Schema and Property Linking.md` in `docs/03 Agent/` detailing the solutions for the schema slicing blocker (using pinned handles and drawer-based connections) and variable-based plan dependencies.
- Created `Config Redesign Task List.md` in `docs/03 Agent/` declaring all database, repository, API, and frontend tasks mapped directly to decisions and blockers.

### 2026-06-07 — Vault rework restored after reset
- Recreated the Obsidian vault structure after the markdown rework was lost in a reset.
- Rewrote the current-state notes around the actual repo surface, including phase-2 backend services and expanded frontend routes.
- Re-established the rule that `docs/` is the single current documentation system and that the old flat docs plus `docs/ai` tree are deprecated.

