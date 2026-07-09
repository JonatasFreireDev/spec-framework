# Readiness Gate

The Readiness Gate answers one question:

Can this product artifact safely move to the next step of the Product Engineering Framework ladder?

Most importantly, it prevents a feature from becoming executable tasks before the Specification, Implementation Plan, and Execution Graph are good enough.

## When To Run

Run this gate before:

- moving a feature into use-case specification;
- approving a specification;
- generating an implementation plan;
- generating an execution graph;
- generating tasks;
- starting implementation work;
- releasing a completed feature.

## Verdicts

### ready

The artifact can move to the next step. Remaining notes are non-blocking.

### ready_with_notes

The artifact can move forward, but there are explicit risks or follow-up items. Notes must be carried into the next context.md or handoff.

### not_ready

The artifact must not move forward. A blocking gap, conflict, missing decision, invalid graph, or incomplete specification exists.

## Minimum Gate For Task Generation

A use case can generate executable tasks only when all of these are true:

- `use-case.md` exists and has observable acceptance criteria.
- `specification.md` exists and covers product, UX, data/API, permissions, analytics, errors, edge cases, rollout, and acceptance criteria.
- `implementation-plan.md` exists and defines phases, dependencies, risks, rollout, rollback, and candidate tasks.
- `execution-graph.json` exists and parses as JSON.
- Every execution graph node has an id, title, type, owner skill, dependency list, source sections, write scope, status, and acceptance checks.
- Every dependency in the graph points to an existing node.
- Every task in `tasks.md` traces to the specification and graph.
- Blocking decisions are either resolved or the affected tasks are marked blocked.

## Required Inputs

- Domain context.
- Goal context.
- Feature context and feature.md.
- Use case context and use-case.md.
- specification.md.
- implementation-plan.md.
- execution-graph.json.
- tasks.md.
- Relevant decisions.

## Output

Use `product/knowledge/templates/readiness-report-template.md`.

Save reports under:

```text
product/audits/readiness/<scope-id>-readiness.md
```

## Validator Command

Run the executable validator with:

```bash
npm run product:readiness -- <use-case-dir>
```

Example:

```bash
npm run product:readiness -- product/domains/events/goals/participate-in-event/features/qr-code-check-in/use-cases/attendee-checks-in-with-qr-code
```

Use JSON output for automation:

```bash
npm run product:readiness -- <use-case-dir> --json
```

The validator checks:

- required files exist;
- required markdown headings exist;
- execution graph parses as JSON;
- execution graph dependencies are valid;
- tasks reference graph node ids;
- blocked tasks have blocking reason;
- no placeholder text remains in approved artifacts.

## Human Review Questions

- Is the user outcome still clear?
- Is the specification detailed enough that two agents would implement the same behavior?
- Are decisions separated from assumptions?
- Are security, privacy, permissions, analytics, and error states explicit?
- Can tasks run in parallel without overlapping write scope?
- Is every blocker visible before implementation starts?