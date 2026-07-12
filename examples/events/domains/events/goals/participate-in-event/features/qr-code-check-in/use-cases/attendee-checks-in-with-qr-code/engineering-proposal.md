# Engineering Proposal: Attendee QR Check-in

## Snapshot

| Field | Value |
| --- | --- |
| ID | `ENGPROP-001` |
| Status | `draft` |
| Source specification | `SPEC-001` |
| Source design | `DES-001` |
| Source discovery | `TD-001` |
| Engineering System | `ENGSYS-EVENTS-001 @ 0.1.0` |
| Owner skill | `engineering-proposal` |
| Next skill | `engineering-review` |

## Delivery

| Field | Value |
| --- | --- |
| Level | `L1` |
| Priority | `P0` |
| Rationale | Inherited from the attendee QR walking-skeleton delivery. |

## Technical Outcome

Propose server-authoritative QR token generation and attendance recording while preserving five-minute expiry, idempotency, permission checks, and non-disclosure of PII. This fixture has no application code, so proposed module names and interfaces are not implementation evidence.

## Proposed Boundaries

| Area | Intended change | Decisions | Evidence gap |
| --- | --- | --- | --- |
| QR token service | Generate short-lived attendee tokens on the server | DEC-001, DEC-002 | No server module exists |
| Event attendance | Record one attendance result per attendee and event | DEC-002 | No schema exists |
| Attendee experience | Present active, refreshing, expired, loading, and error states | N/A | No client exists |
| Observability | Record generation and validation outcomes without token or PII leakage | DEC-002 | No runtime exists |

## Quality And Operations

- Integration and security tests must cover expiration, tampering, replay, permissions, and idempotency.
- Rollout should use a product-controlled feature flag; its concrete owner and gate remain unknown.
- Rollback must disable generation without deleting recorded attendance.

## Blockers

- Real module, data, deployment, test, and gate evidence cannot be inspected in this documentation fixture.
- The proposal cannot advance to an implementation-ready verdict in this repository.

## Handoff

Next: `engineering-review` for an explicit blocked verdict over the documentation-only proposal.
