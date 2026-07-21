# Data Contract: Organizer Validates QR Code

| Field | Value |
| --- | --- |
| Status | draft |
| Source specification | [SPEC-002](../specification.md) |
| Contract version | 2 |

## Ownership And Entities

| Entity/store | Owner | Fields/relationships | Writer/readers | Source of truth |
| --- | --- | --- | --- | --- |
| event_attendance | Events domain. | event_id, attendee_user_id, status, checked_in_at, checked_in_by_user_id. | Validation writes; authorized product/operations read. | Events database. |
| check_in_audit_log | Events/Security operations. | event, organizer, nullable attendee, result, safe reason, created_at. | Validation writes; restricted audit consumers read. | Audit store. |

## Constraints And Invariants

| Constraint | Enforcement | Concurrency/consistency | Failure behavior |
| --- | --- | --- | --- |
| One event_attendance row per event and attendee; checked_in_at first-write wins. | Unique constraint and transactional update/readback. | Concurrent organizers converge on one authoritative record. | Return already_checked_in with existing timestamp. |
| Invalid proof does not require resolved attendee identity in audit. | Nullable attendee field and safe reason enum. | Rejected attempts cannot create attendance. | Store safe event/actor context only when authorized. |

## Lifecycle Migration And Retention

| Data | Create/update/delete lifecycle | Migration/backfill | Retention/deletion | Rollback impact |
| --- | --- | --- | --- | --- |
| Attendance check-in fields | Added before validation activation; set once on success. | No fabricated historical check-ins. | Attendance policy applies. | Disable writers; retain valid records. |
| Audit records | Created per authorized attempt/result. | No historical synthesis. | Security/operations retention policy. | Preserve incident evidence through rollback. |

## Privacy And Classification

| Data | Classification | Purpose/legal basis | Exposure | Encryption/masking |
| --- | --- | --- | --- | --- |
| Attendance and attendee identity | Personal/internal. | Event attendance operations. | Attendee and authorized event staff. | Protected storage/transport; omitted from rejection output/logs. |
| Organizer identity and audit reason | Internal/security evidence. | Accountability and abuse diagnosis. | Restricted support/security roles. | Structured safe reason; no raw token. |

## Requirements

| ID | Requirement | Source | Acceptance criteria | Dependencies |
| --- | --- | --- | --- | --- |
| REQ-105 | Validation data must enforce unique first-write attendance, privacy-safe nullable audit identity, governed retention, no synthetic backfill, and attendance-preserving rollback. | [Use Case](../use-case.md) | AC-105 | REQ-104 |
