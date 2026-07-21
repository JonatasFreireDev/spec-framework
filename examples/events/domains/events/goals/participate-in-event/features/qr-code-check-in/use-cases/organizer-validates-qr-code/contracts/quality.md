# Quality Contract: Organizer Validates QR Code

| Field | Value |
| --- | --- |
| Status | draft |
| Source specification | [SPEC-002](../specification.md) |
| Contract version | 2 |

## Quality Risks

| Risk | Impact/likelihood | Preventive contract | Detection method | Owner |
| --- | --- | --- | --- | --- |
| Unauthorized or duplicate check-in. | High/medium. | REQ-102, REQ-104, REQ-105, REQ-106. | Permission, transaction, retry, and concurrency tests. | QA and Security. |
| Scanner result is inaccessible or operationally slow. | High during venue/medium. | REQ-103, REQ-108. | Accessibility and latency/queue scenario review. | Design, QA, Operations. |
| Failure leaks attendee or token data. | High/medium. | REQ-104, REQ-106, REQ-108. | Negative API, UI content, and telemetry schema tests. | Security and QA. |

## Acceptance Traceability

| Acceptance criterion | Requirement | Risk | Test/evidence method | Expected evidence |
| --- | --- | --- | --- | --- |
| AC-101 | REQ-101 | Wrong product outcome. | TEST-101 managed-event end-to-end scan. | Exactly one authorized check-in. |
| AC-102 | REQ-102 | Incorrect state/failure behavior. | TEST-102 behavior table and concurrent scans. | Deterministic outcomes and no duplicate write. |
| AC-103 | REQ-103 | Inaccessible or unsafe scanner UX. | TEST-103 camera, keyboard, screen-reader, content checks. | Reviewed states and focus. |
| AC-104 | REQ-104 | Unsafe API boundary. | TEST-104 auth/schema/error/idempotency integration tests. | Safe responses and current authorization. |
| AC-105 | REQ-105 | Corrupt/private data. | TEST-105 constraint, audit, migration, retention tests. | Data lifecycle evidence. |
| AC-106 | REQ-106 | Permission/token/privacy abuse. | TEST-106 security negative and replay suite. | Controls deny safely. |
| AC-107 | REQ-107 | Incomplete evidence. | TEST-107 coverage/gate audit. | All risks have owner and evidence target. |
| AC-108 | REQ-108 | Undiagnosable venue incident. | TEST-108 signal/privacy/alert review. | Actionable safe telemetry. |
| AC-109 | REQ-109 | Unsafe activation/recovery. | TEST-109 flag, migration, pilot, rollback rehearsal. | Owned release and recovery record. |

## Test Levels And Environments

| Coverage | Level | Environment/platform | Test data | Isolation/dependencies |
| --- | --- | --- | --- | --- |
| State, authorization, API, data, concurrency. | Unit, integration, security. | Configured server/database; fixture until available. | Synthetic roles, events, tokens, simultaneous organizers. | Controlled clock, isolated transaction store. |
| Scanner and result experience. | End-to-end, accessibility, manual visual. | Web and mobile-camera. | Synthetic attendees; no production identity. | Real camera or reviewed simulator limitations. |
| Venue operations and rollout. | Scenario/rehearsal. | Internal or low-risk pilot. | Approved pilot event. | Support and rollback owners present. |

## Evidence And Exit Conditions

| Gate | Pass condition | Evidence owner | Failure route |
| --- | --- | --- | --- |
| Specification coverage | REQ-101..REQ-109 map to AC-101..AC-109 and TEST-101..TEST-109. | Specification and QA. | Return to concern owner. |
| QA/accessibility/security | No blocker; same implementation diff; concrete logs/screenshots/results. | QA and Security Review. | Return to implementation or decision owner. |
| Pilot readiness | Permission policy, monitoring, support, pause, and rollback accepted. | Product and Operations. | Keep feature disabled. |

## Requirements

| ID | Requirement | Source | Acceptance criteria | Dependencies |
| --- | --- | --- | --- | --- |
| REQ-107 | Verification must cover product, states, authorization, interfaces, data, accessibility, telemetry, and rollout with configured environments, owners, exit conditions, and concrete evidence. | [Quality System](../../../../../../../../../engineering/quality/quality-system.md) | AC-107 | REQ-101, REQ-102, REQ-103, REQ-104, REQ-105, REQ-106 |
