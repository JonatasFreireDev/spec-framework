# UX Contract: [use case]

## Generation And Agent Self-Check

| Field | Value |
| --- | --- |
| Generated on | `YYYY-MM-DD` |
| Purpose | Define the verifiable interaction states that Design must realize. |
| Required inputs and evidence | `[use case, behavior contract, Design System, research, adjacent screens]` |
| Ready when | Entry, navigation, states, feedback, recovery, content, and accessibility are explicit. |
| Current status | `draft` |

## Snapshot

| Field | Value |
| --- | --- |
| Status | `draft` |
| Rationale | `[required when not_applicable]` |
| Source specification | `[SPEC-XXX](../specification.md)` |
| Contract version | `2` |

## Entry Points And Navigation

| Entry point | Preconditions | Destination | Exit/cancel behavior |
| --- | --- | --- | --- |
| `[entry]` | `[conditions]` | `[screen/state]` | `[behavior]` |

## Interaction States

| State | Trigger | Displayed data/actions | Next states |
| --- | --- | --- | --- |
| `[idle/loading/empty/success/error]` | `[trigger]` | `[content/actions]` | `[states]` |

## Feedback And Recovery

| Condition | User feedback | Recovery action | Persistence/retry behavior |
| --- | --- | --- | --- |
| `[condition]` | `[message/cue]` | `[action]` | `[behavior]` |

## Accessibility And Content

| Area | Requirement | Locale/content owner | Verification |
| --- | --- | --- | --- |
| `[focus/label/contrast/motion/copy]` | `[requirement]` | `[owner]` | `[method]` |

## Requirements

| ID | Requirement | Source | Acceptance criteria | Dependencies |
| --- | --- | --- | --- | --- |
| [`REQ-001`](#req-001) | `[testable UX contract]` | `[link]` | [`AC-001`](../tests.md#ac-001) | `[links or None]` |

## Agent Verification Checklist

- [ ] Every behavior state has visible or assistive feedback where applicable.
- [ ] Recovery, cancellation, accessibility, and content ownership are explicit.
- [ ] Design can proceed without inventing interaction behavior.
