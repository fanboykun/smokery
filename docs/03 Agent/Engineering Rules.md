---
type: agent-note
status: current
tags:
  - agent
  - engineering
  - rules
related:
  - "[[Project State]]"
  - "[[Architecture]]"
  - "[[Tech Stack]]"
---

# Engineering Rules

> [!summary]
> These are implementation rules for agents working on Smokery. Preserve boundaries, keep changes coherent, and update [[Project State]] when the work changes project context.

## General Rules

- Keep code explicit and readable.
- Prefer simple, testable functions.
- Avoid broad rewrites unrelated to the task.
- Preserve module boundaries.
- Add or update tests where practical.
- Keep API DTOs stable and documented.
- All Go binaries output to `tmp/`.
- Use `make` targets for dev commands.

## Layering Rules

```text
delivery → app → port ← adapter
              │
              ▼
            domain
```

Rules:

- `domain` imports nothing from this project except other domain packages.
- `port` defines interfaces and imports `model` only.
- `adapter` implements ports using infrastructure SDKs.
- `app` orchestrates use cases and never imports adapters.
- `delivery` is the entry layer for HTTP and CLI.
- `cmd/*` is dependency injection only.

## Frontend Rules

- Use SvelteKit 2 with Svelte 5 runes and TypeScript strict mode.
- Use `openapi-fetch` for type-safe API calls.
- Use `@tanstack/svelte-query` v6 for server state.
- Avoid handwritten `any`.
- Use `bun` as package manager and runtime.
- Be explicit about whether a page is backed by real APIs, browser-local state, or mock data.

## Documentation Rules

- Any Markdown under `docs/` must use Obsidian-compatible syntax.
- Any change that alters project context must update [[Project State]].
- New docs must be linked from [[Smokery Documentation Index]] or a relevant parent note.
- Do not recreate `docs/ai` or flat duplicate notes when canonical vault notes already exist.
- If a note becomes stale, migrate its useful content into the canonical vault note and delete the duplicate.

## Related Notes

- [[Architecture]]
- [[Compiler Pipeline]]
- [[Runner Execution]]
- [[Project State]]

