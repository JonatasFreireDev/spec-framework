# Tasks: [use case name]

## Context

- Source specification: [SPEC-XXX]
- Source design: [DESIGN-XXX | Not applicable]
- Source implementation plan: [PLAN-XXX]
- Source execution graph: [GRAPH-XXX]
- Status: [draft | proposed | approved]

## Delivery

- Level: [L0 | L1 | L2 | L3 | L4 | L5]
- Priority: [P0 | P1 | P2 | P3]
- Rationale:

## Rules

- Tasks must trace back to `specification.md`.
- UI tasks must trace back to `design.md`.
- Tasks must carry Delivery Level and Priority from the graph unless an approved exception exists.
- Tasks must follow dependency order from `execution-graph.json`.
- A task must be independently reviewable and testable.
- Parallel tasks must have disjoint write scopes.

## Task: [TK-XXX] [title]

- Type: [database | backend | frontend | test | analytics | docs | security]
- Status: [pending | in_progress | implemented | validated | blocked]
- Delivery: [Lx/Px]
- Depends on: [TK-XXX]
- Source specification sections:
  - [section]
- Write scope:
  - [path/module]
- Objective:
  - [what this task accomplishes]
- Implementation notes:
  - [constraints, patterns, or files to inspect]
- Acceptance criteria:
  - [ ] [observable behavior]
- Validation:
  - [command/test/manual check]
- Handoff:
  - Next task(s): [TK-XXX]
  - Risks: [risk]
