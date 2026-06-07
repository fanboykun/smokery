---
type: technical-note
status: current
tags:
  - technical
  - runner
related:
  - "[[Architecture]]"
  - "[[Compiler Pipeline]]"
---

# Runner Execution

> [!summary]
> The runner is a pure execution library. It accepts a compiled `SmokePlan`, runs it through hooks, and returns a structured `RunResult`.

## Runner Contract

- Input: compiled `SmokePlan`
- Optional input: runtime options and hooks
- Output: structured `RunResult`

## Hard Constraints

- The runner does not import `libopenapi`.
- The runner does not import database packages.
- The runner does not import blob-storage adapters.
- The runner does not orchestrate jobs or persistence.

## Hook Model

The runner is extended through pre- and post-processors for:

- auth injection
- interpolation
- capture
- redaction
- trace extraction
- assertions

## Current-State Notes

- The same runner is used by the API server and CLI.
- Live server execution is orchestrated outside the runner by in-process worker infrastructure.

## Related Notes

- [[Architecture]]
- [[Compiler Pipeline]]
- [[Product Vision]]

