# Context: Example Feature

```yaml
id: FT-EXAMPLE
type: feature
name: Example Feature
status: draft
owner_skill: feature
slug: _example-feature
last_updated: 2026-07-09
delivery:
  level: N/A
  priority: N/A
  depends_on:
    - GOAL-EXAMPLE
  rationale: Structural example only; not product scope.
```

## Purpose

This context demonstrates how a feature links to a parent user goal and child use cases.

## Parent Artifacts

- GOAL-EXAMPLE - ../../context.md - example parent goal.

## Child Artifacts

- UC-EXAMPLE - use-cases/_example-use-case/context.md - example use case.

## Dependencies

- GOAL-EXAMPLE - parent example goal - blocking: yes.

## Related Artifacts

- framework/template/feature-template.md - feature structure.
- framework/template/context-template.md - context structure.

## Canonical Documents

- Primary: feature.md
- Specification: use-cases/_example-use-case/specification.md
- Design: use-cases/_example-use-case/design.md
- Implementation plan: use-cases/_example-use-case/implementation-plan.md
- Execution graph: use-cases/_example-use-case/execution-graph.json
- Tasks: use-cases/_example-use-case/tasks.md

## Decisions

- N/A - example only.

## Assumptions

- Example feature content is illustrative and must be replaced for real product work.

## Open Questions

- N/A.

## Handoff

Next recommended skill: use-case
Required reading before next step:
- feature.md
- use-cases/_example-use-case/context.md
