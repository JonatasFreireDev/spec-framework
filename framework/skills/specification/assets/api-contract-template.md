# API Contract: [use case]

## Generation And Agent Self-Check

| Field | Value |
| --- | --- |
| Generated on | `YYYY-MM-DD` |
| Purpose | Define stable interface operations and failure semantics. |
| Required inputs and evidence | `[behavior, technical landscape, API standards, decisions, existing code]` |
| Ready when | Operations, authorization, schemas, errors, idempotency, and compatibility are explicit. |
| Current status | `draft` |

## Snapshot

| Field | Value |
| --- | --- |
| Status | `draft` |
| Rationale | `[required when not_applicable]` |
| Source specification | `[SPEC-XXX](../specification.md)` |
| Contract version | `2` |

## Operations

| Operation | Kind/path | Consumer | Side effect | Owner |
| --- | --- | --- | --- | --- |
| `[name]` | `[command/query/endpoint/event]` | `[consumer]` | `[effect or None]` | `[owner]` |

## Authorization

| Operation | Authentication | Permission/tenant rule | Enforcement point | Denial response |
| --- | --- | --- | --- | --- |
| `[operation]` | `[mechanism]` | `[rule]` | `[server boundary]` | `[response]` |

## Request And Response Schemas

| Operation | Request fields and constraints | Success response | Sensitive fields |
| --- | --- | --- | --- |
| `[operation]` | `[schema]` | `[schema]` | `[classification/exposure]` |

## Errors Idempotency And Compatibility

| Operation | Error/status | Retryable | Idempotency/concurrency | Versioning compatibility |
| --- | --- | --- | --- | --- |
| `[operation]` | `[error]` | `[yes/no]` | `[rule/key]` | `[policy]` |

## Requirements

| ID | Requirement | Source | Acceptance criteria | Dependencies |
| --- | --- | --- | --- | --- |
| [`REQ-001`](#req-001) | `[testable interface contract]` | `[link]` | [`AC-001`](../tests.md#ac-001) | `[links or None]` |

## Agent Verification Checklist

- [ ] Every operation has an owner, consumer, authorization rule, and explicit errors.
- [ ] Schemas distinguish required, optional, sensitive, and constrained fields.
- [ ] Retry, idempotency, concurrency, and compatibility behavior agree with the behavior contract.
