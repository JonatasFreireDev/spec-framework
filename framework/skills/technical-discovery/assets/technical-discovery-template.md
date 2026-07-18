# Technical Discovery: [use case]

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
| ID | `TD-XXX` |
| Status | `draft` |
| Source specification | `[SPEC-*]` |
| Source design | `[DES-* or Not applicable]` |
| Owner skill | `technical-discovery` |

## Requirement-to-codebase map

| Requirement | Existing module/API/data owner | Tests/convention | Expected change | Risk |
| --- | --- | --- | --- | --- |
| `REQ-001` | `[real path]` | `[path]` | `[change]` | `[risk]` |

## Change surface

| Area | Existing evidence | Probable paths | Shared resources |
| --- | --- | --- | --- |
| `[area]` | `[path/command]` | `[paths]` | `[resource or none]` |

## Architecture Gate

| Field | Value |
| --- | --- |
| Verdict | `Decision required | Not required` |
| Decision | `[DEC-* or N/A]` |
| Rationale | `[concrete rationale]` |

A referenced DEC must be indexed, approved, covered by a current approval record, and scope-compatible. `Not required` must name the existing pattern or decision that already covers the change.

## Planning blockers

- `[blocker or None]`

## Handoff

Next: `engineering-proposal` after the Architecture Gate is resolved.

## ✅ Agent Verification Checklist

- [ ] Every relevant requirement maps to existing code, tests, conventions, owners, or an evidenced gap.
- [ ] Change surfaces identify modules, APIs, data, migrations, integrations, operations, and risks.
- [ ] The Architecture Gate lists applicable decisions and unresolved questions with owners.
- [ ] Planning is blocked when evidence or architecture decisions are insufficient.
