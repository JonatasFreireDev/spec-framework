# Agent Instructions

## Purpose

This repository is a product repository adopting the Spec Framework. Product documentation is execution infrastructure.

## Source Of Truth

Activate only when `product/.product/framework.json` is valid, then resolve and read the pinned framework root's `FRAMEWORK.md` before generating downstream artifacts.

## Product-Owned Areas

- `product/.product/`: product state, approval history, artifact registry, derivations, and framework adoption metadata.
- `product/foundation/`: problem, vision, and strategy.
- `product/domains/`: domains, goals, features, use cases, specifications, tasks, and validation artifacts.
- `product/knowledge/decisions/`: product decisions only.
- `product/audits/`: product audits and readiness reports.

## Framework-Owned Assets

Framework-owned assets may be copied or installed into the product repository:

- the pinned framework runtime's `skills/`
- the pinned framework runtime's `templates/`
- the pinned framework runtime's validators
- `.github/workflows/framework-validation.yml`

Do not edit framework-owned assets to encode product scope. Put product-specific rules in `product/`.

## First Handoff

Start with Framework Guide when the next command or gate is unclear. For a new product, continue through Product Orchestrator, Problem Discovery, Vision, and Strategy before creating domains or features.

## Canonical Delivery Flow

```text
Problem -> Vision -> Strategy -> Domain -> User Goal -> Feature -> Use Case
-> Specification -> Design -> Technical Discovery -> Architecture Gate
-> Engineering Proposal -> Engineering Review -> Implementation Plan
-> Execution Graph -> Tasks -> Code Runner
-> Code Review -> QA -> Commit Crafter -> PR Finalizer -> Release
```

Do not skip approval gates. A downstream artifact may remain `draft` while its parent is incomplete, but it must not advance to `proposed`, `approved`, or later.

## Design And Design System

- Shared foundations, tokens, components, and patterns live under `product/design/system/` and are owned by the Design System skill.
- Use-case `design.md` remains the canonical experience contract and is owned by UX/UI.
- When a Design System is declared, proposed-or-later Design must pin its approved id/version and record consumed tokens, components, patterns, and deviations.
- Impeccable, Figma, Penpot, images, and other tools are optional adapters or sources; they never replace Specification, `design.md`, UX Review, or human approval.
- Adapter installation is explicit and supervised. Skipping Impeccable does not block the framework.

## Implementation Evidence

- Before implementation, require approved Specification, approved Design or structured `not_applicable`, Technical Discovery, resolved Architecture Gate, applicable Engineering Proposal and passed Engineering Review, Implementation Plan, Execution Graph, Tasks, configured gates, and an active lease when graph runtime is used.
- Engineering Proposal and Review apply to every Tier L delivery and to Tier S/M when `context.md` declares a supported `engineering_trigger`; never infer triggers from prose.
- `implemented` requires immutable working-tree evidence: branch, base commit, changed paths, diff hash, tests, and applicable gate results. It does not require an early commit.
- Code Review and task QA must approve the same current diff hash before Commit Crafter creates local commits.
- After task commits are integrated, require Integrated QA where applicable.
- `validated` or `released` requires structured test and review evidence, code paths, commits, PR when repository policy requires it, and concrete CI, gate, screenshot, or QA evidence.
- Agents must not create or repair approval records unless a human explicitly authorizes a named migration.

## Bootstrap Path

1. Read `product/context.md`.
2. Complete the foundation contexts and documents under `product/foundation/`.
3. Replace `TBD` commands in `product/knowledge/conventions/gates.md` before implementation work.
4. Copy `product/domains/_template-domain/` when creating the first real domain, then use Domain Evolution to select a delivery slice before New Feature orchestration.
5. Use `spec-framework work`, `status`, and `next` to navigate a selected feature without global mutable focus.
6. Use `spec-framework guide` and `dashboard` to inspect the current gate, blockers, decisions, Design System, graph, tasks, and next action.
7. Do not invoke Code Runner while `spec-framework gates` reports applicable `TBD` commands.
8. Rename template folders to stable slugs and update every `slug` field in `context.md`.
9. Keep product scope inside `product/`; do not encode product-specific rules inside the external runtime cache.
