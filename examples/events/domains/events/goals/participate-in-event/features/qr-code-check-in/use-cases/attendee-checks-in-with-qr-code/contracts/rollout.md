# Rollout Contract: Attendee Checks In With QR Code

| Field | Value |
| --- | --- |
| Status | draft |
| Source specification | [SPEC-001](../specification.md) |
| Contract version | 2 |

## Release Strategy

| Stage/audience | Activation | Entry condition | Exit condition | Owner |
| --- | --- | --- | --- | --- |
| Internal fixture and test events. | event_qr_check_in flag for selected events. | Schema, permissions, tests, signals, and rollback verified. | No critical data, authorization, privacy, or accessibility finding. | Product, Engineering, QA, Security. |
| Limited live-event pilot. | Event allowlist plus flag. | Internal evidence accepted and support owner present. | Stable success/failure signals and accepted venue feedback. | Product and Operations. |

## Compatibility Migration And Backfill

| Change | Compatibility window | Migration/backfill | Validation | Cleanup trigger |
| --- | --- | --- | --- | --- |
| Attendance check-in fields and uniqueness. | Deploy nullable fields before writers; old readers ignore them. | No historical check-in backfill. | Existing participation regression plus concurrent write test. | After all supported readers tolerate fields. |
| Token storage and cleanup. | Create before generation API activation. | No token backfill. | Expiry, hashing, prune, and rollback tests. | Scheduled cleanup proven in pilot. |

## Monitoring And Decision Points

| Signal | Expected range | Pause/rollback threshold | Decision owner | Observation window |
| --- | --- | --- | --- | --- |
| Authorization/data integrity | No unauthorized or duplicate mutation. | Any confirmed violation. | Security and Product. | Continuous during pilot. |
| Validation reliability | Establish evidence-backed pilot baseline. | Sustained abnormal failure or latency that harms venue flow. | Operations. | Per event and five-minute windows. |

## Rollback And Recovery

| Failure mode | Rollback action | Data recovery | Maximum tolerated impact | Verification |
| --- | --- | --- | --- | --- |
| Authorization, privacy, or integrity defect. | Disable flag and validation endpoint; preserve audit evidence. | Review affected attendance; never bulk-rewrite without approved remediation. | No continued exposure or mutation after disable. | Flag/endpoint check and incident audit. |
| Reliability degradation. | Stop pilot and return event UI to unavailable guidance. | Keep valid attendance; invalidate active tokens if needed. | Current pilot event only. | Signal recovery and support confirmation. |

## Requirements

| ID | Requirement | Source | Acceptance criteria | Dependencies |
| --- | --- | --- | --- | --- |
| REQ-009 | Activation must be staged and observable with compatible schema order, explicit pause thresholds, token invalidation, attendance-preserving rollback, and owned recovery. | [Use Case](../use-case.md) | AC-009 | REQ-005, REQ-006, REQ-008 |
