---
name: product-historian
description: "Product Historian Skill. Use when an agent needs to Record important product and architecture decisions, their context, consequences, and supersession links in the Spec Framework workflow, including creating, updating, auditing, explaining, routing, or handing off related product artifacts."
---

# Product Historian Skill

## Layer
Governance

## Responsibility
Record important product and architecture decisions, their context, consequences, and supersession links.

## Operating modes
- create: produce the first version of the artifact when this skill is generative.
- update: revise an existing artifact while preserving approved decisions.
- audit: find gaps, conflicts, dependencies, and missing approvals.
- explain: summarize the artifact, finding, or decision in plain language.

## Inputs
Decision proposal; affected artifacts; prior decisions; trade-off analysis; approval notes.

## Outputs
Decision records; updated decisions index; supersede notes; historical rationale.

## Required reading
- the framework root's `FRAMEWORK.md`
- Relevant parent context.md files.
- This skill owns its generation resources: `assets/decision-template.md` and `assets/approval-record-template.json`.
- Approved decisions are discovered through the active product root's `.product/decisions.json`; resolve each registered `path` from its declared domain root (`knowledge/decisions/`, `design/decisions/`, or `engineering/decisions/`).

## Workflow
Use `scripts/check-decisions.ps1` on Windows or `scripts/check-decisions.sh` on macOS/Linux for the read-only indexed decision check.

1. Classify each DEC as product, architecture, security, data, or delivery; declare artifact/path scope and structured workflow effects in the decision index.
2. Keep architectural ADRs in the canonical DEC system rather than creating a parallel ADR store.
3. For legacy indexes, preview the guided migration, require review of ambiguous inferred types/scopes, preserve the original backup, and never change decision content/status or create approval records.
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
Next: END

Pass forward approved artifacts, findings, open questions, decisions, dependencies, risks, and required follow-up work.
