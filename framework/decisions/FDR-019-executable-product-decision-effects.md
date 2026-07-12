# FDR-019: Executable Product Decision Effects

## Snapshot

| Field | Value |
| --- | --- |
| ID | `FDR-019` |
| Status | `approved` |
| Origin EV | `Governance baseline` |
| Date | `2026-07-12` |
| Owner | `Product Historian` |

## Context

Product decisions can unblock the Architecture Gate, but their type, scope, approval integrity, downstream propagation, and operational effects are not one mechanically enforceable contract. Introducing a separate ADR store would create competing decision sources.

## Decision

Keep `DEC-*` as the single product decision identity. An architectural ADR is a DEC with `type: architecture`.

| Concern | Contract |
| --- | --- |
| Type | `product`, `architecture`, `security`, `data`, or `delivery`. |
| Scope | Artifact IDs or product-relative path prefixes; effects apply only within that scope. |
| Validity | A referenced DEC must exist, be indexed, be approved, and have a current approval record matching its content hash. |
| Architecture Gate | `DEC-*` references unblock the gate only when valid and scope-compatible; otherwise the workflow remains blocked. |
| Effects | Structured `workflowEffects` may require task types, gates, evidence, write scopes, and shared resources. Empty effects are valid. |
| Propagation | Technical Discovery, Plan, Graph, and affected Tasks reference applicable decisions. Missing propagation blocks proposed-or-later work. |
| Staleness | Derivations must include the DEC as a source; changing it makes affected downstream artifacts stale. |
| Commands | Decision text never executes. Command Planner consumes only configured gates required by validated structured effects. |
| Impact | `impact --decision DEC-NNN` reports validity, affected/referencing/stale artifacts, effects, and propagation gaps. |

## Consequences

| Type | Consequence |
| --- | --- |
| Positive | Architecture and design constraints become inspectable workflow inputs without a second ADR hierarchy. |
| Positive | Task readiness can reject incomplete decision consequences before code execution. |
| Negative | Existing decision indexes may require gradual enrichment with type, scope, and empty effects. |
| Follow-up | Automated task generation from decisions remains prohibited. |

## FRAMEWORK.md Amendments

| Section | Amendment |
| --- | --- |
| `7` | Strengthen Architecture Gate decision validity and propagation. |
| `11-12` | Define typed scoped decisions and structured workflow effects. |
| `15` | Add decision impact inspection. |

## Related Artifacts

| Artifact | Link |
| --- | --- |
| Feature-to-Task closure | [FDR-018](FDR-018-feature-to-task-operational-closure.md) |
| Framework method | [FRAMEWORK.md](../../FRAMEWORK.md) |

## Supersedes

- N/A
