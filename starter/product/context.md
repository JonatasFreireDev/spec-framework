# Product Context

```yaml
id: PRODUCT-TBD
type: product
name: TBD Product
status: draft
owner_skill: Product Orchestrator
slug: product

parents: []

children:
  - foundation/problem
  - foundation/vision
  - foundation/strategy
  - domains

depends_on: []
used_by:
  - releases
related:
  - .product/framework.json

documents:
  canonical: README.md
  problem: foundation/problem/problem.md
  vision: foundation/vision/vision.md
  strategy: foundation/strategy/strategy.md

delivery:
  level: L0
  priority: P0
  depends_on: []
  rationale: Product foundation must exist before domains, features, and implementation planning.

open_questions:
  - What product name, primary audience, and first domain should replace the starter placeholders?

decisions: []
next_recommended_skill: Problem Discovery AI
```

## Purpose

This file orients agents at the product root. Replace TBD values as soon as the product has a real name, owner, and first validated problem.

## Current Flow

```mermaid
flowchart LR
  problem["Problem"] --> vision["Vision"] --> strategy["Strategy"] --> domain["Domain"] --> goal["Goal"] --> feature["Feature"] --> usecase["Use Case"] --> spec["Specification"]

  classDef current fill:#fff3bf,stroke:#f59f00,color:#1f1f1f;
  classDef pending fill:#f1f3f5,stroke:#adb5bd,color:#495057;
  class problem current;
  class vision,strategy,domain,goal,feature,usecase,spec pending;
```

## Bootstrap Rule

Do not create domains or features until `foundation/problem/problem.md`, `foundation/vision/vision.md`, and `foundation/strategy/strategy.md` contain product-specific content.
