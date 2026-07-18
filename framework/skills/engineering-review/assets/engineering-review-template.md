# Engineering Review: [use case]

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
| ID | `ENGREV-XXX` |
| Status | `draft` |
| Proposal | [engineering-proposal.md](engineering-proposal.md) |
| Proposal hash | `[sha256 or N/A while draft]` |
| Reviewer skill | `engineering-review` |
| Verdict | `passed | required_fix | blocked | not_reviewed` |

## Delivery

| Field | Value |
| --- | --- |
| Level | `[L0 | L1 | L2 | L3 | L4 | L5]` |
| Priority | `[P0 | P1 | P2 | P3]` |
| Rationale | `[inherited delivery rationale]` |

## Review Matrix

| Concern | Evidence | Result | Notes |
| --- | --- | --- | --- |
| Requirement coverage | `[proposal section]` | `[passed/required_fix/blocked/not_reviewed]` | `[notes]` |
| Architecture boundaries | `[system/proposal/DEC]` | `[result]` | `[notes]` |
| Data ownership | `[system/proposal/DEC]` | `[result]` | `[notes]` |
| Dependencies and integrations | `[evidence]` | `[result]` | `[notes]` |
| Security and privacy | `[evidence]` | `[result]` | `[notes]` |
| Quality and observability | `[evidence]` | `[result]` | `[notes]` |
| Migration, rollout, rollback | `[evidence]` | `[result]` | `[notes]` |
| Testability | `[evidence]` | `[result]` | `[notes]` |

## Findings

| Severity | Finding | Route | Owner |
| --- | --- | --- | --- |
| `[blocker/required_fix/note]` | `[finding or None]` | `[proposal/DEC/specification]` | `[owner]` |

## Decision Coverage

| Decision need | Applicable DEC | Approval evidence | Result |
| --- | --- | --- | --- |
| `[need or None]` | `[DEC-* or N/A]` | `[path or N/A]` | `[covered/blocked]` |

## Handoff

Next: `implementation-planner` only when the verdict is `passed` and every required decision is approved.

## ✅ Agent Verification Checklist

- [ ] The reviewed proposal and immutable hash identify the exact review target.
- [ ] Architecture, ownership, dependencies, quality attributes, operations, security, and testability are independently assessed.
- [ ] Findings include evidence, severity, owner, and required correction.
- [ ] The verdict and handoff preserve read-only review and do not grant product approval.
