# Use Case Context

```yaml
id: UC-TEMPLATE
type: use-case
name: TBD Use Case
status: draft
owner_skill: Use Case AI
slug: _template-use-case
rigor_tier: N/A
engineering_triggers: []

parents:
  - FT-TEMPLATE

children:
  - SPEC-TEMPLATE
  - DESIGN-TEMPLATE
  - PLAN-TEMPLATE
  - GRAPH-TEMPLATE
  - TASK-TEMPLATE-001

depends_on:
  - ../../feature.md
used_by: []
related:
  - use-case.md
  - specification.md
  - design.md
  - engineering-proposal.md
  - engineering-review.md
  - implementation-plan.md
  - execution-graph.json
  - tasks.md
  - tests.md

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
  canonical: use-case.md
  specification: specification.md
  design: design.md
  engineering_proposal: engineering-proposal.md
  engineering_review: engineering-review.md
  implementation_plan: implementation-plan.md
  execution_graph: execution-graph.json
  tasks_index: tasks.md
  tests: tests.md
  analytics: analytics.md
  qa_evidence: qa-evidence.md
  security_review: security-review.md
  code_review: code-review.md
  audit: audit.md

delivery:
  level: L1
  priority: P0
  depends_on:
    - FT-TEMPLATE
  rationale: Replace with why this interaction is required for the delivery level.

open_questions:
  - What tier applies: S, M, or L?
  - Which parent artifact version/hash generated each downstream artifact?

decisions: []
next_recommended_skill: Specification AI
```

## Progress

```mermaid
flowchart LR
  usecase["Use Case"] --> spec["Specification"] --> design["Design"] --> plan["Implementation Plan"] --> graph["Execution Graph"] --> tasks["Tasks"] --> code["Code"] --> validation["Validation"] --> audit["Audit"]

  classDef current fill:#fff3bf,stroke:#f59f00,color:#1f1f1f;
  classDef pending fill:#f1f3f5,stroke:#adb5bd,color:#495057;
  class usecase current;
  class spec,design,plan,graph,tasks,code,validation,audit pending;
```

## Handoff

| Field | Value |
| --- | --- |
| Next skill | specification |
