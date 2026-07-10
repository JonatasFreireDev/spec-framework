# Vision Context

```yaml
id: VISION-TBD
type: vision
name: TBD Vision
status: draft
owner_skill: Vision AI
slug: vision

parents:
  - PROBLEM-TBD

children:
  - principles.md
  - north-star.md

depends_on:
  - foundation/problem/problem.md
used_by:
  - foundation/strategy/strategy.md
related:
  - foundation/problem/opportunities.md

documents:
  canonical: vision.md
  principles: principles.md
  north_star: north-star.md

delivery:
  level: L0
  priority: P0
  depends_on:
    - foundation/problem/problem.md
  rationale: Vision translates the validated problem into a product direction before strategy and domains.

open_questions:
  - What product promise should stay stable while features evolve?

decisions: []
next_recommended_skill: Strategy AI
```

## Purpose

Use this context to keep vision work anchored to the problem instead of feature lists.

