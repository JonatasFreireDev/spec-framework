# Implementation Plan: Attendee Checks In With QR Code

## Context

- ID: PLAN-001
- Status: draft
- Source specification: SPEC-001
- Context file: context.md
- Delivery Level: L1 Walking Skeleton
- Priority: P0
- Rationale: The plan sequences the minimum technical work needed to prove QR check-in in the walking skeleton.

## Technical Objective

Build server-authoritative QR generation and validation with idempotent attendance updates, basic UI states, analytics, and tests.

## Constraints

- Product: QR expiration window is 5 minutes per DEC-001.
- UX: QR must have text status and refresh affordance.
- Architecture: validation must happen server-side.
- Data: attendance must be unique per event and attendee.
- Security/privacy: QR payload must not expose PII.
- Delivery: v1 excludes offline validation.

## Proposed Phases

### Phase 1 - Data Model

- Goal: support attendance check-in state and token validation.
- Work: add fields/tables for attendance and QR token tracking.
- Dependencies: existing event and user models.
- Exit criteria: migration and constraints are testable.

### Phase 2 - Server Mutations

- Goal: generate and validate QR tokens server-side.
- Work: implement generation, validation, permission checks, idempotent update.
- Dependencies: Phase 1.
- Exit criteria: server tests cover success and failure states.

### Phase 3 - UI States

- Goal: expose attendee QR and organizer validation feedback.
- Work: QR display, expiration state, validation result states.
- Dependencies: Phase 2 API contract.
- Exit criteria: UI handles loading, success, expired, invalid, permission denied.

### Phase 4 - Analytics And QA

- Goal: instrument and verify the complete flow.
- Work: analytics events, logs, test matrix, QA pass.
- Dependencies: Phases 2 and 3.
- Exit criteria: acceptance criteria are covered.

## Candidate Tasks

- TK-001 - Add attendance check-in persistence - type: database - depends on: none
- TK-002 - Implement QR token generation - type: backend - depends on: TK-001
- TK-003 - Implement QR token validation - type: backend - depends on: TK-001,TK-002
- TK-004 - Build attendee QR UI states - type: frontend - depends on: TK-002
- TK-005 - Build organizer validation UI states - type: frontend - depends on: TK-003
- TK-006 - Add analytics, tests, and QA evidence - type: test/analytics - depends on: TK-003,TK-004,TK-005

## Files Or Modules Likely Touched

- supabase/migrations - data model changes.
- mobile services/actions - QR generation and validation calls.
- mobile event screens - attendee and organizer UI states.
- docs/product artifacts - update context and audit notes.

## Test Strategy

- Unit: token generation and validation helpers.
- Integration: mutation validates permissions and idempotency.
- E2E/manual: attendee presents QR and organizer validates.
- Security/RLS: attendee cannot validate arbitrary event; organizer cannot validate unmanaged event.
- Regression: duplicate scan behavior.

## Rollout And Rollback

- Rollout: behind `event_qr_check_in` flag.
- Feature flag: required.
- Rollback: disable flag and leave attendance data intact.
- Data recovery: no destructive migration in v1.

## Risks

- Token expiration is approved by DEC-001; implementation must enforce 5-minute expiry.
- Permission mismatch - likelihood: medium - mitigation: test organizer authorization explicitly.

## Approved Decisions

- DEC-001 QR expiration duration - owner: product/security - status: approved.
- DEC-002 QR token strategy - owner: architecture/security - status: approved.

## Approval

- Approved by:
- Date:
- Notes:
