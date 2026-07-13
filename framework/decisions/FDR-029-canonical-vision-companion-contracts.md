# FDR-029: Canonical Vision Companion Contracts

## Snapshot

| Field | Value |
| --- | --- |
| ID | `FDR-029` |
| Status | `proposed` |
| Origin EV | `Governance baseline` |
| Date | `2026-07-12` |
| Owner | `Vision AI` |

## Context

The starter and Vision template placed product principles and north-star content inside `vision.md` while also shipping dedicated `principles.md` and `north-star.md` artifacts. The Vision context already identifies those files as distinct children and documents, but the duplicated sections leave agents and adopters without a reliable source of truth when the copies diverge.

## Decision

| Boundary | Canonical Contract |
| --- | --- |
| Vision | `vision.md` owns product direction, target users, non-goals, and decision boundaries. |
| Product Principles | `principles.md` exclusively owns principles, meanings, trade-offs, examples, and anti-principles. |
| North Star | `north-star.md` exclusively owns the durable value outcome, candidate metric, measurement notes, and guardrails. |
| Linking | `vision.md` links to both companion contracts and does not repeat their canonical content. |
| Existing products | Existing duplicated content is preserved until intentionally edited; the next Vision revision consolidates it into the owning companion artifact without changing approval history automatically. |

## Consequences

| Type | Consequence |
| --- | --- |
| Positive | Each Vision surface has one source of truth and a clear owner. |
| Positive | Principles and metrics can evolve without creating conflicting copies in `vision.md`. |
| Negative | Reading the complete Vision package requires following two explicit companion links. |
| Follow-up | Validators may later diagnose duplicated canonical sections after a compatibility period. |

## FRAMEWORK.md Amendments

| Section | Amendment |
| --- | --- |
| `3. Conceptual Model` | Define exclusive ownership for Vision, Product Principles, and North Star. |

## Related Artifacts

| Artifact | Link |
| --- | --- |
| Canonical method | [../../FRAMEWORK.md](../../FRAMEWORK.md) |
| Vision skill | [../skills/vision/SKILL.md](../skills/vision/SKILL.md) |
| Vision template | [../template/vision-template.md](../template/vision-template.md) |
| Starter Vision | [../../starter/product/foundation/vision/vision.md](../../starter/product/foundation/vision/vision.md) |
| Starter principles | [../../starter/product/foundation/vision/principles.md](../../starter/product/foundation/vision/principles.md) |
| Starter north star | [../../starter/product/foundation/vision/north-star.md](../../starter/product/foundation/vision/north-star.md) |

## Supersedes

- N/A
