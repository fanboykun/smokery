---
type: technical-note
status: current
tags:
  - technical
  - compiler
related:
  - "[[Architecture]]"
  - "[[Runner Execution]]"
  - "[[Product Rules]]"
---

# Compiler Pipeline

> [!summary]
> The compiler turns canonical project configuration into a self-contained `SmokePlan`. It is the boundary between OpenAPI knowledge and runtime execution.

## Required Flow

```text
UI / CLI input
  → canonical project configuration
  → compiler validation
  → executable SmokePlan
  → runner execution
```

## Rules

- The compiler bakes everything the runner needs into the plan.
- The compiler emits structured errors and warnings.
- The compiler blocks unsafe destructive execution unless explicitly allowed.
- The compiler is where OpenAPI-driven reasoning happens, not the runner.

## Current-State Notes

- Builder, plan preview, and run-start flows already use the backend compile path.
- Current documentation must distinguish between browser-local config editing and canonical compiled execution.

## Related Notes

- [[Architecture]]
- [[Runner Execution]]
- [[Product Rules]]

