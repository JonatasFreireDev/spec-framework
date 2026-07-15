# Engineering Proposal: [use case]

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

## ✅ Agent Verification Checklist

- [ ] The proposal traces to approved specification, design, discovery, decisions, and system version.
- [ ] Module, data, integration, dependency, and ownership boundaries are explicit.
- [ ] Quality, security, operations, testing, rollout, rollback, and observability are implementable.
- [ ] Deviations and unresolved architecture decisions block handoff when required.
