---
name: distribution-impact-analyzer
description: Analyze the distribution impact of a Spec Framework asset or packaging change across embedding, install, upgrade, documentation, agent targets, and release verification.
---

# Distribution Impact Analyzer

## Purpose

Make distribution effects explicit before a framework asset or packaging change is committed.

## Required reading

- `FRAMEWORK.md`
- `assets.go`
- affected `starter/` and `framework/` assets
- installer, release, and smoke-test code

## Workflow

1. Determine whether the change is embedded, generated, copied at init, materialized at runtime, or documentation-only.
2. Complete an impact matrix for assets, init, upgrade, Codex/Cursor/Claude targets, documentation, CI, archives, checksums, and release smoke.
3. Identify managed versus adopter-owned paths and verify that removals remain safe after upgrade.
4. Require the smallest relevant install, upgrade, target, and release tests.
5. State whether a version bump, migration note, or rollback plan is required.

## Output

Return the impact matrix, affected paths, required verification, release impact, compatibility verdict, and migration/rollback notes.
