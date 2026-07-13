# FDR-032: Existing Implementation Assessment

## Snapshot

| Field | Value |
| --- | --- |
| ID | `FDR-032` |
| Status | `proposed` |
| Origin EV | `Adopter first-run review` |
| Date | `2026-07-12` |
| Owner | `Product Orchestrator` |

## Context

The `existing-implementation` starting point currently changes bootstrap wording but enters the same Foundation path without first recording what the code, tests, runtime configuration, and repository history actually demonstrate. Deriving product intent directly from implementation risks converting accidents, legacy behavior, and unsupported assumptions into approved product truth.

## Decision

| Boundary | Contract |
| --- | --- |
| Applicability | `existing-implementation` begins with one canonical `knowledge/assessments/implementation-assessment.md`. |
| Required content | The assessment inventories observed behavior, architecture, data and integrations, test evidence, operational constraints, security and privacy signals, documentary gaps, conflicts, risks, and candidate product claims. |
| Evidence boundary | Observed implementation is evidence. Candidate intent remains an assumption until it is written and approved in the owning Foundation artifact. |
| Registry | The assessment is registered as `implementation-assessment` and becomes a parent of Problem for this starting point. |
| Approval | The assessment requires an individual hash-matching approval before Problem approval or workspace creation. |
| Foundation | Problem, Vision, Product Principles, North Star, and Strategy remain required after the assessment and before workspace creation. |
| Mutation | Assessment work must not modify application code or fabricate product approvals. |
| Escalation | Material conflicts, unknown data handling, or unsafe runtime behavior become explicit blockers or decisions before downstream delivery. |

## Consequences

| Type | Consequence |
| --- | --- |
| Positive | Existing code becomes traceable evidence instead of implicit product authority. |
| Positive | The first session has a concrete artifact and gate before Foundation drafting. |
| Negative | Existing implementations require an additional approval before the full Foundation sequence. |
| Follow-up | Consider a dedicated CLI scanner only after the manual assessment contract is stable. |

## FRAMEWORK.md Amendments

| Section | Amendment |
| --- | --- |
| `3. Conceptual Model` | Define Implementation Assessment as the evidence boundary for existing implementations. |
| `15. How To Use With Codex` | Route `existing-implementation` through assessment approval before Foundation and workspace creation. |

## Related Artifacts

| Artifact | Link |
| --- | --- |
| Starting points | [FDR-015](FDR-015-starting-points-and-source-import.md) |
| Feature-scoped bootstrap | [FDR-031](FDR-031-feature-scoped-bootstrap.md) |
| Assessment template | [../template/implementation-assessment-template.md](../template/implementation-assessment-template.md) |
| Canonical method | [../../FRAMEWORK.md](../../FRAMEWORK.md) |

## Supersedes

- N/A; this decision specializes FDR-015 for `existing-implementation`.
