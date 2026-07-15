# Metric: [metric name]

## 🧾 Generation And Agent Self-Check

> Complete this section when materializing the artifact. Keep unresolved items explicit in the relevant scope, findings, risks, or handoff section.

| Field | Value |
| --- | --- |
| Generated on | `YYYY-MM-DD` |
| Purpose | `[decision, evidence, contract, or handoff this artifact supports]` |
| Use when | `[workflow stage, trigger, or condition]` |
| Prepared by | `[owning skill, role, or accountable person]` |
| Scope covered | `[artifact, product area, use case, or review boundary]` |
| Required inputs and evidence | `[links to approved parents, documents, code, decisions, or observations]` |
| Ready when | `[artifact-specific completion, evidence, and gate conditions]` |
| Current status | `[status allowed by this artifact's owning workflow]` |


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

## ✅ Agent Verification Checklist

- [ ] The metric has a precise definition, formula, unit, population, window, and owner.
- [ ] Source events, data quality, segmentation, cadence, and instrumentation are traceable.
- [ ] Target, baseline, interpretation, and guardrails support a real product decision.
- [ ] Approval and change history do not imply unavailable data or fabricated measurement.
