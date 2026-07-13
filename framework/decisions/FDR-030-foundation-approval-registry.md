# FDR-030: Foundation Approval Registry

## Snapshot

| Field | Value |
| --- | --- |
| ID | `FDR-030` |
| Status | `proposed` |
| Origin EV | `Governance baseline` |
| Date | `2026-07-12` |
| Owner | `Framework Maintainer` |

## Context

The starter artifact registry contained domain templates and omitted the L0 Foundation ladder. `spec-framework validate` could accept structurally valid Problem, Vision, and Strategy files, while `spec-framework approve` rejected those same paths as unregistered. Manual status edits therefore appeared to advance Foundation without producing current approval evidence in `.product/history/`.

## Decision

| Boundary | Contract |
| --- | --- |
| Initial registry | New products register Problem, Vision, Product Principles, North Star, and Strategy during initialization. |
| Regeneration | Canonical Foundation files declare ID, type, status, and parent IDs so `validate --write-registry` preserves the same entries and relationships. |
| Ordering | Vision requires approved Problem; Principles and North Star require approved Vision; Strategy requires approved Vision, Principles, and North Star. |
| Proportional scope | Registration does not force a product-wide definition. An `existing-feature` adoption may approve feature-scoped Foundation contracts that justify and bound only that delivery while preserving the same parent and evidence gates. |
| Status synchronization | Approving Problem, Vision, or Strategy synchronizes the canonical file and its local `context.md` status before updating the registry and history. |
| Evidence | Foundation statuses at `approved` or later are invalid without a current hash-matching record in `.product/history/`. Manual Markdown edits never constitute approval. |
| Compatibility | Existing products first add Foundation ID/type/parent metadata, then run `validate --write-registry`; no approval records are generated automatically. Human approval is still required. |

## Consequences

| Type | Consequence |
| --- | --- |
| Positive | The official CLI can produce every required L0 approval with traceable evidence. |
| Positive | Registry regeneration no longer drops Foundation artifacts or their parent gates. |
| Positive | A feature-only adoption retains traceability without fabricating a complete product strategy. |
| Negative | Existing products with cosmetic approved statuses become correctly blocked until a human re-approves them. |
| Follow-up | Add a guided metadata migration for adopters whose Foundation files predate the registry contract. |

## FRAMEWORK.md Amendments

| Section | Amendment |
| --- | --- |
| `11. Approval Gates` | Require registered Foundation artifacts, parent ordering, atomic status synchronization, and history evidence. |

## Related Artifacts

| Artifact | Link |
| --- | --- |
| Canonical method | [../../FRAMEWORK.md](../../FRAMEWORK.md) |
| Starter registry | [../../starter/product/.product/artifacts.json](../../starter/product/.product/artifacts.json) |
| Starting points | [FDR-015](FDR-015-starting-points-and-source-import.md) |
| Approval workflow | [../../internal/workflow/workflow.go](../../internal/workflow/workflow.go) |
| Registry generator | [../../internal/validator/output.go](../../internal/validator/output.go) |

## Supersedes

- N/A
