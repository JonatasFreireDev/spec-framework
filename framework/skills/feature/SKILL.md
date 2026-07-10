---
name: feature
description: "Feature Skill. Use when Codex needs to Define features that serve a user goal while keeping scope, Delivery Level, Priority, and testability explicit in the Spec Framework workflow, including creating, updating, auditing, explaining, routing, or handing off related product artifacts."
---

# Feature Skill

## Layer
Product Design

## Responsibility
Define features that serve a user goal while keeping scope, Delivery Level, Priority, and testability explicit.

## Operating modes
- create: produce the first version of the artifact.
- update: revise an existing artifact while preserving approved decisions.
- audit: find gaps, conflicts, dependencies, and missing approvals.
- explain: summarize the artifact and why it exists.

## Inputs
Approved goal and journey; opportunity list; constraints; related features.

## Outputs
feature.md; feature context.md; Delivery Level; Priority; scope/non-goals; related use-case candidates.

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
- [ ] Carries Delivery Level and Priority from roadmap, or marks the missing prioritization as a blocker.
- [ ] Detects gaps, conflicts, and dependencies.
- [ ] Records meaningful decisions or decision candidates.
- [ ] Leaves a clear handoff for the next skill.

## Handoff
Next: 08-use-case.md

Pass forward approved artifacts, open questions, decisions, dependencies, risks, and any remaining audit findings.
