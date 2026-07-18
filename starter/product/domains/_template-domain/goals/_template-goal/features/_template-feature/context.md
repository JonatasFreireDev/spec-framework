# Feature Context

```yaml
id: FT-TEMPLATE
type: feature
name: TBD Feature
status: draft
owner_skill: Feature AI
slug: _template-feature

parents:
  - GOAL-TEMPLATE

children:
  - UC-TEMPLATE

depends_on:
  - ../../goal.md
used_by: []
related:
  - feature.md

relations:
  extends: null
  reuses: []
  depends_on: []
  impacts: []
  supersedes: []

traceability:
  source_demand: null
  source_documents: []
  source_decisions: []

evolution:
  type: new
  previous_version: null
  change_summary: null

documents:
  canonical: feature.md
  use_cases: use-cases/

delivery:
  level: L1
  priority: P0
  depends_on:
    - GOAL-TEMPLATE
  rationale: Replace with why this feature is needed for the delivery level.

open_questions:
  - Which use cases are required for this feature to be valuable?

decisions: []
next_recommended_skill: Use Case AI
```

## Handoff

| Field | Value |
| --- | --- |
| Next skill | use-case |
