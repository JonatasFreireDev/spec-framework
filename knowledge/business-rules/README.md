# Business Rules

## Purpose

Record cross-cutting rules that affect multiple domains, goals, features, or use cases.

## When To Use

Use this folder when a rule is broader than a single use case or feature. Local rules can live in their artifact, but shared rules should be documented here and referenced from specifications.

## Expected Files

- `README.md`: folder purpose and usage.
- Future files named by rule area, for example `attendance.md`, `permissions.md`, or `privacy.md`.

## Responsible Skill

Primary owner: Specification AI.

Supporting skills: Domain Architect AI, Product Historian AI, QA AI.

## Next Step

When a rule changes product behavior, record a decision in `knowledge/decisions/` and reference it from affected specifications.
