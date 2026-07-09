# Tasks: Attendee Checks In With QR Code

## Context

- Source specification: SPEC-001
- Source implementation plan: PLAN-001
- Source execution graph: GRAPH-001
- Status: draft

## Rules

- Tasks must trace back to `specification.md`.
- Tasks must follow dependency order from `execution-graph.json`.
- A task must be independently reviewable and testable.
- Parallel tasks must have disjoint write scopes.

## Task: TK-001 Add Attendance Check-in Persistence

- Type: database
- Status: pending
- Depends on: none
- Source specification sections:
  - Data Contract
  - Permissions And Security
- Write scope:
  - supabase/migrations
- Objective:
  - Add persistence needed for checked-in attendance and QR validation.
- Acceptance criteria:
  - [ ] Attendance check-in timestamp can be stored once per event attendee.
  - [ ] Token data is event-scoped and attendee-scoped.
- Validation:
  - migration tests or SQL checks
- Handoff:
  - Next task(s): TK-002, TK-003
  - Risks: DEC-002 token strategy is approved; implement opaque server-stored token persistence.

## Task: TK-002 Implement QR Token Generation

- Type: backend
- Status: blocked
- Depends on: TK-001
- Source specification sections:
  - API Contract
  - Business Rules
- Write scope:
  - mobile services/actions
  - server functions
- Objective:
  - Generate opaque QR token for authenticated attendee and event.
- Acceptance criteria:
  - [ ] Token contains no raw PII.
  - [ ] Token expires.
  - [ ] Non-attendees cannot generate token.
- Validation:
  - unit/integration tests
- Handoff:
  - Next task(s): TK-004
  - Risks: DEC-001 expiration duration is approved; enforce 5-minute expiry.

## Task: TK-003 Implement QR Token Validation

- Type: backend
- Status: pending
- Depends on: TK-001, TK-002
- Source specification sections:
  - Functional Behavior
  - Permissions And Security
- Write scope:
  - mobile services/actions
  - server functions
- Objective:
  - Validate scanned QR server-side and mark attendance idempotently.
- Acceptance criteria:
  - [ ] Organizer permission is enforced.
  - [ ] Invalid, expired, and wrong-event tokens are rejected.
  - [ ] Duplicate scans are idempotent.
- Validation:
  - integration and security tests
- Handoff:
  - Next task(s): TK-005, TK-006
  - Risks: permission model must match existing app roles.

## Task: TK-004 Build Attendee QR UI States

- Type: frontend
- Status: pending
- Depends on: TK-002
- Source specification sections:
  - UX Contract
- Write scope:
  - mobile event screens
- Objective:
  - Let attendee view QR, expiration, refresh, loading, and error states.
- Acceptance criteria:
  - [ ] QR active and expired states are visible.
  - [ ] Text fallback is available for accessibility.
- Validation:
  - component tests or manual QA
- Handoff:
  - Next task(s): TK-006
  - Risks: screen ownership needs codebase confirmation.

## Task: TK-005 Build Organizer Validation UI States

- Type: frontend
- Status: pending
- Depends on: TK-003
- Source specification sections:
  - UX Contract
- Write scope:
  - mobile organizer event screens
- Objective:
  - Let organizer scan or submit QR payload and see validation results.
- Acceptance criteria:
  - [ ] Success, already checked in, invalid, expired, and permission denied states are shown.
- Validation:
  - component tests or manual QA
- Handoff:
  - Next task(s): TK-006
  - Risks: scanner dependency needs selection.

## Task: TK-006 Add Analytics, Tests, And QA Evidence

- Type: test/analytics
- Status: pending
- Depends on: TK-003, TK-004, TK-005
- Source specification sections:
  - Analytics And Observability
  - Acceptance Criteria
- Write scope:
  - tests
  - analytics instrumentation
  - product audit notes
- Objective:
  - Verify acceptance criteria and instrument success/failure events.
- Acceptance criteria:
  - [ ] qr_check_in_generated, qr_check_in_validated, and qr_check_in_failed are emitted.
  - [ ] QA evidence covers happy path and failure cases.
- Validation:
  - test run and QA report
- Handoff:
  - Next task(s): release readiness audit
  - Risks: analytics provider/event naming conventions need confirmation.