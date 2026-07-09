---
name: task-generator
description: "Task Generator Skill. Use when Codex needs to Generate small, executable, testable tasks from the execution graph and source specification in this repository's Spec Framework workflow, including creating, updating, auditing, explaining, routing, or handing off related product artifacts."
---

# Task Generator Skill

## Layer
Planning

## Responsibility
Generate small, executable, testable tasks from the execution graph and source specification.

## Operating modes
- create: produce the first version of the artifact.
- update: revise an existing artifact while preserving approved decisions.
- audit: find gaps, conflicts, dependencies, and missing approvals.
- explain: summarize the artifact and why it exists.

## Inputs
Approved execution graph; implementation plan; specification; design artifact when applicable; Delivery Level; Priority; repo conventions; test strategy.

## Outputs
tasks.md; task files or task records with Delivery Level/Priority; acceptance checks; handoff notes for implementers.

## Required reading
- FRAMEWORK.md
- Relevant parent context.md files.
- Relevant templates in knowledge/templates/.
- Approved decisions in knowledge/decisions/ and .product/decisions.json.

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
- [ ] Each task carries Delivery Level and Priority inherited from the graph or an approved exception.
- [ ] Detects gaps, conflicts, and dependencies.
- [ ] Records meaningful decisions or decision candidates.
- [ ] Leaves a clear handoff for the next skill.

## Handoff
Next: 14-qa.md

Pass forward approved artifacts, open questions, decisions, dependencies, risks, and any remaining audit findings.
