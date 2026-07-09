# Decision: [decision title]

## 🧭 Snapshot

| Field | Value |
| --- | --- |
| ID | `[DEC-XXX]` |
| Status | `[proposed | approved | superseded | rejected]` |
| Date | `[YYYY-MM-DD]` |
| Scope | `[product/security/architecture/data/UX/release]` |
| Owner | `[role/person]` |

## ✅ Decision

[State the decision clearly.]

## 🧠 Why

[Explain the product, technical, security, or operational reason.]

## ⚖️ Options Considered

| Option | Pros | Cons | Result |
| --- | --- | --- | --- |
| `[option]` | `[pros]` | `[cons]` | `[chosen/rejected]` |

## 🗺️ Decision Impact Flow

```mermaid
flowchart LR
  D["Decision"] --> S["Specifications"]
  D --> P["Implementation Plans"]
  D --> G["Execution Graphs"]
  D --> T["Tasks"]
  D --> R["Release Readiness"]

  classDef done fill:#dcfce7,stroke:#16a34a,color:#14532d;
  classDef current fill:#fef3c7,stroke:#d97706,color:#78350f,stroke-width:3px;
  classDef pending fill:#f8fafc,stroke:#94a3b8,color:#334155;
  classDef blocked fill:#fee2e2,stroke:#dc2626,color:#7f1d1d;

  class D current;
  class S,P,G,T,R pending;
```

## 📌 Consequences

| Type | Consequence | Follow-up |
| --- | --- | --- |
| Positive | `[benefit]` | `[action]` |
| Negative | `[cost/risk]` | `[action]` |

## 📂 Affected Artifacts

| Artifact | Required Update |
| --- | --- |
| `[path/id]` | `[update]` |

## 🔁 Supersedes

- `[DEC-XXX or N/A]`

## 🏁 Approval

| Field | Value |
| --- | --- |
| Approved by |  |
| Date |  |
| Approval record | `[.product/history/approval-...]` |
| Notes |  |
