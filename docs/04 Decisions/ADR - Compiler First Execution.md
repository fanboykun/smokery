---
type: decision
status: accepted
date: 2026-06-07
tags:
  - adr
  - compiler
related:
  - "[[Compiler Pipeline]]"
  - "[[Runner Execution]]"
---

# ADR - Compiler First Execution

> [!summary]
> All smoke execution goes through compiled plans. The runner executes compiled plans only.

## Decision

Smokery uses compiler-first execution. UI and CLI inputs must be normalized and compiled into a `SmokePlan` before execution.

## Related Notes

- [[Compiler Pipeline]]
- [[Runner Execution]]

