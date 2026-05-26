# Product Rules for AI Agents

These rules define product behavior. Do not violate them during implementation.

## 1. The product is UI-first

The primary product is a UI-managed smoke testing platform.

The CLI/runner is useful, but it is not the main product experience.

The intended user journey is:

```text
Create project
  → import OpenAPI
  → review analysis
  → configure environment/auth
  → compose suites and flows
  → preview plan
  → run
  → review/share reports
```

## 2. The compiler is the product backbone

All smoke execution must go through compiled plans.

The UI can be friendly and visual, but the backend must preserve canonical configuration and deterministic compilation.

## 3. The platform must be useful for multiple roles

Do not design only for backend developers.

Required audiences:

- Backend developers: debug details, request/response, traces, logs.
- Frontend developers: flow sequence, response shapes, dependencies.
- Analysts: readable business flow, expected outcomes, pass/fail evidence.
- QA: scenario evidence, assertions, repeatable proof.
- DevOps: health, latency, environment readiness, trends.
- CI: machine-readable pass/fail and artifacts.

## 4. OpenAPI is helpful but not absolute

OpenAPI can be incomplete or wrong.

The platform must support:

- Operation override
- Response shape override
- Auth override
- Pagination/search/filter override
- Destructive flag override
- Empty result policy override

## 5. Plan preview is mandatory

Users must be able to see what will run before it runs.

Plan preview should show:

- Selected flows
- Selected suites
- Matched operations
- Generated cases
- Destructive steps
- Skipped operations
- Warnings
- Errors
- Environment target
- Concurrency

## 6. Report output is not only pass/fail

Run output should answer:

```text
What was tested?
Why was it tested?
What happened?
What changed?
Who needs to act?
Where can we debug it?
What can other teams learn from it?
```

## 7. Collaboration is first-class

The platform should persist comments, failure labels, issue states, and history.

Supported failure states should include:

```text
new
acknowledged
investigating
expected
fixed
ignored
needs-config-update
```

## 8. Safety is first-class

Destructive operations require explicit permission.

Sensitive data must be redacted.

Secrets must never be stored in plain exported config.
