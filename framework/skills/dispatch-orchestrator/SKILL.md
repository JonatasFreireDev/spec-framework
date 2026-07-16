---
name: dispatch-orchestrator
description: "Dispatch Orchestrator Skill. Use when Codex needs to plan, assign, observe, reconcile, or supervise bounded subagent work in the Spec Framework workflow."
---

# Dispatch Orchestrator Skill

## Layer
Execution

## Responsibility
Owns dispatch envelopes, waves, capacity observation, and handoff sequencing. It does not author product artifacts, approve work, or deliver remotely.

## Operating modes
- create: create an explicit assignment envelope.
- update: return or reconcile persisted dispatch state.
- audit: inspect leases, transcripts, scope conflicts, and stale returns.
- explain: summarize eligible work and blockers.

## Inputs
Approved Execution Graph; ready task or import chunk; workspace state; dispatch configuration; current product decisions.

## Outputs
Persisted envelopes; dispatch plan; wave observation; reconciliation findings; handoffs.

## Required reading
- the framework root's `FRAMEWORK.md`
- Relevant parent context.md files.
- Relevant templates in framework/template/.
- Approved product decisions in the active product root's `knowledge/decisions/` and `.product/decisions.json`.

## Workflow
1. Read current graph/chunk readiness and dispatch configuration.
2. Plan only canonical units with no dependency or scope conflict.
3. Require explicit human confirmation before assignment or execution.
4. Persist envelope, lease, required reading, hashes, scope, forbidden operations, and expected evidence.
5. Dispatch QA, Code Review, and Security Review only from the returned Code Runner diff hash.
6. Reconcile and route blockers; never repair approval or product state.

## Quality checklist
- [ ] Preserves traceability to affected artifacts.
- [ ] Uses the correct template and naming conventions.
- [ ] Never assigns overlapping write scopes or resources.
- [ ] Execution uses an enabled harness and explicit confirmation.
- [ ] Detects gaps, conflicts, and dependencies.
- [ ] Records or requests decisions for meaningful changes.
- [ ] Leaves a clear handoff for the next skill or orchestrator.

## Handoff
Next: code-runner, artifact-importer, qa, code-review, security-review, or product-historian.

Pass forward envelope, hashes, evidence, blockers, risks, and required follow-up work.
