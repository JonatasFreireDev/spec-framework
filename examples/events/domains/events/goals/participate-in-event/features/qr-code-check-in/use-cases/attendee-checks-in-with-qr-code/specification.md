# Specification: Attendee Checks In With QR Code

## Context

- ID: SPEC-001
- Status: draft
- Source use case: UC-001
- Source feature: FT-001
- Context file: context.md

## Delivery

- Level: L1 Walking Skeleton
- Priority: P0
- Rationale: Attendee QR generation and presentation are required for the first end-to-end check-in flow.
- Depends on:
  - FT-001
  - DEC-001
  - DEC-002

## Product Context

QR check-in supports the user goal of participating in an event and gives the product a stronger signal of real attendance than RSVP alone.

## Scope

### In Scope

- Generate an opaque QR token for an authenticated attendee and event.
- Validate a scanned QR token server-side.
- Mark an attendance record as checked in exactly once.
- Return clear states for success, expired, invalid, already checked in, and permission denied.
- Emit analytics and logs for generation, success, and failure.

### Non-Goals

- Offline validation.
- Paid ticketing.
- Advanced fraud scoring.

## Functional Behavior

### Main Flow

1. Attendee requests QR code for an event they joined.
2. Server creates an opaque, event-scoped, attendee-scoped token with expiration.
3. Client renders QR code and text status.
4. Organizer scans QR code from an organizer check-in surface.
5. Server validates token, expiration, event, attendee, and organizer permission.
6. Server updates attendance record with checked_in_at if not already set.
7. Server returns success or already_checked_in.

### Alternate Flows

- Expired QR: reject and let attendee refresh.
- Already checked in: return idempotent already_checked_in without changing timestamp.
- Wrong event: reject with invalid_for_event.

### Error States

- invalid_token - no attendee details returned - log validation failure.
- expired_token - user can refresh QR - analytics event emitted.
- permission_denied - organizer cannot validate this event - security log emitted.
- network_error - client shows retryable state.

### Edge Cases

- Multiple scans in quick succession must be idempotent.
- Token replay after successful check-in must not create new attendance records.
- QR payload must not expose raw user PII.

## Business Rules

- QR token must be opaque and validated by the server.
- QR token must expire after 5 minutes, per DEC-001.
- Check-in mutation must be idempotent.
- Only authorized organizers can validate QR codes for their event.

## UX Contract

- Entry points: attendee event detail; organizer check-in scanner.
- UI states: QR active, QR expired, generating, validating, success, invalid, permission denied, already checked in.
- Accessibility requirements: QR code has textual fallback status and refresh action label.
- Copy/content requirements: errors should not expose whether a token belongs to a specific person.

## API Contract

### Commands / Mutations

- `generateEventCheckInQr(eventId)`
  - Request: authenticated attendee session and event id.
  - Response: opaque token, expiresAt, qrPayload.
  - Errors: not_joined, event_not_found, check_in_not_open.

- `validateEventCheckInQr(eventId, token)`
  - Request: authenticated organizer session, event id, token.
  - Response: status success | already_checked_in, checkedInAt.
  - Errors: invalid_token, expired_token, permission_denied, invalid_for_event.

### Queries

- `getEventCheckInStatus(eventId)`
  - Parameters: event id, authenticated attendee.
  - Response: joined status, checked-in status, checkedInAt.

## Data Contract

- Tables/entities: event_attendance, event_check_in_tokens.
- Fields: event_id, attendee_id, checked_in_at, token_hash, expires_at, consumed_at.
- Constraints: unique event_id + attendee_id attendance record.
- Retention/privacy: token records should expire or be pruned.
- Migration needed: yes.

## Permissions And Security

- Who can read: attendee can read own QR/check-in status; organizer can read check-in result for managed event.
- Who can write: server mutation can mark checked_in_at after validating organizer permission.
- Server-authoritative checks: event membership, organizer permission, token expiration, token event binding.
- Abuse cases: token replay, forged QR, organizer scanning wrong event, PII exposure.
- Privacy/LGPD notes: QR payload must not include raw attendee personal data.

## Analytics And Observability

- Analytics events: qr_check_in_generated, qr_check_in_validated, qr_check_in_failed.
- Logs: permission denied, invalid token, expired token, idempotent duplicate scan.
- Metrics: check-in success rate, check-in failure rate, duplicate scan rate.
- Alerts: unusual invalid token spike.

## Performance And Reliability

- Latency expectations: validation should complete in under two seconds under normal conditions.
- Offline/retry behavior: no offline mutation in v1; retry network failures.
- Concurrency/idempotency: duplicate scans must not change checked_in_at after first success.

## Rollout

- Feature flag: event_qr_check_in.
- Migration/backfill: add check-in fields before enabling validation.
- Rollback: disable flag; keep attendance records intact.

## Acceptance Criteria

- [ ] Attendee can generate a non-PII QR token for an event they joined.
- [ ] Organizer can validate a valid token for an event they manage.
- [ ] Invalid, expired, wrong-event, and permission-denied tokens are rejected.
- [ ] Duplicate scans are idempotent.
- [ ] Analytics and logs are emitted for success and failure.

## Open Questions

- DEC-001 QR expiration duration is approved.
- DEC-002 QR token strategy is approved.

## Approval

- Approved by:
- Date:
- Notes:
