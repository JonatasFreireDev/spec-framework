# Product Starter

## Purpose

This folder is the clean skeleton for a new product repository using Specification Driven Development.

Copy the contents of this folder into a new product repository. The copied repository will contain:

```text
.spec-framework/
  # how to operate the framework

product/
  # product-owned SDD tree
```

Then install or copy the framework assets described in [../framework/adoption.md](../framework/adoption.md).

## Included Product-Owned Areas

| Area | Purpose |
| --- | --- |
| `product/.product/` | Product state, framework adoption metadata, decisions index, roadmap, artifact registry, derivations, and approval history. |
| `product/foundation/` | Problem, vision, and strategy documents. |
| `product/knowledge/` | Product-specific knowledge, decisions, business rules, conventions, and examples. |
| `product/domains/` | Product domain tree. |
| `product/audits/` | Product validation reports, readiness checks, and security/threat tracking. |
| `product/releases/` | Release planning and release evidence. |
| `product/design/` | Product design references. |
| `product/engineering/` | Product engineering notes, architecture, runbooks, and product-specific gates. |

## First Product Steps

1. Fill `product/context.md` with product identity.
2. Fill `product/foundation/problem/context.md`, `product/foundation/problem/problem.md`, and `product/foundation/problem/opportunities.md`.
3. Fill `product/foundation/vision/context.md`, `product/foundation/vision/vision.md`, `product/foundation/vision/principles.md`, and `product/foundation/vision/north-star.md`.
4. Fill `product/foundation/strategy/context.md`, `product/foundation/strategy/strategy.md`, `personas.md`, `metrics.md`, and `roadmap.md`.
5. Copy `product/domains/_template-domain/` to the first real domain slug and update every `context.md`.
6. Continue through Goal -> Feature -> Use Case -> Specification -> Design -> Implementation Plan -> Execution Graph -> Tasks.
7. Replace product gate placeholders in `product/knowledge/conventions/gates.md` before marking executable work as implemented or validated.
8. Use `spec-framework work`, then `resume`, leases, checkpoints, and handoffs so execution can be continued without conversation history.
9. Use `spec-framework dashboard --work WORK-001` for the consolidated flow and `spec-framework decisions migrate` to preview legacy decision metadata upgrades.

## Boundary Rule

Do not put framework method decisions in `product/knowledge/decisions/`. Product decisions go there; framework decisions live under `.spec-framework/decisions/`.

## Starter Templates

| Template | Use |
| --- | --- |
| `product/domains/_template-domain/` | Copy for each real domain. |
| `product/domains/_template-domain/goals/_template-goal/` | Copy for each user goal. |
| `product/domains/_template-domain/goals/_template-goal/features/_template-feature/` | Copy for each feature. |
| `product/domains/_template-domain/goals/_template-goal/features/_template-feature/use-cases/_template-use-case/` | Copy for each use case. |
