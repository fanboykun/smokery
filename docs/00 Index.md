---
type: index
status: current
tags:
  - docs
  - obsidian-vault
  - smokery
related:
  - "[[Product Vision]]"
  - "[[Architecture]]"
  - "[[Project State]]"
  - "[[Engineering Rules]]"
---

# Smokery Documentation Index

> [!summary]
> This `docs/` directory is the single current documentation vault for Smokery. Flat root markdown notes and the old `docs/ai` tree are deprecated and should not be recreated.

## 01 Product

- [[Product Vision]] — product concept, current shipped reality, and target direction.
- [[Product Rules]] — product behavior rules that implementation work must not violate.
- [[Smoke Testing Lifecycle]] — pre-smoke, execution, and post-smoke lifecycle with current implementation status.

## 02 Technical

- [[Architecture]] — actual runtime architecture, adapter defaults, and delivery surfaces.
- [[Tech Stack]] — approved libraries, actual defaults, and current installed frontend tooling.
- [[Compiler Pipeline]] — deterministic config-to-plan compiler rules.
- [[Runner Execution]] — runner library contract, hook pipeline, and execution behavior.
- [[Repository Layout]] — current repository shape and important subtrees.
- [[Database Migrations]] — migration workflow and local-vs-Postgres operational notes.

## 03 Agent

- [[Project State]] — latest project state, current phase, open work, and documentation guidance.
- [[Agent Context]] — condensed implementation context for agents.
- [[Engineering Rules]] — architecture, frontend, backend, and documentation rules.
- [[Config Persistence Blockers]] — technical blockers and redesign plan for configuration persistence.
- [[Canvas Schema and Property Linking]] — technical note on handling deep schemas and cross-plan variables.
- [[Config Redesign Task List]] — detailed task checklist for implementing configuration persistence and canvas refactoring.
- [[MVP Task List]] — what is truly left for MVP based on the current repo.
- [[Phase 2 Task List]] — post-MVP roadmap.
- [[Implementation Workflows]] — current implementation flow from spec import through preview, run, analytics, governance, and report follow-up.

## 04 Decisions

- [[ADR - Use Obsidian Vault Docs]]
- [[ADR - Compiler First Execution]]
- [[ADR - Two Products One Core]]
- [[ADR - Hexagonal Architecture]]
- [[ADR - Config Builder Redesign]]

## Non-Markdown Artifacts

- `docs/tech-spec.html` — legacy technical spec artifact kept for reference.

> [!warning]
> Do not create isolated Markdown files in `docs/`. Every new note must include YAML properties, use Obsidian-compatible internal links, and be reachable from this index or a parent note.

> [!note]
> Minimum note pattern for every `docs/**/*.md` write:
> frontmatter with `type`, `status`, `tags`, and `related`;
> a single `#` note title;
> a `> [!summary]` callout near the top;
> and at least one vault link from this index or a relevant parent note.

> [!important]
> Deprecated note sets have been removed. Do not recreate `docs/ai`, root duplicate notes like `docs/architecture.md`, or any parallel markdown tree outside the vault structure unless a new decision note explicitly changes the documentation architecture.

