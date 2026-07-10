---
name: user-goal
description: "User Goal Skill. Use when Codex needs to Model stable user goals inside a domain, replacing vague capabilities with user-centered intent in the Spec Framework workflow, including creating, updating, auditing, explaining, routing, or handing off related product artifacts."
---

# User Goal Skill

## Layer
Product Design

## Responsibility
Model stable user goals inside a domain, replacing vague capabilities with user-centered intent.

## Operating modes
- create: produce the first version of the artifact.
- update: revise an existing artifact while preserving approved decisions.
- audit: find gaps, conflicts, dependencies, and missing approvals.
- explain: summarize the artifact and why it exists.

## Inputs
Approved domain; journeys; personas; problem evidence; business rules.

## Outputs
goal.md; goal context.md; candidate journeys; feature opportunity list.

## Required reading
- the framework root's `FRAMEWORK.md`
- Relevant parent context.md files.
- Relevant templates in knowledge/templates/.
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
- [ ] Detects gaps, conflicts, and dependencies.
- [ ] Records meaningful decisions or decision candidates.
- [ ] Leaves a clear handoff for the next skill.

## Handoff
Next: 06-journey.md

Pass forward approved artifacts, open questions, decisions, dependencies, risks, and any remaining audit findings.