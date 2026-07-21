# API Contract: Organizer Validates QR Code

| Field | Value |
| --- | --- |
| Status | draft |
| Source specification | [SPEC-002](../specification.md) |
| Contract version | 2 |

## Operations

| Operation | Kind/path | Consumer | Side effect | Owner |
| --- | --- | --- | --- | --- |
| getEventCheckInSession | query | Organizer scanner. | None. | Events service. |
| validateEventCheckIn | command | Organizer scanner. | Records attendance and audit on authorized success. | Events service. |

## Authorization

| Operation | Authentication | Permission/tenant rule | Enforcement point | Denial response |
| --- | --- | --- | --- | --- |
| getEventCheckInSession | Organizer session. | Caller has current check-in permission for eventId. | Server before scanner session data. | permission_denied. |
| validateEventCheckIn | Organizer session. | Permission is rechecked for every attempt; eventId must match proof. | Server before identity response and write. | permission_denied or wrong_event without attendee details. |

## Request And Response Schemas

| Operation | Request fields and constraints | Success response | Sensitive fields |
| --- | --- | --- | --- |
| getEventCheckInSession | eventId required. | eventDisplayName, permissionStatus, windowStatus, scanner guidance. | No attendee data. |
| validateEventCheckIn | eventId and qrToken required; scanDeviceId/clientTimestamp optional telemetry. | status, checkInId/checkedInAt where applicable, retryable; attendeeDisplayName only for authorized success/existing. | Raw token and attendee identity are restricted. |

## Errors Idempotency And Compatibility

| Operation | Error/status | Retryable | Idempotency/concurrency | Versioning compatibility |
| --- | --- | --- | --- | --- |
| getEventCheckInSession | authentication_required, permission_denied, outside_window, unavailable. | Only unavailable. | Query has no side effect. | New guidance fields remain optional. |
| validateEventCheckIn | checked_in, already_checked_in, expired, invalid, wrong_event, permission_denied, outside_window, unavailable. | Only unavailable. | Event-attendee uniqueness and transaction return one stable result. | Rejection responses never gain identity fields. |

## Requirements

| ID | Requirement | Source | Acceptance criteria | Dependencies |
| --- | --- | --- | --- | --- |
| REQ-104 | Session and validation operations must reauthorize each attempt, constrain schemas and sensitive fields, enumerate safe outcomes, and remain idempotent under retry/concurrency. | [Use Case](../use-case.md) | AC-104 | REQ-102, REQ-103 |
