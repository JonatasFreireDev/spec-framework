# Data Contract: [use case]

## Generation And Agent Self-Check

| Field | Value |
| --- | --- |
| Generated on | `YYYY-MM-DD` |
| Purpose | Define ownership, invariants, lifecycle, migration, and privacy for affected data. |
| Required inputs and evidence | `[behavior, technical catalog, data standards, decisions, current schema]` |
| Ready when | Ownership, entities, constraints, lifecycle, migration, retention, and classification are explicit. |
| Current status | `draft` |

## Snapshot

| Field | Value |
| --- | --- |
| Status | `draft` |
| Rationale | `[required when not_applicable]` |
| Source specification | `[SPEC-XXX](../specification.md)` |
| Contract version | `2` |

## Ownership And Entities

| Entity/store | Owner | Fields/relationships | Writer/readers | Source of truth |
| --- | --- | --- | --- | --- |
| `[entity]` | `[owner]` | `[shape]` | `[actors/components]` | `[system]` |

## Constraints And Invariants

| Constraint | Enforcement | Concurrency/consistency | Failure behavior |
| --- | --- | --- | --- |
| `[constraint]` | `[database/service]` | `[model]` | `[result]` |

## Lifecycle Migration And Retention

| Data | Create/update/delete lifecycle | Migration/backfill | Retention/deletion | Rollback impact |
| --- | --- | --- | --- | --- |
| `[data]` | `[lifecycle]` | `[plan or None]` | `[policy]` | `[impact]` |

## Privacy And Classification

| Data | Classification | Purpose/legal basis | Exposure | Encryption/masking |
| --- | --- | --- | --- | --- |
| `[field/entity]` | `[public/internal/PII/secret]` | `[purpose]` | `[who/where]` | `[control]` |

## Requirements

| ID | Requirement | Source | Acceptance criteria | Dependencies |
| --- | --- | --- | --- | --- |
| [`REQ-001`](#req-001) | `[testable data contract]` | `[link]` | [`AC-001`](../tests.md#ac-001) | `[links or None]` |

## Agent Verification Checklist

- [ ] Ownership and source-of-truth boundaries are unambiguous.
- [ ] Constraints, concurrency, lifecycle, migration, and rollback agree.
- [ ] Sensitive data has purpose, exposure, retention, and protection rules.
