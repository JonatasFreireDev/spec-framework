---
name: commit-crafter
description: "Commit Crafter Skill. Use when Codex needs to turn verified working-tree changes into small atomic local commits without pushing."
---

# Commit Crafter Skill

## Layer

Delivery

## Responsibility

Create local commits from ready changes when explicitly asked.

Commit Crafter packages completed work into atomic commits by concern, following the active product root's `knowledge/conventions/commits.md`. It never pushes, merges, creates approval records, or hides failing gates.

## Operating Modes

- `commit`: create one or more local commits from a cleanly verified change set.
- `plan`: propose commit grouping without committing.
- `audit`: inspect whether changes are ready to commit.

## Required Reading

- the framework root's `FRAMEWORK.md`.
- the active product root's `knowledge/conventions/commits.md`.
- the active product root's `knowledge/conventions/gates.md`.
- Relevant task files and code evidence.
- QA Evidence, Code Review, and Security Review when the commit is part of validation/release work.
- Matching Code Review and QA diff hashes for working-tree implementation; refuse to commit when the diff changed after either gate.
- For parallel delivery, the current isolated task worktree and its active lease/checkpoint.
- Delivery commit policy in `FRAMEWORK.md`.

## Preconditions

- User explicitly asks to commit or the active orchestration step explicitly invokes Commit Crafter.
- The intended change set is understood and belongs to the current task/use case.
- Technical gates required by the active product root's `knowledge/conventions/gates.md` are green or the user explicitly accepts a draft exception.
- No secrets, credentials, tokens, private keys, or sensitive local artifacts are staged.
- The working tree has been reviewed for unrelated changes.

## Workflow

1. Inspect `git status` and the diff.
2. Identify unrelated changes and leave them unstaged unless the user explicitly includes them.
3. Group changes into atomic commits by concern:
   - framework docs;
   - skill contracts;
   - validator/tooling;
   - templates/conventions;
   - product documentation;
   - application code.
4. Separate documentation-only changes from application code when practical.
5. Scan staged content for obvious secrets and sensitive local artifacts.
6. Run or confirm required gates from the active product root's `knowledge/conventions/gates.md`.
7. Create local commits with messages following the active product root's `knowledge/conventions/commits.md`.
8. Record commit hashes back into task files only when the task status transition requires it and the user asked for that artifact update.
9. Return the immutable task commit to `integration-orchestrator`; do not cherry-pick it into the integration branch directly.
9. Do not push.

## Boundaries

- Do not stage unrelated user changes.
- Do not commit secrets or local environment files.
- Do not push or merge.
- Do not create, edit, or repair approval records.
- Do not mark QA, Code Review, Security Review, or release as approved.

## Quality Checklist

- [ ] Commit scope is clear and atomic.
- [ ] Unrelated changes were left alone.
- [ ] Required gates were run or limitations were reported.
- [ ] Staged diff was checked for secrets.
- [ ] Commit message follows the convention.
- [ ] No push or merge was performed.

## Handoff

Next: pr-finalizer or the human owner.

Pass forward commit hashes, commit messages, files included, gates run, skipped gates or limitations, and any remaining uncommitted changes.
