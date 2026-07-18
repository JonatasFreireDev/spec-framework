# Product Architecture

Shared topology and boundary records are owned by `technical-landscape`; the
Engineering System aggregate references them without duplicating ownership.

## Purpose

Store stable product architecture knowledge and product ADR references:

- `system-context.md`: actors, systems, and external boundaries.
- `modules.md`: module responsibilities and dependencies.
- `data-ownership.md`: sources of truth, consistency, and authorization boundaries.
- `integration-map.md`: APIs, events, queues, and third-party contracts.

Delivery-specific change analysis lives beside its use case in `technical-discovery.md` and links back here.

## Boundary Rule

Framework method decisions do not live here.

Engineering ADRs and engineering-owned decisions live in [`../decisions/`](../decisions/). Index every record in `.product/decisions.json` with `domain: engineering`; this architecture folder remains the stable technical baseline and reference material.
