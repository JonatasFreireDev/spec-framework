# Product Contract: Attendee Checks In With QR Code

| Field | Value |
| --- | --- |
| Status | draft |
| Source specification | [SPEC-001](../specification.md) |
| Contract version | 2 |

## Outcome And Actors

| Actor | Evidenced need | Observable outcome | Success signal | Source |
| --- | --- | --- | --- | --- |
| Attendee | Prove presence without exposing personal data in the QR payload. | An authorized scan marks the attendee present once. | Check-in succeeds or returns the existing state. | [Use Case](../use-case.md) |
| Organizer | Validate proof for the correct event. | The scanner returns a safe, actionable result. | Validation result matches permission and token state. | [Use Case](../use-case.md) |

## Scope And Non-Goals

| In scope | Non-goal | Boundary rationale | Source |
| --- | --- | --- | --- |
| Generate, refresh, present, and validate event-scoped attendee proof. | Offline authority, payments, ticket transfer, and advanced fraud scoring. | L1 proves attendance with the smallest server-authoritative flow. | [Feature](../../../feature.md) |

## Product Rules

| Rule | Actor impact | Failure consequence | Source decision/evidence |
| --- | --- | --- | --- |
| Only joined attendees receive proof; only authorized organizers validate it. | Prevents unrelated users from creating or consuming attendance proof. | Request is denied without mutating attendance or revealing identity. | [Use Case](../use-case.md) |
| Proof expires after five minutes and may be refreshed. | Limits replay while retaining venue usability. | Expired proof is rejected and attendee receives a refresh path. | [DEC-001](../../../../../../../../../knowledge/decisions/DEC-001-qr-expiration-duration.md) |

## Requirements

| ID | Requirement | Source | Acceptance criteria | Dependencies |
| --- | --- | --- | --- | --- |
| REQ-001 | A joined attendee can obtain non-PII proof whose authorized validation records one event check-in. | [Use Case](../use-case.md) | AC-001 | [DEC-001](../../../../../../../../../knowledge/decisions/DEC-001-qr-expiration-duration.md), [DEC-002](../../../../../../../../../knowledge/decisions/DEC-002-qr-token-strategy.md) |
