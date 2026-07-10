# Goal Context

```yaml
id: GOAL-TEMPLATE
type: user-goal
name: TBD Goal
status: draft
owner_skill: User Goal AI
slug: _template-goal

parents:
  - DOMAIN-TEMPLATE

children:
  - FT-TEMPLATE

depends_on:
  - ../../domain.md
used_by: []
related:
  - goal.md
  - journeys.md

documents:
  canonical: goal.md
  journeys: journeys.md
  features: features/

delivery:
  level: L1
  priority: P0
  depends_on:
    - DOMAIN-TEMPLATE
  rationale: Replace with why this user goal matters for the delivery level.

open_questions:
  - What user outcome proves this goal has been achieved?

decisions: []
next_recommended_skill: Feature AI
```

## Handoff

| Field | Value |
| --- | --- |
| Next skill | feature |
