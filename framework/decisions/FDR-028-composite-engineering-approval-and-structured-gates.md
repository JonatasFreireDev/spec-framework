# FDR-028: Composite Engineering Approval And Structured Gates

## Snapshot

| Field | Value |
| --- | --- |
| ID | `FDR-028` |
| Status | `approved` |
| Origin EV | `Governance baseline` |
| Date | `2026-07-12` |
| Owner | `Framework Maintainer` |

## Context

After FDR-027, the Engineering System still had three lifecycle surfaces with no atomic approval path, catalog changes did not invalidate the approval of `engineering-system.md`, operational Architecture Gate resolution accepted any DEC token, and `Not applicable` was detected through free-text substrings. Some distributed guides also retained the pre-Engineering-Proposal flow.

## Decision

| Boundary | Contract |
| --- | --- |
| Composite approval | Engineering System approval hashes every product-owned contract file under `engineering/` in deterministic path order. Any component change invalidates the approval. |
| Atomic status | Approving the Engineering System updates status in `context.md`, `engineering-system.md`, and `engineering-system.yaml` as one rollback-capable operation before recording the composite hash. |
| Architecture Gate | Operational resolution accepts exact structured `Not required` with non-placeholder rationale, or a DEC that is indexed, approved, hash-current, and scope-compatible. |
| Not applicable | Delivery artifacts use structured status `not_applicable` and a non-placeholder `Rationale` field. Free-text mentions never satisfy the gate. |
| Distribution | Framework agent instructions, generated guides, templates, skills, and canonical flow remain synchronized. |
| Migration | Existing draft prose is unchanged. Artifacts intentionally declared N/A must adopt structured status and rationale before advancing; no approval record is generated. |

## Consequences

| Type | Consequence |
| --- | --- |
| Positive | Engineering System approval covers the actual shared technical contract rather than one Markdown file. |
| Positive | CLI navigation and validator agree on Architecture Gate resolution. |
| Positive | Incidental prose cannot bypass Design or specification-contract gates. |
| Negative | Any shared engineering contract change requires human re-approval. |
| Negative | Older N/A artifacts need a small metadata migration before they can advance. |

## FRAMEWORK.md Amendments

| Section | Amendment |
| --- | --- |
| `3. Engineering System` | Define composite approval boundary. |
| `6. Specification Driven Development` | Define structured N/A contract. |
| `7. Implementation Plan` | Align Architecture Gate operational validity. |
| `11. Approval Gates` | Require atomic composite Engineering System approval. |

## Related Artifacts

| Artifact | Link |
| --- | --- |
| Canonical method | [../../FRAMEWORK.md](../../FRAMEWORK.md) |
| Engineering integrity | [FDR-027](FDR-027-engineering-gate-integrity.md) |
| Decision effects | [FDR-019](FDR-019-executable-product-decision-effects.md) |

## Supersedes

- N/A; this decision amends FDR-019, FDR-026, and FDR-027.
