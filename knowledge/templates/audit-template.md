# Audit: [scope]

## 🧭 Executive Snapshot

| Field | Value |
| --- | --- |
| Scope | `[domain/goal/feature/use-case/release]` |
| Auditor skill | `[skill]` |
| Date | `[YYYY-MM-DD]` |
| Verdict | `[✅ approved | 🟡 approved_with_notes | 🔴 blocked]` |
| Next owner | `[skill/person]` |

## 🗺️ Audit Flow

```mermaid
flowchart LR
  A["Collect Context"] --> B["Check Traceability"]
  B --> C["Find Gaps"]
  C --> D["Find Conflicts"]
  D --> E["Check Dependencies"]
  E --> F["Assess Impact"]
  F --> G["Verdict"]

  G -->|approved| H["Proceed"]
  G -->|approved_with_notes| I["Proceed With Follow-ups"]
  G -->|blocked| J["Fix Required"]

  classDef done fill:#dcfce7,stroke:#16a34a,color:#14532d;
  classDef current fill:#fef3c7,stroke:#d97706,color:#78350f,stroke-width:3px;
  classDef pending fill:#f8fafc,stroke:#94a3b8,color:#334155;
  classDef blocked fill:#fee2e2,stroke:#dc2626,color:#7f1d1d;

  class A,B,C,D,E,F done;
  class G current;
  class H,I pending;
  class J blocked;
```

## 📌 Summary

[Short summary of what was checked and result.]

## 🚦 Verdict Matrix

| Area | Result | Evidence | Notes |
| --- | --- | --- | --- |
| Traceability | `[✅/🟡/🔴]` | `[path/section]` | `[note]` |
| Completeness | `[✅/🟡/🔴]` | `[path/section]` | `[note]` |
| Consistency | `[✅/🟡/🔴]` | `[path/section]` | `[note]` |
| Dependencies | `[✅/🟡/🔴]` | `[path/section]` | `[note]` |
| Security/privacy | `[✅/🟡/🔴/➖]` | `[path/section]` | `[note]` |
| UX/accessibility | `[✅/🟡/🔴/➖]` | `[path/section]` | `[note]` |
| Release readiness | `[✅/🟡/🔴/➖]` | `[path/section]` | `[note]` |

## 🔎 Findings

| Severity | Finding | Evidence | Impact | Required Fix | Owner |
| --- | --- | --- | --- | --- | --- |
| `[🔴 blocker/🟡 warning/🔵 note]` | `[finding title]` | `[file/path/section]` | `[why it matters]` | `[fix]` | `[role]` |

## 🧱 Gaps

| Gap | Blocks | Required Fix | Owner |
| --- | --- | --- | --- |
| `[gap]` | `[artifact/task]` | `[fix]` | `[role]` |

## ⚔️ Conflicts

| Conflict | Artifacts | Impact | Resolution Needed |
| --- | --- | --- | --- |
| `[conflict]` | `[paths]` | `[impact]` | `[decision/fix]` |

## 🔗 Dependencies

| Dependency | Required By | Status | Risk |
| --- | --- | --- | --- |
| `[dependency]` | `[artifact/task]` | `[open/ready/blocked]` | `[risk]` |

## 🔐 Decisions

| Decision | Status | Blocks | Owner |
| --- | --- | --- | --- |
| `[decision]` | `[open/proposed/approved]` | `[artifact/task]` | `[role]` |

## 🌡️ Residual Risk

| Risk | Likelihood | Impact | Mitigation |
| --- | --- | --- | --- |
| `[risk]` | `[low/medium/high]` | `[low/medium/high]` | `[mitigation]` |

## 🏁 Approval

| Field | Value |
| --- | --- |
| Approved by |  |
| Date |  |
| Notes |  |
