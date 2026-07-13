# FDR-031: Feature-Scoped Bootstrap

## Snapshot

| Field | Value |
| --- | --- |
| ID | `FDR-031` |
| Status | `proposed` |
| Origin EV | `Adopter first-run feedback` |
| Date | `2026-07-12` |
| Owner | `Product Orchestrator` |

## Context

The full five-artifact Foundation package is proportionate for a new or uncertain product but creates duplicated, low-value documents when an adopter starts from one bounded feature. Keeping only Strategy would remove the user problem, desired outcome, constraints, and success signal that justify delivery.

## Decision

| Boundary | Contract |
| --- | --- |
| Applicability | `existing-feature` uses one canonical `foundation/feature-brief.md` instead of requiring Problem, Vision, Product Principles, North Star, and Strategy approvals. |
| Required content | The brief owns problem, desired outcome, scope, non-goals, principles and constraints, success signal, delivery strategy, evidence, and decisions needed. |
| Feature binding | The brief declares one target Feature id or path. Its approval cannot unlock a different registered feature. |
| Registry | The brief is registered as `feature-brief`; full Foundation artifacts remain available as reference skeletons but are excluded from the active registry for this starting point. |
| Approval | The brief requires an individual hash-matching approval record. Manual status edits do not satisfy the gate. |
| Workspace gate | `spec-framework work` is blocked for `existing-feature` until the current Feature Brief is approved. |
| Escalation | Cross-domain direction, business-model change, broad security policy, or uncertain product/audience scope routes back to the full Foundation package. |
| Other starting points | Their own starting-point contracts apply. This decision does not impose full Foundation on modes specialized by later FDRs. |

## Consequences

| Type | Consequence |
| --- | --- |
| Positive | A bounded feature retains traceability with one proportional approval contract. |
| Positive | The bootstrap no longer asks adopters to create five near-empty Foundation files. |
| Negative | Promotion from feature-scoped to product-scoped work requires explicit Foundation creation rather than silent reuse. |
| Follow-up | Add a guided promotion command if repeated feature briefs reveal product-wide decisions. |

## FRAMEWORK.md Amendments

| Section | Amendment |
| --- | --- |
| `3. Conceptual Model` | Define Feature Brief as the proportional existing-feature entry contract. |
| `15. How To Use With Codex` | Route existing-feature bootstrap through Feature Brief approval before workspace creation. |

## Related Artifacts

| Artifact | Link |
| --- | --- |
| Starting points | [FDR-015](FDR-015-starting-points-and-source-import.md) |
| Foundation registry | [FDR-030](FDR-030-foundation-approval-registry.md) |
| Feature Brief template | [../template/feature-brief-template.md](../template/feature-brief-template.md) |
| Canonical method | [../../FRAMEWORK.md](../../FRAMEWORK.md) |

## Supersedes

- N/A; this decision specializes FDR-015 and FDR-030 for `existing-feature`.
