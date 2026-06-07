---
type: product-note
status: current
tags:
  - product
  - rules
related:
  - "[[Product Vision]]"
  - "[[Smoke Testing Lifecycle]]"
---

# Product Rules

> [!summary]
> These rules define product behavior. Implementation work must not violate them, even if the current repo only partially realizes the full product surface.

## The Product Is UI-Managed But Not UI-Only

The primary product is a UI-managed smoke testing platform, but the CLI is also first-class and must stay aligned with the same domain behavior.

## The Compiler Is The Backbone

All smoke execution must go through compiled plans. UI friendliness does not justify bypassing canonical configuration or deterministic compilation.

## Reports Must Serve Multiple Roles

Do not design only for backend developers.

Required audiences:

- backend developers
- frontend developers
- analysts
- QA
- DevOps
- CI

## OpenAPI Is Helpful But Not Absolute

OpenAPI can be incomplete or wrong. The platform must support overrides for operation classification, response shape, auth, pagination, search, filter, destructive flags, and empty-result policy.

## Plan Preview Is Mandatory

Users must be able to see what will run before it runs, including errors, warnings, destructive steps, selected flows, selected suites, and target environment details.

## Safety Is First-Class

Destructive operations require explicit permission. Sensitive data must be redacted. Secrets must never be stored in plain exported config.

## Related Notes

- [[Product Vision]]
- [[Smoke Testing Lifecycle]]
- [[Compiler Pipeline]]

