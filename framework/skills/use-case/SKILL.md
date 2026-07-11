---
name: use-case
description: "Use Case Skill. Use when Codex needs to Break a feature into concrete, prioritized, verifiable interactions that can become specifications in the Spec Framework workflow, including creating, updating, auditing, explaining, routing, or handing off related product artifacts."
---

# Use Case Skill

## Layer
Product Design

## Responsibility
Break a feature into concrete, prioritized, verifiable interactions that can become specifications.

## Operating modes
- create: produce the first version of the artifact.
- update: revise an existing artifact while preserving approved decisions.
- audit: find gaps, conflicts, dependencies, and missing approvals.
- explain: summarize the artifact and why it exists.

## Inputs
Approved feature; UX states; business rules; technical constraints; edge cases.

## Outputs
use-case.md files; use-case context.md files; inherited Delivery Level and Priority; acceptance intent; open questions.

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
- [ ] Preserves or explicitly adjusts Delivery Level/Priority with rationale and approval.
- [ ] Detects gaps, conflicts, and dependencies.
- [ ] Records meaningful decisions or decision candidates.
- [ ] Leaves a clear handoff for the next skill.

## Handoff
Next: specification.

Pass forward approved artifacts, open questions, decisions, dependencies, risks, and any remaining audit findings.
