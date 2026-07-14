---
name: integration-orchestrator
description: "Integration Orchestrator Skill. Use when Codex needs to plan ordered local integration of task commits that passed code review and task qa, then require integrated qa in the Spec Framework runtime."
---

# Integration Orchestrator Skill

## Layer
Governance

## Responsibility
Plan ordered local integration of task commits that passed Code Review and Task QA, then require Integrated QA. It never resolves conflicts automatically, pushes, or merges remotely.

## Operating modes
- create: produce the first runtime artifact or route.
- update: refresh state while preserving approved contracts.
- audit: detect stale, unsafe, conflicting, or unauthorized runtime state.
- explain: summarize state, evidence, blockers, and next action.

## Inputs
Workspace; DAG; validated task commits; task QA; integration base; rollback point.

## Outputs
INTEGRATION-NNN.json; ordered cherry-pick plan; integrated diff hash; conflict report; Integrated QA handoff.

## Required reading
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
Next: qa in integrated-qa mode, then pr-finalizer or release-orchestrator.

Pass forward workspace, task, hashes, evidence, decisions, dependencies, risks, attempts, blockers, and required reading.
