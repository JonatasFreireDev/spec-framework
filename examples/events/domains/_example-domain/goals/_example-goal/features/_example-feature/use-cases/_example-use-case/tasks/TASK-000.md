# Task: Example task

## Snapshot

| Field | Value |
| --- | --- |
| ID | TASK-000 |
| Status | draft |
| Source graph | GRAPH-EXAMPLE |
| Source specification | UC-EXAMPLE:specification |
| Source node | TASK-000 |
| Owner skill | task-generator |
| Next skill | code-runner |

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
| Level | N/A |
| Priority | N/A |
| Depends on | none |
| Rationale | Structural example only; not product scope. |

## Task Contract

| Field | Value |
| --- | --- |
| Title | Example task |
| Type | example |
| Depends on | none |
| Source sections | Specification > Acceptance Criteria |
| Write scope | example-only |
| Graph node status | pending |

## Objective

Preserve and execute the migrated task intent from the previous task set.

## Scope And Boundaries

### Included Behavior

- Preserve the structural example task and its navigable task contract.

### Non-Goals

- Do not model real product implementation or split this example into file-level work.

### Assumptions And Constraints

- This is a structural example only; it does not grant delivery readiness.

## Implementation Strategy

- Keep the example as one complete vertical task record, even though it has no executable product scope.

## Acceptance Checks

- Example task has a visible acceptance check.

## Test And Evidence Strategy

- Validate navigability and template conformance; no product test is applicable.

## Implementation Links

| Field | Value |
| --- | --- |
| Branch | N/A until implementation |
| Commits | N/A until implementation |
| PR | N/A until implementation |
| Monorepo model | docs and code live in the same product repository |
| Code paths | example-only |

## Validation Evidence

| Field | Value |
| --- | --- |
| Test status | pending |
| Gate logs | N/A until validation |
| CI URL | N/A until validation |
| Screenshots | N/A until validation |
| QA evidence | N/A until validation |
| Security review | N/A until validation |

## Working Tree Evidence

| Field | Value |
| --- | --- |
| Changed paths | N/A for structural example |
| Diff hash | N/A for structural example |
| Narrow test | template/fixture validation |
| Applicable gates | N/A |
| Code Review diff hash | N/A |
| QA diff hash | N/A |

## Migrated Notes

## TASK-000 Example task

- Type: example
- Depends on: none
- Source: `execution-graph.json`
- Status: draft
- Acceptance criteria:

## Blockers

- None recorded in the canonical task file.

## Handoff

| Field | Value |
| --- | --- |
| Ready for implementation | no |
| Required next skill | task-generator |
| Notes | This file is the canonical task record. Update this file, not tasks.md, when status, links, contract, or evidence changes. |
