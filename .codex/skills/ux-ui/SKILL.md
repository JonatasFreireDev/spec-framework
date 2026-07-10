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
design.md; UX flow; UI states; accessibility notes; empty/loading/error states; design risks; Not applicable rationale when there is no UI.

## Required reading
- the framework root's `FRAMEWORK.md`
- Relevant parent context.md files.
- Relevant templates in knowledge/templates/.
- Approved product decisions in the active product root's `knowledge/decisions/` and `.product/decisions.json`.

## Workflow
1. Read the relevant context and identify artifact status.
2. Compare the artifact against the framework, template, and approved decisions.
3. Separate verified facts from assumptions and recommendations.
4. Report gaps, conflicts, dependencies, and risks with file-level evidence when possible.
5. Ask for approval before changing canonical product artifacts.
6. Update context.md or decision indexes only when the change is approved.

## Quality checklist
- [ ] Preserves traceability to affected artifacts.
- [ ] Uses the correct template and naming conventions.
- [ ] Produces approved design before Implementation Planner, or records `Not applicable` for non-UI work.
- [ ] Preserves Delivery Level and Priority from the source specification.
- [ ] Distinguishes blockers from suggestions.
- [ ] Detects gaps, conflicts, and dependencies.
- [ ] Records or requests decisions for meaningful changes.
- [ ] Leaves a clear handoff for the next skill or orchestrator.

## Handoff
Next: 10-implementation-planner.md

Pass forward approved artifacts, findings, open questions, decisions, dependencies, risks, and required follow-up work.
