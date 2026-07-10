# FDR-001: Decision governance: product vs framework

## Snapshot

| Field | Value |
| --- | --- |
| ID | `FDR-001` |
| Status | `approved` |
| Origin EV | `governance baseline` |
| Date | `2026-07-10` |
| Owner | `Evolution Orchestrator` |

## Context

This repository is both the framework laboratory and a template for product repositories. Without an explicit split, method decisions can be mixed with product decisions and then inherited by adopting products as if they were business facts.

Product decisions must remain reserved for product scope: domain behavior, business rules, data, security, privacy, payment, permissions, and delivery commitments.

Framework decisions govern the method: skill contracts, gates, validator behavior, artifact movement, safe parallelism, QA independence, commit policy, and workflow routing.

## Decision

Use two decision layers:

| Layer | Location | Scope |
| --- | --- | --- |
| Product decision | `knowledge/decisions/DEC-*` | Product domain, feature, use case, business rule, data, permissions, privacy, payment, security, or delivery commitment. |
| Framework decision | `engineering/decisions/FDR-*` or explicit `FRAMEWORK.md` amendment | Framework method, skill contract, gates, validator behavior, orchestration, and workflow policy. |

No framework-method decision may be recorded in `knowledge/decisions/`.

When the level is unclear, ask: does this trace to a product problem, domain, feature, or use case? If not, it is an FDR or a `FRAMEWORK.md` amendment, not a DEC.

## Consequences

| Type | Consequence |
| --- | --- |
| Positive | Adopting products inherit a clean product decision log. |
| Positive | Framework evolution remains auditable without polluting product history. |
| Negative | Contributors must classify decisions before writing them. |
| Follow-up | Update `FRAMEWORK.md` section 12 and keep this FDR index current. |

## FRAMEWORK.md Amendments

| Section | Amendment |
| --- | --- |
| `12. Decisoes` | Clarify that product decisions live in `knowledge/decisions/`, while framework decisions live in `engineering/decisions/FDR-*` or as amendments to `FRAMEWORK.md`. |

## Related Artifacts

| Artifact | Link |
| --- | --- |
| FRAMEWORK | [../../FRAMEWORK.md](../../FRAMEWORK.md) |
| FDR index | [README.md](README.md) |

## Supersedes

- `N/A`
