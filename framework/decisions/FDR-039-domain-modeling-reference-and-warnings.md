# FDR-039: Domain Modeling Reference And Warnings

## Snapshot

| Field | Value |
| --- | --- |
| ID | `FDR-039` |
| Status | `proposed` |
| Origin EV | `Adopter domain-modeling feedback` |
| Date | `2026-07-13` |
| Owner | `Domain Architect AI` |

## Context

Starting-point-specific bootstrap guidance can reach its first domain change without directing an agent to the worked domain model. The distributed starter domain scaffold has fewer boundary prompts than the canonical framework template. As a result, a product can validate structurally while using its product name as one catch-all domain, placing identity inside an unrelated business domain, or stopping before a goal, feature, and use case exist.

## Decision

| Boundary | Contract |
| --- | --- |
| Modeling reference | The versioned runtime distributes `examples/events/`. Its domain tree is the canonical reference for initial domain modeling, while remaining outside adopter-owned product scope. |
| Domain contract | Domains model coherent business areas, declare `Owns`, `Does Not Own`, and cross-domain dependencies, and are not named for the product or a UI section. |
| Walking skeleton | The first modeled domain continues through one User Goal, Feature, and Use Case before workspace creation. |
| Bootstrap | Every starting point that creates or revises domains requires the modeling reference before its first domain change. `audit-only` uses it to assess existing boundaries without mutation. |
| Validator | The validator emits non-blocking warnings for a domain matching the product name, missing explicit non-ownership, an approved domain with no goals, and apparent authentication ownership inside a non-identity domain. These are explainable heuristics, not semantic approval gates. |
| Identity exception | A product may intentionally own authentication outside a `users`/identity domain when the domain document states the boundary and rationale; the validator still warns so a human can review it. |

## Consequences

| Type | Consequence |
| --- | --- |
| Positive | Agents receive an installed, concrete reference before any starting point creates or revises a domain. |
| Positive | The starter and framework template ask for the same boundary information. |
| Positive | Structural validity no longer silently implies that obvious domain-modeling anti-patterns were considered. |
| Negative | Warnings can require a documented intentional exception for compact products. |
| Follow-up | Reassess heuristic precision with adopter feedback; do not promote these warnings to errors without evidence of stable, low-false-positive rules. |

## FRAMEWORK.md Amendments

| Section | Amendment |
| --- | --- |
| `3. Conceptual Model` | Clarify domain boundaries and the required walking skeleton. |
| `15. How To Use With Codex` | Require the installed Events reference before initial domain modeling. |

## Related Artifacts

| Artifact | Link |
| --- | --- |
| Feature-scoped bootstrap | [FDR-031](FDR-031-feature-scoped-bootstrap.md) |
| Domain Architect skill | [../skills/domain-architect/SKILL.md](../skills/domain-architect/SKILL.md) |
| Domain template | [../template/domain-template.md](../template/domain-template.md) |
| Events domain example | [../../examples/events/domains/events/domain.md](../../examples/events/domains/events/domain.md) |
| Canonical method | [../../FRAMEWORK.md](../../FRAMEWORK.md) |

## Supersedes

- N/A
