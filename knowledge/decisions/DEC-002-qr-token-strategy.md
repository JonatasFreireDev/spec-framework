# Decision: QR Token Strategy

## Context

- ID: DEC-002
- Status: approved
- Date: 2026-07-09
- Scope: architecture/security

## Decision

Use opaque, server-stored QR check-in tokens for v1.

This decision is approved. Downstream artifacts may use this token strategy.

## Why

QR check-in tokens must not expose attendee personal data and must be revocable, event-scoped, attendee-scoped, and auditable. An opaque server-stored token keeps validation server-authoritative and makes expiration, replay handling, and audit logging straightforward.

This is simpler and safer for v1 than signed stateless tokens because the product still needs to settle details around organizer permissions, replay handling, and operational logging.

## Options Considered

### Option A - Opaque server-stored token

- Pros: easy revocation, server-authoritative validation, simple audit trail, no PII in payload.
- Cons: requires token persistence and cleanup.

### Option B - Signed stateless token

- Pros: less database storage, simpler generation path.
- Cons: harder revocation, careful signing/key rotation needed, replay and audit behavior require extra design.

### Option C - Raw encoded attendance payload

- Pros: simplest to generate.
- Cons: exposes sensitive data and should not be used.

## Consequences

- Positive: validation can enforce current event, attendee, expiration, and organizer permission server-side.
- Negative: requires token table or equivalent persistence plus cleanup policy.
- Follow-up work: update data contract and task TK-001 after approval.

## Affected Artifacts

- product/domains/events/goals/participate-in-event/features/qr-code-check-in/use-cases/attendee-checks-in-with-qr-code/specification.md
- product/domains/events/goals/participate-in-event/features/qr-code-check-in/use-cases/attendee-checks-in-with-qr-code/implementation-plan.md
- product/domains/events/goals/participate-in-event/features/qr-code-check-in/use-cases/attendee-checks-in-with-qr-code/execution-graph.json
- product/domains/events/goals/participate-in-event/features/qr-code-check-in/use-cases/attendee-checks-in-with-qr-code/tasks.md

## Supersedes

- N/A

## Approval

- Approved by:
- Date:
- Notes: Proposed by framework example; awaiting architecture/security approval.