---
type: task-note
status: planned
tags:
  - agent
  - tasks
  - phase-2
related:
  - "[[Project State]]"
  - "[[MVP Task List]]"
---

# Phase 2 Task List

> [!summary]
> Phase 2 starts after MVP behavior is truly complete. These tasks expand advanced reporting, collaboration, analytics hardening, runner deployment, governance persistence, observability, and AI-assisted features.

## Reporting And Collaboration

- [ ] Fully integrated Contract report view.
- [ ] Fully integrated Analyst Flow report view.
- [ ] Fully integrated QA Evidence report view.
- [ ] Fully integrated Correlation report view.
- [ ] Report template builder UI.
- [ ] PDF report export.
- [ ] Richer failure-classification workflows.

## Spec Evolution And Analytics

- [ ] Spec diff refinement and breaking-change semantics.
- [ ] Impact analysis grounded in persisted project configuration.
- [ ] Run comparison UX.
- [ ] Latency trend validation against real run history.
- [ ] Flakiness scoring hardening.

## Runner And Infrastructure

- [x] CLI runner binary shipped as `cmd/smokery`.
- [ ] Private or self-hosted runner mode.
- [ ] Runner registration tokens.
- [ ] Scheduled smoke runs.
- [ ] Redis pub/sub for horizontally scaled live events.
- [ ] Slack or Teams notifications.

## Security, Governance, And AI

- [ ] Persist governance data beyond placeholders or in-memory defaults where still needed.
- [ ] Role-based access control.
- [ ] Project-level permissions.
- [ ] Approval workflow for destructive flows.
- [ ] Audit trail hardening.
- [ ] AI-assisted operation classification review.
- [ ] AI-assisted flow suggestion from OpenAPI.
- [ ] AI-assisted failure explanation.

## Related Notes

- [[MVP Task List]]
- [[Project State]]
- [[Product Vision]]

