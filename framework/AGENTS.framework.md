# Framework Agent Instructions

## Purpose

This file describes how agents should use the installed Spec Framework assets inside an adopter repository.

The framework teaches the process. Product artifacts remain owned by the adopter product.

## Repository Boundary

Use these roots:

| Root | Purpose |
| --- | --- |
| `.spec-framework/` | Installed framework method, skills, templates, validators, tools, and framework decisions. |
| `product/` | Product-owned state, foundation, domains, decisions, audits, releases, and evidence. |
| `.agents/skills/`, `.cursor/skills/`, `.claude/skills/` | Generated agent-specific skill trees. Treat these as derived copies of `.spec-framework/skills/`. |
| `product/knowledge/imports/` | Source evidence, immutable inventories, proposed mappings, conflicts, and import reports. |

Do not write product scope, product decisions, approval records, or delivery evidence into `.spec-framework/`.

Do not write framework-method decisions into `product/knowledge/decisions/`.

Use `.product/workspaces/WORK-NNN.json` for concurrent feature focus; never invent a global active feature. Before implementation, require modular specification contracts by rigor, approved Design, Technical Discovery, a resolved Architecture Gate, configured gates, graph/tasks, and an exclusive claim when graph runtime is used. `implemented` uses working-tree evidence and diff hash; commit only after Code Review and QA approve that same hash.

## Required Reading

Before creating or updating framework-governed work:

1. Read `.spec-framework/FRAMEWORK.md`.
2. Read the relevant `product/**/context.md` files.
3. Read the matching template in `.spec-framework/templates/` when creating or normalizing an artifact.
4. Read approved product decisions in `product/knowledge/decisions/` and `product/.product/decisions.json` when relevant.
5. Read framework decisions in `.spec-framework/decisions/` when the work touches method, gates, validators, skills, or workflow policy.

## Active Product Root

The active product root is `product/`.

When a skill mentions product-relative paths such as `knowledge/conventions/gates.md`, `.product/decisions.json`, `domains/`, `audits/`, or `releases/`, resolve them under `product/`.

When a skill mentions framework-relative paths such as `FRAMEWORK.md`, `templates/`, `skills/`, or framework decisions, resolve them under `.spec-framework/`. Run executable operations through the installed `spec-framework` CLI.

## Gates

Run the product validator from the repository root:

```bash
spec-framework validate
```

If a product artifact is `approved` or later and its approval record is missing or inconsistent, report the blocker and stop. Agents must not create, edit, or repair approval records unless a human explicitly approves that migration.

## Reports

Save product reports under `product/audits/`.

Save framework-upgrade or installation diagnostics under `.spec-framework/` only when they are framework metadata, not product evidence.
