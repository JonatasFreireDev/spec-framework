# Security Contract: Organizer Validates QR Code

| Field | Value |
| --- | --- |
| Status | draft |
| Source specification | [SPEC-002](../specification.md) |
| Contract version | 2 |

## Assets And Trust Boundaries

| Asset | Sensitivity | Boundary crossed | Actor/component | Security objective |
| --- | --- | --- | --- | --- |
| Organizer authority | Privileged event capability. | Organizer client to permission service to validation mutation. | Organizer and Events service. | Current authorization and least privilege. |
| QR bearer proof | Short-lived secret. | Attendee display through scanner to validation service. | Attendee, organizer, server. | Confidentiality, integrity, replay limitation. |
| Attendance/audit identity | Personal and security evidence. | Validation service to data/audit stores and restricted UI. | Events, Security, authorized staff. | Integrity, privacy, accountability. |

## Threats And Abuse Cases

| Threat/abuse | Preconditions | Impact | Existing mitigation | Gap owner |
| --- | --- | --- | --- | --- |
| User invokes validation without event role. | Authenticated account and endpoint access. | False attendance and identity disclosure. | Server permission check on session and every validation. | Product must define role mapping. |
| Copied/forged/wrong-event token. | Attacker presents proof. | False check-in or attendee enumeration. | Opaque high-entropy token, expiry, event binding, generic rejection. | Security review tests replay/rate behavior. |
| Sensitive scanner/log output. | Failure path logs request or displays identity. | Privacy or secret leak. | Allowlisted response and telemetry fields. | QA/security schema review. |

## Controls And Audit

| Control | Enforcement point | Deny/fail behavior | Audit evidence | Verification |
| --- | --- | --- | --- | --- |
| Current organizer permission. | Server session query and validation mutation. | Deny before token resolution response/write. | Safe actor/event denial event. | Negative authorization tests. |
| Token hash, scope, expiry, and single check-in. | Token service and database transaction. | Invalid/expired/wrong-event/existing result. | Safe result and correlation id. | Replay, expiry, wrong-event, concurrency tests. |
| Identity-safe serialization/logging. | API response and telemetry boundary. | Omit identity/token fields on rejection. | Schema review evidence. | Security/observability tests. |

## Residual Risk

| Risk | Likelihood/impact | Acceptance owner | Decision/evidence | Expiry/review trigger |
| --- | --- | --- | --- | --- |
| Organizer roles are not yet canonically defined. | Medium/high. | Product and Security. | [Context](../context.md) | Blocks implementation design. |
| Five-minute copied-token window before legitimate validation. | Low/medium. | Product and Security. | [DEC-001](../../../../../../../../../knowledge/decisions/DEC-001-qr-expiration-duration.md) | Abuse evidence or threat change. |

## Requirements

| ID | Requirement | Source | Acceptance criteria | Dependencies |
| --- | --- | --- | --- | --- |
| REQ-106 | Validation must reauthorize organizer authority and enforce token, event, privacy, audit, and concurrency controls without exposing secrets or identity on failure. | [Security Baseline](../../../../../../../../../knowledge/conventions/security-baseline.md) | AC-106 | REQ-104, REQ-105 |
