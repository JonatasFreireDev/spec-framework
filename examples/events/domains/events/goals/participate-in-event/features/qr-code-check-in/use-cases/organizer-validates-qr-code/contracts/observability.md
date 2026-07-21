# Observability Contract: Organizer Validates QR Code

| Field | Value |
| --- | --- |
| Status | draft |
| Source specification | [SPEC-002](../specification.md) |
| Contract version | 2 |

## Signals And Correlation

| Signal | Trigger | Fields/dimensions | Correlation | Consumer |
| --- | --- | --- | --- | --- |
| qr_check_in_scanner_opened | Authorized scanner session opens. | event_id, permission/window status, client capability; no attendee. | scanner_session_id and request_id. | Product and Operations. |
| qr_check_in_scan_submitted | Payload submitted. | event_id, client capability, attempt sequence; no token. | scanner_session_id and validation_id. | Operations. |
| qr_check_in_validation_succeeded/failed | Server returns result. | safe result/reason, latency, event_id when authorized. | validation_id and request_id. | Product, Operations, Security. |
| qr_check_in_duplicate_detected | Existing attendance returned. | event_id and safe concurrency/retry dimension. | validation_id. | Engineering and Operations. |

## Metrics And Service Expectations

| Metric | Definition | Target/threshold | Window | Decision supported |
| --- | --- | --- | --- | --- |
| Scanner readiness latency | Time from open to ready/rejection state. | Evidence-backed pilot baseline; no invented SLO. | Per event and five-minute rolling. | Diagnose camera/session/permission delay. |
| Validation latency and result rate | Server time/count by safe result. | Interactive expectation under normal conditions. | Per event and five-minute rolling. | Pause rollout when venue flow degrades. |
| Denied/invalid/wrong-event rate | Count by safe reason. | Alert on abnormal deviation from pilot baseline. | Five-minute window. | Detect abuse, configuration, or role regression. |

## Alerts And Diagnosis

| Alert/symptom | Condition | Severity | Runbook/diagnostic path | Owner |
| --- | --- | --- | --- | --- |
| Failure or latency spike during active event. | Sustained abnormal rate against approved pilot baseline. | High. | Check service/dependency health, expiry distribution, permission denials, and rollout. | Events operations. |
| Permission-denied spike. | Material change for one event/role. | High security/operations. | Verify role source and recent policy changes; inspect safe audit evidence. | Security and Product. |

## Logging Privacy And Retention

| Signal/data | Allowed content | Prohibited content | Retention/access | Redaction |
| --- | --- | --- | --- | --- |
| Scanner/client signal | Session/correlation id, event when authorized, capability, safe result. | Camera frame, QR payload, raw token, attendee identity on failure. | Operational policy; restricted diagnostics. | Allowlists at client and ingestion. |
| Server validation/audit | Authorized organizer, event, safe reason, latency; attendee only for valid authorized audit. | Raw token/hash and unnecessary profile fields. | Security baseline and audit policy. | Structured serializer and log filter. |

## Requirements

| ID | Requirement | Source | Acceptance criteria | Dependencies |
| --- | --- | --- | --- | --- |
| REQ-108 | Scanner and server validation must emit correlated, actionable, privacy-safe readiness, attempt, result, latency, denial, and duplicate signals with owned diagnosis and alert paths. | [Security Baseline](../../../../../../../../../knowledge/conventions/security-baseline.md) | AC-108 | REQ-104, REQ-106 |
