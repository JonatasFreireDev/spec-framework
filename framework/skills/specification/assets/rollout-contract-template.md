# Rollout Contract: [use case]

## Generation And Agent Self-Check

| Field | Value |
| --- | --- |
| Generated on | `YYYY-MM-DD` |
| Purpose | Define a reversible release with compatibility, monitoring, and recovery boundaries. |
| Required inputs and evidence | `[behavior, API/data compatibility, operations, decisions, dependencies]` |
| Ready when | Release, migration, backfill, monitoring, decision points, rollback, and recovery are explicit. |
| Current status | `draft` |

## Snapshot

| Field | Value |
| --- | --- |
| Status | `draft` |
| Rationale | `[required when not_applicable]` |
| Source specification | `[SPEC-XXX](../specification.md)` |
| Contract version | `2` |

## Release Strategy

| Stage/audience | Activation | Entry condition | Exit condition | Owner |
| --- | --- | --- | --- | --- |
| `[stage]` | `[flag/config/deploy]` | `[condition]` | `[condition]` | `[owner]` |

## Compatibility Migration And Backfill

| Change | Compatibility window | Migration/backfill | Validation | Cleanup trigger |
| --- | --- | --- | --- | --- |
| `[change]` | `[window]` | `[plan or None]` | `[check]` | `[trigger]` |

## Monitoring And Decision Points

| Signal | Expected range | Pause/rollback threshold | Decision owner | Observation window |
| --- | --- | --- | --- | --- |
| `[metric/event]` | `[range]` | `[threshold]` | `[owner]` | `[window]` |

## Rollback And Recovery

| Failure mode | Rollback action | Data recovery | Maximum tolerated impact | Verification |
| --- | --- | --- | --- | --- |
| `[failure]` | `[action]` | `[plan]` | `[boundary]` | `[check]` |

## Requirements

| ID | Requirement | Source | Acceptance criteria | Dependencies |
| --- | --- | --- | --- | --- |
| [`REQ-001`](#req-001) | `[testable rollout contract]` | `[link]` | [`AC-001`](../tests.md#ac-001) | `[links or None]` |

## Agent Verification Checklist

- [ ] Activation, compatibility, migration, and cleanup are reversible or explicitly governed.
- [ ] Monitoring connects to concrete pause and rollback decisions.
- [ ] Recovery covers code, configuration, and data consequences.
