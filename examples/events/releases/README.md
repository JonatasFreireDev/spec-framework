# Releases

## Purpose

Releases stores release readiness, scope, validation, and handoff information. It helps confirm that product, design, engineering, QA, audit, and decision artifacts agree before shipping.

## When To Use

Use this folder when grouping validated use cases into a release candidate, preparing a release gate, or recording released scope.

## Expected Files

- `README.md`: release process and expectations.
- Future files such as `RELEASE-001.md`, release notes, release checklists, and rollback summaries.

## Responsible Skill

Primary owner: Release Orchestrator.

Supporting skills: QA AI, Gap Finder AI, Conflict AI, Dependency AI, Impact Analysis AI, Documentation Writer AI.

## Release Gate

A release candidate should list:

- Included domains, goals, features, and use cases.
- Approved specifications and designs.
- Implementation plans, execution graphs, and tasks.
- Test evidence and QA results.
- Open risks and residual risk acceptance.
- Rollback plan and owner.

## Next Step

Before a release candidate is approved, run audit checks for gaps, conflicts, dependencies, impact, documentation freshness, tests, QA, review, and security when applicable.
