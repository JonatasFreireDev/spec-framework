# Product Starter

## Purpose

This folder is the clean skeleton for a new product repository using Specification Driven Development.

The CLI composes product-owned assets from `starter/product/` according to the selected versioned contract under `framework/init/contracts/`. The adopter repository receives:

```text
product/
  # product-owned SDD tree
```

Framework assets live in the versioned user cache and are resolved from `product/.product/framework.json`. Do not copy `.spec-framework/` or agent skill trees into the adopter repository.

## Included Product-Owned Areas

| Area | Purpose |
| --- | --- |
| `product/.product/` | Product state, framework adoption metadata, decisions index, roadmap, artifact registry, derivations, and approval history. |
| `product/foundation/` | Full Foundation plus proportional starting-point contracts when selected. |
| `product/knowledge/` | Product-specific knowledge, decisions, business rules, conventions, and examples. |
| `product/domains/` | Product domain tree. |
| `product/audits/` | Product validation reports, readiness checks, and security/threat tracking. |
| `product/releases/` | Release planning and release evidence. |
| `product/design/` | Product design references. |
| `product/engineering/` | Product engineering notes, architecture, runbooks, and product-specific gates. |

## First Product Steps

These are the default `new-product` steps. For every initialized repository, read `product/BOOTSTRAP.md` first; it names the active starting-point contract and may replace or prepend this sequence.

1. Fill `product/context.md` with product identity.
2. Fill `product/foundation/problem/context.md`, `product/foundation/problem/problem.md`, and `product/foundation/problem/opportunities.md`.
3. Fill `product/foundation/vision/context.md`, `product/foundation/vision/vision.md`, `product/foundation/vision/principles.md`, and `product/foundation/vision/north-star.md`.
4. Fill `product/foundation/strategy/context.md`, `product/foundation/strategy/strategy.md`, `personas.md`, `metrics.md`, and `roadmap.md`.
5. Copy `product/domains/_template-domain/` to the first real domain slug and update every `context.md`.
6. Continue through Goal -> Feature -> Use Case -> Specification -> Design -> Technical Discovery -> applicable Engineering Proposal and Engineering Review -> Implementation Plan -> Execution Graph -> Tasks.
7. Replace product gate placeholders in `product/knowledge/conventions/gates.md` before marking executable work as implemented or validated.
8. Use `spec-framework work`, then `resume`, leases, checkpoints, and handoffs so execution can be continued without conversation history.
9. Use `spec-framework engineering-system inspect` for the shared technical baseline, `spec-framework dashboard --work WORK-001` for the consolidated flow, and `spec-framework decisions migrate` to preview legacy decision metadata upgrades.
10. For `existing-documents`, ask the Artifact Importer agent to read each source and fill `traceability.json` plus proposed `mapping.json` entries before materialization; review unmapped gaps explicitly.

## Boundary Rule

Do not put framework maintenance history in any product decision domain. Product decisions are indexed in `product/.product/decisions.json` and stored under `product/knowledge/decisions/`, `product/design/decisions/`, or `product/engineering/decisions/` according to domain; the framework method lives in the pinned external runtime.

## Starter Templates

| Template | Use |
| --- | --- |
| `product/domains/_template-domain/` | Copy for each real domain. |
| `product/domains/_template-domain/goals/_template-goal/` | Copy for each user goal. |
| `product/domains/_template-domain/goals/_template-goal/features/_template-feature/` | Copy for each feature. |
| `product/domains/_template-domain/goals/_template-goal/features/_template-feature/use-cases/_template-use-case/` | Copy for each use case. |
