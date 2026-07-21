# Quality Contract: [use case]

## Generation And Agent Self-Check

| Field | Value |
| --- | --- |
| Generated on | `YYYY-MM-DD` |
| Purpose | Map delivery risks and acceptance to proportionate verification and evidence. |
| Required inputs and evidence | `[all applicable contracts, Quality System, test strategy, gates]` |
| Ready when | Risks, criteria, levels, environments, evidence, and exit conditions are complete. |
| Current status | `draft` |

## Snapshot

| Field | Value |
| --- | --- |
| Status | `draft` |
| Source specification | `[SPEC-XXX](../specification.md)` |
| Contract version | `2` |

## Quality Risks

| Risk | Impact/likelihood | Preventive contract | Detection method | Owner |
| --- | --- | --- | --- | --- |
| `[risk]` | `[rating]` | `[REQ/link]` | `[method]` | `[owner]` |

## Acceptance Traceability

| Acceptance criterion | Requirement | Risk | Test/evidence method | Expected evidence |
| --- | --- | --- | --- | --- |
| [`AC-001`](../tests.md#ac-001) | [`REQ-001`](#req-001) | `[risk]` | `[TEST-001/manual/review]` | `[artifact/log/screenshot]` |

## Test Levels And Environments

| Coverage | Level | Environment/platform | Test data | Isolation/dependencies |
| --- | --- | --- | --- | --- |
| `[behavior]` | `[unit/integration/e2e/security]` | `[configured target]` | `[class]` | `[constraints]` |

## Evidence And Exit Conditions

| Gate | Pass condition | Evidence owner | Failure route |
| --- | --- | --- | --- |
| `[gate]` | `[objective condition]` | `[owner]` | `[skill/human]` |

## Requirements

| ID | Requirement | Source | Acceptance criteria | Dependencies |
| --- | --- | --- | --- | --- |
| [`REQ-001`](#req-001) | `[testable quality contract]` | `[link]` | [`AC-001`](../tests.md#ac-001) | `[links or None]` |

## Agent Verification Checklist

- [ ] Every acceptance criterion maps to risk, method, owner, and evidence.
- [ ] Levels, environments, data, and platforms conform to the configured Quality System.
- [ ] Exit conditions are objective and route failures without fabricating evidence.
