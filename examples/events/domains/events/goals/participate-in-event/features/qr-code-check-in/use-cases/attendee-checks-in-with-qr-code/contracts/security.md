# Security Contract: Attendee Checks In With QR Code

| Field | Value |
| --- | --- |
| Status | draft |
| Source specification | [SPEC-001](../specification.md) |
| Contract version | 2 |

## Assets And Trust Boundaries

| Asset | Sensitivity | Boundary crossed | Actor/component | Security objective |
| --- | --- | --- | --- | --- |
| Raw QR token | Short-lived bearer secret. | Attendee display to organizer scanner to server. | Attendee, organizer, validation service. | Confidentiality and replay resistance. |
| Attendance identity/state | Personal/internal. | Client to server to database. | Authorized actors and Events service. | Authorization, integrity, and privacy. |

## Threats And Abuse Cases

| Threat/abuse | Preconditions | Impact | Existing mitigation | Gap owner |
| --- | --- | --- | --- | --- |
| Screenshot or replay. | Attacker obtains active proof. | Unauthorized presence claim. | Five-minute expiry, event scope, server validation, single check-in. | Security review verifies residual replay window. |
| Forged, wrong-event, or enumerated token. | Attacker submits guesses or copied proof. | Identity disclosure or false check-in. | Opaque high-entropy token, generic rejection, rate monitoring. | Engineering defines rate control. |
| Unauthorized organizer. | Authenticated user lacks event permission. | Attendance corruption and privacy breach. | Server permission check before lookup response and mutation. | Product defines canonical roles. |

## Controls And Audit

| Control | Enforcement point | Deny/fail behavior | Audit evidence | Verification |
| --- | --- | --- | --- | --- |
| Token binding, expiry, and hashing. | Token and validation services. | Generic invalid/expired result. | Safe result reason without token or identity. | Security and integration tests. |
| Organizer authorization. | Server mutation boundary. | Deny without attendance change or attendee disclosure. | Permission-denied event with event and actor identifiers. | Negative authorization test. |
| Idempotent attendance write. | Transaction/database constraint. | Return existing state. | Duplicate result metric. | Concurrent validation test. |

## Residual Risk

| Risk | Likelihood/impact | Acceptance owner | Decision/evidence | Expiry/review trigger |
| --- | --- | --- | --- | --- |
| Copied token remains usable inside five-minute window before legitimate check-in. | Low/medium. | Product and Security. | [DEC-001](../../../../../../../../../knowledge/decisions/DEC-001-qr-expiration-duration.md) | Abuse evidence or threat-register change. |
| Organizer role policy is not yet canonical. | Medium/high for implementation. | Product. | [Context](../context.md) | Blocks Engineering Proposal. |

## Requirements

| ID | Requirement | Source | Acceptance criteria | Dependencies |
| --- | --- | --- | --- | --- |
| REQ-006 | Server-side token and permission controls must prevent forgery, replay side effects, unauthorized writes, secret logging, and identity disclosure on rejected validation. | [Security Baseline](../../../../../../../../../knowledge/conventions/security-baseline.md) | AC-006 | REQ-004, REQ-005 |
