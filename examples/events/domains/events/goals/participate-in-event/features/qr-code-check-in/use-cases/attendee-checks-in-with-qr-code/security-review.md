# Security Review: Attendee Checks In With QR Code

## Snapshot

| Field | Value |
| --- | --- |
| ID | SEC-001 |
| Status | draft |
| Source use case | UC-001 |
| Source specification | SPEC-001 |
| Source QA evidence | QA-001 |
| Owner skill | Security Review AI |
| Next skill | QA AI |

## Navigation

| Artifact | Link |
| --- | --- |
| Context | [context.md](context.md) |
| Specification | [specification.md](specification.md) |
| Design | [design.md](design.md) |
| Implementation Plan | [implementation-plan.md](implementation-plan.md) |
| Execution Graph | [execution-graph.json](execution-graph.json) |
| Tasks | [tasks.md](tasks.md) |
| Tests | [tests.md](tests.md) |
| QA Evidence | [qa-evidence.md](qa-evidence.md) |
| Audit | [audit.md](audit.md) |

## Delivery

| Field | Value |
| --- | --- |
| Level | L1 |
| Priority | P0 |
| Depends on | SPEC-001, UC-001:tests, QA-001, DEC-001, DEC-002, DEC-008 |
| Rationale | Tier L is required because the flow uses authentication, organizer permissions, opaque tokens, and attendance writes. |

## Security Scope

| Area | In Scope | Out Of Scope |
| --- | --- | --- |
| Authentication | Attendee and organizer sessions. | Anonymous check-in. |
| Authorization | Attendee can generate own QR; organizer can validate managed event. | Offline organizer authorization. |
| Data and privacy | QR payload, token storage, attendance status, logs, analytics. | Payment or ticket ownership. |
| Abuse prevention | Token replay, duplicate scans, forged QR, wrong-event scan. | Advanced fraud scoring. |
| Observability | Safe failure logs and analytics events. | Full production alert tuning. |

## Threat Model Summary

| Threat | Actor | Impact | Required Control | Evidence |
| --- | --- | --- | --- | --- |
| Forged QR token | Attacker | Unauthorized check-in attempt | Server-side opaque token validation | Pending implementation evidence |
| Organizer without permission validates | Unauthorized user | Attendance state tampering | Server-authoritative organizer permission check | Pending security test |
| QR payload leaks PII | Any scanner | Privacy exposure | Opaque non-PII payload | Pending payload review |
| Token replay | Attendee or observer | Duplicate or stale check-in attempt | Expiration, consumed state, idempotency | Pending integration test |
| Unsafe logs reveal attendee data | Operator or log reader | Privacy exposure | Safe enumerable failure reasons | Pending log review |

## Control Checklist

| Control | Expected Evidence | Result | Notes |
| --- | --- | --- | --- |
| Server-side authorization | API/security test | blocked | No implementation evidence yet. |
| Least privilege | Role matrix and permission tests | blocked | Organizer roles need implementation confirmation. |
| Sensitive data minimization | QR payload/log review | blocked | Must prove no raw PII in payload. |
| Input validation | Invalid token tests | blocked | Must reject malformed and wrong-event payloads. |
| Abuse/replay/rate limits | Replay and duplicate tests | blocked | Idempotency and token expiry must be demonstrated. |
| Secrets and tokens | Token storage/review evidence | blocked | Opaque token strategy is approved but not implemented. |
| Safe logging and analytics | Log/event review | blocked | Must avoid sensitive attendee fields. |
| Rollback and monitoring | Plan/runbook evidence | blocked | Feature flag exists in plan, not validated. |

## Findings

| Severity | Finding | Evidence | Required Fix | Owner |
| --- | --- | --- | --- | --- |
| blocker | No implementation or QA evidence exists yet. | [qa-evidence.md](qa-evidence.md) | Implement and run the required security and QA checks. | Code Runner AI + QA AI |

## Residual Risks

| Risk | Severity | Mitigation | Approval Needed | Owner |
| --- | --- | --- | --- | --- |
| Online-only check-in can fail under poor connectivity. | medium | Keep offline validation out of scope and show retryable network state. | no | Product + UX |

## Security Verdict

| Field | Value |
| --- | --- |
| Verdict | blocked |
| Blocks validation | yes |
| Blocks release | yes |
| Required decisions | N/A |
| Next owner | Code Runner AI, then QA AI |
