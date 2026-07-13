# Tests: Organizer Validates QR Code

## Context

- ID: TEST-002
- Status: draft
- Source specification: SPEC-002
- Delivery Level: L1 Walking Skeleton
- Priority: P0

## Quality Policy

| Field | Value |
| --- | --- |
| Engineering System | `ENGSYS-EVENTS-001 @ 0.1.0` |
| Policy | [Engineering Quality System](../../../../../../../../engineering/quality/quality-system.md) |
| Applicable risks | permissions, token lifecycle, data mutation, visual surface, accessibility, observability |
| Environments | documentation-fixture |
| Test data | synthetic-event, synthetic-user, synthetic-token |
| Platforms | web, mobile-camera |
| Deviations or exceptions | None; absent runtime evidence blocks validation |

## Behavioral Tests

- Valid organizer scans valid QR and receives checked-in result.
- Duplicate scan returns already-checked-in state without duplicate write.
- Expired token returns expired state.
- Invalid token returns invalid state without attendee details.
- Wrong-event token returns wrong-event state without unnecessary attendee details.
- Outside-window token returns outside-window state.

## Permission Tests

- Unauthenticated organizer cannot validate.
- Authenticated user without event permission cannot validate.
- Organizer permission is checked server-side on every validation.

## Data Tests

- `event_id + attendee_user_id` uniqueness prevents duplicate check-ins.
- Audit log records validation result and reason.
- Failed invalid scans do not require attendee ID.

## UX Tests

- Scanner ready state appears after session load.
- Camera permission state has clear recovery.
- Validating state prevents repeated submissions.
- Result states have scan next or retry action.
- Result text does not rely on color alone.

## Analytics Tests

- Scanner opened event fires once per scanner session.
- Scan submitted event fires for each scan attempt.
- Success and failure events include safe reason fields.
- Duplicate scans emit duplicate-detected event.

## Accessibility Tests

- Result state receives focus after validation.
- Screen reader can announce success and failure states.
- All actions are keyboard reachable.

## Residual Risk

- Offline venue operation is not covered unless approved as in scope.
- Manual fallback is not covered unless approved as in scope.
