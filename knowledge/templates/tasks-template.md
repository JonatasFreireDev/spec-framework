# Tasks: [use case name]

## 🧭 Snapshot

| Field | Value |
| --- | --- |
| ID | `[TASKSET-XXX]` |
| Status | `[draft | proposed | approved]` |
| Source graph | `[GRAPH-XXX]` |
| Source specification | `[SPEC-XXX]` |
| Owner skill | Task AI |
| Next skill | Code Runner AI or QA AI |

## 🔗 Navigation

| Artifact | Link |
| --- | --- |
| Context | [context.md](context.md) |
| Specification | [specification.md](specification.md) |
| Implementation Plan | [implementation-plan.md](implementation-plan.md) |
| Execution Graph | [execution-graph.json](execution-graph.json) |
| Tests | [tests.md](tests.md) |
| Audit | [audit.md](audit.md) |

## 🚚 Delivery

| Field | Value |
| --- | --- |
| Level | `[L0 | L1 | L2 | L3 | L4 | L5]` |
| Priority | `[P0 | P1 | P2 | P3]` |
| Depends on | `[artifact ids/paths]` |
| Rationale | `[why this task set belongs here]` |

## 🕸️ Task Graph

```mermaid
flowchart LR
  T1["TK-XXX-001"] --> T2["TK-XXX-002"]
  T1 --> T3["TK-XXX-003"]
  T2 --> T4["TK-XXX-004"]
  T3 --> T4

  classDef done fill:#dcfce7,stroke:#16a34a,color:#14532d;
  classDef current fill:#fef3c7,stroke:#d97706,color:#78350f,stroke-width:3px;
  classDef pending fill:#f8fafc,stroke:#94a3b8,color:#334155;
  classDef blocked fill:#fee2e2,stroke:#dc2626,color:#7f1d1d;

  class T1 current;
  class T2,T3,T4 pending;
```

## ✅ Task List

| Task | Type | Depends On | Write Scope | Status | Acceptance |
| --- | --- | --- | --- | --- | --- |
| `TK-XXX-001` `[title]` | `[database/backend/frontend/test/analytics/docs/security]` | `[]` | `[paths/modules]` | `[pending/blocked/done]` | `[observable check]` |

## 🚧 Blocked Tasks

| Task | Blocking Reason | Decision/Dependency Needed | Owner |
| --- | --- | --- | --- |
| `[task]` | `[reason]` | `[decision/dependency]` | `[role]` |

## 🧪 Validation Methods

| Task | Validation |
| --- | --- |
| `[task]` | `[test/check/review]` |

## ⚠️ Parallelism Notes

- [Which tasks can run in parallel and why their write scopes do not overlap.]

## 🏁 Handoff

| Field | Value |
| --- | --- |
| Ready for implementation | `[yes/no]` |
| Required next skill | `[skill]` |
| Notes | `[notes]` |
