# Implementation Plan: Organizer Validates QR Code

## Context

- ID: PLAN-002
- Status: draft
- Source specification: SPEC-002
- Source design: DES-002
- Delivery Level: L1 Walking Skeleton
- Priority: P0

## Objective

Plan the implementation of organizer QR validation without writing application code in this framework repo.

## Scope

- Organizer scanner session query.
- Server-side validation mutation.
- Attendance idempotency.
- Permission checks.
- Result states in organizer UI.
- Analytics and audit logs.
- Tests for valid, duplicate, invalid, expired, wrong-event, permission-denied, and network paths.

## Dependencies

- Approved Specification.
- Approved Design.
- DEC-001 QR expiration duration.
- DEC-002 QR token strategy.
- User identity and organizer permission model.
- Event attendance data model.

## Phases

1. Data model and constraints for attendance and audit logging.
2. Server-side validation service and permission checks.
3. API query and mutation contract.
4. Organizer scanner UI states.
5. Analytics and observability instrumentation.
6. Automated and manual test coverage.
7. QA, audit, and release readiness.

## Sequencing

Implement server-authoritative validation before UI marks success. UI may render scanner states only after API contracts are stable enough to test.

## Risks

- Role permissions may be underspecified.
- Offline validation may be requested later and would change security assumptions.
- Token expiration rules may affect venue operations.
- Poor scanner feedback could slow queues.

## Test Plan

- Unit tests for token validation outcomes.
- Integration tests for permission and idempotent attendance mutation.
- UI state tests for each scanner result.
- Accessibility checks for status changes.
- Analytics event assertions.

## Rollback Plan

- Disable feature flag if available.
- Hide organizer scanner entry point.
- Preserve attendance records already created.
- Keep audit logs for investigation.

## Probable Application Areas

These are placeholders for a future application repository, not files to create here:

- event attendance persistence.
- check-in validation service.
- organizer event management API.
- organizer scanner UI.
- analytics and audit logging.

## Decisions Needing ADR Or Approval

- Offline validation support.
- Exact organizer permission roles.
- Manual fallback for camera failure.
- Feature flag naming and rollout owner.

## Candidate Tasks

See `tasks.md` and `execution-graph.json`.
