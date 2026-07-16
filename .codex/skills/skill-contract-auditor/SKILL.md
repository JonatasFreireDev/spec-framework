---
name: skill-contract-auditor
description: Audit repository-maintenance and shipped Spec Framework skills for valid metadata, ownership, portable paths, handoffs, and consistency with the canonical method.
---

# Skill Contract Auditor

## Purpose

Keep skills actionable and non-overlapping. This skill audits contracts; it does not silently change a skill's product ownership or workflow.

## Required reading

- `FRAMEWORK.md`
- `AGENTS.md`
- `.codex/README.md` for maintenance skills
- `framework/skills/README.md` for shipped skills

## Workflow

1. Validate each `SKILL.md` has only `name` and `description` frontmatter and a precise activation description.
2. Verify the skill has one clear responsibility, explicit exclusions, required reading, and a safe handoff.
3. Check paths work from the framework runtime and active product root; reject lab-specific or generated-tree source assumptions.
4. Check for stale concepts, nonexistent commands, duplicated ownership, bypassed gates, or approval authority violations.
5. Route method-level corrections to `FRAMEWORK.md`, contracts, validators, and tests in the same change.

## Output

Report compliant skills, findings by severity, affected owners, exact corrective surface, and verification required after correction.
