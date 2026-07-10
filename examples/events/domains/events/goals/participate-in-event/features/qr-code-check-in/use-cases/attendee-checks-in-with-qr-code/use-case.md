# Use Case: Attendee Checks In With QR Code

## Context

- ID: UC-001
- Status: draft
- Feature: FT-001
- Context file: context.md

## Actor

Primary actor: attendee.
Secondary actor: organizer.

## Goal

The attendee is marked as present at the event after an organizer validates their QR code.

## Preconditions

- Attendee is authenticated.
- Attendee has joined or RSVP'd to the event.
- Organizer is authenticated and has permission to manage check-in for the event.
- Event exists and is within the check-in window.

## Main Flow

1. Attendee opens the event check-in screen.
2. System generates or retrieves a valid QR token for that attendee and event.
3. Attendee presents the QR code to the organizer.
4. Organizer scans the QR code.
5. System validates token, event, attendee, expiration, and organizer permission.
6. System marks attendance as checked in.
7. Attendee and organizer see success confirmation.

## Alternate Flows

### QR Expired

1. Organizer scans an expired QR code.
2. System rejects validation and shows an expired state.
3. Attendee can refresh the QR code.

### Already Checked In

1. Organizer scans a valid QR for an attendee already checked in.
2. System returns an idempotent already-checked-in state.

## Error And Edge Cases

- Invalid token - show invalid QR error and do not reveal attendee details.
- Wrong event - reject validation.
- Organizer lacks permission - reject validation and log permission failure.
- Network failure - show retryable error; do not mark attendance locally.

## Business Rules

- QR token must be event-scoped and attendee-scoped.
- QR token must expire.
- Check-in mutation must be idempotent.
- Organizer must have explicit permission for the event.

## UX States

- Default: attendee sees QR code and expiration hint.
- Loading: QR generation or validation in progress.
- Empty: attendee has not joined event.
- Error: invalid, expired, permission denied, network failure.
- Success: attendance marked as checked in.

## Acceptance Criteria

- [ ] Attendee can generate a QR code for an event they joined.
- [ ] Organizer can validate a valid QR code for an event they manage.
- [ ] Expired or invalid QR codes are rejected.
- [ ] Duplicate scans do not create duplicate side effects.
- [ ] Validation emits analytics and audit logs.

## Specification Readiness

- [x] Scope is clear.
- [x] Business rules are linked.
- [x] UI states are known.
- [x] Data and permission needs are known.
- [x] Analytics and observability needs are known.