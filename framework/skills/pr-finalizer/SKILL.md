---
name: pr-finalizer
description: "PR Finalizer Skill. Use when Codex needs to prepare or open a pull request after gates, QA, Code Review, and Security Review are ready, without merging."
---

# PR Finalizer Skill

## Layer

Delivery

## Responsibility

Prepare or open a pull request for validated work without merging it.

PR Finalizer verifies hard preconditions, links evidence, writes PR metadata back to task files when appropriate, and stops before merge.

## Operating Modes

- `prepare`: produce a PR title/body and checklist without opening a PR.
- `open`: open a PR when `gh` or the configured provider is available and authenticated.
- `audit`: inspect whether a branch is ready for PR.

## Required Reading

- the framework root's `FRAMEWORK.md`.
- the active product root's `knowledge/conventions/pull-requests.md`.
- the active product root's `knowledge/conventions/gates.md`.
- Relevant task files with branch, commits, and code paths.
- QA Evidence.
- Code Review.
- Security Review when applicable.
- Audit/release readiness reports when present.
- The integration record and Integrated QA evidence when multiple task branches were combined.
- `framework/decisions/FDR-008-delivery-commits-and-prs.md`.

## Hard Preconditions

- Branch is known and contains the intended commits.
- Required gates are green.
- QA Evidence is approved or the PR is explicitly draft/prototype.
- Code Review is approved or the PR is explicitly draft/prototype.
- Security Review has no blocker when applicable.
- Integrated QA passed over the current integrated diff hash when integration was required.
- Blocking findings have route and owner.
- PR body links evidence: tasks, QA Evidence, Code Review, Security Review when applicable, gate logs, screenshots, and CI URL when available.

If any hard precondition fails, stop and report the blocker.

## Workflow

1. Inspect branch, commits, diff, and task evidence.
2. Re-run or confirm gates required by the active product root's `knowledge/conventions/gates.md`.
3. Verify QA Evidence, Code Review, and Security Review readiness.
4. Build a PR title and body from the active product root's `knowledge/conventions/pull-requests.md`.
5. Open the PR only when the user asked for it and the provider tool is available.
6. Record the PR URL/id back into relevant task files when required by status transition.
7. Do not merge.

## Boundaries

- Do not merge.
- Do not bypass red QA, Code Review, or Security Review gates.
- Do not create approval records.
- Do not invent missing evidence.
- Do not push unless the user explicitly asks and the branch is the intended delivery branch.

## Quality Checklist

- [ ] PR includes links to tasks and evidence.
- [ ] Gates are green or limitation is explicit for draft/prototype.
- [ ] QA Evidence is ready.
- [ ] Code Review is ready.
- [ ] Security Review has no blocker when applicable.
- [ ] Task PR fields are updated only when evidence is real.
- [ ] Merge is left to the human/repository policy.

## Handoff

Next: human reviewer, release-orchestrator, or product-historian for missing decisions.

Pass forward PR URL/id or prepared PR body, branch, commits, evidence links, unresolved blockers, and release readiness notes.
