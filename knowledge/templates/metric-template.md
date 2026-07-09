# Metric: [metric name]

## 🧭 Snapshot

| Field | Value |
| --- | --- |
| ID | `[MET-XXX]` |
| Status | `[draft | proposed | approved]` |
| Source artifact | `[STRAT/GOAL/FT/UC id or path]` |
| Owner | `[role/person]` |

## 📊 Definition

[Precise metric definition.]

## 🧮 Formula

| Component | Definition |
| --- | --- |
| Numerator | `[definition]` |
| Denominator | `[definition]` |
| Window | `[time window]` |
| Filters/exclusions | `[filters]` |

## 🗺️ Measurement Flow

```mermaid
flowchart LR
  A["User/System Action"] --> B["Event or Log"]
  B --> C["Metric Computation"]
  C --> D["Dashboard/Review"]
  D --> E["Decision"]

  classDef done fill:#dcfce7,stroke:#16a34a,color:#14532d;
  classDef current fill:#fef3c7,stroke:#d97706,color:#78350f,stroke-width:3px;
  classDef pending fill:#f8fafc,stroke:#94a3b8,color:#334155;
  classDef blocked fill:#fee2e2,stroke:#dc2626,color:#7f1d1d;

  class A done;
  class B,C current;
  class D,E pending;
```

## 🎯 Why It Matters

[Product, quality, safety, or operational reason.]

## 🧯 Guardrails

| Guardrail | Why Needed |
| --- | --- |
| `[guardrail]` | `[reason]` |

## 🔌 Instrumentation

| Source | Name | Notes |
| --- | --- | --- |
| Event | `[event]` | `[notes]` |
| Log | `[log]` | `[notes]` |
| Data source | `[source]` | `[notes]` |

## 🏁 Approval

| Field | Value |
| --- | --- |
| Approved by |  |
| Date |  |
| Notes |  |
