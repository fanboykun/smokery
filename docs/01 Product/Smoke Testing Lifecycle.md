---
type: product-note
status: current
tags:
  - product
  - lifecycle
related:
  - "[[Product Vision]]"
  - "[[Compiler Pipeline]]"
  - "[[Runner Execution]]"
---

# Smoke Testing Lifecycle

> [!summary]
> Smokery has three lifecycle stages: pre-smoke configuration, compiled execution, and post-smoke reporting or collaboration. In the current repo, all three exist, but several post-smoke and governance experiences are still unevenly integrated across backend and frontend.

## Current Lifecycle Status

- Pre-smoke: active through spec import, operation review, environment/auth config, flow editing, suite editing, plan preview, and canvas-first builder work.
- Execution: active through the shared runner, CLI commands, API run endpoints, and websocket event stream.
- Post-smoke: active through run history, run detail, comments, artifacts, debug/CI/Mermaid reports, failure classification, analytics, and spec-evolution APIs, though some frontend views still rely on mock integrations.

## Pre-Smoke

Pre-smoke includes spec import, operation analysis, environment and auth configuration, suites, flows, validation, and plan preview.

## Execution

```text
Canonical Project Config
        ↓
Compiler Validation
        ↓
Self-contained SmokePlan
        ↓
Runner Execution
        ↓
Structured RunResult
```

> [!important]
> The UI must never directly execute arbitrary form state.

## Post-Smoke

Post-smoke output should produce collaborative artifacts: run summary, flow results, step results, suite results, generated case results, request and response metadata, assertions, captures, traces, comments, failure classification, analytics, governance follow-up, report artifacts, diagrams, and CI summaries.

## Related Notes

- [[Product Vision]]
- [[Product Rules]]
- [[Compiler Pipeline]]
- [[Runner Execution]]

