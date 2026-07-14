# Feature: [feature name]

## 🧭 Snapshot

| Field | Value |
| --- | --- |
| ID | `[FT-XXX]` |
| Status | `[draft | proposed | approved]` |
| Domain | [`DOMAIN-XXX`](<path-to-domain.md>#domain-xxx) |
| User goal | [`GOAL-XXX`](<path-to-goal.md>#goal-xxx) |
| Owner skill | Feature AI |
| Next skill | Use Case AI |

## 🚚 Delivery

| Field | Value |
| --- | --- |
| Level | `[L0 | L1 | L2 | L3 | L4 | L5]` |
| Priority | `[P0 | P1 | P2 | P3]` |
| Depends on | `[artifact ids/paths]` |
| Rationale | `[why this belongs here]` |

## 📌 Summary

[Describe the concrete solution and how it supports the parent goal.]

## 🎯 Problem Fit

| User Problem | Business Reason | Evidence |
| --- | --- | --- |
| `[problem]` | `[reason]` | `[path/source]` |

## 🧱 Scope

| In Scope | Non-Goals |
| --- | --- |
| `[behavior]` | `[excluded behavior]` |

## Delivery Slice

| Field | Value |
| --- | --- |
| User value | `[observable value]` |
| Entry point | `[where the slice begins]` |
| End state | `[observable completion]` |
| Independently releasable | `[yes/no]` |
| Reversible | `[yes/no and how]` |
| Deferred | `[explicitly postponed behavior]` |

## 🎬 Use Cases

| Use Case | Status | Delivery | Priority | Notes |
| --- | --- | --- | --- | --- |
| [`UC-XXX`](<path-to-use-case.md>#uc-xxx) `[name]` | `[status]` | `[L0-L5]` | `[P0-P3]` | `[notes]` |

## 🗺️ Feature Flow

```mermaid
flowchart TD
  G["[GOAL-XXX] Goal"] --> F["[FT-XXX] Feature"]
  F --> U1["[UC-XXX] Use Case"]
  U1 --> S1["[SPEC-XXX] Specification"]

  classDef done fill:#dcfce7,stroke:#16a34a,color:#14532d;
  classDef current fill:#fef3c7,stroke:#d97706,color:#78350f,stroke-width:3px;
  classDef pending fill:#f8fafc,stroke:#94a3b8,color:#334155;
  classDef blocked fill:#fee2e2,stroke:#dc2626,color:#7f1d1d;

  class G done;
  class F current;
  class U1,S1 pending;
```

## 🎨 UX Notes

| Area | Notes |
| --- | --- |
| Entry points | `[entry points]` |
| Core states | `[states]` |
| Empty/loading/error states | `[states]` |
| Accessibility | `[requirements]` |

## 🔐 Data, Permissions, And Risks

| Topic | Detail |
| --- | --- |
| Data touched | `[entities/fields]` |
| Permission model | `[who can do what]` |
| Sensitive data or abuse risks | `[risk]` |

## 📊 Analytics

| Event | Meaning |
| --- | --- |
| `[event]` | `[meaning]` |

## ⚠️ Open Questions

| Question | Blocks | Owner |
| --- | --- | --- |
| `[question]` | `[artifact]` | `[role]` |

## 🏁 Approval

| Field | Value |
| --- | --- |
| Approved by |  |
| Date |  |
| Notes |  |
