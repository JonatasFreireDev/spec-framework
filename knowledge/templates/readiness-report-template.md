# Readiness Report: [feature/use case]

## Context

- Scope: [DOMAIN/GOAL/FT/UC/SPEC]
- Auditor: [skill/orchestrator]
- Date: [YYYY-MM-DD]
- Verdict: [ready | ready_with_notes | not_ready]

## Summary

[Short explanation of whether this artifact can move to the next step.]

## Required Artifacts

| Artifact | Path | Status | Result |
| --- | --- | --- | --- |
| Domain context | [path] | [status] | [pass/fail] |
| Goal context | [path] | [status] | [pass/fail] |
| Feature | [path] | [status] | [pass/fail] |
| Use case | [path] | [status] | [pass/fail] |
| Specification | [path] | [status] | [pass/fail] |
| Implementation plan | [path] | [status] | [pass/fail] |
| Execution graph | [path] | [status] | [pass/fail] |
| Tasks | [path] | [status] | [pass/fail] |

## Gate Checks

### Traceability

- [ ] Every child artifact links to a parent.
- [ ] Every task links to a specification section.
- [ ] The execution graph points to the source specification and implementation plan.

### Specification Completeness

- [ ] Scope and non-goals are explicit.
- [ ] Functional behavior is specified.
- [ ] Business rules are listed.
- [ ] UX states are listed.
- [ ] API/data contracts are present or intentionally N/A.
- [ ] Permissions and security are covered.
- [ ] Analytics and observability are covered.
- [ ] Acceptance criteria are observable.

### Planning Completeness

- [ ] Implementation phases are sequenced.
- [ ] Dependencies are explicit.
- [ ] Risks are documented.
- [ ] Rollout and rollback are documented.
- [ ] Decisions needed are listed.

### Execution Graph Completeness

- [ ] Graph JSON parses.
- [ ] Nodes have ids, titles, types, owners, dependencies, source sections, write scopes, status, and acceptance checks.
- [ ] Dependencies reference existing nodes.
- [ ] Blocked nodes are explained.
- [ ] Parallel lanes do not imply overlapping write scopes.

### Task Readiness

- [ ] Tasks are small enough to implement and review independently.
- [ ] Tasks have acceptance criteria and validation method.
- [ ] Blocked tasks name the blocking decision or dependency.
- [ ] No implementation task starts from an unapproved or incomplete specification unless explicitly marked exploratory.

## Findings

### [Severity] [Finding]

- Evidence: [path/section]
- Impact: [why it matters]
- Required fix: [fix]
- Owner: [role]

## Blocking Decisions

- [Decision needed] - blocks: [artifact/task]

## Result

- Verdict: [ready | ready_with_notes | not_ready]
- Can generate/execute tasks: [yes/no]
- Required next step: [skill/orchestrator]