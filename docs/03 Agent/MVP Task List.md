---
type: task-note
status: current
tags:
  - agent
  - tasks
  - mvp
related:
  - "[[Project State]]"
  - "[[Phase 2 Task List]]"
---

# MVP Task List

> [!summary]
> Backend, compiler, runner, CLI, and most core frontend routes exist today. The remaining MVP gap is no longer “build the pages”, but “finish the end-to-end product behavior behind them”.

## Status Snapshot

| Phase | Description | Status |
|---|---|---|
| 0 | Repository foundation | Done |
| 1 | Backend and database foundation | Done |
| 2 | OpenAPI ingestion | Done |
| 3 | Canonical config model | Done |
| 4 | Compiler MVP | Done |
| 5 | Runner MVP | Done |
| 6 | Jobs and live run streaming | Done |
| 7 | Report MVP backend surface | Partial |
| 8 | Frontend route and builder foundation | Done |
| 9 | Quality hardening | Partial |
| A | Ports and adapters refactor | Done |
| B | CLI smoke runner | Done |
| C | End-to-end MVP completion | In progress |

## Current Remaining MVP Work

- [ ] Persist project configuration beyond browser-local storage so environments, auth profiles, flows, suites, and canvas metadata are not local-only.
- [ ] Align frontend project-configuration flows with backend canonical config and persistence expectations.
- [ ] Replace mock-backed report, analytics, or governance pages with real backend-powered views where endpoints already exist.
- [ ] Verify and harden end-to-end compile, run, comments, artifact, failure-classification, and report flows across API and frontend.
- [ ] Continue UI polish, loading states, responsive fixes, and final builder usability cleanup.
- [ ] Run and stabilize `make test`, `make check`, and `make lint` as the minimum quality gate.

## Related Notes

- [[Project State]]
- [[Engineering Rules]]
- [[Phase 2 Task List]]
- [[ADR - Config Builder Redesign]]
- [[Config Redesign Task List]]

