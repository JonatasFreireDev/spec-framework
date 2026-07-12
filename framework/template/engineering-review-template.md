# Engineering Review: [use case]

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
