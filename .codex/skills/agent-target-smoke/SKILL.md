---
name: agent-target-smoke
description: Smoke-test Spec Framework initialization and upgrade output for Codex, Cursor, and Claude Code agent targets and their Guide-first dispatchers.
---

# Agent Target Smoke

## Purpose

Verify the generated harness entry points without treating any target-specific tree as the source of truth.

## Required reading

- `FRAMEWORK.md`
- `internal/install/`
- `internal/dispatcher/`
- relevant init, upgrade, and dispatcher tests

## Workflow

1. Initialize an isolated temporary product for each requested target: Codex, Cursor, and Claude Code.
2. Verify the target receives only its supported dispatcher and that the dispatcher requires a valid `product/.product/framework.json`.
3. Verify that framework skills and templates resolve from the versioned runtime, not from repository-local copied trees.
4. Upgrade the fixture and verify dispatcher refresh while product content remains unchanged.
5. Check target instructions for portable paths and absence of unsupported harness-specific instructions.

## Output

Report target, generated paths, activation result, upgrade result, unsupported-target findings, and exact test evidence.
