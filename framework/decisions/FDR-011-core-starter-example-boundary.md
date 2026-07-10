# FDR-011: Framework core, product starter, and examples boundary

## Snapshot

| Field | Value |
| --- | --- |
| ID | `FDR-011` |
| Status | `approved` |
| Origin EV | `EV-015` |
| Date | `2026-07-10` |
| Owner | `Documentation Orchestrator` |

## Context

This repository started as both the framework laboratory and a product-shaped documentation tree. That was useful for bootstrapping, but it makes adoption confusing: framework method artifacts, product starter files, and worked examples can look like one product.

New product repositories need a clean SDD skeleton without inheriting framework FDR history, lab approval records, or example product scope.

## Decision

Separate repository intent into three areas:

| Area | Responsibility |
| --- | --- |
| `framework/` | Reusable framework core documentation, adoption model, and future packaging boundary. |
| `starter/` | Clean new-repository skeleton containing `.spec-framework/` for method assets and `product/` for product artifacts. |
| `examples/` | Worked examples that teach the framework but are not copied into products by default. |

During the transition, existing root-level operational paths remain in place so validators, CI, and Codex skills continue to work. Future work should add a bootstrap CLI that copies `starter/`, installs framework assets into `.spec-framework/`, optionally mirrors skills into `.codex/skills/`, and avoids copying lab history.

## Consequences

| Type | Consequence |
| --- | --- |
| Positive | New product repos get a cleaner starting point with a visible boundary between how the framework works and what the product is. |
| Positive | Framework evolution history no longer needs to masquerade as product history. |
| Positive | Examples can become optional learning material. |
| Negative | The repo temporarily has both old root-level lab paths and new adoption folders. |
| Follow-up | Build `scripts/init-product.mjs`, add validator support for `--product-root`, and migrate example domains into `examples/` once validators support example scope. |

## FRAMEWORK.md Amendments

| Section | Amendment |
| --- | --- |
| `4. Estrutura De Pastas` | Clarify that new repos use `.spec-framework/` for method assets and `product/` for product artifacts. |
| `15. Como Usar Com Codex` | Clarify that `spec-framework` is a lab/framework repo, while product repos start from `starter/`. |

## Related Artifacts

| Artifact | Link |
| --- | --- |
| Framework core README | [../../framework/README.md](../../framework/README.md) |
| Adoption guide | [../../framework/adoption.md](../../framework/adoption.md) |
| Product starter | [../../starter/README.md](../../starter/README.md) |
| Examples | [../../examples/README.md](../../examples/README.md) |

## Supersedes

- `N/A`
