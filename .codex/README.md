# Project Codex Configuration

## Purpose

This folder stores Codex-facing assets that belong to this repository. It keeps the framework skills inside the project so an agent can read and use them with the product documents they govern.

## When To Use

Use this folder when maintaining repository-local skills, skill metadata, or future Codex configuration. Do not store product artifacts here.

In this framework repository, worked product artifacts live under `examples/events/`. In adopter repositories, product artifacts live under `product/`. Framework method assets are resolved from the versioned external runtime.

## Expected Files

- `skills/`: repository-local maintenance skills for Codex. They are not mirrored into adopter repositories.

## Ownership and routing

| Surface | Canonical owner | Use it for |
| --- | --- | --- |
| `FRAMEWORK.md` | Framework method | Gates, lifecycle, authority, and architecture. |
| `framework/skills/` | Distributed framework skills | Product-facing specialist and orchestrator contracts. |
| `.codex/skills/` | Repository maintenance skills | Reviewing, testing, packaging, and preserving this repository. |
| `starter/` | Clean adopter skeleton | New-product content only. |
| `examples/events/` | Worked product fixture | Learning and validation; never a starter source. |
| Git history | Maintenance record | Evolution of the framework method and contracts. |

Use `framework-change-reviewer` before or after a method-level change; use `distribution-impact-analyzer` and `upgrade-compatibility` whenever embedded, starter, installer, or packaging assets change. Use `runtime-contract-tester` for CLI or runtime changes, `agent-target-smoke` for harness output, `fixture-governor` for starter/example boundaries, `skill-contract-auditor` for skill contracts, and `ci-matrix-maintainer` to select and maintain mechanical gates.

## Maintenance rules

- Read `FRAMEWORK.md` before changing framework behavior.
- Update every affected contract, validator, template, starter asset, documentation surface, and test in the same change.
- Do not edit generated agent trees as canonical sources or create product artifacts under `.codex/`.
- Preserve adopter-owned content and approval records during `upgrade`; never use an overwrite to repair a migration.
- Keep optional capabilities explicit, disabled by default when they change execution behavior, and covered by combination tests.
- Report unavailable environment gates as limitations, never as passing evidence.
- Record framework method evolution directly in `FRAMEWORK.md`, affected contracts, validators, and tests; Git is the maintenance record.

## Responsible Skill

Primary owner: Documentation Orchestrator. Use the specialist maintenance skill matching the changed surface before handoff.

## Next Step

When a skill changes, run `skill-contract-auditor` or at minimum validate that its `SKILL.md` has frontmatter with `name` and `description`, that its instructions resolve product paths through the active product root, and that its guidance remains consistent with `FRAMEWORK.md`.
