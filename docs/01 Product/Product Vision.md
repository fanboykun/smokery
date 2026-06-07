---
type: product-note
status: current
tags:
  - product
  - vision
  - smokery
related:
  - "[[Smoke Testing Lifecycle]]"
  - "[[Product Rules]]"
  - "[[Compiler Pipeline]]"
---

# Product Vision

> [!summary]
> Smokery is a collaborative OpenAPI-driven smoke testing platform. The repo already ships a working API server, standalone CLI, and broad Svelte frontend surface, but the product is still closing the gap between browser-local builder state, persisted collaboration, and fully backend-powered report, analytics, and governance experiences.

## Product Identity

```text
OpenAPI-driven compiler
+ composable flow and suite runner
+ persistent smoke evidence and collaboration
+ standalone CLI smoke runner
```

## Core Pipeline

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

## Current Product Reality

- The backend and CLI share the same compiler and runner.
- The frontend already includes project overview, spec import, operations, builder, plan preview, run detail, analytics, spec evolution, impact, and governance routes.
- Builder configuration is still partly browser-local today through frontend state, while preview and run execution already use the backend compile-and-run path.
- The backend already exposes phase-2 surfaces for analytics, spec evolution, governance, and failure classification.
- Some frontend pages remain mock-backed or only partially integrated with the newer backend endpoints.

## Core Value

Smokery should become living API behavior documentation. Every run should answer what was tested, why it was tested, what happened, where failures can be debugged, and what other teams can learn from it.

> [!note]
> When product docs describe a capability, they should make clear whether it is shipping now, partially integrated, or still planned.

## Related Notes

- [[Smoke Testing Lifecycle]]
- [[Product Rules]]
- [[Architecture]]
- [[Compiler Pipeline]]
- [[Runner Execution]]

