# Specification: Organizer Validates QR Code

## Snapshot

| Field | Value |
| --- | --- |
| ID | SPEC-002 |
| Status | draft |
| Source use case | [UC-002](use-case.md) |
| Source feature | [FT-001](../../feature.md) |
| Contract version | 2 |
| Delivery | L1 / P0 |

## Contract Applicability

| Contract | Applies | Status | Rationale |
| --- | --- | --- | --- |
| Product | yes | draft | Organizer validation closes the L1 attendance proof. |
| Behavior | yes | draft | Permission, token, check-in-window, and idempotency states govern the result. |
| UX | yes | draft | Scanner, camera, validation, and result states are user-visible. |
| API | yes | draft | Session and validation operations are server-authoritative. |
| Data | yes | draft | Validation writes attendance and audit records. |
| Security | yes | draft | The operation authorizes a sensitive attendance mutation. |
| Quality | yes | draft | Concurrent writes, denial, privacy, and accessibility require evidence. |
| Observability | yes | draft | Venue failures and abuse must be diagnosable without leaking attendee data. |
| Rollout | yes | draft | Permission and data changes need staged, reversible activation. |

## Evidence And Boundary

| Kind | Evidence or statement | Source | Confidence/decision status |
| --- | --- | --- | --- |
| Approved decision | Validation resolves opaque five-minute tokens on the server. | [DEC-001](../../../../../../../../knowledge/decisions/DEC-001-qr-expiration-duration.md), [DEC-002](../../../../../../../../knowledge/decisions/DEC-002-qr-token-strategy.md) | approved |
| Scope | An authenticated organizer validates attendee proof for one managed event and records one check-in. | [Use Case](use-case.md) | proposed use-case scope |
| Non-goal | Offline authority, role-management UI, payments, and manual attendee search. | [Use Case](use-case.md) | proposed use-case scope |
| Unresolved choice | Exact organizer roles and manual camera fallback are not approved. | [Context](context.md) | blocking for implementation |

## Cross-Contract Synthesis

| Concern | Implementable outcome | Contract | Blocking dependency |
| --- | --- | --- | --- |
| Product and behavior | Authorized validation returns checked-in or an explicit non-mutating result. | [Product](contracts/product.md), [Behavior](contracts/behavior.md) | Organizer role decision. |
| Experience and interfaces | Scanner and result states recover from camera, token, permission, and network failures. | [UX](contracts/ux.md), [API](contracts/api.md) | Manual fallback remains out of scope until approved. |
| Data and trust | Unique attendance and privacy-safe audit records make concurrent validation idempotent. | [Data](contracts/data.md), [Security](contracts/security.md) | Permission ownership must be resolved. |
| Quality and operations | Tests and signals cover authorization, concurrency, accessibility, venue latency, rollout, and recovery. | [Quality](contracts/quality.md), [Observability](contracts/observability.md), [Rollout](contracts/rollout.md) | Runtime environments are not configured. |

## Traceability Summary

| Requirement range | Acceptance range | Source contracts | Test/evidence destination |
| --- | --- | --- | --- |
| REQ-101 through REQ-109 | AC-101 through AC-109 | [contracts](contracts/) | [Tests](tests.md), [QA Evidence](qa-evidence.md) |

## Adversarial Review

| Check | Result | Evidence or routed correction |
| --- | --- | --- |
| Contradictions and duplicated requirements | passed | Concern ownership is unique and shared dependencies are linked. |
| Alternate, error, edge, and abuse coverage | passed for draft | Duplicate, expired, invalid, wrong-event, outside-window, denial, camera, and network cases are covered. |
| Unsafe assumptions and missing decisions | blocked for implementation | Organizer roles, permission source, and manual fallback require product direction. |
| Cross-contract terminology and ownership | passed | Scanner session, token validation, attendance, and audit terms are consistent. |

## Open Questions And Decisions

| Question/Decision | Owner | Blocks |
| --- | --- | --- |
| Which organizer roles grant check-in permission and which system owns that mapping? | Product and Engineering | Engineering Proposal and implementation |
| Is a manual attendee-search fallback allowed when the camera is unavailable? | Product and Security | UX finalization and implementation |
| Is online-only validation acceptable for L1 venue operations? | Product and Operations | Release approval |

## Approval

| Field | Value |
| --- | --- |
| Approved by |  |
| Date |  |
| Notes | Returned to draft for explicit Specification Contract v2 migration; no approval created. |
