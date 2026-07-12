# Framework Core

## Purpose

This folder describes the reusable Spec Framework core: method, operational contracts, validators, skills, templates, and adoption rules.

The `framework/` directory hosts the executable framework core. New product repositories should not copy the whole repository root; they should start from `starter/` and consume the framework core through the documented adoption path.

## Ownership Boundary

| Area | Owner | Product Repo Copies It? | Notes |
| --- | --- | --- | --- |
| `FRAMEWORK.md` | Framework core | Yes, into the versioned user cache | Canonical method contract. |
| `framework/skills/` | Framework core | Yes, into the versioned user cache | Operational agent contracts resolved by the global dispatcher. |
| `framework/template/` | Framework core | Yes, into the versioned user cache | Reusable artifact templates. |
| Go CLI | Framework core | Installed as a release binary | Mechanical gates and migration tools. |
| `framework/tests/` | Framework core | No | Tests the framework laboratory and distribution flow. |
| `framework/decisions/FDR-*` | Framework core | Yes, into the versioned user cache | Framework method history, not product history. |
| `starter/product/` | Product starter | Yes | Clean product-owned skeleton. |
| `examples/` | Examples | Optional | Learning material, not production source of truth. |

## Adoption Models

| Model | Status | Best Use |
| --- | --- | --- |
| Template copy | Supported manually through `starter/`. | First real adopters and fast experiments. |
| CLI bootstrap | Supported through versioned Go release binaries. | Repeatable product creation with versioned framework assets. |
| Submodule | Future option. | Larger teams that need strict framework versioning. |

## Current Boundary Rule

Treat the root repository as the framework laboratory and treat `starter/` as the copyable product skeleton.

Do not use `examples/events/` as the canonical starter. It contains worked product history and example artifacts. New products contain product artifacts under `product/`; method assets stay in the external user cache.

## Next Step

Use `spec-framework init` to copy `starter/product/`, cache embedded framework assets, install namespaced user dispatchers, and record the adopted version.

Use `spec-framework upgrade` to refresh the external runtime and manifest without overwriting adopter-owned product content.

For framework development, use `go run ./cmd/spec-framework`; adopters use the precompiled release binary.
