---
name: code-review
description: "Code Review Skill. Use when an agent needs to perform a read-only implementation review for completeness, adherence, and quality before validation or release."
---

# Code Review Skill

## Layer

Validation

## Responsibility

Review implemented code and evidence without editing files.

Code Review verifies that implementation matches the Specification, tasks, architecture, and project conventions. It reports findings and routes fixes through the framework's fixed failure-routing policy. It does not repair code, approve QA, create approval records, commit, push, merge, or open PRs.

## Operating Modes

- `review`: review a completed task, use case, or PR surface.
- `audit`: inspect a bundle for review readiness.
- `explain`: summarize review findings and routing.

## Required Reading
- [`lifecycle-and-approvals.md`](../../docs/lifecycle-and-approvals.md) for diff-hash, authority, and failure-routing boundaries.

- the framework root's `FRAMEWORK.md`.
- This skill owns `assets/code-review-template.md` and `assets/review-finding-template.json`.
- Relevant `context.md`.
- Specification, Design, Implementation Plan, Execution Graph, Tasks, Tests, QA Evidence, and Security Review when present.
- the active product root's `knowledge/conventions/gates.md`.
- Product decisions are discovered through the active product root's `.product/decisions.json`; resolve each registered `path` from its declared domain root.
- Failure routing in `FRAMEWORK.md`.
- The diff, branch, commits, PR, or code paths referenced by the task files.

## Review Lenses

| Lens | Question |
| --- | --- |
| Completeness | Did the implementation deliver everything required by the Specification, tasks, acceptance criteria, and evidence contract? |
| Adherence | Does the implementation follow the approved architecture, contracts, data model, permissions, routing, and non-goals without drift? |
| Quality | Is the code maintainable, minimal, tested, reusable where appropriate, safe around errors, performant enough for the delivery level, and free of obvious dead code or unsafe behavior? |

## Finding Severity

| Severity | Meaning |
| --- | --- |
| `blocker` | Must be fixed before validation/release. |
| `required_fix` | Must be fixed before the reviewed artifact can be considered approved. |
| `note` | Non-blocking observation or follow-up. |

## Workflow

1. Confirm the review target and status.
2. Read the source artifacts and code evidence.
3. Record the base commit and normalized diff hash; review exactly that immutable working-tree snapshot and mark prior review stale when it changes.
4. For an isolated task, verify the lease owner, worktree branch, checkpoint input hash, and declared `writeScope` before reviewing its diff.
4. Inspect the diff/code paths against the three lenses.
5. Classify findings as `blocker`, `required_fix`, or `note`.
5. Cite file and line when possible.
6. Route each blocking or required finding using the fixed failure-routing policy:
   - defect/regression -> `bug-fixer`;
   - incomplete implementation or drift from task contract -> `code-runner`;
   - missing tests/evidence -> `qa`;
   - missing decision -> `product-historian` plus human approval.
7. Produce or update `code-review.md` with verdict, findings, route, owner, and residual risk.
8. Stay read-only. Do not fix code.

## Verdict Rules

- `passed`: no blocker or required fix remains.
- `passed_with_notes`: no blocker remains, but notes or non-blocking follow-ups exist.
- `blocked`: any blocker or required fix remains.

## Quality Checklist

- [ ] Review is read-only.
- [ ] Completeness, adherence, and quality lenses were applied.
- [ ] Findings cite file/line when possible.
- [ ] Blockers and required fixes have route and owner.
- [ ] Missing decisions are routed to Product Historian and human approval.
- [ ] Verdict matches remaining findings.

## Handoff

Next: qa when passed; otherwise bug-fixer, code-runner, or product-historian according to findings.

Pass forward the review target, verdict, findings, routes, owners, files reviewed, residual risks, and whether validation can proceed.
