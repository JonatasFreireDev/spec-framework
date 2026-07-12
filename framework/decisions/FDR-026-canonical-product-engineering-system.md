# FDR-026: Canonical Product Engineering System

## Snapshot

| Field | Value |
| --- | --- |
| ID | `FDR-026` |
| Status | `approved` |
| Origin EV | `Governance baseline` |
| Date | `2026-07-12` |
| Owner | `Framework Maintainer` |

## Context

The framework gives product Design a shared, versioned system and an independent review gate, while stable engineering knowledge is only suggested by starter README files. Delivery-specific Technical Discovery, architecture decisions, Implementation Plans, Code Review, and QA exist, but they do not define a single contract for how a product is built and operated. Agents therefore rediscover architecture from code, mix current-state discovery with proposed changes, and can reach implementation planning without an independent review of the intended technical solution.

## Decision

| Boundary | Contract |
| --- | --- |
| Shared system | Products may declare a versioned Engineering System under `engineering/`. `engineering-system.md` is the human contract and `engineering-system.yaml` is the mechanical catalog. |
| Origin and maturity | The system declares `origin_mode: generate | evolve | adopt`, semantic version, and maturity by area using `baseline | mapped | governed | verified | operated`. Maturity describes available evidence; it is not approval. |
| Stable knowledge | Architecture, standards, quality attributes, fitness functions, runbooks, and evidence belong under `engineering/`. Product decisions remain `DEC-*` records under `knowledge/decisions/`. |
| Delivery proposal | A delivery-specific `engineering-proposal.md` describes the intended technical change after Technical Discovery and a resolved Architecture Gate. It links rather than duplicates stable Engineering System knowledge. |
| Independent review | `engineering-review.md` is owned by a read-only Engineering Review specialist. It evaluates architecture, data ownership, dependencies, operations, tests, and decision coverage before Implementation Planning. |
| Applicability | Tier L requires Engineering Proposal and Engineering Review. Tier S and M require them when `context.md` declares a supported structured `engineering_trigger`; automation must not infer triggers from prose. |
| Approval | Engineering artifacts do not approve architecture. Structural, data, security, privacy, or other governed choices require an applicable approved `DEC-*` with a current approval record. |
| Compatibility | Existing products remain valid until a use case is migrated or governed by the Tier L rule. Upgrade adds no product approval records and never overwrites adopter-owned engineering content. |
| Runtime | Skills and validators ship through the versioned external runtime selected by `product/.product/framework.json`; only product-owned engineering artifacts live in the adopter repository. |

## Consequences

| Type | Consequence |
| --- | --- |
| Positive | Stable engineering knowledge becomes reusable and traceable across deliveries. |
| Positive | Discovery, proposed solution, implementation sequencing, and code review have distinct owners and gates. |
| Positive | Tier L receives an independent technical design review before executable work. |
| Negative | Tier L gains two planning artifacts and an additional approval checkpoint. |
| Negative | Early Engineering Systems may be incomplete until products map their real code and operations. |
| Positive | Structured technical triggers extend the review gate to applicable Tier S and M work without prose inference. |
| Follow-up | Add CLI inspection and product-configured fitness-function execution without hardcoding a technology stack. |

## FRAMEWORK.md Amendments

| Section | Amendment |
| --- | --- |
| `3. Conceptual Model` | Define Engineering System, Engineering Proposal, and Engineering Review. |
| `4. Folder Structure` | Expand the shared `engineering/` product area and delivery artifacts. |
| `6. Specification Driven Development` | Insert proposal and independent review after the Architecture Gate. |
| `7. Implementation Plan` | Require an applicable approved Engineering Review before planning. |
| `9. Skills` | Register Engineering System and Engineering Review specialists. |
| `15. How To Use With Codex` | Expose Engineering System inspection and the canonical trigger list through the versioned CLI runtime. |

## Related Artifacts

| Artifact | Link |
| --- | --- |
| Canonical method | [../../FRAMEWORK.md](../../FRAMEWORK.md) |
| Technical Discovery | [../skills/technical-discovery/SKILL.md](../skills/technical-discovery/SKILL.md) |
| Engineering System | [../skills/engineering-system/SKILL.md](../skills/engineering-system/SKILL.md) |
| Engineering Review | [../skills/engineering-review/SKILL.md](../skills/engineering-review/SKILL.md) |
| External runtime | [FDR-025](FDR-025-external-runtime-and-manifest-only-activation.md) |
| Prior operational flow | [FDR-018](FDR-018-feature-to-task-operational-closure.md) |

## Supersedes

- N/A
