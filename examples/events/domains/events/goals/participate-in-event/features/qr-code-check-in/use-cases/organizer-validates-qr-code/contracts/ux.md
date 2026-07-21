# UX Contract: Organizer Validates QR Code

| Field | Value |
| --- | --- |
| Status | draft |
| Source specification | [SPEC-002](../specification.md) |
| Contract version | 2 |

## Entry Points And Navigation

| Entry point | Preconditions | Destination | Exit/cancel behavior |
| --- | --- | --- | --- |
| Event management check-in action. | Organizer session and current event permission. | Scanner ready or permission/window state. | Return to event management and release camera. |

## Interaction States

| State | Trigger | Displayed data/actions | Next states |
| --- | --- | --- | --- |
| Camera permission required | Scanner opens without camera access. | Rationale, grant/retry action, and exit. | Scanner ready or exit. |
| Scanner ready/scanning | Camera available. | Event name, scan target, torch/camera affordance where supported. | Validating or camera error. |
| Validating | Payload captured. | Progress, disabled repeat submission, safe cancel policy. | Success, existing, rejection, or retry. |
| Result | Server returns outcome. | Status text, permitted attendee display on success only, timestamp where applicable, scan-next/retry. | Scanner ready or exit. |

## Feedback And Recovery

| Condition | User feedback | Recovery action | Persistence/retry behavior |
| --- | --- | --- | --- |
| Expired proof | Tell organizer attendee must refresh. | Scan refreshed proof. | Expired proof remains invalid. |
| Invalid/wrong event/denied | Generic safe reason without attendee identity. | Verify event/permission or scan different proof. | No automatic repeat. |
| Network failure | Retryable service message; no success indication. | Retry online. | Preserve no local attendance. |

## Accessibility And Content

| Area | Requirement | Locale/content owner | Verification |
| --- | --- | --- | --- |
| Result announcements | Move focus/announce outcome; do not rely on color, sound, or animation alone. | Design and Product. | Keyboard and screen-reader review. |
| Camera controls | Controls have accessible names and adequate targets. | Design. | Automated semantics and manual target review. |
| Privacy-safe copy | Rejections do not confirm attendee identity. | Product and Security. | Content/security review. |

## Requirements

| ID | Requirement | Source | Acceptance criteria | Dependencies |
| --- | --- | --- | --- | --- |
| REQ-103 | The organizer experience must provide accessible scanner, permission, progress, result, retry, scan-next, and exit states with identity-safe rejection content. | [Use Case](../use-case.md) | AC-103 | REQ-102 |
