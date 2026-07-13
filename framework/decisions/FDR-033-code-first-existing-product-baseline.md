# FDR-033: Code-First Existing Product Baseline

## Snapshot

| Field | Value |
| --- | --- |
| ID | `FDR-033` |
| Status | `proposed` |
| Origin EV | `Adopter starting-point review` |
| Date | `2026-07-12` |
| Owner | `Product Orchestrator` |

## Context

An existing product may have active code, tests, deployments, users, and operational history while having little canonical product documentation. Requiring separate Problem, Vision, Product Principles, and North Star reconstruction encourages repetitive documents and may present inferred intent as historical fact. Treating the repository as an `existing-implementation` also loses the important fact that a product is already operating and has observable users, value, constraints, and outcomes.

## Decision

| Boundary | Contract |
| --- | --- |
| Applicability | `existing-product` assumes an operating product whose repository and runtime evidence may be more complete than its documentation. |
| Baseline | One canonical `foundation/product-baseline.md` consolidates current audience, evidenced needs, delivered value, capabilities, operating model, signals, constraints, decision rules, risks, and unknowns. |
| Evidence | Code, tests, configuration, telemetry, support history, and releases may support baseline claims. Observed behavior does not prove original intent. |
| Strategy | `foundation/strategy/strategy.md` remains a separate required artifact and has Product Baseline as its parent. It owns future bets, trade-offs, priorities, metrics, and change direction. |
| Registry | Problem, Vision, Product Principles, and North Star are excluded from the active registry for this starting point. Their skeletons remain available if uncertainty requires promotion to the full Foundation path. |
| Approval | Product Baseline and Strategy each require individual current hash-matching approval before workspace creation. |
| Escalation | Unknown audience, unclear value, major product repositioning, or conflicting evidence promotes the repository to the full Problem, Vision, Principles, North Star, and Strategy path. |

## Consequences

| Type | Consequence |
| --- | --- |
| Positive | A code-first existing product gets a concise current-state contract without pretending documentation already existed. |
| Positive | Current reality remains separate from future Strategy. |
| Negative | Some concepts normally separated across Foundation artifacts are consolidated in Product Baseline. |
| Follow-up | Revisit promotion tooling after multiple real adopter migrations. |

## FRAMEWORK.md Amendments

| Section | Amendment |
| --- | --- |
| `3. Conceptual Model` | Define Product Baseline for code-first existing products. |
| `15. How To Use With Codex` | Route `existing-product` through Product Baseline and Strategy approvals. |

## Related Artifacts

| Artifact | Link |
| --- | --- |
| Starting points | [FDR-015](FDR-015-starting-points-and-source-import.md) |
| Existing implementation | [FDR-032](FDR-032-existing-implementation-assessment.md) |
| Product Baseline template | [../template/product-baseline-template.md](../template/product-baseline-template.md) |
| Canonical method | [../../FRAMEWORK.md](../../FRAMEWORK.md) |

## Supersedes

- N/A; this decision specializes FDR-015 for `existing-product`.
