# API Contract: Attendee Checks In With QR Code

| Field | Value |
| --- | --- |
| Status | draft |
| Source specification | [SPEC-001](../specification.md) |
| Contract version | 2 |

## Operations

| Operation | Kind/path | Consumer | Side effect | Owner |
| --- | --- | --- | --- | --- |
| generateEventCheckInQr | command | Attendee client | Creates or rotates an opaque token. | Events service. |
| validateEventCheckInQr | command | Organizer scanner | Records check-in once when authorized. | Events service. |
| getEventCheckInStatus | query | Attendee client | None. | Events service. |

## Authorization

| Operation | Authentication | Permission/tenant rule | Enforcement point | Denial response |
| --- | --- | --- | --- | --- |
| generateEventCheckInQr | Attendee session. | Caller joined the requested event. | Server before token creation. | not_joined or event_not_found. |
| validateEventCheckInQr | Organizer session. | Caller can manage check-in for eventId. | Server before token lookup response or write. | permission_denied without attendee data. |

## Request And Response Schemas

| Operation | Request fields and constraints | Success response | Sensitive fields |
| --- | --- | --- | --- |
| generateEventCheckInQr | eventId required. | token, expiresAt, qrPayload. | Token is bearer proof; payload excludes PII. |
| validateEventCheckInQr | eventId and token required. | status and checkedInAt; attendee display only on authorized success. | Identity is conditional and never returned on rejection. |
| getEventCheckInStatus | eventId required. | joined, checkedIn, checkedInAt. | Caller may read only own status. |

## Errors Idempotency And Compatibility

| Operation | Error/status | Retryable | Idempotency/concurrency | Versioning compatibility |
| --- | --- | --- | --- | --- |
| generateEventCheckInQr | not_joined, check_in_not_open, unavailable. | Only unavailable. | Rotation invalidates prior active token according to service policy. | Adding statuses must preserve existing safe failure handling. |
| validateEventCheckInQr | expired, invalid, wrong_event, permission_denied, already_checked_in, unavailable. | Only unavailable. | Event-attendee key makes concurrent calls idempotent. | Success fields remain optional for rejection responses. |

## Requirements

| ID | Requirement | Source | Acceptance criteria | Dependencies |
| --- | --- | --- | --- | --- |
| REQ-004 | Generation, status, and validation operations must enforce caller scope, explicit schemas, safe errors, and idempotent concurrency semantics. | [Use Case](../use-case.md) | AC-004 | REQ-002, REQ-003 |
