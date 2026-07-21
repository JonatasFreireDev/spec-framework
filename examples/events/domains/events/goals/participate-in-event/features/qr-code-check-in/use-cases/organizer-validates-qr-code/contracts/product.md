# Product Contract: Organizer Validates QR Code

| Field | Value |
| --- | --- |
| Status | draft |
| Source specification | [SPEC-002](../specification.md) |
| Contract version | 2 |

## Outcome And Actors

| Actor | Evidenced need | Observable outcome | Success signal | Source |
| --- | --- | --- | --- | --- |
| Organizer | Process attendee proof quickly for the correct managed event. | Scanner produces a trustworthy result and next-scan action. | Exactly one server-authoritative check-in or explicit rejection. | [Use Case](../use-case.md) |
| Attendee | Receive correct attendance confirmation without identity exposure on invalid proof. | Valid proof is acknowledged; invalid proof reveals no unnecessary data. | Attendance state matches authorized validation. | [Use Case](../use-case.md) |

## Scope And Non-Goals

| In scope | Non-goal | Boundary rationale | Source |
| --- | --- | --- | --- |
| Scanner session, online token validation, permission check, result feedback, audit, and next scan. | Offline authority, role-management UI, payments, manual search, and attendee proof generation. | This use case owns organizer consumption of proof; UC-001 owns generation. | [Use Case](../use-case.md) |

## Product Rules

| Rule | Actor impact | Failure consequence | Source decision/evidence |
| --- | --- | --- | --- |
| Every validation checks current organizer permission and event window server-side. | An authenticated but unauthorized user cannot mutate attendance. | Denied result, no identity disclosure, safe audit signal. | [Use Case](../use-case.md) |
| Opaque proof expires after five minutes and resolves only on the server. | Organizer cannot trust client decoding or stale proof. | Expired/invalid result with attendee refresh guidance. | [DEC-001](../../../../../../../../../knowledge/decisions/DEC-001-qr-expiration-duration.md), [DEC-002](../../../../../../../../../knowledge/decisions/DEC-002-qr-token-strategy.md) |

## Requirements

| ID | Requirement | Source | Acceptance criteria | Dependencies |
| --- | --- | --- | --- | --- |
| REQ-101 | An authorized organizer can validate attendee proof for one managed event and receive a safe result that records at most one check-in. | [Use Case](../use-case.md) | AC-101 | [UC-001](../../attendee-checks-in-with-qr-code/use-case.md) |
