# UX Contract: Attendee Checks In With QR Code

| Field | Value |
| --- | --- |
| Status | draft |
| Source specification | [SPEC-001](../specification.md) |
| Contract version | 2 |

## Entry Points And Navigation

| Entry point | Preconditions | Destination | Exit/cancel behavior |
| --- | --- | --- | --- |
| Attendee event detail. | Joined event and authenticated session. | QR proof state. | Back returns to event without invalidating active proof. |
| Organizer event management. | Check-in permission. | Scanner state. | Exit stops camera use and returns to event management. |

## Interaction States

| State | Trigger | Displayed data/actions | Next states |
| --- | --- | --- | --- |
| Generating or active proof | Attendee opens check-in. | Progress or QR, event name, expiry text, refresh action. | Active, expired, or error. |
| Validating | Organizer submits scan. | Progress and disabled repeated submission. | Success, already checked in, or rejection. |
| Rejection | Expired, invalid, wrong event, denied, or network failure. | Safe reason and applicable refresh/retry action. | Active proof, scanner ready, or exit. |

## Feedback And Recovery

| Condition | User feedback | Recovery action | Persistence/retry behavior |
| --- | --- | --- | --- |
| Expired proof | Attendee is told proof expired. | Generate a new token. | Old token remains invalid. |
| Network failure | Organizer sees retryable failure; attendee sees no false success. | Retry scan online. | No local attendance mutation. |
| Already checked in | Both actors see existing attendance state. | Continue or scan next. | Timestamp remains unchanged. |

## Accessibility And Content

| Area | Requirement | Locale/content owner | Verification |
| --- | --- | --- | --- |
| QR alternative | Text communicates active/expired state and refresh action. | Product content. | Screen-reader and keyboard review. |
| Result feedback | Focus moves to result; status does not rely on color alone. | Design. | Automated roles plus manual assistive check. |
| Privacy | Failure copy never confirms attendee identity. | Product and Security. | UX/security review. |

## Requirements

| ID | Requirement | Source | Acceptance criteria | Dependencies |
| --- | --- | --- | --- | --- |
| REQ-003 | Attendee and organizer interfaces must expose accessible progress, success, rejection, expiry, refresh, retry, and already-checked-in states without leaking identity. | [Use Case](../use-case.md) | AC-003 | REQ-002 |
