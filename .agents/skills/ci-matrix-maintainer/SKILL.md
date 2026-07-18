---
name: ci-matrix-maintainer
description: Maintain the Spec Framework verification matrix by mapping changed surfaces to unit, integration, race, validation, packaging, and release-smoke gates.
---

# CI Matrix Maintainer

## Purpose

Keep mechanical confidence proportional to the change. Do not suppress a failing, skipped, or environment-blocked gate.

## Required reading

- `AGENTS.md`
- `.github/workflows/`
- `.agents/skills/verify/SKILL.md`
- changed code, assets, and tests

## Workflow

1. Map each changed surface to its required checks: formatting, unit, integration, example validation, race, target smoke, init/upgrade, packaging, and release smoke.
2. Verify the workflow contains the required gate or document a justified local-only check.
3. Add combination tests when optional runtime capabilities interact with imports, dispatch, reviews, approvals, or upgrade.
4. Run the applicable checks and preserve exact output for failures.
5. Report App Control, unavailable binaries, or other environment limits as blocked evidence, never as passing evidence.

## Output

Return a check matrix with status, command or workflow, coverage rationale, blockers, and recommended next gate owner.
