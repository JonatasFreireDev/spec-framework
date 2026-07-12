# Engineering Proposal: [use case]

## Snapshot

| Field | Value |
| --- | --- |
| ID | `ENGPROP-XXX` |
| Status | `draft` |
| Source specification | `[SPEC-*]` |
| Source design | `[DES-* or Not applicable]` |
| Source discovery | `[TD-*]` |
| Engineering System | `[ENGSYS-* @ version or Not configured]` |
| Owner skill | `engineering-proposal` |
| Next skill | `engineering-review` |

## Delivery

| Field | Value |
| --- | --- |
| Level | `[L0 | L1 | L2 | L3 | L4 | L5]` |
| Priority | `[P0 | P1 | P2 | P3]` |
| Rationale | `[inherited delivery rationale]` |

## Technical Outcome

[Describe the intended system change without sequencing implementation tasks or writing code.]

## Existing Evidence

| Area | Current contract or code | Constraint |
| --- | --- | --- |
| `[area]` | `[path]` | `[constraint]` |

## Proposed Change

| Area | Intended change | Requirements | Decisions |
| --- | --- | --- | --- |
| `[module/data/API/integration]` | `[change]` | `[REQ-*]` | `[DEC-* or N/A]` |

## Boundaries And Ownership

| Resource | Owner before | Owner after | Interface |
| --- | --- | --- | --- |
| `[resource]` | `[module/system]` | `[module/system]` | `[API/event/data contract]` |

## Quality And Operations

| Attribute | Expected behavior | Verification or evidence |
| --- | --- | --- |
| `[reliability/security/observability/performance]` | `[expectation]` | `[test/gate/runbook]` |

## Rollout And Recovery

| Concern | Proposal |
| --- | --- |
| Compatibility | `[contract]` |
| Migration | `[contract or Not applicable]` |
| Rollout | `[contract]` |
| Rollback | `[contract]` |

## Deviations And Open Decisions

| Item | Impact | Required owner or decision |
| --- | --- | --- |
| `[deviation/question or None]` | `[impact]` | `[DEC-*/human/none]` |

## Handoff

Next: `engineering-review`.
