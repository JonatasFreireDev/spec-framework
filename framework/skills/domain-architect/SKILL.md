---
name: domain-architect
description: "Domain Architect Skill. Use when Codex needs to Model product domains, ownership boundaries, dependencies, and cross-domain risks in the Spec Framework workflow, including creating, updating, auditing, explaining, routing, or handing off related product artifacts."
---

# Domain Architect Skill

## Layer
Foundation

## Responsibility
Model product domains, ownership boundaries, dependencies, and cross-domain risks.

## Operating modes
- create: produce the first version of the artifact.
- update: revise an existing artifact while preserving approved decisions.
- audit: find gaps, conflicts, dependencies, and missing approvals.
- explain: summarize the artifact and why it exists.

## Inputs
Approved strategy; glossary; business rules; existing code/docs; domain candidates.

## Outputs
domain.md files; domain context.md files; ubiquitous language; invariants; commands/events; data ownership; source-of-truth and authorization boundaries; dependency notes; boundary decisions.

## Required reading
- the framework root's `FRAMEWORK.md`
- Relevant parent context.md files.
- Relevant templates in framework/template/.
- The pinned framework runtime's `examples/events/domains/README.md` and `examples/events/domains/events/domain.md` before creating the first domain.
- Approved decisions are discovered through the active product root's `.product/decisions.json`; resolve each registered `path` from its declared domain root (`knowledge/decisions/`, `design/decisions/`, or `engineering/decisions/`).

## Discovery and challenge

Follow the shared [Discovery And Challenge contract](../discovery-and-challenge.md) before substantive creation or material revision.

## Workflow
1. Read the parent context and confirm the artifact status.
2. Identify missing information, assumptions, conflicts, and dependencies.
3. Partition responsibilities into coherent business areas; do not use the product name, a UI navigation section, or a catch-all domain as the boundary.
4. Propose the artifact or revision using the matching template, including Owns, Does Not Own, and cross-domain contracts.
5. For the first delivery slice, materialize Domain -> User Goal -> Feature -> Use Case before routing to workspace creation.
6. If authentication or identity is owned by a non-identity business domain, either split or explicitly justify the boundary and dependency.
7. Record decision candidates for high-impact or hard-to-reverse choices.
8. Ask for approval before moving the artifact to the next ladder step.
9. Update context.md with new links, dependencies, questions, and status changes.

## Quality checklist
- [ ] Preserves traceability to the parent artifact.
- [ ] Uses the correct template and naming conventions.
- [ ] States scope, non-goals, assumptions, and open questions.
- [ ] Defines what the domain does not own and its cross-domain dependencies.
- [ ] Does not put authentication inside an unrelated business domain without an explicit boundary decision.
- [ ] Materializes a first goal, feature, and use case rather than stopping at `domain.md`.
- [ ] Detects gaps, conflicts, and dependencies.
- [ ] Records meaningful decisions or decision candidates.
- [ ] Leaves a clear handoff for the next skill.

## Handoff
Next: user-goal, then domain-evolution-orchestrator when goals and journeys are ready.

Pass forward approved artifacts, open questions, decisions, dependencies, risks, and any remaining audit findings.
