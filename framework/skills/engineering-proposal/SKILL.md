---
name: engineering-proposal
description: "Engineering Proposal Skill. Use when Codex needs to translate approved delivery contracts and Technical Discovery into an intended technical solution before independent Engineering Review in the Spec Framework workflow."
---

# Engineering Proposal Skill

## Layer
Planning

## Responsibility
Own delivery-specific `engineering-proposal.md`. Describe the intended technical change and its alignment with stable engineering contracts without sequencing tasks, approving architecture, or writing implementation code.

## Operating modes
- create: produce the first technical proposal after a resolved Architecture Gate.
- update: revise the proposal in response to approved contract changes or review findings.
- audit: detect missing requirements, ownership, decisions, quality attributes, and operational consequences.
- explain: summarize the proposed change and why it is or is not reviewable.

## Inputs
Approved Specification; approved Design or `Not applicable`; approved Technical Discovery; resolved Architecture Gate; pinned Engineering System when configured; approved decisions; quality and operations constraints.

## Outputs
`engineering-proposal.md`; intended boundaries and ownership; dependency and integration changes; quality and operations contract; deviations; decision candidates; review handoff.

## Required reading
- the framework root's `FRAMEWORK.md`
- Relevant parent context.md files.
- Relevant templates in framework/template/.
- Approved product decisions in the active product root's `knowledge/decisions/` and `.product/decisions.json`.

## Discovery and challenge

Follow the shared [Discovery And Challenge contract](../discovery-and-challenge.md) before substantive creation or material revision.

## Workflow
1. Confirm Technical Discovery is approved, its Architecture Gate is resolved, and applicability comes from Tier L or declared structured engineering triggers rather than prose inference.
2. Pin the configured Engineering System id/version or explicitly record `Not configured`.
3. Separate existing evidence from the intended change and map every applicable `REQ-*` to affected boundaries.
4. Define module and data ownership, interfaces, dependencies, quality attributes, tests, observability, migration, rollout, and rollback consequences.
5. Link deviations and governed choices to applicable approved decisions; stop and request a decision when required.
6. Keep implementation sequence and task slicing out of the proposal.
7. Hand the proposal to independent Engineering Review without advancing its status or creating approval evidence.

## Quality checklist
- [ ] Preserves traceability to affected artifacts.
- [ ] Uses the correct template and naming conventions.
- [ ] Distinguishes current evidence from the proposed change.
- [ ] Maps requirements to boundaries, ownership, quality, and operations.
- [ ] Pins the Engineering System or records `Not configured` explicitly.
- [ ] Detects gaps, conflicts, and dependencies.
- [ ] Records or requests decisions for meaningful changes.
- [ ] Leaves a clear handoff for the next skill or orchestrator.

## Handoff
Next: `engineering-review`.

Pass forward the proposal, source contracts, pinned system, decisions, deviations, risks, and open blockers.
