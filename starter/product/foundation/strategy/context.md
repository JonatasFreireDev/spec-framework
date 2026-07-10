# Strategy Context

```yaml
id: STRATEGY-TBD
type: strategy
name: TBD Strategy
status: draft
owner_skill: Strategy AI
slug: strategy

parents:
  - VISION-TBD

children:
  - personas.md
  - competitors.md
  - metrics.md
  - roadmap.md

depends_on:
  - foundation/problem/problem.md
  - foundation/vision/vision.md
used_by:
  - domains/
related:
  - foundation/vision/north-star.md

documents:
  canonical: strategy.md
  personas: personas.md
  competitors: competitors.md
  metrics: metrics.md
  roadmap: roadmap.md

delivery:
  level: L0
  priority: P0
  depends_on:
    - foundation/vision/vision.md
  rationale: Strategy decides the first product bets before domain modeling begins.

open_questions:
  - Which segment and delivery level should the first domain serve?

decisions: []
next_recommended_skill: Domain Architect AI
```

## Purpose

Use this context to ensure domains and goals come from explicit strategic choices.

