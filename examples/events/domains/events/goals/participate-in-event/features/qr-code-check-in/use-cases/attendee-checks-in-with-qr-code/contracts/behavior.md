# Behavior Contract: Attendee Checks In With QR Code

| Field | Value |
| --- | --- |
| Status | draft |
| Source specification | [SPEC-001](../specification.md) |
| Contract version | 2 |

## Preconditions And Triggers

| Trigger | Preconditions | Actor/system | Rejection behavior |
| --- | --- | --- | --- |
| Attendee opens check-in proof. | Authenticated, joined event, check-in window open. | Attendee and token service. | Return not-joined, event-not-found, or window-closed without issuing proof. |
| Organizer submits scanned proof. | Authenticated organizer, managed event, valid scanner session. | Organizer and validation service. | Deny before attendance mutation. |

## State Transitions

| From | Event/command | Guard | To | Side effects |
| --- | --- | --- | --- | --- |
| no active token | generate proof | attendee joined | active token | Store opaque token hash with five-minute expiry. |
| active token | validate proof | token, event, attendee, and permission valid | checked in | Set checked_in_at once and consume token. |
| checked in | validate again | same attendance | checked in | Return existing timestamp; write no duplicate attendance. |

## Alternate Error And Edge Flows

| Type | Condition | Expected behavior | Recovery/retry | Signal |
| --- | --- | --- | --- | --- |
| Alternate | Token expired before scan. | Reject as expired without identity disclosure. | Attendee refreshes proof. | qr_check_in_failed reason expired. |
| Error | Wrong event, invalid token, or denied organizer. | Reject before mutation with safe result. | Correct event or permission; invalid proof is not retried automatically. | Safe denial log. |
| Edge | Concurrent scans validate the same attendee. | One write wins; all callers receive checked-in or already-checked-in. | No client reconciliation write. | Duplicate metric. |

## Invariants

| Invariant | Concurrency/idempotency impact | Enforcement point | Source |
| --- | --- | --- | --- |
| One attendance row exists per event and attendee and its first checked_in_at is immutable. | Unique key and transactional mutation prevent duplicates. | Service and database. | [Use Case](../use-case.md) |
| QR payload contains no raw attendee PII. | Replay or screenshot cannot directly disclose identity. | Token generation. | [DEC-002](../../../../../../../../../knowledge/decisions/DEC-002-qr-token-strategy.md) |

## Requirements

| ID | Requirement | Source | Acceptance criteria | Dependencies |
| --- | --- | --- | --- | --- |
| REQ-002 | Validation must produce deterministic success, already-checked-in, expired, invalid, wrong-event, denied, window-closed, and retryable-network outcomes without duplicate writes. | [Use Case](../use-case.md) | AC-002 | REQ-001 |
