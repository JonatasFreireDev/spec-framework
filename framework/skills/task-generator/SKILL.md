---
name: task-generator
description: "Task Generator Skill. Use when Codex needs to generate complete, executable, testable vertical task contracts from the execution graph and source specification in the Spec Framework workflow, including creating, updating, auditing, explaining, routing, or handing off related product artifacts."
---

# Task Generator Skill

## Layer
Planning

## Responsibility
Generate complete, executable, testable vertical task contracts from the execution graph and source specification.

Task Generator materializes exactly the nodes of a reviewed proposed graph. It does not invent additional tasks during materialization and never overwrites an existing canonical task.

Prefer fewer, richer tasks. A task is not a thin checklist or a file-count partition: it is the complete contract for one coherent, reviewable outcome. Split only at a real dependency, a safe-parallelism boundary with disjoint write scopes, a distinct ownership/toolchain boundary, or a rollback/risk boundary. A task may span frontend, backend, data, tests, and many files when those changes close one vertical outcome.

## Operating modes
- create: produce the first version of the artifact.
- update: revise an existing artifact while preserving approved decisions.
- audit: find gaps, conflicts, dependencies, and missing approvals.
- explain: summarize the artifact and why it exists.

## Inputs
Approved execution graph; implementation plan; specification; applicable Design, Technical Discovery, Engineering Proposal, and Engineering Review; Delivery Level; Priority; repo conventions; test strategy.

## Outputs
tasks.md; task files or task records with Delivery Level/Priority; acceptance checks; handoff notes for implementers.

## Required reading
- [`execution-runtime.md`](../../docs/execution-runtime.md) for canonical task paths, graph materialization, write scopes, and shared resources.
- the framework root's `FRAMEWORK.md`
- Relevant parent context.md files.
- Relevant templates in framework/template/.
- Approved decisions are discovered through the active product root's `.product/decisions.json`; resolve each registered `path` from its declared domain root (`knowledge/decisions/`, `design/decisions/`, or `engineering/decisions/`).

## Discovery and challenge

Follow the shared [Discovery And Challenge contract](../discovery-and-challenge.md) before substantive creation or material revision.

## Workflow
1. Read the parent context and confirm the artifact status.
2. Identify missing information, assumptions, conflicts, and dependencies.
3. Size nodes before materialization. Merge coupled checklists, layers, or file groups into one task unless a real split boundary exists. Never split solely to make tasks appear small.
4. For every graph node, declare concrete `writeScope` paths/modules and `sharedResources` when generated files, indexes, locales, local database state, schema, contracts, or other shared assets are touched.
5. Give every task a complete contract: outcome, included behavior, non-goals, assumptions, implementation approach, acceptance checks, planned tests/evidence, rollback or operational notes, and handoff.
6. Guarantee safe parallelism by construction: nodes with no dependency path between them must have disjoint `writeScope` values and must not share the same `sharedResources`.
7. When scopes or shared resources overlap, choose exactly one: merge the tasks, add a dependency to serialize the work, or assign the shared resource to a single node. Split scopes only when the resulting tasks still close independent outcomes.
8. Materialize the proposed graph atomically into `tasks/*.md` and generated `tasks.md`, changing graph status to `materialized`.
9. Record decision candidates for high-impact or hard-to-reverse choices.
10. Ask for approval of the validated Graph + Tasks set before moving either to the next ladder step.
11. Update context.md with new links, dependencies, questions, and status changes.
12. Map every task to concrete `REQ-*` and `AC-*` identifiers and at least one planned `TEST-*` or explicit non-test evidence method.
13. Propagate applicable DEC references and ensure every structured workflow effect is covered by Graph/Task contracts; do not create unapproved scope to satisfy an effect.

## Quality checklist
- [ ] Preserves traceability to the parent artifact.
- [ ] Uses the correct template and naming conventions.
- [ ] States scope, non-goals, assumptions, and open questions.
- [ ] Is a complete vertical outcome, not a file-count, layer, or checklist partition.
- [ ] Includes the implementation approach, tests/evidence, rollback/operational notes, and an implementer handoff.
- [ ] Each task carries Delivery Level and Priority inherited from the graph or an approved exception.
- [ ] Every graph node has concrete `writeScope`.
- [ ] Parallel graph nodes have disjoint `writeScope`.
- [ ] Parallel graph nodes do not compete for the same `sharedResources`.
- [ ] Detects gaps, conflicts, and dependencies.
- [ ] Records meaningful decisions or decision candidates.
- [ ] Leaves a clear handoff for the next skill.

## Handoff
Next: code-runner after tasks and graph are approved and implementation gates are configured.

Pass forward approved artifacts, open questions, decisions, dependencies, risks, and any remaining audit findings.
