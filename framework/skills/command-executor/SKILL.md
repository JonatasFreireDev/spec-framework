---
name: command-executor
description: "Command Executor Skill. Use when Codex needs to execute validated shell-free r0/r1 command plans with direct argv, confined cwd, sanitized environment, timeout, limited attempts, and evidence capture in the Spec Framework runtime."
---

# Command Executor Skill

## Layer
Execution

## Responsibility
Execute validated shell-free R0/R1 command plans with direct argv, confined cwd, sanitized environment, timeout, limited attempts, and evidence capture. It never invents commands or performs remote/destructive operations.

## Operating modes
- create: produce the first runtime artifact or route.
- update: refresh state while preserving approved contracts.
- audit: detect stale, unsafe, conflicting, or unauthorized runtime state.
- explain: summarize state, evidence, blockers, and next action.

## Inputs
Approved non-stale command plan; active lease; isolated worktree; policy.

## Outputs
Command evidence JSON; stdout/stderr digest or safe log reference; result; routed failure.

## Required reading
- [`execution-runtime.md`](../../docs/execution-runtime.md) for command execution limits, stale inputs, and scope controls.
- the framework root's `FRAMEWORK.md`
- Relevant parent context.md files.
- Relevant templates in framework/template/.
- Approved product decisions in the active product root's `knowledge/decisions/` and `.product/decisions.json`.

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
- [ ] Enforces the runtime authority boundary.
- [ ] Detects gaps, conflicts, and dependencies.
- [ ] Records or requests decisions for meaningful changes.
- [ ] Leaves a clear handoff for the next skill or orchestrator.

## Handoff
Next: code-runner, code-review, qa, or delivery-orchestrator according to the plan.

Pass forward workspace, task, hashes, evidence, decisions, dependencies, risks, attempts, blockers, and required reading.
