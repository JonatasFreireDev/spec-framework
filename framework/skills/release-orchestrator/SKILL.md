---
name: release-orchestrator
description: "Release Orchestrator. Use when Codex needs to Verify readiness before a release or merge milestone in the Spec Framework workflow, including creating, updating, auditing, explaining, routing, or handing off related product artifacts."
---

# Release Orchestrator

## Mission
Verify readiness before a release or merge milestone.

## Type
Orchestrator. Controls workflow, gates, handoffs, and approval checkpoints. It should not invent canonical product content when a specialist skill owns that content.

## Required reading
- the framework root's `FRAMEWORK.md`
- Relevant context.md files for the requested scope.
- the active product root's `.product/state.json`, `.product/roadmap.json`, and `.product/decisions.json` when present.
- Approved product decisions in the active product root's `knowledge/decisions/`.

## Default sequence
Gap Finder -> Conflict Finder -> UX/UI audit when UI exists -> QA -> Code Review -> Security Review when sensitive -> Documentation Writer

## Operating rules
1. Identify the current artifact status before routing work.
2. Route work to the smallest specialist skill that owns the next artifact.
3. Stop at approval gates when scope, architecture, data, security, roadmap, or release risk changes.
4. Preserve traceability from parent artifacts to child artifacts.
5. Keep audit findings separate from product decisions until approved.
6. Update context and indexes only after the source artifact is approved.
7. Route red gates with FDR-006: defect to bug-fixer, missing test to QA/tests, incomplete implementation to code-runner, and missing decision to Product Historian plus human approval.
8. Enforce QA re-entry after any code change. Never advance over a red QA or Security Review gate.

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
- [ ] Blocking findings have route and owner.
- [ ] No red QA or Security Review gate is bypassed.
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
