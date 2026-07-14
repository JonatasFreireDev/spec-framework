---
name: execution-scheduler
description: "Execution Scheduler Skill. Use when Codex needs to calculate deterministic parallel waves from the dag, write scopes, shared resources, capabilities, leases, priority, and capacity in the Spec Framework runtime."
---

# Execution Scheduler Skill

## Layer
Planning

## Responsibility
Calculate deterministic parallel waves from the DAG, write scopes, shared resources, capabilities, leases, priority, and capacity. It does not spawn agents or execute tasks.

## Operating modes
- create: produce the first runtime artifact or route.
- update: refresh state while preserving approved contracts.
- audit: detect stale, unsafe, conflicting, or unauthorized runtime state.
- explain: summarize state, evidence, blockers, and next action.

## Inputs
Approved execution graph; workspace; active leases; resource capacities; parallelism limit.

## Outputs
Versioned wave plan; ready/serialized tasks; conflict rationale; capability gaps.

## Required reading
- [`execution-runtime.md`](../../docs/execution-runtime.md) for leases, scheduling, conflict detection, and non-execution boundaries.
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
Next: delivery-orchestrator or task agents for the ready wave.

Pass forward workspace, task, hashes, evidence, decisions, dependencies, risks, attempts, blockers, and required reading.
