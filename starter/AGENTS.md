# Agent Instructions

## Purpose

This repository is a product repository adopting the Spec Framework. Product documentation is execution infrastructure.

## Source Of Truth

Read `.spec-framework/FRAMEWORK.md` first when it exists. If this starter has not yet installed framework assets, follow `README.md` and install the framework before generating downstream artifacts.

## Product-Owned Areas

- `product/.product/`: product state, approval history, artifact registry, derivations, and framework adoption metadata.
- `product/foundation/`: problem, vision, and strategy.
- `product/domains/`: domains, goals, features, use cases, specifications, tasks, and validation artifacts.
- `product/knowledge/decisions/`: product decisions only.
- `product/audits/`: product audits and readiness reports.

## Framework-Owned Assets

Framework-owned assets may be copied or installed into the product repository:

- `.spec-framework/skills/`
- `.spec-framework/templates/`
- `.spec-framework/validators/`
- `.github/workflows/framework-validation.yml`

Do not edit framework-owned assets to encode product scope. Put product-specific rules in `product/`.

## First Handoff

Start with Problem Discovery AI, Vision AI, and Strategy AI before creating domains or features.
