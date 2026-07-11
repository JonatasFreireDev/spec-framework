# Adoption Guide

## Goal

Create a new Specification Driven Development product repository without mixing product-owned documents with framework development history.

## Recommended Path

Automated bootstrap:

```bash
spec-framework init --target ../my-product
```

CLI-style bootstrap from the framework repository:

```bash
spec-framework init --target ../my-product --agents codex,cursor,claude --yes
```

Install a versioned release binary as described in [install.md](install.md). Go and Node.js are not runtime requirements for adopters.

Manual bootstrap:

1. Create an empty product repository.
2. Copy the contents of `starter/` into the product repository root.
3. Install framework assets into `.spec-framework/`:
   - `.spec-framework/FRAMEWORK.md`
   - `.spec-framework/AGENTS.framework.md`
   - `.spec-framework/decisions/FDR-*`
   - `.spec-framework/skills/`
   - `.spec-framework/templates/`
4. Generate one or more agent skill trees: `.agents/skills/`, `.cursor/skills/`, and `.claude/skills/`.
5. Replace `product/` starter placeholders with product-specific content.
6. Run the validation wrapper against the product root.

```bash
spec-framework validate
```

Direct validator form when debugging:

```bash
spec-framework validate --product-root product --framework-root .spec-framework --write-registry --write-report
```

Upgrade an initialized product from the framework repository:

```bash
spec-framework upgrade --target ../my-product --agents codex --yes
```

## What Belongs To The Product

| Product-Owned Area | Purpose |
| --- | --- |
| `product/.product/` | Product state, registry, derivations, approval records, and adopted framework metadata. |
| `product/foundation/` | Problem, vision, and strategy for the product. |
| `product/domains/` | Product domains, goals, features, use cases, specifications, and tasks. |
| `product/knowledge/decisions/` | Product decisions only. |
| `product/knowledge/business-rules/` | Product business rules. |
| `product/audits/` | Product audits, readiness reports, QA evidence references, and threat register. |
| `product/releases/` | Product release notes and release readiness. |
| `product/design/` | Product design artifacts and mockups. |

## What Belongs To The Framework

| Framework-Owned Area | Purpose |
| --- | --- |
| `.spec-framework/FRAMEWORK.md` | Method contract. |
| `.spec-framework/AGENTS.framework.md` | Agent instructions for resolving framework and product roots. |
| `.spec-framework/decisions/FDR-*` | Framework method decisions. |
| `.spec-framework/skills/` | Operational skills. |
| `.spec-framework/templates/` | Reusable artifact templates. |
| Installed `spec-framework` binary | Validation, bootstrap, upgrade, and migration tooling. |

## Non-Goals For Starter Repositories

- Do not copy framework FDRs into `product/knowledge/decisions/`.
- Do not inherit example domains as real product scope.
- Do not inherit retroactive approval records from the framework lab.
- Do not edit `.spec-framework/` internals to encode product scope.

## Upgrade Direction

Stable commands:

```bash
spec-framework init --target ../my-product
spec-framework validate
spec-framework upgrade --target ../my-product
```

Adoption is backed by the validator, package smoke tests, and the `.spec-framework/` / `product/` boundary.
