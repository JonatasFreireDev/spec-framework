# Context: Product Problem

```yaml
id: PROB-001
type: problem
name: Product Problem
status: draft
owner_skill: problem-discovery
parents: []
children:
  - VIS-001
depends_on: []
used_by:
  - foundation/vision/context.md
documents:
  canonical: problem.md
  opportunities: opportunities.md
delivery:
  level: L0
  priority: P0
  rationale: The product needs a validated problem before vision, strategy, domains, or implementation can be safely derived.
open_questions:
  - Which user pain is the first product bet?
  - Which evidence is strong enough to approve the problem statement?
decisions: []
```

## Purpose

This context explains the root problem layer. It tells an agent what evidence to inspect before proposing vision, strategy, domains, or features.

## Required Reading

- `FRAMEWORK.md`
- `problem.md`
- `opportunities.md`
- Notes in `researches/` and `interviews/`

## Handoff

Next recommended skill: Vision AI.

Before handoff, the problem should clearly state the target user, the painful situation, the current workaround, evidence quality, assumptions, non-goals, and approval state.
