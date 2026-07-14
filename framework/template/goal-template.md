# User Goal: [goal name]

## 🧭 Snapshot

| Field | Value |
| --- | --- |
| ID | `[GOAL-XXX]` |
| Status | `[draft | proposed | approved]` |
| Domain | [`DOMAIN-XXX`](<path-to-domain.md>#domain-xxx) |
| Owner skill | User Goal AI |
| Next skill | Journey AI or Feature AI |

## 🎯 User Intent

As a `[user]`, I want to `[goal]`, so I can `[outcome]`.

## 💡 Why This Goal Matters

[Explain the product value and how it traces to strategy.]

## 🗺️ Journey Summary

```mermaid
flowchart LR
  A["Trigger"] --> B["User tries to progress"]
  B --> C["System supports action"]
  C --> D["Success moment"]
  B --> E["Failure or recovery"]

  classDef done fill:#dcfce7,stroke:#16a34a,color:#14532d;
  classDef current fill:#fef3c7,stroke:#d97706,color:#78350f,stroke-width:3px;
  classDef pending fill:#f8fafc,stroke:#94a3b8,color:#334155;
  classDef blocked fill:#fee2e2,stroke:#dc2626,color:#7f1d1d;

  class A done;
  class B,C current;
  class D pending;
  class E blocked;
```

## 🧱 Candidate Features

| Feature | Status | Delivery | Priority | Notes |
| --- | --- | --- | --- | --- |
| [`FT-XXX`](<path-to-feature.md>#ft-xxx) `[name]` | `[status]` | `[L0-L5]` | `[P0-P3]` | `[notes]` |

## 📏 Rules And Constraints

| Rule/Constraint | Source | Impact |
| --- | --- | --- |
| `[rule]` | `[decision/path]` | `[impact]` |

## 📊 Metrics

| Metric | Meaning |
| --- | --- |
| `[metric]` | `[meaning]` |

## ⚠️ Risks And Open Questions

| Item | Blocks | Owner |
| --- | --- | --- |
| `[risk/question]` | `[artifact]` | `[role]` |

## 🏁 Approval

| Field | Value |
| --- | --- |
| Approved by |  |
| Date |  |
| Notes |  |
