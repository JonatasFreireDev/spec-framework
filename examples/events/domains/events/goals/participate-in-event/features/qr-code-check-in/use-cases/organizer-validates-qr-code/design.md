# Design: Organizer Validates QR Code

## Context

- ID: DES-002
- Status: draft
- Source specification: SPEC-002
- Use case: UC-002
- Delivery Level: L1 Walking Skeleton
- Priority: P0

## UX Goal

Let an organizer scan attendees quickly, understand the result immediately, and recover safely from invalid, expired, duplicate, permission, camera, or network states.

## Entry Points

- Event management screen.
- Event check-in action.
- Optional release or operations checklist link if an event is in progress.

## Primary Flow

1. Organizer opens event.
2. Organizer selects check-in scanner.
3. System confirms event name and check-in window.
4. Organizer grants camera permission if needed.
5. Organizer scans QR.
6. Scanner pauses while validation runs.
7. Result appears with next scan action.

## UI Regions

- Header: event name, check-in window state, close/back control.
- Scanner area: camera preview and scanning target.
- Status area: ready, validating, or result message.
- Result detail: attendee display data only on successful or already checked-in state.
- Action area: scan next, retry, request QR refresh, or exit.

## States

### Ready

Scanner is active and the organizer can scan the attendee QR.

### Camera Permission Required

Explain that camera access is needed. Provide retry permission action and a non-committal note that manual fallback requires approval.

### Validating

Pause repeated scans and show progress. Do not show attendee details yet.

### Success

Show checked-in confirmation, attendee display name, timestamp, and scan next action.

### Already Checked In

Show already checked-in state, prior timestamp when available, and scan next action.

### Expired QR

Show that the attendee should refresh their QR code. Do not expose additional attendee details.

### Invalid QR

Show generic invalid QR state. Do not reveal whether the token resembles a real attendee.

### Wrong Event

Show wrong event state and ask attendee to open the correct event QR. Avoid unnecessary attendee details.

### Permission Denied

Show that the organizer lacks permission and must contact an event owner or admin.

### Network Retry

Show retryable failure. Do not mark local check-in as complete.

## Accessibility

- Provide text equivalents for scanner and validation states.
- Do not rely on color alone for success or failure.
- Move focus to result state after validation.
- Ensure result text is short, explicit, and screen-reader friendly.
- Keep scan next action reachable without precise gestures.

## Content Guidelines

- Use action-oriented messages.
- Avoid exposing technical token details to organizers.
- Avoid attendee identity details for rejected scans.

## Data Displayed

- Event name.
- Check-in window state.
- Attendee display name only for successful or already checked-in result.
- Check-in timestamp when applicable.
- Result reason in user-safe terms.

## Open Questions

- Should there be a manual fallback when camera access fails?
- What exact organizer roles can access the scanner?

## Approval

- UX approved by:
- Date:
- Notes:
