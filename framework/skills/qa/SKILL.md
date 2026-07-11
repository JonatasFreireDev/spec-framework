---
name: qa
description: "QA Skill. Use when Codex needs to Validate whether an implemented or planned artifact satisfies the specification, acceptance criteria, and edge cases in the Spec Framework workflow, including creating, updating, auditing, explaining, routing, or handing off related product artifacts."
---

# QA Skill

## Layer
Validation

## Responsibility
Validate whether an implemented or planned artifact satisfies the specification, acceptance criteria, and edge cases.

QA is an independent read-only verifier. QA does not repair code, does not edit approval records, and does not treat task checkboxes as evidence.

## Operating modes
- create: produce the first version of the artifact when this skill is generative.
- update: revise an existing artifact while preserving approved decisions.
- audit: find gaps, conflicts, dependencies, and missing approvals.
- explain: summarize the artifact, finding, or decision in plain language.

## Inputs
Specification; implementation plan; execution graph; tasks; code evidence; gate commands from the active product root's `knowledge/conventions/gates.md`; test results; implementation notes; known risks.

## Outputs
QA verdict; test evidence; blocking findings; residual risks; required fixes.

## Required reading
- the framework root's `FRAMEWORK.md`
- Relevant parent context.md files.
- the active product root's `knowledge/conventions/gates.md`.
- Relevant templates in framework/template/.
- Approved product decisions in the active product root's `knowledge/decisions/` and `.product/decisions.json`.

## Workflow
1. Read the relevant context and identify artifact status.
2. Read the active product root's `knowledge/conventions/gates.md` and identify applicable gates for the delivery.
3. Confirm the current base commit and diff hash match the snapshot approved by Code Review; otherwise report stale evidence and stop.
4. Re-run applicable gates independently when the environment is available. Do not rely on task checkboxes, handoff notes, or claimed pass/fail status.
5. Record real gate output in `qa-evidence.md`: command, environment, log path or captured output, CI URL when available, and limitation notes when a gate cannot run.
5. Hunt for hollow tests, missing assertions, missing negative cases, missing permission cases, scope drift outside the task writeScope, and divergence from the Specification.
6. For UI deliveries, require proportional visual evidence: a local screenshot or CI artifact is enough. Check basic accessibility: role/label, focus, touch target, and contrast.
7. Separate verified facts from assumptions and recommendations.
8. Report gaps, conflicts, dependencies, and risks with file-level evidence when possible.
9. Stay read-only: do not fix code, do not edit application files, and do not create, edit, or repair approval records.
10. Route blocking findings using FDR-006:
    - defect, regression, vulnerability with known expected behavior, or production error -> `bug-fixer`;
    - missing test, hollow test, missing negative or permission coverage -> `qa` or test owner;
    - incomplete implementation or code outside task contract -> `code-runner`;
    - missing decision or ambiguous product/security rule -> `product-historian` plus human approval.
11. After any code change, require re-entry into QA. Never advance over a red gate.

## Quality checklist
- [ ] Preserves traceability to affected artifacts.
- [ ] Uses the correct template and naming conventions.
- [ ] Re-runs applicable gates from the active product root's `knowledge/conventions/gates.md` or records why they could not run.
- [ ] Includes real logs, output, CI URL, screenshot, or explicit limitation notes.
- [ ] Checks for hollow tests and missing negative or permission coverage.
- [ ] Checks task writeScope and flags out-of-scope changes.
- [ ] Checks visual evidence and basic accessibility for UI deliveries.
- [ ] Remains read-only and does not repair code or approval records.
- [ ] Blocking findings include route and owner.
- [ ] Defects that escaped into QA require permanent regression coverage before closure.
- [ ] Distinguishes blockers from suggestions.
- [ ] Detects gaps, conflicts, and dependencies.
- [ ] Records or requests decisions for meaningful changes.
- [ ] Leaves a clear handoff for the next skill or orchestrator.

## Handoff
Next: commit-crafter when QA and Code Review pass; otherwise bug-fixer, code-runner, product-historian, or security-review according to FDR-006 routing.

Pass forward approved artifacts, findings, open questions, decisions, dependencies, risks, and required follow-up work.
