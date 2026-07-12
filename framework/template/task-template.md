# Task: [task title]

## Snapshot

| Field | Value |
| --- | --- |
| ID | `[TK-XXX-001]` |
| Status | `[draft | proposed | approved | in_progress | implemented | validated | released]` |
| Source graph | `[GRAPH-XXX]` |
| Source specification | `[SPEC-XXX]` |
| Source node | `[TK-XXX-001]` |
| Owner skill | `[code-runner | documentation-writer | qa | security-review]` |
| Next skill | `[next skill]` |

## Navigation

| Artifact | Link |
| --- | --- |
| Context | [../context.md](../context.md) |
| Specification | [../specification.md](../specification.md) |
| Implementation Plan | [../implementation-plan.md](../implementation-plan.md) |
| Execution Graph | [../execution-graph.json](../execution-graph.json) |
| Tasks Index | [../tasks.md](../tasks.md) |
| Tests | [../tests.md](../tests.md) |

## Delivery

| Field | Value |
| --- | --- |
| Level | `[L0 | L1 | L2 | L3 | L4 | L5 | N/A]` |
| Priority | `[P0 | P1 | P2 | P3 | N/A]` |
| Depends on | `[task ids or artifact ids]` |
| Rationale | `[why this task belongs to this level and priority]` |

## Task Contract

| Field | Value |
| --- | --- |
| Title | `[task title]` |
| Type | `[database/backend/frontend/test/analytics/docs/security]` |
| Depends on | `[task ids or none]` |
| Source sections | `[Specification sections or plan sections]` |
| Requirements | `[REQ-* ids]` |
| Acceptance criteria | `[AC-* ids]` |
| Planned tests | `[TEST-* ids or explicit evidence method]` |
| Applicable decisions | `[DEC-* ids or N/A]` |
| Write scope | `[paths/modules this task may touch]` |
| Shared resources | `[generated indexes/locales/local database/schema/contracts or none]` |
| Graph node status | `[graph node operational status]` |

## Objective

[Explain the smallest useful outcome this task must produce.]

## Acceptance Checks

- [Observable, reviewable check.]

## Implementation Links

| Field | Value |
| --- | --- |
| Branch | `N/A until implementation` |
| Base commit | `N/A until implementation` |
| Diff hash | `N/A until implementation` |
| Commits | `N/A until QA and Code Review pass` |
| PR | `N/A until implementation` |
| Code paths | `N/A until implementation` |
| Monorepo model | `docs and code live in the same product repository` |
| Commit convention | `knowledge/conventions/commits.md` |
| PR convention | `knowledge/conventions/pull-requests.md` |

## Working Tree Evidence

| Field | Value |
| --- | --- |
| Changed paths | `[repo-relative paths or N/A]` |
| Diff hash | `[sha256 of reviewed diff or N/A]` |
| Narrow test | `[command and result or N/A]` |
| Applicable gates | `[gate ids and results or N/A]` |
| Code Review diff hash | `[hash or pending]` |
| QA diff hash | `[hash or pending]` |

## Validation Evidence

| Field | Value |
| --- | --- |
| Test status | `pending` |
| Gate logs | `N/A until validation` |
| CI URL | `N/A until validation` |
| Screenshots | `N/A until validation` |
| QA evidence | `N/A until validation` |
| Security review | `N/A until validation` |

## Blockers

- [Open blocker or `None`.]

## Handoff

| Field | Value |
| --- | --- |
| Ready for implementation | `[yes/no; requires configured gates]` |
| Readiness command | `spec-framework task readiness --graph ../execution-graph.json --task [TK-XXX-001]` |
| Required next skill | `[skill]` |
| Notes | `[handoff notes]` |
