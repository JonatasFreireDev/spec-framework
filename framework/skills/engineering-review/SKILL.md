---
name: engineering-review
description: "Independent read-only Engineering Review. Use to verify an Engineering Proposal against the Specification, Technical Discovery, Engineering System, approved decisions, quality attributes, and operational constraints before implementation planning."
---

# Engineering Review Skill

## Layer
Validation

## Responsibility
Own delivery-specific `engineering-review.md` as an independent read-only verdict. Review the intended solution without editing it, approving decisions, implementing code, or replacing Code Review, QA, or Security Review.

## Operating modes
- review: produce a first verdict for the current proposal content.
- recheck: review a revised proposal and record the new proposal hash.
- audit: identify stale verdicts, missing evidence, and unresolved governed choices.
- explain: summarize findings, routes, and the planning gate.

## Inputs
Approved Specification; Design; approved Technical Discovery; `engineering-proposal.md`; pinned Engineering System; approved decisions and approval records; quality, security, operations, and test contracts.

## Outputs
`engineering-review.md`; proposal hash; `passed`, `required_fix`, `blocked`, or `not_reviewed` verdict; routed findings; planning handoff.

## Required reading
- [`engineering-systems.md`](../../docs/engineering-systems.md) for shared contract versioning, hashes, and review boundaries.
- the framework root's `FRAMEWORK.md`
- Relevant parent context.md files.
- This skill owns its generation resources: `assets/engineering-review-template.md`.
- Approved decisions are discovered through the active product root's `.product/decisions.json`; resolve each registered `path` from its declared domain root (`knowledge/decisions/`, `design/decisions/`, or `engineering/decisions/`).

## Workflow
1. Read the proposal and record its normalized SHA-256 content hash so later edits make the review stale.
2. Verify requirement coverage, architecture boundaries, ownership, integrations, dependencies, quality attributes, security, observability, migration, rollout, rollback, and testability.
3. Verify every governed choice has a scope-compatible approved `DEC-*` and current approval record.
4. Set `blocked` for missing decisions or unsafe unresolved contracts, `required_fix` for correctable proposal gaps, and `passed` only when no blocking finding remains.
5. Route proposal gaps to `engineering-proposal`, specification conflicts to `specification`, and decision gaps to `product-historian` plus a human.
6. Do not edit reviewed inputs, implementation code, or approval records.
7. Hand a passed review to Implementation Planner; all other verdicts return to their named owner.

## Quality checklist
- [ ] Preserves traceability to affected artifacts.
- [ ] Uses the correct template and naming conventions.
- [ ] Is independent and read-only with respect to reviewed artifacts.
- [ ] Records the reviewed proposal hash and concrete evidence.
- [ ] Distinguishes proposal fixes from missing decisions and Specification conflicts.
- [ ] Detects gaps, conflicts, and dependencies.
- [ ] Records or requests decisions for meaningful changes.
- [ ] Leaves a clear handoff for the next skill or orchestrator.

## Handoff
Next: `implementation-planner` only after a `passed` verdict over the current proposal hash.

Pass forward the verdict, proposal hash, evidence, findings, routes, approved decisions, risks, and required follow-up work.
