---
name: subagent-return-reviewer
description: "Subagent Return Reviewer Skill. Use when Codex needs to validate a bounded subagent return against its envelope, hashes, evidence, and route in the Spec Framework workflow."
---

# Subagent Return Reviewer Skill

## Layer
Validation

## Responsibility
Owns `dispatch-return.md` validation and routing. It does not fix code, approve artifacts, accept residual risk, or close external reviews.

## Operating modes
- create: record a validated return.
- update: refresh a return after new evidence.
- audit: find stale hashes, missing evidence, and scope violations.
- explain: summarize return readiness and blockers.

## Inputs
Dispatch envelope; transcript; task/chunk context; hashes; evidence; findings.

## Outputs
`dispatch-return.md`; route; blockers; stale-return verdict.

## Required reading
- the framework root's `FRAMEWORK.md`
- Relevant parent context.md files.
- Relevant templates in framework/template/.
- Approved product decisions in the active product root's `knowledge/decisions/` and `.product/decisions.json`.

## Workflow
1. Verify agent, unit, input hash, diff hash when applicable, scope, and evidence.
2. Reject returns with missing evidence, stale diff, or forbidden operation.
3. Preserve unresolved gaps and route to the independent owner.
4. Record the return without changing canonical task approval or lifecycle.

## Quality checklist
- [ ] Preserves traceability to affected artifacts.
- [ ] Uses the correct template and naming conventions.
- [ ] Return matches its envelope and current diff.
- [ ] Findings have route and owner.
- [ ] Detects gaps, conflicts, and dependencies.
- [ ] Records or requests decisions for meaningful changes.
- [ ] Leaves a clear handoff for the next skill or orchestrator.

## Handoff
Next: qa, code-review, security-review, artifact-importer, code-runner, or product-historian.

Pass forward return, hashes, evidence, blockers, findings, risks, and required follow-up work.
