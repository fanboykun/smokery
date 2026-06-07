---
type: decision
status: accepted
date: 2026-06-07
tags:
  - adr
  - product
related:
  - "[[Product Vision]]"
  - "[[Architecture]]"
---

# ADR - Two Products One Core

> [!summary]
> Smokery ships both an API server and a CLI runner on top of one shared domain core.

## Decision

Features that live in the domain or app layer should work for both products unless explicitly delivery-specific.

## Related Notes

- [[Product Vision]]
- [[Architecture]]

