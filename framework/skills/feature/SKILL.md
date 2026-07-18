---
name: feature
description: "Feature Skill. Use when an agent needs to Define features that serve a user goal while keeping scope, Delivery Level, Priority, and testability explicit in the Spec Framework workflow, including creating, updating, auditing, explaining, routing, or handing off related product artifacts."
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
Approved goal and journey; opportunity list; constraints; related features; approved demand classification; existing Feature context.

## Outputs
feature.md; feature context.md; Delivery Level; Priority; explicit delivery slice; scope/non-goals; related use-case candidates.

## Required reading
- the framework root's `FRAMEWORK.md`
- Relevant parent context.md files.
- This skill owns its generation resources: `assets/feature-template.md`.
- Approved decisions are discovered through the active product root's `.product/decisions.json`; resolve each registered `path` from its declared domain root (`knowledge/decisions/`, `design/decisions/`, or `engineering/decisions/`).

## Discovery and challenge

Follow the shared [Discovery And Challenge contract](../discovery-and-challenge.md) before substantive creation or material revision.

## Workflow
1. Read the parent context and confirm the artifact status.
2. Read sibling Features and related Use Cases before creating or revising scope.
3. Identify missing information, assumptions, conflicts, and dependencies.
4. For an incoming demand, confirm that it belongs to this Goal and record whether it extends this Feature or creates a sibling Feature.
5. Propose the artifact or revision using the matching template.
6. Record decision candidates for high-impact or hard-to-reverse choices.
7. Ask for approval before moving the artifact to the next ladder step.
8. Update context.md with `relations`, `traceability`, links, dependencies, questions, and status changes.

## Quality checklist
- [ ] Preserves traceability to the parent artifact.
- [ ] Uses the correct template and naming conventions.
- [ ] States scope, non-goals, assumptions, and open questions.
- [ ] Carries Delivery Level and Priority from roadmap, or marks the missing prioritization as a blocker.
- [ ] Declares user value, entry point, observable end state, independent releasability, reversibility, and deferred behavior.
- [ ] Detects gaps, conflicts, and dependencies.
- [ ] Records meaningful decisions or decision candidates.
- [ ] Leaves a clear handoff for the next skill.

## Handoff
Next: use-case after the feature and delivery slice are approved.

Pass forward approved artifacts, open questions, decisions, dependencies, risks, and any remaining audit findings.
