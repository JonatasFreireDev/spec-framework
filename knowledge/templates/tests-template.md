# Tests: [use case name]

## 🧭 Snapshot

| Field | Value |
| --- | --- |
| ID | `[TEST-XXX]` |
| Status | `[draft | proposed | approved]` |
| Source specification | `[SPEC-XXX]` |
| Owner skill | QA AI |
| Next skill | Audit Orchestrator or Release Orchestrator |

## 🎯 Test Goal

[Describe what confidence this test plan must provide.]

## 🧪 Coverage Matrix

| Area | Required Coverage | Status |
| --- | --- | --- |
| Behavioral | `[main and alternate flows]` | `[draft/proposed/approved]` |
| Permissions/security | `[checks]` | `[status]` |
| Data | `[constraints and mutations]` | `[status]` |
| UX states | `[states]` | `[status]` |
| Accessibility | `[requirements]` | `[status]` |
| Analytics/observability | `[events/logs/metrics]` | `[status]` |
| Performance/reliability | `[expectations]` | `[status]` |

## 🗺️ Test Flow

```mermaid
flowchart LR
  A["Fixtures"] --> B["Behavior tests"]
  B --> C["Permission tests"]
  C --> D["UX and accessibility checks"]
  D --> E["Analytics assertions"]
  E --> F["QA verdict"]

  classDef done fill:#dcfce7,stroke:#16a34a,color:#14532d;
  classDef current fill:#fef3c7,stroke:#d97706,color:#78350f,stroke-width:3px;
  classDef pending fill:#f8fafc,stroke:#94a3b8,color:#334155;
  classDef blocked fill:#fee2e2,stroke:#dc2626,color:#7f1d1d;

  class A done;
  class B current;
  class C,D,E,F pending;
```

## ✅ Test Cases

| Test | Preconditions | Steps | Expected Result |
| --- | --- | --- | --- |
| `[test name]` | `[preconditions]` | `[steps]` | `[result]` |

## ⚠️ Residual Risk

| Risk | Why It Remains | Mitigation |
| --- | --- | --- |
| `[risk]` | `[reason]` | `[mitigation]` |

## 🏁 QA Result

| Field | Value |
| --- | --- |
| Verdict | `[passed | passed_with_notes | blocked]` |
| Required fixes | `[fixes]` |
| Next owner | `[role/skill]` |
