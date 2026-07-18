---
name: evolution-orchestrator
description: "Evolution Orchestrator. Use when an agent needs to Turn feedback, metrics, and audit findings into approved improvement candidates in the Spec Framework workflow, including creating, updating, auditing, explaining, routing, or handing off related product artifacts."
---

# Evolution Orchestrator

## Mission
Turn feedback, metrics, and audit findings into approved improvement candidates.

## Type
Orchestrator. Controls workflow, gates, handoffs, and approval checkpoints. It should not invent canonical product content when a specialist skill owns that content.

## Required reading
- the framework root's `FRAMEWORK.md`
- Relevant context.md files for the requested scope.
- the active product root's `.product/state.json`, `.product/roadmap.json`, and `.product/decisions.json` when present.
- Approved decisions are discovered through the active product root's `.product/decisions.json`; resolve each registered `path` from its declared domain root.

## Default sequence
Demand intake/import -> Evolution -> Impact Analyzer -> affected generative skill -> Documentation Writer

## Discovery and challenge

Follow the shared [Discovery And Challenge contract](../discovery-and-challenge.md) when comparing candidates or requesting human selection.

## Operating rules
1. Identify the current artifact status before routing work.
2. Load the nearest context.md, its parent chain, sibling delivery artifacts, approved decisions, and shared Engineering and Design baselines.
3. Route the demand to the smallest specialist skill that owns the proposed destination.
4. Stop for human selection when the demand could be an extension, new Use Case, new Feature, new Goal, or new Domain.
5. Stop at approval gates when scope, architecture, data, security, roadmap, or release risk changes.
6. Preserve traceability from the source demand through context relations and child artifacts.
7. Keep audit findings separate from product decisions until approved.
8. Update context and indexes only after the source artifact is approved.

## Outputs
ranked proposals; approval questions; updated specs or roadmap candidates; rejected ideas log.

## Gate checklist
- [ ] Inputs are approved or explicitly marked as draft.
- [ ] Parent and child artifacts are linked.
- [ ] Demand classification and reuse/impact relations are explicit.
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
