# Framework Agent Instructions

## Purpose

This file describes how agents should use the installed Spec Framework assets inside an adopter repository.

The framework teaches the process. Product artifacts remain owned by the adopter product.

## Activation Boundary

Activate these instructions only when the current repository contains a valid `product/.product/framework.json` that identifies `framework: spec-framework`, pins a concrete version, and declares `activation.mode: manifest-only`. If the manifest is absent or invalid, do not load specialized framework contracts or change product files; only explain the state or perform an explicitly requested bootstrap or `init` route.

## Repository Boundary

Use these roots:

| Root | Purpose |
| --- | --- |
| Versioned user cache | Installed framework method, skills, templates, validators, and tools. |
| `product/` | Product-owned state, foundation, domains, decisions, audits, releases, and evidence. |
| User-scoped `spec-framework` dispatcher | Manifest-gated resolver for the pinned specialized skill contracts. |
| `product/knowledge/imports/` | Source evidence, immutable inventories, proposed mappings, per-source traceability, conflicts, and import reports. |

Do not write product scope, product decisions, approval records, or delivery evidence into the external runtime cache.

Do not write framework-method decisions into any product decision domain.

Operational boundaries:

- Use `.product/workspaces/WORK-NNN/` for concurrent focus; never invent a global active feature.
- Resume from `state.json`, the latest checkpoint, and the latest handoff. Legacy `WORK-NNN.json` is read-only until explicit migration.
- Follow the detailed implementation prerequisites, evidence, status, and approval gates in `FRAMEWORK.md` and the owning skills; this file defines only their cross-agent boundaries.
- Do not write outside the current skill's declared scope. Record working-tree evidence, diff hashes, commits, and integrated QA only when the applicable delivery contract requires them.

When the product declares `product/design/system/`, route shared foundations, tokens, components, patterns, versions, and sources to the Design System skill. Detailed Design System and UX/UI gates remain in `FRAMEWORK.md` and their owning skills. External visual tools are optional adapters; their installation, output, and availability never grant product approval or replace framework contracts.

## Required Reading

Authority order:

```text
FRAMEWORK.md → owning skill → matching template → product context and decisions → current CLI evidence
```

`FRAMEWORK.md` defines the method and gates. This file defines cross-agent behavior. The owning skill defines specialized responsibility and write scope. The template defines artifact shape. Product context and approved decisions define product intent. CLI output defines current mechanical state. Later sources specialize earlier sources within their scope but cannot weaken framework gates, human authority, or approval requirements.

## Common Agent Rules

These rules apply to every agent and every skill. Read and apply them before following a specialized `SKILL.md`.

- Inspect the current repository and CLI evidence before asking for facts that can be discovered safely.
- Ask the human whenever uncertainty could change scope, architecture, compatibility, migration, or a consequential decision. Do not invent requirements, decisions, data, approvals, or evidence.
- Keep references to existing artifacts, documents, sections, code, and evidence navigable with relative Markdown links or the repository's canonical link format. External references and not-yet-materialized artifacts must declare their type and state. Treat an unlinked or broken repository reference as a gap.
- After an implementation, review the complete change for blockers, functional, technical, documentation, testing, compatibility, CI, installation, upgrade, and distribution gaps. Fix findings only when the current skill has write authority and the correction is within scope; otherwise route the finding to its owner. Repeat the review after each correction until no known blocker or gap remains.
- Preserve product scope, history, approval records, delivery evidence, and adopter-owned content. Never repair approval records without an explicitly authorized migration.
- When a change has variants, optional integrations, endpoints, starting points, or configuration modes, evaluate a modular and configurable design. Define optional modules, activation parameters, dependencies, default behavior, safe plug-in/plug-out behavior, compatibility expectations, and combination tests. For changes without meaningful variation, record that modularity was evaluated and is not applicable. Keep optional capabilities decoupled from the core when practical.
- Revalidate mechanical state after commands. Do not assume a command, write, test, approval, or migration succeeded without checking its result.
- For implementations, audits, and configuration changes, report the summary, decisions, validations, blockers and gaps found, corrections applied, pending questions, and residual risks. Explicitly state when no known blockers or gaps remain. If a finding cannot be resolved, report its reason, impact, and recommended action. For explanations and read-only diagnostics, report evidence, current state, risks, and next action without inventing a mutation report.

Before creating or updating framework-governed work:

1. Resolve and read the pinned framework root's `FRAMEWORK.md`.
2. Read the relevant `product/**/context.md` files.
3. Read the matching template in the pinned framework root's `templates/` when creating or normalizing an artifact.
4. Read approved decisions from `product/.product/decisions.json` and resolve each record from its registered `path`; default domain roots are `knowledge/decisions/`, `design/decisions/`, and `engineering/decisions/`.

Before writing, confirm the owning skill, artifact scope, and allowed write scope. Do not edit artifacts owned by another skill without an explicit handoff or route.

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

When a human requests approval, use the batch-capable approval route whenever the scope contains more than one artifact or names a Foundation/stage. First run `spec-framework approve-batch` without `--yes` to preview the exact files, IDs, hashes, ignored items, blockers, and next gate. Ask for explicit human confirmation and approver identity; only then rerun with `--yes`. Never infer approval from conversational agreement. The command supports `--artifact`, `--ids`, `--foundation`, `--stage`, `--all-eligible`, and `--until`; stale artifacts, invalid parents, and blockers must remain excluded.

Whenever a workflow stage is complete and the next stage would require an approval gate, stop after the preview and ask the human whether the listed artifacts should be approved. Do not continue to the next stage until the human confirms the exact listed scope.

## Stop Conditions

Stop and ask for human direction when any of the following applies:

- a blocking question could change scope, architecture, compatibility, migration, or a consequential decision;
- an approval is missing, inconsistent, stale, or required for the requested mutation;
- the requested scope or owning skill is ambiguous or conflicting;
- a finding requires a decision, approval-record migration, destructive action, remote mutation, or authority outside the current skill;
- required gate configuration, parent evidence, or source traceability is missing;
- a command result conflicts with the expected mechanical state.

## Reports

Save product reports under `product/audits/` when the operation is authorized to create a product report. `audit-only` and read-only diagnostics must not create reports unless explicitly requested and permitted by their contract.

Keep framework-upgrade or installation diagnostics in the external cache only when they are framework metadata, not product evidence.
