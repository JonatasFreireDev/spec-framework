---
name: qa
description: "QA Skill. Use when an agent needs to Validate whether an implemented or planned artifact satisfies the specification, acceptance criteria, and edge cases in the Spec Framework workflow, including creating, updating, auditing, explaining, routing, or handing off related product artifacts."
---

# QA Skill

## Layer
Validation

## Responsibility
Validate whether an implemented or planned artifact satisfies the specification, acceptance criteria, and edge cases.

QA is an independent read-only verifier. QA does not repair code, does not edit approval records, and does not treat task checkboxes as evidence.

## Operating modes
- task-qa: independently validate one isolated task diff before its local commit is eligible for integration.
- integrated-qa: validate the combined integration diff and cross-task behavior before PR finalization; task-level passes do not substitute for this mode.
- create: produce the first version of the artifact when this skill is generative.
- update: revise an existing artifact while preserving approved decisions.
- audit: find gaps, conflicts, dependencies, and missing approvals.
- explain: summarize the artifact, finding, or decision in plain language.

## Inputs
Specification; pinned Engineering Quality System when configured; delivery `tests.md`; implementation plan; execution graph; tasks; code evidence; gate commands from the active product root's `knowledge/conventions/gates.md`; test results; implementation notes; known risks.

## Outputs
QA verdict; test evidence; blocking findings; residual risks; required fixes.

## Required reading
- [`engineering-systems.md`](../../docs/engineering-systems.md) and [`lifecycle-and-approvals.md`](../../docs/lifecycle-and-approvals.md).
- the framework root's `FRAMEWORK.md`
- Relevant parent context.md files.
- the active product root's `knowledge/conventions/gates.md`.
- This skill owns its generation resources: `assets/tests-template.md`, `assets/qa-evidence-template.md`, and `assets/analytics-template.md`.
- Approved decisions are discovered through the active product root's `.product/decisions.json`; resolve each registered `path` from its declared domain root (`knowledge/decisions/`, `design/decisions/`, or `engineering/decisions/`).

## Workflow
1. Read the relevant context and identify artifact status.
2. Read the active product root's `knowledge/conventions/gates.md` and identify applicable gates for the delivery.
3. When the Engineering Quality System is configured, verify that `tests.md` pins its current id/version, selects configured environment, data, and platform values, applies its risk, coverage, and evidence policies, and declares `None` or only open, unexpired, in-scope exceptions.
4. Confirm the current base commit and diff hash match the snapshot approved by Code Review; otherwise report stale evidence and stop.
5. In `integrated-qa`, confirm every integrated commit passed task QA and Code Review, then test the combined diff, dependency boundaries, migrations, shared resources, and cross-task regressions.
6. Re-run applicable gates independently when the environment is available. Do not rely on task checkboxes, handoff notes, or claimed pass/fail status.
7. Record real gate output in `qa-evidence.md`: command, environment, log path or captured output, CI URL when available, and limitation notes when a gate cannot run.
8. Hunt for hollow tests, missing assertions, missing negative cases, missing permission cases, scope drift outside the task writeScope, and divergence from the Specification.
9. For UI deliveries, require proportional visual evidence: a local screenshot or CI artifact is enough. Check basic accessibility: role/label, focus, touch target, and contrast.
10. Separate verified facts from assumptions and recommendations.
11. Report gaps, conflicts, dependencies, and risks with file-level evidence when possible.
12. Stay read-only: do not fix code, do not edit application files, and do not create, edit, or repair approval records.
13. Route blocking findings using the fixed failure-routing policy:
    - defect, regression, vulnerability with known expected behavior, or production error -> `bug-fixer`;
    - missing test, hollow test, missing negative or permission coverage -> `qa` or test owner;
    - incomplete implementation or code outside task contract -> `code-runner`;
    - missing decision or ambiguous product/security rule -> `product-historian` plus human approval.
14. After any code change, require re-entry into QA. Never advance over a red gate.

## Quality checklist
- [ ] Preserves traceability to affected artifacts.
- [ ] Uses the correct template and naming conventions.
- [ ] Re-runs applicable gates from the active product root's `knowledge/conventions/gates.md` or records why they could not run.
- [ ] Verifies the pinned Engineering Quality System and any declared deviations when configured.
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
Next: commit-crafter when QA and Code Review pass; otherwise bug-fixer, code-runner, product-historian, or security-review according to the fixed failure-routing policy.

Pass forward approved artifacts, findings, open questions, decisions, dependencies, risks, and required follow-up work.
