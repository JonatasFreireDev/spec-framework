# Decision: QR Expiration Duration

## Context

- ID: DEC-001
- Status: approved
- Date: 2026-07-09
- Scope: feature/security

## Decision

Use short-lived QR check-in tokens with a proposed expiration window of 5 minutes.

This decision is approved. Downstream artifacts may use this expiration duration.

## Why

QR check-in tokens prove attendance for a specific event and attendee. A short expiration window reduces replay risk if a QR code is screenshotted, shared, or scanned later outside the intended check-in moment.

Five minutes is a starting proposal because it balances venue friction with replay resistance. It gives attendees enough time to open the screen and present the QR while limiting the usefulness of copied tokens.

## Options Considered

### Option A - 1 minute

- Pros: stronger replay resistance.
- Cons: high friction in venue lines, more refresh failures, worse accessibility for slower interactions.

### Option B - 5 minutes

- Pros: practical for live event check-in while still limiting replay risk.
- Cons: copied QR codes remain usable briefly.

### Option C - 15 minutes

- Pros: lower user friction and fewer refreshes.
- Cons: larger replay window and weaker abuse resistance.

## Consequences

- Positive: QR tokens become safer to display in a public venue context.
- Negative: attendee UI must show expiration and support refresh.
- Follow-up work: update specification, implementation plan, graph, and tasks after approval.

## Affected Artifacts

- domains/events/goals/participate-in-event/features/qr-code-check-in/use-cases/attendee-checks-in-with-qr-code/specification.md
- domains/events/goals/participate-in-event/features/qr-code-check-in/use-cases/attendee-checks-in-with-qr-code/implementation-plan.md
- domains/events/goals/participate-in-event/features/qr-code-check-in/use-cases/attendee-checks-in-with-qr-code/execution-graph.json
- domains/events/goals/participate-in-event/features/qr-code-check-in/use-cases/attendee-checks-in-with-qr-code/tasks.md

## Supersedes

- N/A

## Approval

- Approved by:
- Date:
- Notes: Proposed by framework example; awaiting product/security approval.
