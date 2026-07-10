# Roadmap Item: [item name]

## 🧭 Snapshot

| Field | Value |
| --- | --- |
| ID | `[ROAD-XXX]` |
| Status | `[draft | proposed | approved]` |
| Source strategy | `[STRAT-XXX]` |
| Owner skill | Strategy AI |

## 🚚 Delivery

| Field | Value |
| --- | --- |
| Level | `[L0 | L1 | L2 | L3 | L4 | L5]` |
| Priority | `[P0 | P1 | P2 | P3]` |
| Depends on | `[artifact ids/paths]` |
| Rationale | `[why this level and priority are assigned]` |

## 🎯 Outcome

[What changes for users or operations when complete.]

## 🧱 Scope

| Includes | Excludes |
| --- | --- |
| `[scope]` | `[non-goal]` |

## 🗺️ Roadmap Flow

```mermaid
flowchart LR
  S["Strategy"] --> R["Roadmap Item"]
  R --> D["Domain"]
  D --> G["Goal"]
  G --> F["Feature"]
  F --> U["Use Cases"]

  classDef done fill:#dcfce7,stroke:#16a34a,color:#14532d;
  classDef current fill:#fef3c7,stroke:#d97706,color:#78350f,stroke-width:3px;
  classDef pending fill:#f8fafc,stroke:#94a3b8,color:#334155;
  classDef blocked fill:#fee2e2,stroke:#dc2626,color:#7f1d1d;

  class S done;
  class R current;
  class D,G,F,U pending;
```

## 📂 Candidate Artifacts

| Type | Artifact | Status |
| --- | --- | --- |
| Domain | `[DOMAIN-XXX]` | `[status]` |
| Goal | `[GOAL-XXX]` | `[status]` |
| Feature | `[FT-XXX]` | `[status]` |
| Use case | `[UC-XXX]` | `[status]` |

## ⚠️ Risks

| Risk | Impact | Mitigation |
| --- | --- | --- |
| `[risk]` | `[impact]` | `[mitigation]` |

## 🏁 Approval

| Field | Value |
| --- | --- |
| Approved by |  |
| Date |  |
| Notes |  |
