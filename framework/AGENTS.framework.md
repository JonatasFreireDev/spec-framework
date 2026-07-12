# Framework Agent Instructions

## Purpose

This file describes how agents should use the installed Spec Framework assets inside an adopter repository.

The framework teaches the process. Product artifacts remain owned by the adopter product.

## Repository Boundary

Use these roots:

| Root | Purpose |
| --- | --- |
| Versioned user cache | Installed framework method, skills, templates, validators, tools, and framework decisions. |
| `product/` | Product-owned state, foundation, domains, decisions, audits, releases, and evidence. |
| User-scoped `spec-framework` dispatcher | Manifest-gated resolver for the pinned specialized skill contracts. |
| `product/knowledge/imports/` | Source evidence, immutable inventories, proposed mappings, conflicts, and import reports. |

Do not write product scope, product decisions, approval records, or delivery evidence into the external runtime cache.

Do not write framework-method decisions into `product/knowledge/decisions/`.

Use `.product/workspaces/WORK-NNN/` for concurrent feature focus and resume from `state.json`, the latest checkpoint, and the latest handoff; never invent a global active feature. A legacy `WORK-NNN.json` is read-only compatible until explicit runtime migration. Before implementation, require modular specification contracts by rigor, approved Design or structured `not_applicable`, Technical Discovery, a resolved Architecture Gate, applicable Engineering Proposal and passed current Engineering Review, configured gates, graph/tasks, and an active lease when graph runtime is used. `implemented` uses working-tree evidence and diff hash; commit only after Code Review and task QA approve that same hash, then require Integrated QA after combining task commits.

When the product declares `product/design/system/`, use the Design System skill for shared foundations, tokens, components, patterns, versions, and sources. UX/UI owns use-case Design and must pin the approved Design System version before proposed-or-later use. External visual tools are optional adapters; their installation, output, and availability never grant product approval or replace framework contracts.

## Required Reading

Before creating or updating framework-governed work:

1. Resolve and read the pinned framework root's `FRAMEWORK.md`.
2. Read the relevant `product/**/context.md` files.
3. Read the matching template in the pinned framework root's `templates/` when creating or normalizing an artifact.
4. Read approved product decisions in `product/knowledge/decisions/` and `product/.product/decisions.json` when relevant.
5. Read framework decisions in the pinned framework root's `decisions/` when the work touches method, gates, validators, skills, or workflow policy.

## Active Product Root

The active product root is `product/`.

When a skill mentions product-relative paths such as `knowledge/conventions/gates.md`, `.product/decisions.json`, `domains/`, `audits/`, or `releases/`, resolve them under `product/`.

When a skill mentions framework-relative paths such as `FRAMEWORK.md`, `templates/`, `skills/`, or framework decisions, resolve them under the versioned runtime returned by the CLI. Run executable operations through the installed `spec-framework` CLI.

Use `spec-framework guide` or `dashboard` when routing is unclear. Use `spec-framework adapters` for supervised optional-adapter discovery or installation; never install an external adapter silently.

## Gates

Run the product validator from the repository root:

```bash
spec-framework validate
```

If a product artifact is `approved` or later and its approval record is missing or inconsistent, report the blocker and stop. Agents must not create, edit, or repair approval records unless a human explicitly approves that migration.

## Reports

Save product reports under `product/audits/`.

Keep framework-upgrade or installation diagnostics in the external cache only when they are framework metadata, not product evidence.
