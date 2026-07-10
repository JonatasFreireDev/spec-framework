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

1. Fill `product/foundation/problem/problem.md`.
2. Fill `product/foundation/vision/vision.md`.
3. Fill `product/foundation/strategy/strategy.md`.
4. Create the first domain under `product/domains/<domain>/`.
5. Create the first user goal, feature, and use case.
6. Generate Specification, Design, Implementation Plan, Execution Graph, and Tasks through the framework gates.

## Boundary Rule

Do not put framework method decisions in `product/knowledge/decisions/`. Product decisions go there; framework decisions live under `.spec-framework/decisions/`.
