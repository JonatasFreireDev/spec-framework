# Design: Attendee Checks In With QR Code

## Context

- ID: UC-001:design
- Source specification: SPEC-001
- Status: draft

## UX Goal

Make check-in feel quick and trustworthy without exposing unnecessary personal data.

## Entry Points

- Attendee event detail screen.
- Organizer event management/check-in screen.

## Screens Or Components

- Attendee QR panel - shows QR, expiration, refresh, and status.
- Organizer scanner/validation result - scans QR and shows result.

## States

- Default: active QR and expiration hint.
- Loading: generating or validating.
- Empty: attendee has not joined event.
- Error: invalid, expired, permission denied, network failure.
- Success: checked in.
- Disabled/permission denied: organizer lacks permission.

## Interaction Rules

- Expired QR should offer refresh.
- Validation errors should avoid revealing attendee personal details.
- Duplicate scan should show already checked in without alarm.

## Accessibility

- Keyboard: refresh and retry actions must be reachable.
- Screen reader: QR status must be announced as text.
- Contrast: success/error states must meet contrast requirements.
- Motion: no required motion for validation feedback.

## Content

- Labels: "Check-in QR", "Refresh QR", "Validate check-in".
- Error messages: "This QR code is no longer valid." "You cannot validate this event."
- Confirmation messages: "Check-in confirmed." "Already checked in."

## Open Questions

- Scanner library and camera permission UX are not selected.
