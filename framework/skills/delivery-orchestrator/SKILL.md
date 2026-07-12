---
name: delivery-orchestrator
description: "Delivery Orchestrator Skill. Use when Codex needs to route one workspace through allowed orchestrators using persisted state, handoffs, checkpoints, blockers, and attempt limits in the Spec Framework runtime."
---

# Delivery Orchestrator Skill

## Layer
Governance

## Responsibility
Route one workspace through allowed orchestrators using persisted state, handoffs, checkpoints, blockers, and attempt limits. It never authors specialist artifacts or executes commands.

## Operating modes
- create: produce the first runtime artifact or route.
- update: refresh state while preserving approved contracts.
- audit: detect stale, unsafe, conflicting, or unauthorized runtime state.
- explain: summarize state, evidence, blockers, and next action.

## Inputs
Objective; workspace; latest checkpoint; handoffs; validator findings; approvals.

## Outputs
Updated workspace state; structured handoff; checkpoint request; routed next orchestrator or blocker.

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
6. For documentary work, update runtime state and create a checkpoint/handoff whenever the canonical gate changes. Use `guide` output to explain rather than execute the next skill.
7. Use the consolidated dashboard model when reporting workflow state so humans and agents see the same stages, blockers, graph/tasks, decisions, leases, and next actions.
8. Stop at human approval, remote mutation, conflict, security, or attempt-limit gates.

## Quality checklist
- [ ] Preserves traceability to affected artifacts.
- [ ] Uses the correct template and naming conventions.
- [ ] Records base commit, hashes, owner, timestamps, attempts, and blockers.
- [ ] Enforces the FDR-017 authority boundary.
- [ ] Detects gaps, conflicts, and dependencies.
- [ ] Records or requests decisions for meaningful changes.
- [ ] Leaves a clear handoff for the next skill or orchestrator.

## Handoff
Next: delivery-orchestrator.

Use the persisted runtime state to route the concrete next specialist or orchestrator. Pass forward workspace, task, hashes, evidence, decisions, dependencies, risks, attempts, blockers, and required reading.
