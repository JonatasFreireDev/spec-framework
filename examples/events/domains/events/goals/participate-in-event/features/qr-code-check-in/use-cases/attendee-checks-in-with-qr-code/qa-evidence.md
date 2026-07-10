# QA Evidence: Attendee Checks In With QR Code

## Snapshot

| Field | Value |
| --- | --- |
| ID | QA-001 |
| Status | draft |
| Source use case | UC-001 |
| Source specification | SPEC-001 |
| Source tests | UC-001:tests |
| Owner skill | QA AI |
| Next skill | Security Review AI |

## Navigation

| Artifact | Link |
| --- | --- |
| Context | [context.md](context.md) |
| Specification | [specification.md](specification.md) |
| Implementation Plan | [implementation-plan.md](implementation-plan.md) |
| Execution Graph | [execution-graph.json](execution-graph.json) |
| Tasks Index | [tasks.md](tasks.md) |
| Tests | [tests.md](tests.md) |
| Security Review | [security-review.md](security-review.md) |
| Audit | [audit.md](audit.md) |

## Code Traceability

| Task | Branch | Commits | PR | Code Paths |
| --- | --- | --- | --- | --- |
| TK-001 | N/A until implementation | N/A until implementation | N/A until implementation | supabase/migrations |
| TK-002 | N/A until implementation | N/A until implementation | N/A until implementation | mobile services/actions, server functions |
| TK-003 | N/A until implementation | N/A until implementation | N/A until implementation | mobile services/actions, server functions |
| TK-004 | N/A until implementation | N/A until implementation | N/A until implementation | mobile event screens |
| TK-005 | N/A until implementation | N/A until implementation | N/A until implementation | mobile organizer event screens |
| TK-006 | N/A until implementation | N/A until implementation | N/A until implementation | tests, analytics instrumentation, product audit notes |

## Gate Evidence

| Field | Value |
| --- | --- |
| Test command | N/A until implementation |
| Gate logs | N/A until validation |
| CI URL | N/A until validation |
| Screenshots | N/A until validation |
| Environment | N/A until validation |

## Acceptance Evidence Matrix

| Acceptance Criterion | Source | Validation Method | Evidence | Result |
| --- | --- | --- | --- | --- |
| Attendee can generate a non-PII QR token for an event they joined. | [specification.md](specification.md) | Integration test | Pending implementation evidence | not run |
| Organizer can validate a valid token for an event they manage. | [tests.md](tests.md) | Integration or E2E test | Pending implementation evidence | not run |
| Invalid, expired, wrong-event, and permission-denied tokens are rejected. | [tests.md](tests.md) | Security and integration tests | Pending implementation evidence | not run |
| Duplicate scans are idempotent. | [tests.md](tests.md) | Integration test | Pending implementation evidence | not run |
| Analytics and logs are emitted for success and failure. | [analytics.md](analytics.md) | Event/log assertion | Pending implementation evidence | not run |

## Test Execution

| Test | Type | Command Or Method | Evidence | Result |
| --- | --- | --- | --- | --- |
| T-001 attendee generates QR | integration | Planned app test | Pending | not run |
| T-002 non-attendee denied | security | Planned API/integration test | Pending | not run |
| T-003 organizer validates valid QR | integration | Planned app test | Pending | not run |
| T-004 duplicate scan idempotency | integration | Planned database/service test | Pending | not run |
| T-005 expired token rejection | integration | Planned API/integration test | Pending | not run |
| T-006 organizer without permission denied | security | Planned API/integration test | Pending | not run |
| T-007 analytics emitted | observability | Planned event/log review | Pending | not run |

## Security And Privacy Evidence

| Control | Evidence | Result | Notes |
| --- | --- | --- | --- |
| Authorization | Pending API/security test | not run | Organizer permission must be server-authoritative. |
| Data privacy | Pending QR payload and log review | not run | QR payload must not expose raw PII. |
| Abuse/edge cases | Pending replay/expired/duplicate tests | not run | Token replay and duplicate scans must remain safe. |
| Safe logging/analytics | Pending analytics and audit log review | not run | Failure logs must avoid sensitive attendee data. |

## Defects And Fix Verification

| Finding | Severity | Fix Evidence | Status |
| --- | --- | --- | --- |
| No implementation evidence exists yet. | blocker | Requires code, tests, and QA run | open |

## Residual Risk

| Risk | Why It Remains | Mitigation | Approval |
| --- | --- | --- | --- |
| Offline check-in remains unsupported. | Offline mutation is out of scope for this use case. | Keep online-only behavior explicit and test network errors. | N/A |

## QA Verdict

| Field | Value |
| --- | --- |
| Verdict | blocked |
| Coverage complete | no |
| Security evidence complete | no |
| Blocks validation | yes |
| Blocks release | yes |
| Next owner | Code Runner AI, then QA AI |
