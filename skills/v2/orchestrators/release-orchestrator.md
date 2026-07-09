# Release Orchestrator

## Mission
Verify readiness before a release or merge milestone.

## Type
Orchestrator. Controls workflow, gates, handoffs, and approval checkpoints. It should not invent canonical product content when a specialist skill owns that content.

## Required reading
- product/FRAMEWORK.md
- Relevant context.md files for the requested scope.
- .product/state.json, .product/roadmap.json, and .product/decisions.json when present.
- Approved decisions in product/knowledge/decisions/.

## Default sequence
Gap Finder -> Conflict Finder -> UX/UI audit when UI exists -> QA -> Code Review -> Security Review when sensitive -> Documentation Writer

## Operating rules
1. Identify the current artifact status before routing work.
2. Route work to the smallest specialist skill that owns the next artifact.
3. Stop at approval gates when scope, architecture, data, security, roadmap, or release risk changes.
4. Preserve traceability from parent artifacts to child artifacts.
5. Keep audit findings separate from product decisions until approved.
6. Update context and indexes only after the source artifact is approved.

## Outputs
release verdict; evidence; blockers; required fixes; release notes inputs.

## Gate checklist
- [ ] Inputs are approved or explicitly marked as draft.
- [ ] Parent and child artifacts are linked.
- [ ] Required design artifacts are approved or explicitly marked Not applicable.
- [ ] Delivery Level and Priority are consistent across released artifacts.
- [ ] Decisions are recorded or queued for approval.
- [ ] Gaps, conflicts, dependencies, and risks are visible.
- [ ] Handoff names the next skill and required reading.
- [ ] No downstream task starts from an unapproved Specification.

## Final handoff format
Status: draft | proposed | approved | blocked
Current step:
Next skill:
Approved artifacts:
Open questions:
Blocking findings:
Decision candidates:
Updated files:
