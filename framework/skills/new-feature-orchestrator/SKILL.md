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
Impact Analyzer -> Feature -> Use Case -> Specification -> UX/UI -> Implementation Planner -> Execution Graph -> Task Generator

## Operating rules
1. Identify the current artifact status before routing work.
2. Route work to the smallest specialist skill that owns the next artifact.
3. Stop at approval gates when scope, architecture, data, security, roadmap, or release risk changes.
4. Preserve traceability from parent artifacts to child artifacts.
5. Keep audit findings separate from product decisions until approved.
6. Require `Delivery Level` and `Priority` before Specification and preserve them through every downstream artifact.
7. For UI-bearing use cases, require approved `design.md` before Implementation Planner. For non-UI use cases, require `design.md` marked `Not applicable`.
8. Update context and indexes only after the source artifact is approved.

## Outputs
approved feature scope; specification; design; implementation plan; execution graph; task set; approval log.

## Gate checklist
- [ ] Inputs are approved or explicitly marked as draft.
- [ ] Parent and child artifacts are linked.
- [ ] Delivery Level and Priority are declared and justified.
- [ ] Design is approved or explicitly marked Not applicable before planning.
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
