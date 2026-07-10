# FDR-002: Gate commands are product conventions

## Snapshot

| Field | Value |
| --- | --- |
| ID | `FDR-002` |
| Status | `approved` |
| Origin EV | `EV-011 / gate governance` |
| Date | `2026-07-10` |
| Owner | `Documentation Orchestrator` |

## Context

Code Runner and QA need to run technical gates, but the framework must not hardcode commands because each adopting product has its own stack, package manager, CI, database, and visual verification tools.

## Decision

Each product declares its technical gates in `knowledge/conventions/gates.md`.

Skills that execute or verify implementation work must read this file instead of embedding command names in their skill contract. In this framework repository, the gates are placeholders because there is no real product stack.

## Consequences

| Type | Consequence |
| --- | --- |
| Positive | The framework stays stack-agnostic. |
| Positive | QA and future Code Runner can use the same source of truth for gates. |
| Negative | A product that does not fill `gates.md` cannot produce strong implementation evidence. |
| Follow-up | Future product repositories must replace placeholders before implementation gates are enforced. |

## FRAMEWORK.md Amendments

| Section | Amendment |
| --- | --- |
| `11. Gates De Aprovacao` | Reference `knowledge/conventions/gates.md` as the product-specific source of technical gate commands. |

## Related Artifacts

| Artifact | Link |
| --- | --- |
| Gates convention | [../../knowledge/conventions/gates.md](../../knowledge/conventions/gates.md) |
| FRAMEWORK | [../../FRAMEWORK.md](../../FRAMEWORK.md) |

## Supersedes

- `N/A`
