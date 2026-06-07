---
type: decision
status: accepted
date: 2026-06-07
tags:
  - adr
  - architecture
related:
  - "[[Architecture]]"
---

# ADR - Hexagonal Architecture

> [!summary]
> Smokery uses hexagonal architecture so delivery and infrastructure can vary while the domain core stays shared.

## Decision

The system follows `delivery → app → port ← adapter`, with domain services remaining infrastructure-free.

## Related Notes

- [[Architecture]]

