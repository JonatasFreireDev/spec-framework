# Tasks: Organizer Validates QR Code

## Context

- ID: TASKSET-002
- Status: draft
- Source graph: GRAPH-002
- Source specification: SPEC-002
- Delivery Level: L1 Walking Skeleton
- Priority: P0
- Rationale: These tasks derive the executable planning units for the L1 organizer validation flow.

## Task List

### TK-002-001 Define Attendance Idempotency And Audit Persistence

- Type: data
- Depends on: none
- Output: documented data constraints for one check-in per attendee per event and audit logging expectations.
- Acceptance: duplicate scans cannot create duplicate attendance records.

### TK-002-002 Define Organizer Permission Validation

- Type: security
- Depends on: none
- Output: permission rule for organizer check-in access.
- Acceptance: validation is denied when organizer lacks event permission.

### TK-002-003 Define QR Token Validation Service Contract

- Type: backend
- Depends on: TK-002-001, TK-002-002
- Output: validation outcomes for valid, duplicate, expired, invalid, wrong event, outside window, and denied cases.
- Acceptance: every outcome maps to specification and UI state.

### TK-002-004 Define Check-in API Contract

- Type: api
- Depends on: TK-002-003
- Output: query and mutation contract for scanner session and validation.
- Acceptance: response shape supports UI states without leaking private data.

### TK-002-005 Define Organizer Scanner UI States

- Type: frontend
- Depends on: TK-002-004
- Output: scanner state model and UX copy requirements.
- Acceptance: all design states are represented.

### TK-002-006 Define Analytics And Observability Instrumentation

- Type: analytics
- Depends on: TK-002-003, TK-002-005
- Output: events, logs, metrics, and alert candidates.
- Acceptance: validation attempts and failure reasons are measurable.

### TK-002-007 Define Validation Test Coverage

- Type: qa
- Depends on: TK-002-001, TK-002-002, TK-002-003, TK-002-004, TK-002-005, TK-002-006
- Output: test plan covering behavior, security, UX states, and observability.
- Acceptance: tests prove acceptance criteria in `specification.md`.

## Notes

These are planning tasks for a future implementation repository. They do not implement application code in this framework repository.
