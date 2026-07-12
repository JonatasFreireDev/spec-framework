# Engineering Proposal: Organizer QR Validation

## Snapshot

| Field | Value |
| --- | --- |
| ID | `ENGPROP-002` |
| Status | `draft` |
| Source specification | `SPEC-002` |
| Source design | `DES-002` |
| Source discovery | `TD-002` |
| Engineering System | `ENGSYS-EVENTS-001 @ 0.1.0` |
| Owner skill | `engineering-proposal` |
| Next skill | `engineering-review` |

## Delivery

| Field | Value |
| --- | --- |
| Level | `L1` |
| Priority | `P0` |
| Rationale | Inherited from the organizer validation walking-skeleton delivery. |

## Technical Outcome

Propose a server-authoritative validation boundary that verifies token integrity, expiry, event scope, organizer permission, and attendance idempotency before returning a scanner result. This fixture has no application code or approved organizer permission model.

## Proposed Boundaries

| Area | Intended change | Decisions | Evidence gap |
| --- | --- | --- | --- |
| Organizer client | Submit captured token and render explicit result states | N/A | No client exists |
| QR validation service | Validate token, event, expiry, replay, and authorization | DEC-001, DEC-002 | No service exists |
| Event attendance | Persist idempotent check-in and audit outcome | DEC-002 | No schema exists |
| Authorization | Allow only approved organizer roles | Decision required | Roles remain open |

## Quality And Operations

- Negative tests must cover wrong event, expired or invalid token, duplicate scan, and permission denial.
- Observability must distinguish operational failures without logging secrets or QR payloads.
- Offline validation and manual camera fallback remain outside the approved technical contract.

## Blockers

- Organizer permission roles require a product decision and human approval.
- Real module, deployment, test, and gate evidence cannot be inspected in this documentation fixture.

## Handoff

Next: `engineering-review` for an explicit blocked verdict.
