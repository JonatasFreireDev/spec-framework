# Quality Contract: Attendee Checks In With QR Code

| Field | Value |
| --- | --- |
| Status | draft |
| Source specification | [SPEC-001](../specification.md) |
| Contract version | 2 |

## Quality Risks

| Risk | Impact/likelihood | Preventive contract | Detection method | Owner |
| --- | --- | --- | --- | --- |
| Unauthorized or duplicate attendance write. | High/medium. | REQ-002, REQ-005, REQ-006. | Permission, transaction, and concurrent integration tests. | Engineering and QA. |
| Inaccessible or misleading result state. | Medium/medium. | REQ-003. | Accessibility automation and assistive manual review. | Design and QA. |
| Missing failure evidence or unsafe telemetry. | High/medium. | REQ-006, REQ-008. | Event/log schema review and negative tests. | Security and QA. |

## Acceptance Traceability

| Acceptance criterion | Requirement | Risk | Test/evidence method | Expected evidence |
| --- | --- | --- | --- | --- |
| AC-001 | REQ-001 | Wrong product outcome. | TEST-001 end-to-end attendee-to-organizer flow. | One authorized attendance record. |
| AC-002 | REQ-002 | State or concurrency error. | TEST-002 behavior table and rapid duplicate scans. | Deterministic status and unchanged timestamp. |
| AC-003 | REQ-003 | Inaccessible or privacy-leaking UX. | TEST-003 accessibility and content review. | Focus, labels, text-equivalent, safe copy. |
| AC-004 | REQ-004 | Unsafe interface boundary. | TEST-004 API auth/schema/error tests. | Contract responses and denials. |
| AC-005 | REQ-005 | Corrupt, retained, or exposed data. | TEST-005 migration, constraint, cleanup, privacy tests. | Schema and lifecycle evidence. |
| AC-006 | REQ-006 | Token or authorization abuse. | TEST-006 security negative/replay tests. | Controls deny safely. |
| AC-007 | REQ-007 | Unverified critical behavior. | TEST-007 gate coverage review. | All required evidence destinations assigned. |
| AC-008 | REQ-008 | Undiagnosable venue failure. | TEST-008 telemetry and alert review. | Safe correlated signals. |
| AC-009 | REQ-009 | Irreversible rollout. | TEST-009 migration/flag/rollback rehearsal. | Recorded activation and recovery result. |

## Test Levels And Environments

| Coverage | Level | Environment/platform | Test data | Isolation/dependencies |
| --- | --- | --- | --- | --- |
| Rules, expiry, and state transitions. | Unit and integration. | documentation-fixture until runtime exists. | Synthetic event, users, and tokens. | Controlled clock and isolated database. |
| Authorization, concurrency, and lifecycle. | Integration and security. | Configured server/database environment. | Authorized and denied organizers; concurrent scans. | Real constraints and transaction semantics. |
| QR and result states. | End-to-end and accessibility. | Web and mobile-camera targets. | Synthetic proof only. | Camera capability or reviewed simulator. |

## Evidence And Exit Conditions

| Gate | Pass condition | Evidence owner | Failure route |
| --- | --- | --- | --- |
| Specification conformance | REQ-001..REQ-009 map to AC-001..AC-009 and TEST-001..TEST-009. | Specification and QA. | Return gaps to owning contract. |
| Runtime QA | Applicable tests pass in configured environments with no fabricated evidence. | QA. | Keep delivery unvalidated. |
| Security and accessibility | No blocker; residual risk has owner; visual evidence is reviewed. | Security Review and QA. | Return to implementation or product decision. |

## Requirements

| ID | Requirement | Source | Acceptance criteria | Dependencies |
| --- | --- | --- | --- | --- |
| REQ-007 | Verification must cover every contract risk with configured levels, environments, data, owners, objective exit conditions, and concrete evidence. | [Quality System](../../../../../../../../../engineering/quality/quality-system.md) | AC-007 | REQ-001, REQ-002, REQ-003, REQ-004, REQ-005, REQ-006 |
