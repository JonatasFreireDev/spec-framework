---
name: specification
description: "Specification Skill. Use when an agent needs to Create the implementation contract that unifies product, UX, rules, data, APIs, analytics, security, and acceptance criteria in the Spec Framework workflow, including creating, updating, auditing, explaining, routing, or handing off related product artifacts."
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
Approved use case; local and parent context.md files; sibling Specifications; business rules; design notes; Engineering System; Design System; decisions; approved demand classification.

## Outputs
`specification.md` root contract; applicable `contracts/*.md`; stable `REQ-*` and `AC-*` coverage; Delivery Level/Priority rationale; unresolved questions; decision candidates; context.md updates.

## Required reading
- the framework root's `FRAMEWORK.md`
- Relevant parent context.md files.
- This skill owns its generation resources: `assets/specification-template.md` and `assets/specification-contract-template.md`.
- Approved decisions are discovered through the active product root's `.product/decisions.json`; resolve each registered `path` from its declared domain root (`knowledge/decisions/`, `design/decisions/`, or `engineering/decisions/`).

## Discovery and challenge

Follow the shared [Discovery And Challenge contract](../discovery-and-challenge.md) before substantive creation or material revision.

## Workflow
1. Read the local Use Case context, the complete parent chain, sibling Use Cases and Specifications, and confirm the artifact status.
2. Load reusable approved decisions, Engineering System contracts, Design System references, and code evidence named by the context.
3. Confirm whether the demand extends an existing contract or creates a new interaction; do not duplicate an existing requirement.
4. Identify missing information, assumptions, conflicts, and dependencies.
5. Propose the artifact or revision using the matching template and preserve source-demand traceability.
6. Record decision candidates for high-impact or hard-to-reverse choices.
7. Ask for approval before moving the artifact to the next ladder step.
8. Update context.md with new relations, reuse references, impacts, links, dependencies, questions, and status changes.
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
