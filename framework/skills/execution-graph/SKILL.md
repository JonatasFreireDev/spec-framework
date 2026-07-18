---
name: execution-graph
description: "Execution Graph Skill. Use when an agent needs to Convert an implementation plan into a DAG of executable tasks with explicit dependencies and parallelization boundaries in the Spec Framework workflow, including creating, updating, auditing, explaining, routing, or handing off related product artifacts."
---

# Execution Graph Skill

## Layer
Planning

## Responsibility
Convert an implementation plan into a DAG of complete vertical task contracts with explicit dependencies and parallelization boundaries. Prefer the fewest nodes that preserve real dependency, safe-parallelism, ownership/toolchain, or rollback/risk boundaries; never partition merely by file count, technical layer, or checklist length.

## Operating modes
- create: produce the first version of the artifact.
- update: revise an existing artifact while preserving approved decisions.
- audit: find gaps, conflicts, dependencies, and missing approvals.
- explain: summarize the artifact and why it exists.

## Inputs
Approved implementation plan; Delivery Level; Priority; candidate tasks; dependency notes; ownership boundaries.

## Outputs
execution-graph.json; ordered DAG; Delivery Level/Priority on graph and nodes; parallel lanes; blocked nodes; graph risks.

## Required reading
- [`execution-runtime.md`](../../docs/execution-runtime.md) for graph lifecycle, write scopes, shared resources, and runtime boundaries.
- the framework root's `FRAMEWORK.md`
- Relevant parent context.md files.
- This skill owns its generation resources: `assets/execution-graph-template.json`.
- Approved decisions are discovered through the active product root's `.product/decisions.json`; resolve each registered `path` from its declared domain root (`knowledge/decisions/`, `design/decisions/`, or `engineering/decisions/`).

## Discovery and challenge

Follow the shared [Discovery And Challenge contract](../discovery-and-challenge.md) before substantive creation or material revision.

## Workflow
1. Read the parent context and confirm the artifact status.
2. Create the graph as `draft`, validate its DAG and contracts, then move it to `proposed` for human review without requiring task files to exist.
3. After review, use confirmed graph materialization to create canonical task files and `tasks.md`; never create them ad hoc or overwrite existing tasks.
4. Require Graph + Tasks validation before the graph advances from `materialized` to `approved`.
5. Reference applicable DEC IDs and cover their required task types, write scopes, shared resources, gates, and evidence contracts.
2. Identify missing information, assumptions, conflicts, and dependencies.
3. Propose the artifact or revision using the matching template.
4. Record decision candidates for high-impact or hard-to-reverse choices.
5. Ask for approval before moving the artifact to the next ladder step.
6. Update context.md with new links, dependencies, questions, and status changes.

## Quality checklist
- [ ] Preserves traceability to the parent artifact.
- [ ] Uses the correct template and naming conventions.
- [ ] States scope, non-goals, assumptions, and open questions.
- [ ] Propagates Delivery Level and Priority to every node unless a node has an approved exception.
- [ ] Detects gaps, conflicts, and dependencies.
- [ ] Records meaningful decisions or decision candidates.
- [ ] Leaves a clear handoff for the next skill.

## Handoff
Next: task-generator through confirmed graph materialization.

Pass forward approved artifacts, open questions, decisions, dependencies, risks, and any remaining audit findings.
