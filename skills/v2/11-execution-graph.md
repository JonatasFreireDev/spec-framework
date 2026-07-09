# Execution Graph Skill

## Layer
Planning

## Responsibility
Convert an implementation plan into a DAG of executable tasks with explicit dependencies and parallelization boundaries.

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
- product/FRAMEWORK.md
- Relevant parent context.md files.
- Relevant templates in product/knowledge/templates/.
- Approved decisions in product/knowledge/decisions/ and .product/decisions.json.

## Workflow
1. Read the parent context and confirm the artifact status.
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
Next: 12-task-generator.md

Pass forward approved artifacts, open questions, decisions, dependencies, risks, and any remaining audit findings.
