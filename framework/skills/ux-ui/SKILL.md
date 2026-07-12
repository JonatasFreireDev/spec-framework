---
name: ux-ui
description: "UX/UI Skill. Use when Codex needs to Translate an approved specification into flows, states, interaction rules, accessibility requirements, design handoff notes, and mockup requirements before implementation planning in the Spec Framework workflow, including creating, updating, auditing, explaining, routing, or handing off related product artifacts."
---

# UX/UI Skill

## Layer
Product Design

## Responsibility
Translate an approved specification into flows, states, interaction rules, accessibility requirements, design handoff notes, and mockup requirements before implementation planning.

## Operating modes
- create: produce the first version of the artifact when this skill is generative.
- update: revise an existing artifact while preserving approved decisions.
- audit: find gaps, conflicts, dependencies, and missing approvals.
- explain: summarize the artifact, finding, or decision in plain language.

## Inputs
Approved specification; source feature/use case; inherited Delivery Level and Priority; journey; design system; UX principles; platform constraints.

## Outputs
design.md; origin mode and visual maturity; versioned source references; screen inventory and requirement coverage; UX flow; UI states; accessibility notes; empty/loading/error states; fidelity policy and deviations; design risks; Not applicable rationale when there is no UI.

## Required reading
- the framework root's `FRAMEWORK.md`
- Relevant parent context.md files.
- Relevant templates in framework/template/.
- Approved product decisions in the active product root's `knowledge/decisions/` and `.product/decisions.json`.

## Workflow
1. Read the relevant context, the declared Design System when present, and identify artifact status.
2. Select or preserve `generate`, `evolve`, or `adopt`; never infer that a reference is canonical.
3. Register every source with authority, location, and immutable version or hash.
4. Compare screens and states against stable REQ/AC identifiers and the design system.
5. Separate verified facts from assumptions and recommendations.
6. Report missing states, conflicts, fidelity deviations, dependencies, and risks before proposing changes.
7. Keep design prototypes under `product/design/` and mark them non-production.
8. Ask for approval before changing canonical product artifacts or modifying an adopted strict-fidelity source.
9. Hand off to independent UX Review before Design approval.

## Quality checklist
- [ ] Preserves traceability to affected artifacts.
- [ ] Uses the correct template and naming conventions.
- [ ] Produces approved design before Implementation Planner, or records `Not applicable` for non-UI work.
- [ ] Preserves Delivery Level and Priority from the source specification.
- [ ] Distinguishes blockers from suggestions.
- [ ] Detects gaps, conflicts, and dependencies.
- [ ] Records or requests decisions for meaningful changes.
- [ ] Versioned sources, authority, maturity, coverage, and deviations are explicit.
- [ ] Pins the Design System id/version and records consumed tokens, components, patterns, and deviations when a system is declared.
- [ ] Does not treat imported or generated visual evidence as approval.
- [ ] Leaves a clear handoff for the next skill or orchestrator.

## Handoff
Next: technical-discovery.

Pass forward approved artifacts, findings, open questions, decisions, dependencies, risks, and required follow-up work.
