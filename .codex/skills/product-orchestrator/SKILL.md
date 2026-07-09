---
name: product-orchestrator
description: "Product Orchestrator. Use when Codex needs to Create a product foundation from zero or reset an incomplete foundation in this repository's Spec Framework workflow, including creating, updating, auditing, explaining, routing, or handing off related product artifacts."
---

# Product Orchestrator

## Mission
Create a product foundation from zero or reset an incomplete foundation.

## Type
Orchestrator. Controls workflow, gates, handoffs, and approval checkpoints. It should not invent canonical product content when a specialist skill owns that content.

## Required reading
- FRAMEWORK.md
- Relevant context.md files for the requested scope.
- .product/state.json, .product/roadmap.json, and .product/decisions.json when present.
- Approved decisions in knowledge/decisions/.

## Default sequence
Problem Discovery -> Vision -> Strategy -> Domain Architect -> User Goal -> Roadmap alignment

## Operating rules
1. Identify the current artifact status before routing work.
2. Route work to the smallest specialist skill that owns the next artifact.
3. Stop at approval gates when scope, architecture, data, security, roadmap, or release risk changes.
4. Preserve traceability from parent artifacts to child artifacts.
5. Keep audit findings separate from product decisions until approved.
6. Update context and indexes only after the source artifact is approved.

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
