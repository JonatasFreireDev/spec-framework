# Vision: [product or area name]

## 🧭 Snapshot

| Field | Value |
| --- | --- |
| ID | `[VIS-XXX]` |
| Type | `vision` |
| Parent IDs | `[PROB-XXX]` |
| Status | `[draft | proposed | approved]` |
| Source problem | `[PROB-XXX/path]` |
| Owner skill | Vision AI |
| Next skill | Strategy AI |

## 🌟 Vision Statement

[Describe the product future, for whom, why now, and what durable outcome it should create.]

## 👥 Target Users

| User | Desired Outcome | Current Friction |
| --- | --- | --- |
| `[user segment]` | `[outcome]` | `[friction]` |

## 🧭 Companion Contracts

| Contract | Canonical Artifact |
| --- | --- |
| Product principles, trade-offs, examples, and anti-principles | `principles.md` |
| North-star outcome, metric, measurement notes, and guardrails | `north-star.md` |

## 🗺️ Vision To Strategy Flow

```mermaid
flowchart LR
  P["Approved Problem"] --> V["Vision"]
  V --> PR["principles.md"]
  V --> NS["north-star.md"]
  PR --> S["Strategy"]
  NS --> S

  classDef done fill:#dcfce7,stroke:#16a34a,color:#14532d;
  classDef current fill:#fef3c7,stroke:#d97706,color:#78350f,stroke-width:3px;
  classDef pending fill:#f8fafc,stroke:#94a3b8,color:#334155;
  classDef blocked fill:#fee2e2,stroke:#dc2626,color:#7f1d1d;

  class P done;
  class V current;
  class PR,NS,S pending;
```

## 🚫 Non-Goals

- [What this vision does not include.]

## 🔐 Decisions Needed

| Decision | Blocks | Owner |
| --- | --- | --- |
| `[decision]` | `[artifact]` | `[role]` |

## 🏁 Approval

| Field | Value |
| --- | --- |
| Approved by |  |
| Date |  |
| Notes |  |
