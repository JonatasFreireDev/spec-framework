# FDR-016: Delivery Closure and Operational Workspaces

## Snapshot

| Field | Value |
| --- | --- |
| ID | `FDR-016` |
| Status | `approved` |
| Origin EV | `Governance baseline` |
| Date | `2026-07-11` |
| Owner | `Evolution Orchestrator` |

## Context

The framework defines a strong Domain-to-Code artifact chain but leaves operational gaps between Domain and feature selection, uses stale numbered skill handoffs, has no friendly work/approval navigation, stores a monolithic specification, lacks mandatory codebase discovery, models but cannot operate an execution graph, and requires commits too early for independent QA of working-tree changes.

## Decision

The framework adopts the following delivery flow:

```text
Domain -> Domain Evolution -> Feature Selection -> New Feature -> Use Cases
-> Specification Contracts -> Design -> Technical Discovery -> Architecture Gate
-> Plan -> Graph -> Tasks -> Code Runner -> Code Review -> QA
-> Commit Crafter -> PR Finalizer -> Security/Release gates
```

| Area | Contract |
| --- | --- |
| Domain evolution | A new orchestrator coordinates goals, journeys, opportunities, candidate features, impact, slicing, and selection. |
| Work selection | Concurrent `.product/workspaces/WORK-NNN.json` records select a scoped Domain/Goal/Feature and current step. No global active feature exists. |
| Navigation | `work`, `status`, and `next` resolve paths or scoped IDs and report readiness without inventing content. |
| Approval | `approve` previews and records an explicit human grant with normalized SHA-256; it cannot bypass parent blockers. |
| Specification | `specification.md` remains the root contract; modular files live in `contracts/` and are required proportionally by rigor. |
| Technical discovery | Delivery-specific `technical-discovery.md` maps requirements to the existing codebase; stable architecture remains in `engineering/`. |
| Architecture | Architectural impact must reference an approved decision or explicitly record `Not required` with rationale before planning. |
| Traceability | Stable `REQ-*`, `AC-*`, task, test, and evidence references form a mechanically validated coverage chain. |
| Graph runtime | `graph ready/claim/release/complete` operates task readiness and exclusive claims without automatically executing code. |
| Implementation state | `implemented` requires working-tree evidence and diff hash, not commits. Any diff change stales review/QA. |
| Validation state | `validated` requires green Code Review and QA plus commits, code paths, tests, gates, and PR when policy requires it. |
| Gates | Applicable `TBD` gate commands block Code Runner readiness; explicit `N/A` with rationale is allowed. |
| Handoffs | Skill references use canonical skill names and are validated against the installed skill registry. |

## Consequences

| Type | Consequence |
| --- | --- |
| Positive | Humans and agents can select, inspect, approve, and advance one feature without guessing the next artifact. |
| Positive | QA can verify immutable working-tree evidence before Commit Crafter packages the change. |
| Positive | Large specifications remain navigable without creating competing sources of truth. |
| Negative | More schemas, CLI commands, and validator rules must remain synchronized. |
| Follow-up | Automatic multi-agent graph execution remains out of scope until claim semantics mature. |

## FRAMEWORK.md Amendments

| Section | Amendment |
| --- | --- |
| `6` | Add modular contracts, requirement IDs, technical discovery, and architecture gate. |
| `8` | Add graph runtime and claims. |
| `9-10` | Add Domain Evolution and Technical Discovery skills and normalize handoffs. |
| `11` | Move commit evidence from implemented to validated and define diff-hash staleness. |
| `15` | Add work, status, next, approve, gates, and graph commands. |

## Related Artifacts

| Artifact | Link |
| --- | --- |
| Framework method | [FRAMEWORK.md](../../FRAMEWORK.md) |
| Domain Evolution | [skill](../skills/domain-evolution-orchestrator/SKILL.md) |
| Technical Discovery | [skill](../skills/technical-discovery/SKILL.md) |
| CLI | [app.go](../../internal/cli/app.go) |
| Validator | [rules.go](../../internal/validator/rules.go) |

## Supersedes

- Amend implementation evidence portions of `FDR-005` and `FDR-008`; all other clauses remain active.
