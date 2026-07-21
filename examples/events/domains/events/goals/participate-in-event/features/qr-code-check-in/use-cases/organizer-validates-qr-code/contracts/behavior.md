# Behavior Contract: Organizer Validates QR Code

| Field | Value |
| --- | --- |
| Status | draft |
| Source specification | [SPEC-002](../specification.md) |
| Contract version | 2 |

## Preconditions And Triggers

| Trigger | Preconditions | Actor/system | Rejection behavior |
| --- | --- | --- | --- |
| Organizer opens event scanner. | Authenticated and permitted for event; window open. | Organizer and session service. | Block scanner or return permission/window state. |
| Organizer scans proof. | Ready scanner session and online service. | Scanner and validation service. | Do not create local attendance on unavailable service. |

## State Transitions

| From | Event/command | Guard | To | Side effects |
| --- | --- | --- | --- | --- |
| scanner ready | submit scan | camera and payload available | validating | Disable duplicate submission and start correlation. |
| validating | authorized valid proof | event, attendee, expiry, permission valid | success | Record check-in once and safe audit. |
| validating | attendance already exists | same event and attendee | already checked in | Return timestamp; no write. |
| validating | any rejection | invalid, expired, wrong event, denied, outside window | rejection | Emit safe reason; reveal no attendee identity. |

## Alternate Error And Edge Flows

| Type | Condition | Expected behavior | Recovery/retry | Signal |
| --- | --- | --- | --- | --- |
| Alternate | Proof already checked in. | Show existing timestamp and next-scan action. | Continue scanning. | Duplicate result event. |
| Error | Camera unavailable or permission denied. | Explain scanner limitation; no unapproved manual search. | Retry camera permission or exit. | Client capability event. |
| Error | Network/service unavailable. | Show retryable result and keep attendance unchanged. | Retry online. | Availability/latency signal. |
| Edge | Several organizers scan the same proof concurrently. | One authoritative write; all receive success/existing state. | Continue scanning. | Concurrency result metric. |

## Invariants

| Invariant | Concurrency/idempotency impact | Enforcement point | Source |
| --- | --- | --- | --- |
| Validation never trusts client-decoded identity, event, permission, expiry, or attendance state. | Every attempt rechecks current server state. | Validation service. | [DEC-002](../../../../../../../../../knowledge/decisions/DEC-002-qr-token-strategy.md) |
| A failed result never mutates attendance or reveals unnecessary attendee data. | Retries cannot convert a rejection into a local success. | Service and response serializer. | [Use Case](../use-case.md) |

## Requirements

| ID | Requirement | Source | Acceptance criteria | Dependencies |
| --- | --- | --- | --- | --- |
| REQ-102 | Scanner and validation state transitions must deterministically cover readiness, validation, success, existing attendance, token failures, denial, window, camera, network, and concurrent scans. | [Use Case](../use-case.md) | AC-102 | REQ-101 |
