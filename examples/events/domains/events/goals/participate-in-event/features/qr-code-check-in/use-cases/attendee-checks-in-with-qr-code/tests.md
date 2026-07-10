# Tests: Attendee Checks In With QR Code

## Context

- Source specification: SPEC-001
- Source tasks: TK-001..TK-006
- Status: draft

## Test Matrix

| ID | Scenario | Type | Source acceptance criteria | Expected result |
| --- | --- | --- | --- | --- |
| T-001 | Attendee generates QR for joined event | integration | AC1 | QR token returned without PII. |
| T-002 | Non-attendee tries to generate QR | security | AC1 | Request rejected. |
| T-003 | Organizer validates valid QR | integration | AC2 | Attendance checked_in_at set. |
| T-004 | Duplicate scan | integration | AC4 | Already checked in returned; timestamp unchanged. |
| T-005 | Expired token | integration | AC3 | Validation rejected with expired_token. |
| T-006 | Organizer without permission validates | security | AC3 | Validation rejected and logged. |
| T-007 | Analytics emitted | integration/manual | AC5 | Expected events are emitted. |

## Required Coverage

- Happy path: attendee generates QR and organizer validates.
- Permission denied: unmanaged organizer cannot validate.
- Validation errors: invalid, expired, wrong event.
- Edge cases: duplicate scan and token replay.
- Analytics/observability: success and failure events.
- Regression: existing event participation remains unchanged.

## Test Data

- Event with organizer.
- Attendee who joined event.
- User who did not join event.
- Organizer without permission.
- Expired token fixture.

## Commands Or Manual Steps

- To be filled after implementation stack is confirmed.

## Risks Not Covered

- Offline scanning is out of scope for v1.