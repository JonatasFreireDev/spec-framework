# Use Case: Organizer Validates QR Code

## Context

- ID: UC-002
- Status: proposed
- Domain: DOMAIN-001 Events
- Goal: GOAL-001 Participate in event
- Feature: FT-001 QR Code Check-in
- Delivery Level: L1 Walking Skeleton
- Priority: P0
- Context file: context.md

## Actor

Primary actor: organizer.

Secondary actor: attendee.

## Goal

The organizer validates an attendee QR code and the system records a single authoritative check-in for the correct event.

## Preconditions

- Organizer is authenticated.
- Organizer has permission to manage check-in for the event.
- Attendee has joined or RSVP'd to the event.
- Attendee can present an event-scoped QR code.
- The event is inside the allowed check-in window.

## Trigger

The organizer opens the event check-in scanner and scans an attendee QR code.

## Main Flow

1. Organizer opens the check-in scanner from the event management area.
2. System confirms organizer permission for the event.
3. Organizer scans the attendee QR code.
4. System extracts the signed QR token.
5. System validates token signature, event scope, attendee scope, expiration, and prior check-in state.
6. System records check-in server-side with an idempotent mutation.
7. System shows attendee name or safe display identifier, event name, and success confirmation.
8. System emits analytics and audit logs.

## Alternate Flows

### Already Checked In

1. Organizer scans a token for an attendee already checked in.
2. System returns the existing check-in state.
3. UI shows `Already checked in` with timestamp and no duplicate side effect.

### Expired QR Code

1. Organizer scans an expired token.
2. System rejects validation.
3. UI tells the organizer the attendee must refresh the QR code.

### Wrong Event

1. Organizer scans a valid token for a different event.
2. System rejects validation and does not reveal unnecessary attendee details.

### Permission Denied

1. Organizer lacks check-in permission.
2. System blocks scanner access or validation.
3. System logs the permission failure.

## Error And Edge Cases

- Invalid token: reject with generic invalid QR state.
- Network failure: show retryable error and do not mark attendance locally.
- Camera unavailable: show manual fallback guidance if a fallback is approved.
- Rapid duplicate scans: return idempotent result.
- Event outside check-in window: reject and explain the allowed state.

## Business Rules

- QR token must be event-scoped and attendee-scoped.
- QR token must expire according to DEC-001.
- Token strategy must follow DEC-002.
- Check-in mutation must be idempotent.
- Organizer permission must be checked server-side.
- Invalid validation must not expose unnecessary attendee data.

## UX States

- Scanner ready.
- Camera permission required.
- Scanning.
- Validating.
- Success.
- Already checked in.
- Expired QR.
- Invalid QR.
- Wrong event.
- Permission denied.
- Network retry.

## Acceptance Criteria

- [ ] Organizer can scan a valid attendee QR for an event they manage.
- [ ] Valid scan records exactly one check-in.
- [ ] Duplicate scans return an idempotent already-checked-in state.
- [ ] Expired, invalid, wrong-event, and permission-denied scans are rejected.
- [ ] Validation emits analytics and audit logs.
- [ ] UI avoids exposing attendee details when validation fails.

## Approval

- Approved by:
- Date:
- Notes:
