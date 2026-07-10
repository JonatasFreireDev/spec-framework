# Analytics: Attendee Checks In With QR Code

## Context

- Source artifact: SPEC-001
- Status: draft

## Measurement Goal

Measure whether QR check-in improves reliable attendance confirmation and identify validation friction.

## Events

| Event | Trigger | Properties | User | Notes |
| --- | --- | --- | --- | --- |
| qr_check_in_generated | Attendee receives active QR | event_id, expires_in_seconds | attendee | No PII in payload. |
| qr_check_in_validated | Organizer validates QR successfully | event_id, result | organizer | result can be success or already_checked_in. |
| qr_check_in_failed | Validation fails | event_id, reason | organizer | reason must not expose attendee PII. |

## Funnels Or Metrics

- RSVP to check-in conversion - checked-in attendees divided by RSVP attendees.
- QR validation success rate - successful validations divided by validation attempts.

## Guardrails

- Invalid token rate should not spike after rollout.
- Permission denied validation attempts should remain low.

## Privacy Notes

- Do not include attendee name, email, phone, or raw token in analytics.

## Validation

- Confirm all three events fire in happy path and failure paths.