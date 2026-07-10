# Adoption Guide

## Goal

Create a new Specification Driven Development product repository without mixing product-owned documents with framework development history.

## Recommended Path

Automated bootstrap:

```bash
node scripts/init-product.mjs --target ../my-product
```

CLI-style bootstrap from the framework repository:

```bash
node scripts/spec-framework.mjs init --target ../my-product
```

Local linked CLI:

```bash
npm link
spec-framework init --target ../my-product
```

Packaged CLI from a local tarball:

```bash
npm pack
mkdir ../my-consumer
cd ../my-consumer
npm install ../spec-framework/spec-framework-0.1.0.tgz --no-save
npx spec-framework init --target ../my-product
cd ../my-product
npx spec-framework validate
```

The package path is currently a controlled local/Git-based adoption path, not a public npm release contract. Keep using `starter/`, `.spec-framework/`, and `product/` as the canonical boundary inside generated repositories.

Manual bootstrap:

1. Create an empty product repository.
2. Copy the contents of `starter/` into the product repository root.
3. Install framework assets into `.spec-framework/`:
   - `.spec-framework/FRAMEWORK.md`
   - `.spec-framework/decisions/FDR-*`
   - `.spec-framework/skills/`
   - `.spec-framework/templates/`
   - `.spec-framework/validators/`
   - `.spec-framework/tools/`
4. Optionally copy skills into `.codex/skills/` for Codex auto-discovery.
5. Replace `product/` starter placeholders with product-specific content.
6. Run the validator against the product root.

```bash
node .spec-framework/validators/framework-validator.mjs --product-root product --write-registry --write-report
```

Preferred validation wrapper after bootstrap:

```bash
node .spec-framework/tools/validate-product.mjs
```

Upgrade an initialized product from the framework repository:

```bash
node scripts/upgrade-product.mjs --target ../my-product
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
| `.spec-framework/decisions/FDR-*` | Framework method decisions. |
| `.spec-framework/validators/` | Mechanical validation gates. |
| `.spec-framework/skills/` | Operational skills. |
| `.spec-framework/templates/` | Reusable artifact templates. |
| `.spec-framework/tools/` | Bootstrap, upgrade, and migration tooling. |

## Non-Goals For Starter Repositories

- Do not copy framework FDRs into `product/knowledge/decisions/`.
- Do not inherit example domains as real product scope.
- Do not inherit retroactive approval records from the framework lab.
- Do not edit `.spec-framework/` internals to encode product scope.

## Upgrade Direction

Future framework versions should support:

```bash
spec-framework init --target ../my-product
spec-framework validate
spec-framework upgrade --target ../my-product
```

Adoption is backed by the validator, package smoke tests, and the `.spec-framework/` / `product/` boundary.
