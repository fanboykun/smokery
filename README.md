# Smokery

A collaborative OpenAPI-driven smoke testing platform delivered as **two products**:

- **API server** — UI-managed platform with persistence, live runs, reports, analytics, and governance surfaces.
- **CLI smoke runner** (`smokery`) — standalone executable for local config and CI.

Both products share the same compiler, runner, and domain code. The CLI is not a stripped-down API.

## Core Idea

```text
OpenAPI Spec + Smoke Configuration
        ↓
Compiler / Composer
        ↓
Executable SmokePlan
        ↓
Runner
        ↓
Persistent RunResult
        ↓
Reports, Diagrams, Debug Views, Collaboration
```

## Documentation

The `docs/` directory is an Obsidian-style documentation vault.

Start here:

- `docs/00 Index.md` — documentation entry point.
- `docs/03 Agent/Project State.md` — latest project state and implementation context.
- `AGENTS.md` — hard rules for AI agents.

Vault structure:

```text
docs/
  00 Index.md
  01 Product/
  02 Technical/
  03 Agent/
  04 Decisions/
  tech-spec.html
```

Markdown under `docs/` intentionally uses Obsidian features such as YAML properties, `[[wikilinks]]`, task lists, and callouts.

## Current Runtime Defaults

- Database: SQLite by default, PostgreSQL optional
- Artifacts: filesystem by default, MinIO optional
- Backend: Go 1.26, Echo, Huma, sqlc, libopenapi, zerolog
- Frontend: SvelteKit 2, Svelte 5, TypeScript 6, TanStack Query, TailwindCSS 4

## Development Principle

This project is expected to be implemented primarily by AI agents. Agents must follow `AGENTS.md`, preserve the compiler-first architecture, and keep `docs/03 Agent/Project State.md` updated when project context changes.

## Development

```bash
make install
make generate
make build
make dev
make test
make check
```
