---
name: specification
description: "Specification Skill. Use when Codex needs to Create the implementation contract that unifies product, UX, rules, data, APIs, analytics, security, and acceptance criteria in the Spec Framework workflow, including creating, updating, auditing, explaining, routing, or handing off related product artifacts."
---

# Specification Skill

## Layer
Specification

## Responsibility
Create the implementation contract that unifies product, UX, rules, data, APIs, analytics, security, and acceptance criteria.

## Operating modes
- create: produce the first version of the artifact.
- update: revise an existing artifact while preserving approved decisions.
- audit: find gaps, conflicts, dependencies, and missing approvals.
- explain: summarize the artifact and why it exists.

## Inputs
Approved use case; feature context; business rules; design notes; engineering constraints; decisions.

## Outputs
`specification.md` root contract; applicable `contracts/*.md`; stable `REQ-*` and `AC-*` coverage; Delivery Level/Priority rationale; unresolved questions; decision candidates; context.md updates.

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
7. Keep `specification.md` as the canonical index; split concerns into product, behavior, UX, API, data, security, quality, observability, and rollout contracts according to rigor.
8. Give every testable requirement and acceptance criterion a stable ID and reject duplicate or orphan IDs.

## Quality checklist
- [ ] Preserves traceability to the parent artifact.
- [ ] Uses the correct template and naming conventions.
- [ ] States scope, non-goals, assumptions, and open questions.
- [ ] Includes Delivery Level and Priority and explains any change from the source feature/use case.
- [ ] Detects gaps, conflicts, and dependencies.
- [ ] Records meaningful decisions or decision candidates.
- [ ] Leaves a clear handoff for the next skill.

## Handoff
Next: ux-ui.

Pass forward approved specification, Delivery Level, Priority, open questions, decisions, dependencies, risks, and any remaining audit findings.
