# Task: [task title]

## 🧾 Generation And Agent Self-Check

> Complete this section when materializing the artifact. Keep unresolved items explicit in the relevant scope, findings, risks, or handoff section.

| Field | Value |
| --- | --- |
| Generated on | `YYYY-MM-DD` |
| Purpose | `[decision, evidence, contract, or handoff this artifact supports]` |
| Use when | `[workflow stage, trigger, or condition]` |
| Prepared by | `[owning skill, role, or accountable person]` |
| Scope covered | `[artifact, product area, use case, or review boundary]` |
| Required inputs and evidence | `[links to approved parents, documents, code, decisions, or observations]` |
| Ready when | `[artifact-specific completion, evidence, and gate conditions]` |
| Current status | `[status allowed by this artifact's owning workflow]` |


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

[Explain the coherent vertical outcome this task must produce.]

## Scope And Boundaries

### Included Behavior

- [Concrete behavior and integration this task closes end-to-end.]

### Non-Goals

- [Adjacent work intentionally excluded from this task.]

### Assumptions And Constraints

- [Approved decision, dependency, compatibility, rollout, or operational constraint.]

## Implementation Strategy

- [Describe the coherent approach across affected modules. Do not divide this task merely by file or technical layer.]
- [Name contract, migration, integration, or rollback boundaries relevant to implementation.]

## Acceptance Checks

- [Observable, reviewable check.]

## Test And Evidence Strategy

- [Planned TEST-* coverage or explicit evidence method, including relevant negative and integration cases.]
- [Applicable validation commands, operational evidence, or visual/accessibility evidence.]

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

## ✅ Agent Verification Checklist

- [ ] The task maps to one graph node, specification scope, delivery level, priority, and owner.
- [ ] Objective, scope/non-goals, implementation strategy, dependencies, acceptance checks, tests/evidence, and validation commands are executable.
- [ ] The task closes one coherent vertical outcome; it was not split merely by file count, layer, or checklist length.
- [ ] Working-tree and validation evidence use the current branch, base commit, paths, diff hash, and gates.
- [ ] Status and handoff follow lifecycle rules and do not claim review, QA, commit, or release prematurely.
