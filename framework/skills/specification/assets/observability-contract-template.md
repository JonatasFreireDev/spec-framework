# Observability Contract: [use case]

## Generation And Agent Self-Check

| Field | Value |
| --- | --- |
| Generated on | `YYYY-MM-DD` |
| Purpose | Define privacy-safe signals that prove behavior and support diagnosis. |
| Required inputs and evidence | `[behavior, API, data, operations baseline, service expectations]` |
| Ready when | Signals, correlation, metrics, expectations, alerts, diagnosis, privacy, and retention are explicit. |
| Current status | `draft` |

## Snapshot

| Field | Value |
| --- | --- |
| Status | `draft` |
| Rationale | `[required when not_applicable]` |
| Source specification | `[SPEC-XXX](../specification.md)` |
| Contract version | `2` |

## Signals And Correlation

| Signal | Trigger | Fields/dimensions | Correlation | Consumer |
| --- | --- | --- | --- | --- |
| `[event/log/trace]` | `[condition]` | `[safe fields]` | `[request/trace/business id]` | `[team/system]` |

## Metrics And Service Expectations

| Metric | Definition | Target/threshold | Window | Decision supported |
| --- | --- | --- | --- | --- |
| `[metric]` | `[formula]` | `[target]` | `[window]` | `[action]` |

## Alerts And Diagnosis

| Alert/symptom | Condition | Severity | Runbook/diagnostic path | Owner |
| --- | --- | --- | --- | --- |
| `[alert]` | `[condition]` | `[severity]` | `[link]` | `[owner]` |

## Logging Privacy And Retention

| Signal/data | Allowed content | Prohibited content | Retention/access | Redaction |
| --- | --- | --- | --- | --- |
| `[signal]` | `[fields]` | `[PII/token/secret]` | `[policy]` | `[control]` |

## Requirements

| ID | Requirement | Source | Acceptance criteria | Dependencies |
| --- | --- | --- | --- | --- |
| [`REQ-001`](#req-001) | `[testable observability contract]` | `[link]` | [`AC-001`](../tests.md#ac-001) | `[links or None]` |

## Agent Verification Checklist

- [ ] Success, failure, latency, and abuse signals support defined decisions.
- [ ] Alerts are actionable and route to a real diagnostic owner.
- [ ] Logs and telemetry exclude or protect sensitive data.
