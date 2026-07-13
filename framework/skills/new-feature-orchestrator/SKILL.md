---
name: new-feature-orchestrator
description: "New Feature Orchestrator. Use when Codex needs to Move a proposed feature from idea to executable task graph in the Spec Framework workflow, including creating, updating, auditing, explaining, routing, or handing off related product artifacts."
---

# New Feature Orchestrator

## Mission
Move a proposed feature from idea to executable task graph.

## Type
Orchestrator. Controls workflow, gates, handoffs, and approval checkpoints. It should not invent canonical product content when a specialist skill owns that content.

## Required reading
- the framework root's `FRAMEWORK.md`
- Relevant context.md files for the requested scope.
- the active product root's `.product/state.json`, `.product/roadmap.json`, and `.product/decisions.json` when present.
- Approved product decisions in the active product root's `knowledge/decisions/`.

## Default sequence
Impact Analyzer -> Feature -> Use Case -> Specification -> UX/UI -> Technical Discovery -> Architecture Gate -> Engineering Proposal -> Engineering Review -> Implementation Planner -> Execution Graph -> Task Generator

## Discovery and challenge

Follow the shared [Discovery And Challenge contract](../discovery-and-challenge.md) when resolving scope, route, or approval questions.

## Operating rules
1. Identify the current artifact status before routing work.
2. Route work to the smallest specialist skill that owns the next artifact.
3. Stop at approval gates when scope, architecture, data, security, roadmap, or release risk changes.
4. Preserve traceability from parent artifacts to child artifacts.
5. Keep audit findings separate from product decisions until approved.
6. Require `Delivery Level` and `Priority` before Specification and preserve them through every downstream artifact.
7. For UI-bearing use cases, require approved `design.md` before Technical Discovery. For non-UI use cases, require structured `not_applicable` status and rationale.
8. Require Engineering Proposal and a passed Engineering Review for Tier L and any delivery whose context declares a supported `engineering_trigger`.
8. Update context and indexes only after the source artifact is approved.

## Outputs
approved feature scope; specification contracts; design; technical discovery; resolved Architecture Gate; applicable Engineering Proposal and Engineering Review; implementation plan; execution graph; task set; approval evidence references.

## Gate checklist
- [ ] Inputs are approved or explicitly marked as draft.
- [ ] Parent and child artifacts are linked.
- [ ] Delivery Level and Priority are declared and justified.
- [ ] Design is approved or uses structured `not_applicable` status with rationale before planning.
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
