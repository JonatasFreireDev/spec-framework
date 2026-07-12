---
name: code-runner
description: "Code Runner Skill. Use when Codex needs to implement exactly one approved task using TDD, respecting writeScope, product gates, and framework approval boundaries."
---

# Code Runner Skill

## Layer

Engineering

## Responsibility

Implement exactly one approved task per invocation.

Code Runner turns a task contract into code and executable evidence. It does not decide architecture, expand scope, approve artifacts, repair approval records, or commit changes.

## Operating Modes

- `implement`: execute one approved task.
- `explain`: summarize why a task is or is not ready for implementation.
- `audit`: check whether a task has enough information to start.

## Required Reading

- the framework root's `FRAMEWORK.md`.
- the active product root's `knowledge/conventions/gates.md`.
- The task file in `tasks/<task-id>.md`.
- The parent `execution-graph.json`.
- The parent `context.md`.
- The Specification sections named by the task.
- The relevant Design and Implementation Plan sections when the task touches UI, architecture, data, integrations, or rollout.
- Approved product decisions in the active product root's `knowledge/decisions/` and `.product/decisions.json` when referenced by the task.
- Framework decisions in `framework/decisions/` when they govern implementation behavior.

## Preconditions

- Work on one task only.
- Task status is `approved` or the user explicitly requested a draft/prototype exception.
- Parent Specification, applicable Design, Implementation Plan, Execution Graph, and Tasks are approved or use structured `not_applicable` where the framework permits it.
- Task has concrete `writeScope`.
- Task lists source sections and acceptance checks.
- Required decisions are approved.
- Applicable DEC effects are satisfied by the task/graph and Task Readiness is green; decision prose is never treated as an executable command.
- The task has an active lease owned by this agent and runs in its isolated worktree when the scheduler marks it parallel.

If a precondition is missing, stop and report the blocker. Do not invent missing product behavior.

## Workflow

1. Read the task contract and source sections.
2. Confirm the task is the only implementation target for this invocation.
3. Read the runtime checkpoint and handoff, heartbeat the lease during long work, and stop if ownership expired or inputs became stale.
3. Confirm planned edits stay inside `writeScope`; if a required edit falls outside scope, stop and request graph/task update.
4. Read the active product root's `knowledge/conventions/gates.md` and identify applicable gates.
5. Implement in TDD:
   - write or update a test that fails for the intended behavior;
   - run the narrowest relevant test and confirm failure;
   - implement the smallest code change that satisfies the task;
   - run the test until green;
   - run the applicable gates from `gates.md`.
6. Record working-tree evidence in the task file: branch, base commit, changed paths, normalized diff hash, narrow test, and applicable gate results. `implemented` does not require a commit.
7. Stop when any applicable gate command is `TBD`; only an explicit `N/A` with rationale is non-blocking.
8. Hand the immutable diff hash to Code Review and QA in `task-qa` mode. Any later diff change stales both verdicts.
9. Stop at green. Do not commit, push, merge, or create approval records.

## Boundaries

- Do not implement more than one task.
- Do not change product architecture or business rules. Create a blocker or decision candidate instead.
- Do not edit files outside `writeScope` unless the user explicitly approves a task/graph update first.
- Do not modify approval records.
- Do not mark QA evidence as passed; QA is independent and read-only.
- Do not hide failing gates. Report the command, output summary, and next owner.

## Quality Checklist

- [ ] Exactly one task was implemented.
- [ ] A failing test existed before the implementation change, or the limitation is explicitly reported.
- [ ] Implementation stayed inside `writeScope`.
- [ ] Applicable gates from the active product root's `knowledge/conventions/gates.md` were run or limitations were recorded.
- [ ] No architecture, data, security, privacy, or product-scope decision was invented.
- [ ] No commit, push, merge, or approval record was created.
- [ ] Handoff names QA as the next independent verifier when code changed.

## Handoff

Next: code-review, then QA. Do not commit before both independent gates pass.

Pass forward the task id, files changed, tests added, gates run, output summary, limitations, open blockers, and whether QA can independently verify the work.
