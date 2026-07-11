---
name: command-planner
description: "Command Planner Skill. Use when Codex needs to own an immutable command plan derived only from approved gates, tasks, runbooks, repository scripts, or explicit human commands in the Spec Framework runtime."
---

# Command Planner Skill

## Layer
Planning

## Responsibility
Own an immutable command plan derived only from approved gates, tasks, runbooks, repository scripts, or explicit human commands. It never executes commands.

## Operating modes
- create: produce the first runtime artifact or route.
- update: refresh state while preserving approved contracts.
- audit: detect stale, unsafe, conflicting, or unauthorized runtime state.
- explain: summarize state, evidence, blockers, and next action.

## Inputs
Approved task; gates; runbooks; base commit; worktree; allowed write scope.

## Outputs
CMDPLAN-NNN.json with argv, cwd, source, risk, timeout, expected exit code, hashes, and evidence contract.

## Required reading
- the framework root's `FRAMEWORK.md`
- Relevant parent context.md files.
- Relevant templates in framework/template/.
- Approved product decisions in the active product root's `knowledge/decisions/` and `.product/decisions.json`.
- FDR-016 and FDR-017.

## Workflow
1. Load the workspace, current commit, approved artifacts, and runtime policy.
2. Validate hashes, ownership, dependencies, attempts, and authority before acting.
3. Refuse stale inputs, scope escapes, conflicting resources, and unsupported risk levels.
4. Produce or update only the runtime artifact owned by this skill.
5. Persist evidence and an explicit structured handoff; never rely on chat history.
6. Stop at human approval, remote mutation, conflict, security, or attempt-limit gates.

## Quality checklist
- [ ] Preserves traceability to affected artifacts.
- [ ] Uses the correct template and naming conventions.
- [ ] Records base commit, hashes, owner, timestamps, attempts, and blockers.
- [ ] Enforces the FDR-017 authority boundary.
- [ ] Detects gaps, conflicts, and dependencies.
- [ ] Records or requests decisions for meaningful changes.
- [ ] Leaves a clear handoff for the next skill or orchestrator.

## Handoff
Next: command-executor after plan validation and required approval.

Pass forward workspace, task, hashes, evidence, decisions, dependencies, risks, attempts, blockers, and required reading.
