# Context: Attendee Checks In With QR Code

```yaml
id: UC-001
type: use-case
name: Attendee checks in with QR code
status: draft
owner_skill: 08-use-case.md
slug: attendee-checks-in-with-qr-code
rigor_tier: L
last_updated: 2026-07-09
delivery:
  level: L1
  priority: P0
  depends_on:
    - FT-001
    - DEC-001
    - DEC-002
  rationale: Attendee QR presentation is required to close the walking skeleton with organizer validation.
```

## Purpose

Defines the concrete interaction where an attendee presents a QR code and an organizer validates it to mark attendance.

## Parent Artifacts

- FT-001 - ../../context.md - parent feature

## Child Artifacts

- SPEC-001 - specification.md - implementation contract
- UC-001:design - design.md - UX flow and states
- PLAN-001 - implementation-plan.md - build strategy
- GRAPH-001 - execution-graph.json - task DAG
- UC-001:tasks - tasks.md - executable task set
- UC-001:tests - tests.md - validation plan
- QA-001 - qa-evidence.md - validation evidence
- SEC-001 - security-review.md - security review
- UC-001:analytics - analytics.md - measurement plan
- UC-001:audit - audit.md - audit evidence
- TK-001..TK-006 - tasks.md - executable work

## Dependencies

- Authenticated attendee identity - blocking: yes
- Organizer permission model - blocking: yes
- Event attendance persistence - blocking: yes

## Rigor Tier

| Field | Value |
| --- | --- |
| Tier | L |
| Trigger checklist | auth, permissions, token privacy, database migration |
| Human approval | Approved by EV-003 policy rollout |
| Rationale | The flow authenticates attendees, validates organizer permissions, handles QR tokens, and writes attendance state. |

## Related Artifacts

- ../../feature.md - feature scope

## Canonical Documents

- Primary: use-case.md
- Specification: specification.md
- Implementation plan: implementation-plan.md
- Execution graph: execution-graph.json
- Tasks: tasks.md
- QA evidence: qa-evidence.md
- Security review: security-review.md

## Decisions

- DEC-001 - QR expiration duration - approved.
- DEC-002 - QR token strategy - approved.
- DEC-008 - Rigor tiers for use cases - approved.

## Assumptions

- QR token is opaque and server-validated.
- QR token expires after a short time window.

## Open Questions

- Exact expiration window is approved by DEC-001.
- Offline validation is out of scope for v1 but may become future work.

## Handoff

Next recommended skill: 09-specification.md
Required reading before next step:
- use-case.md
- specification.md
