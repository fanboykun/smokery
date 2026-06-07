---
type: agent-note
status: current
tags:
  - agent
  - context
related:
  - "[[Project State]]"
  - "[[Architecture]]"
  - "[[Product Vision]]"
---

# Agent Context

> [!summary]
> Smokery is an OpenAPI-driven compiler plus composable smoke-test platform with a real API server, a real standalone CLI, and a broad but unevenly integrated frontend. The biggest current gap is between local or mock-backed UI surfaces and the fuller persisted platform the product docs describe.

## Product Identity

```text
OpenAPI-driven compiler
+ composable smoke-test platform
+ persistent smoke evidence and collaboration
+ standalone CLI smoke runner
```

## Core Pipeline

```text
OpenAPI Spec + Smoke Config
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

## Main Architecture Rule

The compiler bakes everything the runner needs into the SmokePlan. The runner does not read the OpenAPI spec, does not talk to the database, and does not persist results.

## Current Implementation Reality

- The backend ships project, spec, operation, plan preview, run, report, comment, artifact, failure-classification, spec-evolution, analytics, and governance endpoints.
- The CLI ships `import-spec`, `compile`, `run`, and `report`.
- The frontend already has routes for builder, report, analytics, spec evolution, impact, and governance flows.
- Some advanced frontend pages still need real endpoint integration or persistence cleanup.
- Canvas-builder work exists and still compiles back down into canonical flows and suites before execution.

## Critical Concepts

- OpenAPI describes possibilities.
- User configuration selects intent.
- Compiler validates and composes a plan.
- Runner executes compiled plans only.
- Hooks are how the runner is extended.
- Destructive operations are blocked by default.
- Cleanup is first-class.
- List endpoint suites are first-class.
- Run results are persistent and collaborative.
- Reports, analytics, and governance must distinguish shipped backend capability from placeholder UI wiring.

## Related Notes

- [[Project State]]
- [[Product Vision]]
- [[Architecture]]
- [[Engineering Rules]]
- [[MVP Task List]]

