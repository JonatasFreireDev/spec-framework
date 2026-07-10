# Context: Example Goal

```yaml
id: GOAL-EXAMPLE
type: goal
name: Example Goal
status: draft
owner_skill: user-goal
slug: _example-goal
last_updated: 2026-07-09
delivery:
  level: N/A
  priority: N/A
  depends_on:
    - DOMAIN-EXAMPLE
  rationale: Structural example only; not product scope.
```

## Purpose

This context demonstrates how a user goal links to a parent domain and child features.

## Parent Artifacts

- DOMAIN-EXAMPLE - ../../context.md - example parent domain.

## Child Artifacts

- FT-EXAMPLE - features/_example-feature/context.md - example feature.

## Dependencies

- DOMAIN-EXAMPLE - parent example domain - blocking: yes.

## Related Artifacts

- knowledge/templates/goal-template.md - goal structure.
- knowledge/templates/context-template.md - context structure.

## Canonical Documents

- Primary: goal.md
- Specification: N/A
- Design: N/A
- Implementation plan: N/A
- Execution graph: N/A
- Tasks: N/A

## Decisions

- N/A - example only.

## Assumptions

- Example goal content is illustrative and must be replaced for real product work.

## Open Questions

- N/A.

## Handoff

Next recommended skill: feature
Required reading before next step:
- goal.md
- features/_example-feature/context.md
