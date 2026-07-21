# Security Contract: [use case]

## Generation And Agent Self-Check

| Field | Value |
| --- | --- |
| Generated on | `YYYY-MM-DD` |
| Purpose | Define delivery-specific threats, controls, audit obligations, and residual risk. |
| Required inputs and evidence | `[security baseline, threat register, API/data contracts, decisions]` |
| Ready when | Assets, boundaries, threats, abuse, controls, audit, and residual risks are explicit. |
| Current status | `draft` |

## Snapshot

| Field | Value |
| --- | --- |
| Status | `draft` |
| Rationale | `[required when not_applicable]` |
| Source specification | `[SPEC-XXX](../specification.md)` |
| Contract version | `2` |

## Assets And Trust Boundaries

| Asset | Sensitivity | Boundary crossed | Actor/component | Security objective |
| --- | --- | --- | --- | --- |
| `[asset]` | `[classification]` | `[boundary]` | `[actor]` | `[confidentiality/integrity/availability]` |

## Threats And Abuse Cases

| Threat/abuse | Preconditions | Impact | Existing mitigation | Gap owner |
| --- | --- | --- | --- | --- |
| `[threat]` | `[conditions]` | `[impact]` | `[control]` | `[owner or None]` |

## Controls And Audit

| Control | Enforcement point | Deny/fail behavior | Audit evidence | Verification |
| --- | --- | --- | --- | --- |
| `[control]` | `[boundary]` | `[behavior]` | `[safe log/event]` | `[test/review]` |

## Residual Risk

| Risk | Likelihood/impact | Acceptance owner | Decision/evidence | Expiry/review trigger |
| --- | --- | --- | --- | --- |
| `[risk or None]` | `[rating]` | `[human owner]` | `[link]` | `[trigger]` |

## Requirements

| ID | Requirement | Source | Acceptance criteria | Dependencies |
| --- | --- | --- | --- | --- |
| [`REQ-001`](#req-001) | `[testable security control]` | `[link]` | [`AC-001`](../tests.md#ac-001) | `[links or None]` |

## Agent Verification Checklist

- [ ] Authorization, privacy, abuse, tokens, logs, and sensitive data are considered.
- [ ] Controls fail safely and have evidence without leaking protected data.
- [ ] Residual risk requires an explicit human owner and decision where material.
