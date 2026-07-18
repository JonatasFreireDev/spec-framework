---
name: product-orchestrator
description: "Product Orchestrator. Use when an agent needs to Create a product foundation from zero or reset an incomplete foundation in the Spec Framework workflow, including creating, updating, auditing, explaining, routing, or handing off related product artifacts."
---

# Product Orchestrator

## Mission
Create a product foundation from zero or reset an incomplete foundation.

## Type
Orchestrator. Controls workflow, gates, handoffs, and approval checkpoints. It should not invent canonical product content when a specialist skill owns that content.

## Required reading
- the framework root's `FRAMEWORK.md`
- Relevant context.md files for the requested scope.
- This skill owns `assets/product-baseline-template.md`, `assets/implementation-assessment-template.md`, and `assets/feature-brief-template.md` for the supported starting points.
- the active product root's `.product/state.json`, `.product/roadmap.json`, and `.product/decisions.json` when present.
- Approved decisions are discovered through the active product root's `.product/decisions.json`; resolve each registered `path` from its declared domain root.

## Default sequence
Problem Discovery -> Vision -> Strategy -> Domain Architect -> User Goal -> Roadmap alignment

For a repository whose canonical manifest declares `starting_point: existing-feature`, replace that sequence with Feature Brief -> explicit individual approval -> existing Feature selection -> workspace. Escalate to the default sequence when the feature cannot be bounded without product-wide decisions.

For `starting_point: existing-implementation`, prepend Implementation Assessment -> explicit individual approval to the default Foundation sequence. Use its observations as evidence and keep inferred product claims unapproved until the owning Foundation specialist resolves them.

For `starting_point: existing-product`, use Product Baseline -> explicit individual approval -> Strategy -> explicit individual approval -> Domain and delivery. Code and operational evidence may establish current state; Strategy remains the explicit future-facing decision contract.

## Discovery and challenge

Follow the shared [Discovery And Challenge contract](../discovery-and-challenge.md) when eliciting foundation choices or resolving a blocking route.

## Operating rules
1. Identify the current artifact status before routing work.
2. Route work to the smallest specialist skill that owns the next artifact.
3. Stop at approval gates when scope, architecture, data, security, roadmap, or release risk changes.
4. Preserve traceability from parent artifacts to child artifacts.
5. Keep audit findings separate from product decisions until approved.
6. Update context and indexes only after the source artifact is approved.
7. Do not require or synthesize full Foundation artifacts for a bounded `existing-feature` route; own `foundation/feature-brief.md` as the proportional entry contract.
8. For `existing-implementation`, own `knowledge/assessments/implementation-assessment.md`, do not modify application code during assessment, and retain the full Foundation sequence before workspace creation.
9. For `existing-product`, own `foundation/product-baseline.md`, keep uncertain intent visible, and promote to the full Foundation sequence when the baseline cannot establish audience and delivered value confidently.

## Outputs
foundation artifacts; domain map; initial goal catalog; roadmap with Delivery Levels and Priorities; open decisions; next recommended feature slices.

## Gate checklist
- [ ] Inputs are approved or explicitly marked as draft.
- [ ] Parent and child artifacts are linked.
- [ ] Roadmap candidates have Delivery Level and Priority before feature planning starts.
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
