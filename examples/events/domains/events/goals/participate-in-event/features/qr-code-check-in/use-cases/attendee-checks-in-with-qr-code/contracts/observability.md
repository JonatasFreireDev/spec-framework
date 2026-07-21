# Observability Contract: Attendee Checks In With QR Code

| Field | Value |
| --- | --- |
| Status | draft |
| Source specification | [SPEC-001](../specification.md) |
| Contract version | 2 |

## Signals And Correlation

| Signal | Trigger | Fields/dimensions | Correlation | Consumer |
| --- | --- | --- | --- | --- |
| qr_check_in_generated | Proof issued or refreshed. | event_id, outcome, expiry bucket; no token or attendee PII. | request_id and event_id. | Product and Operations. |
| qr_check_in_validated | Authorized success or existing check-in. | event_id, result, latency bucket. | validation_id and request_id. | Product and Operations. |
| qr_check_in_failed | Expired, invalid, wrong event, denial, or service failure. | safe reason, event_id when authorized, latency bucket. | validation_id and request_id. | Operations and Security. |

## Metrics And Service Expectations

| Metric | Definition | Target/threshold | Window | Decision supported |
| --- | --- | --- | --- | --- |
| Validation success rate | Successful or already-checked-in results divided by authorized attempts. | Baseline during pilot; no invented production target. | Five-minute and event windows. | Pause rollout on abnormal degradation. |
| Validation latency | Server duration by result. | Under two seconds in normal documented expectation. | Five-minute rolling view. | Investigate venue queue impact. |
| Invalid/denied spike | Count by safe reason. | Alert on statistically abnormal pilot baseline. | Five-minute window. | Detect abuse or permission regression. |

## Alerts And Diagnosis

| Alert/symptom | Condition | Severity | Runbook/diagnostic path | Owner |
| --- | --- | --- | --- | --- |
| Validation failures spike. | Failure ratio exceeds approved pilot baseline. | High during live event. | Check service health, token expiry, permission denials, and recent rollout. | Events operations. |
| Duplicate results spike. | Duplicate ratio changes materially. | Medium. | Inspect scanner retries and concurrency without rewriting attendance. | Events engineering. |

## Logging Privacy And Retention

| Signal/data | Allowed content | Prohibited content | Retention/access | Redaction |
| --- | --- | --- | --- | --- |
| Validation log | request/validation id, safe reason, event id when authorized, latency. | Raw token, token hash, attendee name, email, or QR payload. | Security baseline and operational access policy. | Structured allowlist before emission. |
| Audit record | Authorized actor, event, result, timestamp; attendee only for valid authorized resolution. | Secret material and unnecessary profile data. | Attendance audit policy. | Field-level allowlist. |

## Requirements

| ID | Requirement | Source | Acceptance criteria | Dependencies |
| --- | --- | --- | --- | --- |
| REQ-008 | Generation and validation must emit correlated, actionable, privacy-safe signals for success, failure, latency, denial, replay, and rollout decisions. | [Security Baseline](../../../../../../../../../knowledge/conventions/security-baseline.md) | AC-008 | REQ-004, REQ-006 |
