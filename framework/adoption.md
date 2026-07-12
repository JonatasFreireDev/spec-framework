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

The wizard also asks for the starting point. All skills remain installed; the choice changes the generated bootstrap. For existing documents, use `--starting-point existing-documents` with `--source-dir` or `--sources`. The command inventories sources under `product/knowledge/imports/` but does not create Domains, User Goals, or Features without explicit approval.

After the Artifact Importer fills `mapping.json`, review the selected mappings and materialize them explicitly:

```bash
spec-framework import materialize --run IMPORT-001 --approved-by "Product Owner" --yes
```

The command rejects missing evidence, paths outside `product/`, duplicate targets, non-draft content, and existing destination files.

Use `spec-framework work --feature <path-or-id>` to create an independent workspace, then `status` and `next` to see blockers and the canonical next skill. Use `approve` for human-reviewed status grants, `gates` before Code Runner, and `graph ready/claim/release/complete` to coordinate task ownership.

See [delivery-closure.md](delivery-closure.md) for the complete operational flow and command examples.

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
spec-framework dashboard --work WORK-001
spec-framework decisions migrate
```

Use `decisions migrate` as a preview first. Existing repositories should use `--interactive` to review ambiguous inferred types and scopes before applying the metadata migration.

Adoption is backed by the validator, package smoke tests, and the `.spec-framework/` / `product/` boundary.
