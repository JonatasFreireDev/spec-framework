---
name: vision
description: "Vision Skill. Use when Codex needs to Turn an approved problem into a product vision, principles, north star, and strategic boundaries in the Spec Framework workflow, including creating, updating, auditing, explaining, routing, or handing off related product artifacts."
---

# Vision Skill

## Layer
Foundation

## Responsibility
Turn an approved problem into a product vision, principles, north star, and strategic boundaries.

## Operating modes
- create: produce the first version of the artifact.
- update: revise an existing artifact while preserving approved decisions.
- audit: find gaps, conflicts, dependencies, and missing approvals.
- explain: summarize the artifact and why it exists.

## Inputs
Approved problem; target audience; evidence; constraints; anti-goals; founder intent.

## Outputs
vision.md; principles.md; north-star.md; context.md updates; decision candidates.

Artifact ownership is exclusive: `vision.md` owns the product direction, target users, non-goals, and decision boundaries; `principles.md` owns principles, trade-offs, examples, and anti-principles; `north-star.md` owns the value outcome, candidate metric, measurement notes, and guardrails. Link these artifacts instead of copying their canonical content into `vision.md`.

## Required reading
- the framework root's `FRAMEWORK.md`
- Relevant parent context.md files.
- Relevant templates in framework/template/.
- Approved decisions are discovered through the active product root's `.product/decisions.json`; resolve each registered `path` from its declared domain root (`knowledge/decisions/`, `design/decisions/`, or `engineering/decisions/`).

## Discovery and challenge

Follow the shared [Discovery And Challenge contract](../discovery-and-challenge.md) before substantive creation or material revision.

## Workflow
1. Read the parent context and confirm the artifact status.
2. Identify missing information, assumptions, conflicts, and dependencies.
3. Propose the artifact or revision using the matching template.
   - Keep principles and north-star details only in their dedicated companion artifacts.
   - In `vision.md`, reference those companions instead of creating a second source of truth.
4. Record decision candidates for high-impact or hard-to-reverse choices.
5. Ask for approval before moving the artifact to the next ladder step.
6. Update context.md with new links, dependencies, questions, and status changes.

## Quality checklist
- [ ] Preserves traceability to the parent artifact.
- [ ] Uses the correct template and naming conventions.
- [ ] States scope, non-goals, assumptions, and open questions.
- [ ] Detects gaps, conflicts, and dependencies.
- [ ] Records meaningful decisions or decision candidates.
- [ ] Does not duplicate principles or north-star content in `vision.md`.
- [ ] Leaves a clear handoff for the next skill.

## Handoff
Next: strategy.

Pass forward approved artifacts, open questions, decisions, dependencies, risks, and any remaining audit findings.
