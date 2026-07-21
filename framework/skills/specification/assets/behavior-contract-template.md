# Behavior Contract: [use case]

## Generation And Agent Self-Check

| Field | Value |
| --- | --- |
| Generated on | `YYYY-MM-DD` |
| Purpose | Make the interaction deterministic across success, alternate, error, and edge paths. |
| Required inputs and evidence | `[use case, product rules, adjacent behavior, decisions]` |
| Ready when | Triggers, transitions, invariants, failures, concurrency, and requirements are explicit. |
| Current status | `draft` |

## Snapshot

| Field | Value |
| --- | --- |
| Status | `draft` |
| Source specification | `[SPEC-XXX](../specification.md)` |
| Contract version | `2` |

## Preconditions And Triggers

| Trigger | Preconditions | Actor/system | Rejection behavior |
| --- | --- | --- | --- |
| `[trigger]` | `[conditions]` | `[owner]` | `[observable result]` |

## State Transitions

| From | Event/command | Guard | To | Side effects |
| --- | --- | --- | --- | --- |
| `[state]` | `[event]` | `[condition]` | `[state]` | `[effects]` |

## Alternate Error And Edge Flows

| Type | Condition | Expected behavior | Recovery/retry | Signal |
| --- | --- | --- | --- | --- |
| `[alternate/error/edge]` | `[condition]` | `[behavior]` | `[recovery]` | `[event/log/metric]` |

## Invariants

| Invariant | Concurrency/idempotency impact | Enforcement point | Source |
| --- | --- | --- | --- |
| `[must always hold]` | `[impact]` | `[server/client/data]` | `[link]` |

## Requirements

| ID | Requirement | Source | Acceptance criteria | Dependencies |
| --- | --- | --- | --- | --- |
| [`REQ-001`](#req-001) | `[testable behavior]` | `[link]` | [`AC-001`](../tests.md#ac-001) | `[links or None]` |

## Agent Verification Checklist

- [ ] Main, alternate, error, and edge behavior is deterministic.
- [ ] State transitions, invariants, concurrency, and retries agree.
- [ ] Every behavior requirement maps to observable acceptance.
