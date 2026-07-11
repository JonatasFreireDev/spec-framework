---
name: impact-analyzer
description: "Impact Analyzer Skill. Use when Codex needs to Analyze what changes when a proposed product, UX, technical, or business decision changes in the Spec Framework workflow, including creating, updating, auditing, explaining, routing, or handing off related product artifacts."
---

# Impact Analyzer Skill

## Layer
Audit

## Responsibility
Analyze what changes when a proposed product, UX, technical, or business decision changes.

## Operating modes
- create: produce the first version of the artifact when this skill is generative.
- update: revise an existing artifact while preserving approved decisions.
- audit: find gaps, conflicts, dependencies, and missing approvals.
- explain: summarize the artifact, finding, or decision in plain language.

## Inputs
Change proposal; affected artifacts; decisions; dependency graph; roadmap; release context.

## Outputs
Impact report; affected artifacts; risk level; migration or update plan; approval questions.

## Required reading
- the framework root's `FRAMEWORK.md`
- Relevant parent context.md files.
- Relevant templates in framework/template/.
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
- [ ] Distinguishes blockers from suggestions.
- [ ] Detects gaps, conflicts, and dependencies.
- [ ] Records or requests decisions for meaningful changes.
- [ ] Leaves a clear handoff for the next skill or orchestrator.

## Handoff
Next: evolution.

Pass forward approved artifacts, findings, open questions, decisions, dependencies, risks, and required follow-up work.
