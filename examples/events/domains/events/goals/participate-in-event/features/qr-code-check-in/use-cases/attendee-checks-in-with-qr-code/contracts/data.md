# Data Contract: Attendee Checks In With QR Code

| Field | Value |
| --- | --- |
| Status | draft |
| Source specification | [SPEC-001](../specification.md) |
| Contract version | 2 |

## Ownership And Entities

| Entity/store | Owner | Fields/relationships | Writer/readers | Source of truth |
| --- | --- | --- | --- | --- |
| event_attendance | Events domain | event_id, attendee_id, checked_in_at, checked_in_by_user_id. | Validation service writes; authorized product flows read. | Events database. |
| event_check_in_tokens | Events domain | token_hash, event_id, attendee_id, expires_at, consumed_at. | Token service writes; validation service reads/consumes. | Events database. |

## Constraints And Invariants

| Constraint | Enforcement | Concurrency/consistency | Failure behavior |
| --- | --- | --- | --- |
| event_id plus attendee_id is unique. | Database unique constraint and transaction. | Concurrent validations converge on one row and timestamp. | Return already_checked_in after conflict/readback. |
| Token hash is unique, opaque, event-scoped, attendee-scoped, and expiring. | Service validation and indexed database fields. | Consumption and attendance write occur atomically where supported. | Invalid or expired result without identity disclosure. |

## Lifecycle Migration And Retention

| Data | Create/update/delete lifecycle | Migration/backfill | Retention/deletion | Rollback impact |
| --- | --- | --- | --- | --- |
| Attendance fields | Created on join; checked_in_at set once. | Add nullable check-in fields and unique constraint before activation; no historical backfill. | Product attendance retention policy applies. | Disabling feature retains valid attendance history. |
| Token rows | Created/rotated on request; consumed or expired. | Create token storage before API activation. | Prune expired/consumed tokens after audit window. | Tokens can be invalidated without deleting attendance. |

## Privacy And Classification

| Data | Classification | Purpose/legal basis | Exposure | Encryption/masking |
| --- | --- | --- | --- | --- |
| attendee_id and attendance | Personal/internal. | Event participation and operations. | Attendee self and authorized event staff. | Protected at rest/in transit; excluded from rejection logs. |
| token_hash | Secret credential material. | Short-lived attendance proof. | Server components only. | Store hash, never raw token; redact telemetry. |

## Requirements

| ID | Requirement | Source | Acceptance criteria | Dependencies |
| --- | --- | --- | --- | --- |
| REQ-005 | Attendance and token stores must enforce ownership, uniqueness, expiry, atomic/idempotent mutation, privacy, cleanup, and reversible migration boundaries. | [DEC-002](../../../../../../../../../knowledge/decisions/DEC-002-qr-token-strategy.md) | AC-005 | REQ-004 |
