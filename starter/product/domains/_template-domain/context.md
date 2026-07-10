# Domain Context

```yaml
id: DOMAIN-TEMPLATE
type: domain
name: TBD Domain
status: draft
owner_skill: Domain Architect AI
slug: _template-domain

parents:
  - STRATEGY-TBD

children:
  - GOAL-TEMPLATE

depends_on:
  - foundation/strategy/strategy.md
used_by: []
related:
  - domain.md
  - decisions.md

documents:
  canonical: domain.md
  decisions: decisions.md
  goals: goals/

delivery:
  level: L1
  priority: P0
  depends_on:
    - foundation/strategy/strategy.md
  rationale: Replace with the reason this domain belongs to the selected delivery level.

open_questions:
  - What responsibilities belong inside this domain, and what belongs elsewhere?

decisions: []
next_recommended_skill: User Goal AI
```

## Bootstrap Notes

Update `slug` to match the renamed folder. Keep the slug stable after creation.

## Handoff

| Field | Value |
| --- | --- |
| Next skill | user-goal |
