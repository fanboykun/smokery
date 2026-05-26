# Workflow 06 — Quality, Testing, and Review

## Goal

Ensure AI-generated implementation is safe, testable, and aligned with the platform architecture.

## Backend testing checklist

Test these areas:

- OpenAPI parse success and failure.
- Operation classification.
- Operation override application.
- Compiler validation errors.
- Destructive operation blocking.
- Flow step compilation.
- Suite operation selection.
- List case generation.
- Variable interpolation.
- JSONPath capture.
- Status assertion.
- JSON Schema assertion.
- List response shape detection.
- Empty result policy.
- Pagination sanity.
- Search-from-response behavior.
- Enum filter behavior.
- Trace context extraction.
- Redaction behavior.
- Persistence integration with PostgreSQL.

## Frontend testing/checking checklist

- TypeScript typecheck passes.
- Forms validate expected payloads.
- Compile errors render clearly.
- Plan preview handles warnings/errors.
- Run event store applies events correctly.
- Report components tolerate missing optional fields.

## Review checklist for agents

Before completing a task, verify:

- The change follows compiler-first architecture.
- Runner does not depend on raw UI config.
- Destructive operations remain protected.
- Secrets are not stored or logged.
- New database changes have migrations.
- New queries are sqlc-compatible.
- New APIs have request validation.
- New UI states handle loading/error/empty.
- New result structures are persistent and reportable.

## Completion summary

Use this format:

```text
Implemented:
- ...

Changed files:
- ...

Validation:
- ...

Notes / follow-ups:
- ...
```
