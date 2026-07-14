---
name: problem-discovery
description: "Problem Discovery Skill. Use when Codex needs to Discover, frame, and validate the core problem before solution work starts in the Spec Framework workflow, including creating, updating, auditing, explaining, routing, or handing off related product artifacts."
---

# Problem Discovery Skill

## Layer
Foundation

## Responsibility
Discover, frame, and validate the core problem before solution work starts.

## Operating modes
- create: produce the first version of the artifact.
- update: revise an existing artifact while preserving approved decisions.
- audit: find gaps, conflicts, dependencies, and missing approvals.
- explain: summarize the artifact and why it exists.

## Inputs
Raw idea; user pain; market signal; stakeholder notes; research snippets; existing product context.

## Outputs
problem.md; opportunities.md; research questions; assumptions; risk notes; context.md updates.

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
Next: vision.

Pass forward approved artifacts, open questions, decisions, dependencies, risks, and any remaining audit findings.
