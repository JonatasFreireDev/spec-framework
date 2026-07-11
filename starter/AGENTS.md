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

Start with Product Orchestrator, Problem Discovery AI, Vision AI, and Strategy AI before creating domains or features.

## Bootstrap Path

1. Read `product/context.md`.
2. Complete the foundation contexts and documents under `product/foundation/`.
3. Replace `TBD` commands in `product/knowledge/conventions/gates.md` before implementation work.
4. Copy `product/domains/_template-domain/` when creating the first real domain, then use Domain Evolution to select a delivery slice before New Feature orchestration.
5. Use `spec-framework work`, `status`, and `next` to navigate a selected feature without global mutable focus.
6. Do not invoke Code Runner while `spec-framework gates` reports applicable `TBD` commands.
7. Rename template folders to stable slugs and update every `slug` field in `context.md`.
8. Keep product scope inside `product/`; do not encode product-specific rules inside `.spec-framework/`.
