---
type: workflow-note
status: current
tags:
  - agent
  - workflow
  - implementation
related:
  - "[[Agent Context]]"
  - "[[Engineering Rules]]"
  - "[[Project State]]"
---

# Implementation Workflows

> [!summary]
> This note replaces the older scattered `docs/ai/workflows` notes with one current implementation workflow note anchored to the actual repo surface.

## Core Implementation Flow

```text
Import spec
  → review operations
  → configure environments/auth/flows/suites/canvas
  → preview compiled plan
  → start run
  → inspect result, comments, artifacts, reports, analytics, and governance follow-up
```

## Current Workflow Boundaries

- Spec import and operation review are backend-backed today.
- Builder and config editing are frontend-real, but some state is still browser-local.
- Plan preview and run creation already go through the backend compiler and runner.
- Backend phase-2 services for analytics, governance, spec evolution, and failure classification exist.
- Not every frontend page is fully wired to those newer backend endpoints yet.

## Agent Rule

> [!important]
> When updating implementation docs, describe the real current workflow first, then list the next missing step. Do not flatten planned behavior into present-tense current-state docs.

## Related Notes

- [[Project State]]
- [[Agent Context]]
- [[MVP Task List]]

