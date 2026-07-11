# Framework Core

## Purpose

This folder describes the reusable Spec Framework core: method, operational contracts, validators, skills, templates, and adoption rules.

The `framework/` directory hosts the executable framework core. New product repositories should not copy the whole repository root; they should start from `starter/` and consume the framework core through the documented adoption path.

## Ownership Boundary

| Area | Owner | Product Repo Copies It? | Notes |
| --- | --- | --- | --- |
| `FRAMEWORK.md` | Framework core | Yes, into `.spec-framework/FRAMEWORK.md` | Canonical method contract. |
| `framework/skills/` | Framework core | Yes, into `.spec-framework/skills/` and selected agent skill trees | Operational agent contracts. |
| `framework/template/` | Framework core | Yes, into `.spec-framework/templates/` | Reusable artifact templates. |
| Go CLI | Framework core | Installed as a release binary | Mechanical gates and migration tools. |
| `framework/tests/` | Framework core | No | Tests the framework laboratory and distribution flow. |
| `framework/decisions/FDR-*` | Framework core | Yes, into `.spec-framework/decisions/` | Framework method history, not product history. |
| `starter/` | Product starter | Yes | Clean product-owned skeleton with `.spec-framework/` and `product/`. |
| `examples/` | Examples | Optional | Learning material, not production source of truth. |

## Adoption Models

| Model | Status | Best Use |
| --- | --- | --- |
| Template copy | Supported manually through `starter/`. | First real adopters and fast experiments. |
| CLI bootstrap | Supported through versioned Go release binaries. | Repeatable product creation with versioned framework assets. |
| Submodule | Future option. | Larger teams that need strict framework versioning. |

## Current Boundary Rule

Treat the root repository as the framework laboratory and treat `starter/` as the copyable product skeleton.

Do not use `examples/events/` as the canonical starter. It contains worked product history and example artifacts. New products should contain product artifacts under `product/` and method assets under `.spec-framework/`.

## Next Step

Use `spec-framework init` to copy `starter/`, install framework assets, generate selected agent skill formats, and record the adopted version.

Use `spec-framework upgrade` to refresh `.spec-framework/` assets without touching `product/`.

For framework development, use `go run ./cmd/spec-framework`; adopters use the precompiled release binary.
