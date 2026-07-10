# Specification: Organizer Validates QR Code

## Context

- ID: SPEC-002
- Status: proposed
- Source use case: UC-002
- Source feature: FT-001
- Context file: context.md

## Delivery

- Level: L1 Walking Skeleton
- Priority: P0
- Rationale: Organizer validation is required to prove the event attendance flow end to end.
- Depends on:
  - DOMAIN-users
  - DOMAIN-events
  - DEC-001
  - DEC-002

## Product Context

The Events domain needs a trusted way for organizers to validate that an attendee is present. QR validation turns attendee intent into an operational attendance record while limiting fraud, duplicate scans, and unnecessary personal data exposure.

## User Goal

Participate in event.

## Feature Scope

### In Scope

- Organizer opens scanner for an event they manage.
- Organizer scans attendee QR code.
- System validates token signature, event, attendee, expiration, and organizer permission.
- System records idempotent check-in.
- System shows success, already checked in, invalid, expired, wrong event, permission denied, and network retry states.
- System emits analytics and audit logs.

### Non-Goals

- Offline authoritative validation.
- Payment, refunds, or ticket transfer.
- Organizer role management UI.
- Attendee QR generation, except as an upstream dependency.
- Manual attendee search fallback unless separately approved.

## Use Cases

- UC-002 Organizer Validates QR Code.
- Related upstream: UC-001 Attendee Checks In With QR Code.

## Business Rules

- QR token must be signed and tamper resistant.
- QR token must include or resolve to event and attendee scope.
- QR token expiration follows DEC-001.
- Token strategy follows DEC-002.
- Check-in mutation must be idempotent.
- Organizer permission must be checked server-side for every validation.
- Failed validation must not expose unnecessary attendee details.

## UX Flow

1. Organizer enters event management.
2. Organizer opens check-in scanner.
3. System requests camera permission if needed.
4. Organizer scans QR.
5. UI shows validating state.
6. UI shows a result state and next scan affordance.

## UI States

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

## API Contracts

### Mutation: Validate QR Check-in

- Name: `validateEventCheckIn`
- Actor: authenticated organizer.
- Request:
  - `eventId`
  - `qrToken`
  - `scanDeviceId` optional
  - `clientTimestamp`
- Response:
  - `status`: `checked_in`, `already_checked_in`, `expired`, `invalid`, `wrong_event`, `permission_denied`, `outside_window`
  - `checkInId` when applicable
  - `attendeeDisplayName` only for successful or already checked-in responses
  - `checkedInAt` when applicable
  - `retryable` boolean
- Errors:
  - authentication required
  - authorization denied
  - malformed request
  - service unavailable

### Query: Event Check-in Session

- Name: `getEventCheckInSession`
- Parameters:
  - `eventId`
- Response:
  - event display name
  - organizer permission status
  - check-in window status
  - scanner availability guidance

## Data Model

- `event_attendance`
  - `event_id`
  - `attendee_user_id`
  - `status`
  - `checked_in_at`
  - `checked_in_by_user_id`
  - unique constraint on `event_id + attendee_user_id`
- `check_in_audit_log`
  - `event_id`
  - `organizer_user_id`
  - `attendee_user_id` nullable for invalid scans
  - `result`
  - `reason`
  - `created_at`

## Events

- `qr_check_in_scanner_opened`
- `qr_check_in_scan_submitted`
- `qr_check_in_validation_succeeded`
- `qr_check_in_validation_failed`
- `qr_check_in_duplicate_detected`

## Analytics

See `analytics.md`.

## Permissions

- Organizer must be authenticated.
- Organizer must have event check-in permission.
- Attendee data should only be revealed after successful validation or already checked-in result for the same event.

## Security

- Treat QR token as bearer proof with short validity.
- Validate server-side; client validation is advisory only.
- Log denied permission attempts.
- Avoid leaking attendee identity for invalid, expired, or wrong-event tokens.

## Performance

- Scanner open should load session state within an acceptable interactive delay.
- Validation should return promptly enough for venue queues.
- Duplicate rapid scans should not produce duplicate writes.

## Accessibility

- Scanner must provide text status for each result.
- Result states must not rely on color alone.
- Camera permission failure must provide clear next action.
- Success and failure states should be readable by assistive technology.

## Error States

- Invalid token.
- Expired token.
- Wrong event.
- Outside check-in window.
- Permission denied.
- Network or service unavailable.
- Camera unavailable.

## Edge Cases

- Same QR scanned repeatedly.
- Organizer opens scanner before event check-in window.
- Attendee refreshes QR while organizer scans older token.
- Multiple organizers scan the same attendee nearly simultaneously.

## Observability

- Log validation result and reason.
- Count failed scans by reason.
- Alert only on abnormal failure spikes or permission-denied spikes.
- Preserve privacy in logs.

## Rollout Strategy

- Gate behind event check-in feature flag if the product has flags.
- Pilot with internal or low-risk events before broad release.
- Monitor validation failure rates and queue-impact feedback.

## Feature Flags

- Candidate flag: `events.qr_check_in.organizer_validation`.
- Flag requirement needs approval before implementation.

## Acceptance Criteria

- [ ] Organizer permission is checked before validation.
- [ ] Valid scan records one check-in.
- [ ] Duplicate scan returns already checked-in without duplicate mutation.
- [ ] Invalid, expired, wrong-event, outside-window, and permission-denied scans are handled.
- [ ] Successful response may show attendee display data; failed responses do not reveal unnecessary attendee data.
- [ ] Analytics and audit logs are emitted for validation attempts.

## Open Questions

- Should a manual fallback exist for camera failure?
- Which event roles count as organizer check-in permission?
- Is online-only validation acceptable for L1?
