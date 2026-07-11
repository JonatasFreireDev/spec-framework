---
name: implementation-planner
description: "Implementation Planner Skill. Use when Codex needs to Think like a tech lead and translate an approved specification plus approved design into a build strategy without writing code in the Spec Framework workflow, including creating, updating, auditing, explaining, routing, or handing off related product artifacts."
---

# Implementation Planner Skill

## Layer
Planning

## Responsibility
Think like a tech lead and translate an approved specification plus approved design into a build strategy without writing code.

## Operating modes
- create: produce the first version of the artifact.
- update: revise an existing artifact while preserving approved decisions.
- audit: find gaps, conflicts, dependencies, and missing approvals.
- explain: summarize the artifact and why it exists.

## Inputs
Approved specification; approved design or `Not applicable` design artifact; Delivery Level; Priority; architecture context; dependencies; risks; codebase constraints.

## Outputs
implementation-plan.md; inherited Delivery Level and Priority; phases; dependency notes; rollback plan; candidate task slices.

## Required reading
- the framework root's `FRAMEWORK.md`
- Relevant parent context.md files.
- Relevant templates in framework/template/.
- Approved product decisions in the active product root's `knowledge/decisions/` and `.product/decisions.json`.

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
- [ ] Refuses to plan UI-bearing work without approved design.
- [ ] Carries Delivery Level and Priority into phases and candidate tasks.
- [ ] Detects gaps, conflicts, and dependencies.
- [ ] Records meaningful decisions or decision candidates.
- [ ] Leaves a clear handoff for the next skill.

## Handoff
Next: execution-graph.

Pass forward approved artifacts, open questions, decisions, dependencies, risks, and any remaining audit findings.
