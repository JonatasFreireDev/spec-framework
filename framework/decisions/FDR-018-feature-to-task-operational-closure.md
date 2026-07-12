# FDR-018: Feature-to-Task Operational Closure

## Snapshot

| Field | Value |
| --- | --- |
| ID | `FDR-018` |
| Status | `approved` |
| Origin EV | `Governance baseline` |
| Date | `2026-07-11` |
| Owner | `Delivery Orchestrator` |

## Context

The canonical Feature-to-Task chain is documented, but proposed graphs cannot be valid before their task paths exist, navigation and validation do not share every discovery gate, traceability remains advisory, and documentary work is not fully resumable. The framework needs operational closure without introducing a PRD artifact.

## Decision

| Concern | Contract |
| --- | --- |
| Product entry | Keep the existing Problem, Vision, Strategy, Domain, Goal, Feature, Use Case, and Specification model. Do not add a PRD or Product Brief source of truth. |
| Graph lifecycle | `draft` and `proposed` graphs may reference task paths not yet materialized; `materialized` and `approved` graphs require every canonical task file. |
| Materialization | A confirmed atomic command creates missing task files and the generated `tasks.md`; it never overwrites existing tasks. |
| Gate order | Feature → Use Case → Specification/contracts → Design → Technical Discovery → Architecture Gate → Implementation Plan → Graph → Tasks. CLI and validator share this order. |
| Traceability | Missing REQ/AC/task/test coverage is advisory in draft and blocking when the downstream artifact is proposed or later. |
| Readiness | A task readiness command returns human and JSON verdicts over gates, dependencies, traceability, scopes, staleness, leases, and technical gates. |
| Resume | Documentary transitions persist state, checkpoints, and handoffs; `guide` explains the current gate without executing an agent. |
| Stage approval | Stage review is a preview; confirmed stage approval creates individual approval records atomically and never approves decisions implicitly. |

## Consequences

| Type | Consequence |
| --- | --- |
| Positive | A Feature can reach an objectively ready task through one consistent state machine. |
| Positive | Graph planning no longer violates the task creation gate. |
| Negative | Graph and stage operations require atomic multi-file rollback and compatibility tests. |
| Follow-up | Automatic agent spawning remains deferred by FDR-017. |

## FRAMEWORK.md Amendments

| Section | Amendment |
| --- | --- |
| `6-8` | Define complete discovery gates and graph materialization lifecycle. |
| `11` | Define progressive traceability and atomic stage approval. |
| `15` | Document materialize, readiness, guide, review, and approve-stage commands. |

## Related Artifacts

| Artifact | Link |
| --- | --- |
| Runtime v2 | [FDR-017](FDR-017-resumable-parallel-runtime.md) |
| Framework method | [FRAMEWORK.md](../../FRAMEWORK.md) |

## Supersedes

- The graph/task creation timing clause of FDR-016; all other clauses remain active.
