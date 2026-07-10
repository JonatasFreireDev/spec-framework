# Framework Core

## Purpose

This folder describes the reusable Spec Framework core: method, operational contracts, validators, skills, templates, and adoption rules.

The repository root still hosts the executable framework while the project transitions to a package/CLI-ready layout. New product repositories should not copy the whole repository root blindly; they should start from `starter/` and consume the framework core through the documented adoption path.

## Ownership Boundary

| Area | Owner | Product Repo Copies It? | Notes |
| --- | --- | --- | --- |
| `FRAMEWORK.md` | Framework core | Yes, into `.spec-framework/FRAMEWORK.md` | Canonical method contract. |
| `.codex/skills/` | Framework core | Yes, into `.spec-framework/skills/` and optionally `.codex/skills/` | Operational agent contracts. |
| `knowledge/templates/` | Framework core | Yes, into `.spec-framework/templates/` | Reusable artifact templates. |
| `engineering/validators/` | Framework core | Yes, into `.spec-framework/validators/` | Mechanical gates. |
| `engineering/tests/` | Framework core | Optional | Required when developing the framework itself. |
| `engineering/decisions/FDR-*` | Framework core | Yes, into `.spec-framework/decisions/` | Framework method history, not product history. |
| `starter/` | Product starter | Yes | Clean product-owned skeleton with `.spec-framework/` and `product/`. |
| `examples/` | Examples | Optional | Learning material, not production source of truth. |

## Adoption Models

| Model | Status | Best Use |
| --- | --- | --- |
| Template copy | Supported manually through `starter/`. | First real adopters and fast experiments. |
| CLI bootstrap | Recommended target. | Repeatable product creation with versioned framework assets. |
| Package/submodule | Future option. | Larger teams that need strict framework versioning. |

## Current Transition Rule

Until the CLI exists, treat the root repository as the framework laboratory and treat `starter/` as the copyable product skeleton.

Do not use root-level `domains/`, `foundation/`, `.product/history/`, or `audits/` as the canonical starter. They contain framework development history and examples. New products should contain product artifacts under `product/` and method assets under `.spec-framework/`.

## Next Step

Use `scripts/init-product.mjs` to copy `starter/`, install framework assets into `.spec-framework/`, mirror skills for Codex discovery, and record the adopted framework version in `product/.product/framework.json`.

Use `scripts/upgrade-product.mjs` to refresh `.spec-framework/` assets in an existing product without touching `product/`.
