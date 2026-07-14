# Threat Register

## Snapshot

| Field | Value |
| --- | --- |
| ID | `[THREAT-REGISTER-XXX]` |
| Status | `[draft | proposed | approved]` |
| Owner skill | Threat Modeler AI |
| Governed by | `FRAMEWORK.md security policy` |
| Scope | `[product/domain/release]` |
| Last reviewed | `YYYY-MM-DD` |

## Threat Flow

```mermaid
flowchart LR
  B["Security baseline"] --> T["Threat scenario"]
  T --> M["Mitigation"]
  M --> E["Evidence"]
  E --> R["Security Review"]
  R --> A["Audit / Release"]

  classDef done fill:#dcfce7,stroke:#16a34a,color:#14532d;
  classDef current fill:#fef3c7,stroke:#d97706,color:#78350f,stroke-width:3px;
  classDef pending fill:#f8fafc,stroke:#94a3b8,color:#334155;
  classDef blocked fill:#fee2e2,stroke:#dc2626,color:#7f1d1d;

  class B,T current;
  class M,E,R,A pending;
```

## Register

| ID | Status | Severity | Likelihood | Scenario | Affected Artifacts | Required Mitigation | Evidence | Route | Owner |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| `[THR-001]` | `[open/mitigated/accepted/blocked]` | `[blocker/high/medium/low]` | `[high/medium/low]` | `[attack or failure scenario]` | `[paths]` | `[control]` | `[path/log/test/decision]` | `[code-runner/qa/bug-fixer/product-historian/security-review]` | `[owner]` |

## Accepted Residual Risks

| Threat ID | Risk | Approval | Compensating Control | Review Date |
| --- | --- | --- | --- | --- |
| `[THR-XXX]` | `[risk]` | `[DEC/approval/PR]` | `[control]` | `YYYY-MM-DD` |

## Review Notes

| Date | Reviewer | Change | Follow-up |
| --- | --- | --- | --- |
| `YYYY-MM-DD` | `[name/skill]` | `[what changed]` | `[next step]` |
