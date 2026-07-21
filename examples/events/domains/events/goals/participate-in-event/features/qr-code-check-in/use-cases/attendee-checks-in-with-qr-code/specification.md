# Specification: Attendee Checks In With QR Code

## Snapshot

| Field | Value |
| --- | --- |
| ID | SPEC-001 |
| Status | draft |
| Source use case | [UC-001](use-case.md) |
| Source feature | [FT-001](../../feature.md) |
| Contract version | 2 |
| Delivery | L1 / P0 |

## Contract Applicability

| Contract | Applies | Status | Rationale |
| --- | --- | --- | --- |
| Product | yes | draft | QR proof closes the attendee side of the L1 attendance outcome. |
| Behavior | yes | draft | Token generation, refresh, validation, and duplicate handling are stateful. |
| UX | yes | draft | Attendee and organizer require visible, accessible states. |
| API | yes | draft | Generation, status, and validation cross a server boundary. |
| Data | yes | draft | Opaque tokens and attendance state require persistence and cleanup. |
| Security | yes | draft | The flow crosses identity, permission, token, and privacy boundaries. |
| Quality | yes | draft | Permission and state mutation risks require explicit verification. |
| Observability | yes | draft | Success, expiry, denial, and replay need privacy-safe diagnosis. |
| Rollout | yes | draft | Schema and flag changes require reversible activation. |

## Evidence And Boundary

| Kind | Evidence or statement | Source | Confidence/decision status |
| --- | --- | --- | --- |
| Approved decision | QR tokens expire after five minutes. | [DEC-001](../../../../../../../../knowledge/decisions/DEC-001-qr-expiration-duration.md) | approved |
| Approved decision | V1 uses opaque, server-stored, event- and attendee-scoped tokens. | [DEC-002](../../../../../../../../knowledge/decisions/DEC-002-qr-token-strategy.md) | approved |
| Scope | Generate and present attendee proof; validate it server-side and record one check-in. | [Use Case](use-case.md) | draft product scope |
| Non-goal | Offline validation, payment/ticketing, and advanced fraud scoring. | [Feature](../../feature.md) | draft product scope |

## Cross-Contract Synthesis

| Concern | Implementable outcome | Contract | Blocking dependency |
| --- | --- | --- | --- |
| Product and behavior | A joined attendee receives expiring proof; authorized validation records attendance exactly once. | [Product](contracts/product.md), [Behavior](contracts/behavior.md) | Organizer role policy remains a downstream decision. |
| Experience and interfaces | QR generation, refresh, validation, and failure states are explicit without revealing identity on rejection. | [UX](contracts/ux.md), [API](contracts/api.md) | None for the documented draft. |
| Data and trust | Token material is opaque and prunable; attendance uniqueness and server authorization protect the mutation. | [Data](contracts/data.md), [Security](contracts/security.md) | Event permission source must be selected before implementation. |
| Quality and operations | Risk-based tests, privacy-safe signals, staged activation, and rollback cover the L1 path. | [Quality](contracts/quality.md), [Observability](contracts/observability.md), [Rollout](contracts/rollout.md) | Runtime stack and environments are not configured. |

## Traceability Summary

| Requirement range | Acceptance range | Source contracts | Test/evidence destination |
| --- | --- | --- | --- |
| REQ-001 through REQ-009 | AC-001 through AC-009 | [contracts](contracts/) | [Tests](tests.md), [QA Evidence](qa-evidence.md) |

## Adversarial Review

| Check | Result | Evidence or routed correction |
| --- | --- | --- |
| Contradictions and duplicated requirements | passed | Each REQ has one concern owner; shared rules are linked rather than copied as requirements. |
| Alternate, error, edge, and abuse coverage | passed for draft | Expiry, invalid token, wrong event, denial, network failure, replay, and concurrent scans are covered. |
| Unsafe assumptions and missing decisions | blocked for implementation | Organizer role source and runtime stack remain explicit downstream blockers. |
| Cross-contract terminology and ownership | passed | Token, attendance, organizer permission, and check-in states use consistent names. |

## Open Questions And Decisions

| Question/Decision | Owner | Blocks |
| --- | --- | --- |
| Which event role source grants organizer check-in permission? | Product and Engineering | Engineering Proposal and implementation |
| Which runtime stack and deployment environment will host the flow? | Engineering | Technical planning and executable tests |

## Approval

| Field | Value |
| --- | --- |
| Approved by |  |
| Date |  |
| Notes | Migrated to Specification Contract v2 as draft; no approval created. |
