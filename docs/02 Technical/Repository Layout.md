---
type: technical-note
status: current
tags:
  - technical
  - repository
related:
  - "[[Architecture]]"
  - "[[Tech Stack]]"
---

# Repository Layout

> [!summary]
> Smokery is a monorepo with one Go workspace for both backend products, one SvelteKit frontend, shared generated types, infrastructure configs, and an Obsidian-style documentation vault. The old `docs/ai` tree and flat root markdown notes are no longer part of the intended layout.

## Expected Shape

```text
.
├── Makefile
├── AGENTS.md
├── README.md
├── configs
├── apps
│   ├── core
│   └── web
├── packages
│   └── types
├── docs
│   ├── 00 Index.md
│   ├── 01 Product
│   ├── 02 Technical
│   ├── 03 Agent
│   ├── 04 Decisions
│   └── tech-spec.html
└── tmp
```

## Current Important Subtrees

- `apps/core/cmd/server` — API server entrypoint.
- `apps/core/cmd/smokery` — standalone CLI entrypoint.
- `apps/core/internal/adapter/sqlite` — default local relational adapter.
- `apps/core/internal/adapter/postgres` — sqlc-backed PostgreSQL adapter.
- `apps/core/internal/adapter/inproc` — current worker, event bus, and retention loop.
- `apps/core/internal/frontend` — embedded frontend assets for production server builds.
- `apps/web/src/routes` — current SvelteKit route tree including builder, reports, analytics, spec evolution, and governance surfaces.

## Core Rule

`apps/core` is one Go module that produces multiple binaries. Do not create separate backend modules for server and CLI unless the architecture is explicitly changed.

## Docs Rule

The `docs/` directory is an Obsidian vault. New Markdown files must follow [[ADR - Use Obsidian Vault Docs]] and must be linked from [[Smokery Documentation Index]] or a parent note.

## Related Notes

- [[Architecture]]
- [[Tech Stack]]
- [[Engineering Rules]]

