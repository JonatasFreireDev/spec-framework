# Tests: Attendee Checks In With QR Code

## Context

- Source specification: SPEC-001
- Source tasks: TK-001..TK-006
- Status: draft

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

## Specification v2 Acceptance Traceability

| Criterion | Requirement | Planned test/evidence | Risk covered |
| --- | --- | --- | --- |
| AC-001 | REQ-001 | TEST-001 end-to-end joined-attendee proof and authorized validation. | Product outcome and scope. |
| AC-002 | REQ-002 | TEST-002 state table, expiry, wrong-event, network, and concurrent scans. | Behavior and idempotency. |
| AC-003 | REQ-003 | TEST-003 QR/result accessibility and privacy-safe content review. | UX and accessibility. |
| AC-004 | REQ-004 | TEST-004 API authorization, schema, error, and retry suite. | Interface safety. |
| AC-005 | REQ-005 | TEST-005 schema, uniqueness, token cleanup, migration, and privacy tests. | Data integrity/lifecycle. |
| AC-006 | REQ-006 | TEST-006 replay, forgery, denial, secret-log, and disclosure negatives. | Security and abuse. |
| AC-007 | REQ-007 | TEST-007 gate and evidence completeness audit. | Quality-system conformance. |
| AC-008 | REQ-008 | TEST-008 signal, correlation, privacy, metric, and alert review. | Operability. |
| AC-009 | REQ-009 | TEST-009 additive migration, flag, pause, and rollback rehearsal. | Release reversibility. |

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
