# Project Codex Configuration

## Purpose

This folder stores Codex-facing assets that belong to this repository. It keeps the framework skills inside the project so an agent can read and use them with the product documents they govern.

## When To Use

Use this folder when maintaining repository-local skills, skill metadata, or future Codex configuration. Do not store product artifacts here.

In this framework repository, worked product artifacts live under `examples/events/`. In adopter repositories, product artifacts live under `product/`. Framework method assets are resolved from the versioned external runtime.

## Expected Files

- `skills/`: repository-local maintenance skills for Codex. They are not mirrored into adopter repositories.

## Responsible Skill

Primary owner: Documentation Orchestrator.

## Next Step

When a skill changes, validate that its `SKILL.md` still has frontmatter with `name` and `description`, and that the instructions resolve product paths through the active product root instead of assuming product folders at the framework repository root.
