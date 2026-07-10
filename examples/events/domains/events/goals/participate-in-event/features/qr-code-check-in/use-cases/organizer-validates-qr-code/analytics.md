# Analytics: Organizer Validates QR Code

## Context

- ID: ANA-002
- Status: draft
- Source use case: UC-002
- Delivery Level: L1 Walking Skeleton
- Priority: P0

## Product Questions

- Can organizers complete check-in validation reliably?
- Which validation failures are most common?
- Are duplicate scans frequent enough to affect venue operations?
- Are permission failures signaling configuration problems?

## Events

### qr_check_in_scanner_opened

- Actor: organizer.
- Properties: `event_id`, `organizer_role`, `check_in_window_state`.
- Purpose: measure scanner entry.

### qr_check_in_scan_submitted

- Actor: organizer.
- Properties: `event_id`, `scan_source`, `client_timestamp`.
- Purpose: count validation attempts.

### qr_check_in_validation_succeeded

- Actor: organizer.
- Properties: `event_id`, `check_in_id`, `time_to_validate_ms`.
- Purpose: measure successful check-ins.

### qr_check_in_validation_failed

- Actor: organizer.
- Properties: `event_id`, `reason`, `retryable`.
- Purpose: diagnose expired, invalid, wrong-event, permission, outside-window, and network failures.

### qr_check_in_duplicate_detected

- Actor: organizer.
- Properties: `event_id`, `checked_in_at_present`.
- Purpose: measure duplicate scan behavior.

## Logs

- Validation result and reason.
- Organizer permission denial.
- Token validation failure category.
- Service or network error category.

## Metrics

- Validation success rate.
- Failed scans by reason.
- Duplicate scan rate.
- Median validation latency.
- Permission denied rate.

## Privacy

- Avoid logging raw QR tokens.
- Avoid logging attendee details for invalid, expired, or wrong-event scans.
- Use stable internal IDs only where necessary for auditability.

## Open Questions

- Which analytics system will receive these events?
- What retention policy applies to check-in audit logs?
