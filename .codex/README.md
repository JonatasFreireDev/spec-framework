# Project Codex Configuration

## Purpose

This folder stores Codex-facing assets that belong to this repository. It keeps the framework skills inside the project so an agent can read and use them with the product documents they govern.

## When To Use

Use this folder when maintaining repository-local skills, skill metadata, or future Codex configuration. Do not store product artifacts here; product artifacts belong in `foundation/`, `knowledge/`, `domains/`, `audits/`, `engineering/`, and `releases/`.

## Expected Files

- `skills/`: repository-local Codex skills in folder-per-skill format.

## Responsible Skill

Primary owner: Documentation Orchestrator.

## Next Step

When a skill changes, validate that its `SKILL.md` still has frontmatter with `name` and `description`, and that the instructions point to real framework paths.
