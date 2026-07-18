---
name: implementation-planner
description: "Implementation Planner Skill. Use when an agent needs to Think like a tech lead and translate an approved specification plus approved design into a build strategy without writing code in the Spec Framework workflow, including creating, updating, auditing, explaining, routing, or handing off related product artifacts."
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
Approved specification; approved design or `Not applicable` design artifact; approved Technical Discovery; applicable Engineering Proposal and passed Engineering Review; Delivery Level; Priority; architecture context; dependencies; risks; codebase constraints.

## Outputs
implementation-plan.md; inherited Delivery Level and Priority; phases; dependency notes; rollback plan; candidate task slices.

## Required reading
- the framework root's `FRAMEWORK.md`
- Relevant parent context.md files.
- This skill owns its generation resources: `assets/implementation-plan-template.md`.
- Approved decisions are discovered through the active product root's `.product/decisions.json`; resolve each registered `path` from its declared domain root (`knowledge/decisions/`, `design/decisions/`, or `engineering/decisions/`).

## Discovery and challenge

Follow the shared [Discovery And Challenge contract](../discovery-and-challenge.md) before substantive creation or material revision.

## Workflow
1. Propagate every applicable DEC reference and its validated workflow effects into phases, gates, evidence, rollout, rollback, and candidate tasks; never infer commands from decision prose.
1. Read the parent context and confirm the artifact status.
2. Identify missing information, assumptions, conflicts, and dependencies.
3. Propose the artifact or revision using the matching template.
4. For Tier L or another applicable delivery, refuse to plan until Engineering Review passed against the current Engineering Proposal.
5. Record decision candidates for high-impact or hard-to-reverse choices.
6. Ask for approval before moving the artifact to the next ladder step.
7. Update context.md with new links, dependencies, questions, and status changes.

## Quality checklist
- [ ] Preserves traceability to the parent artifact.
- [ ] Uses the correct template and naming conventions.
- [ ] States scope, non-goals, assumptions, and open questions.
- [ ] Refuses to plan UI-bearing work without approved design.
- [ ] Refuses applicable work without a passed, current Engineering Review.
- [ ] Carries Delivery Level and Priority into phases and candidate tasks.
- [ ] Detects gaps, conflicts, and dependencies.
- [ ] Records meaningful decisions or decision candidates.
- [ ] Leaves a clear handoff for the next skill.

## Handoff
Next: execution-graph.

Pass forward approved artifacts, open questions, decisions, dependencies, risks, and any remaining audit findings.
