# Framework Agent Instructions

## Purpose

This file describes how agents should use the installed Spec Framework assets inside an adopter repository.

The framework teaches the process. Product artifacts remain owned by the adopter product.

## Repository Boundary

Use these roots:

| Root | Purpose |
| --- | --- |
| Versioned user cache | Installed framework method, skills, templates, validators, and tools. |
| `product/` | Product-owned state, foundation, domains, decisions, audits, releases, and evidence. |
| User-scoped `spec-framework` dispatcher | Manifest-gated resolver for the pinned specialized skill contracts. |
| `product/knowledge/imports/` | Source evidence, immutable inventories, proposed mappings, per-source traceability, conflicts, and import reports. |

Do not write product scope, product decisions, approval records, or delivery evidence into the external runtime cache.

Do not write framework-method decisions into `product/knowledge/decisions/`.

Operational boundaries:

- Use `.product/workspaces/WORK-NNN/` for concurrent focus; never invent a global active feature.
- Resume from `state.json`, the latest checkpoint, and the latest handoff. Legacy `WORK-NNN.json` is read-only until explicit migration.
- Before implementation, require the rigor-appropriate Specification contracts, approved Design or structured `not_applicable`, Technical Discovery, resolved Architecture Gate, applicable Engineering Proposal and current passed Engineering Review, configured gates, Graph, Tasks, and an active lease when used.
- Record working-tree evidence and diff hash at `implemented`. Commit only after Code Review and task QA approve that same hash; require Integrated QA after task integration.

When the product declares `product/design/system/`, use the Design System skill for shared foundations, tokens, components, patterns, versions, and sources. UX/UI owns use-case Design and must pin the approved Design System version before proposed-or-later use. External visual tools are optional adapters; their installation, output, and availability never grant product approval or replace framework contracts.

## Required Reading

Authority order:

```text
FRAMEWORK.md → owning skill → matching template → product context and decisions → current CLI evidence
```

Later sources specialize earlier ones within their scope but cannot weaken framework gates or human approval requirements.

Before creating or updating framework-governed work:

1. Resolve and read the pinned framework root's `FRAMEWORK.md`.
2. Read the relevant `product/**/context.md` files.
3. Read the matching template in the pinned framework root's `templates/` when creating or normalizing an artifact.
4. Read approved product decisions in `product/knowledge/decisions/` and `product/.product/decisions.json` when relevant.

## Active Product Root

The active product root is `product/`.

When a skill mentions product-relative paths such as `knowledge/conventions/gates.md`, `.product/decisions.json`, `domains/`, `audits/`, or `releases/`, resolve them under `product/`.

When a skill mentions framework-relative paths such as `FRAMEWORK.md`, `templates/`, or `skills/`, resolve them under the versioned runtime returned by the CLI. Run executable operations through the installed `spec-framework` CLI.

Use Framework Guide as the default route unless current CLI output or an explicit human request names both the specialist and concrete scope. Revalidate persisted handoffs/checkpoints with `dashboard`, `status`, `next`, or `guide` before routing directly. A skill name without scope is only a hint. Use `spec-framework adapters` for supervised optional-adapter discovery or installation; never install an external adapter silently.

For definition and planning work, follow the pinned runtime's `skills/discovery-and-challenge.md`. Inspect evidence first and use the harness-native structured question tool whenever it is available. The canonical capability is `native_user_question`; adapters map it to the harness default. Do not silently replace an available question tool with assumptions or questions embedded in a final response.

## Gates

Run the product validator from the repository root:

```bash
spec-framework validate
```

If a product artifact is `approved` or later and its approval record is missing or inconsistent, report the blocker and stop. Agents must not create, edit, or repair approval records unless a human explicitly approves that migration.

## Reports

Save product reports under `product/audits/`.

Keep framework-upgrade or installation diagnostics in the external cache only when they are framework metadata, not product evidence.
