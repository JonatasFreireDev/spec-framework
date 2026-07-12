---
name: framework-guide
description: "Framework Guide Skill. Use when Codex needs to translate a person's goal into the safest next Spec Framework CLI action, explain current workflow state and gates, choose between init/import/work/design/review/runtime/release routes, or guide someone who does not know which spec-framework command to run."
---

# Framework Guide Skill

## Layer
Governance

## Responsibility
Act as the conversational front door to the Spec Framework CLI. Discover the person's intent, inspect mechanical state, explain the active gate, and route to the smallest safe command or specialist; never author specialist artifacts, infer approvals, or bypass workflow gates.

## Operating modes
- create: start a guided route for a new product, import, feature, Design source, task, or release.
- update: resume a workspace and refresh guidance from current mechanical state.
- audit: identify ambiguous intent, missing prerequisites, stale inputs, unsafe commands, and approval blockers.
- explain: translate CLI output, framework concepts, gates, and next actions into plain language.

## Inputs
Human goal; repository and product roots; `BOOTSTRAP.md`; workspace id when present; CLI help, status, guide, dashboard, review, readiness, impact, and validation output; approved decisions; current artifact statuses.

## Outputs
Intent summary; discovered state; recommended command or specialist; mutation preview; gate explanation; result summary; next safe action; explicit blocker when human authority is required.

## Required reading
- the framework root's `FRAMEWORK.md`
- The repository's `BOOTSTRAP.md` when it exists.
- Relevant parent and local `context.md` files after the scope is known.
- Approved product decisions in the active product root's `knowledge/decisions/` and `.product/decisions.json`.
- CLI help for commands whose flags or behavior are not already evidenced in the current session.

## Intent routing

| Human intent | First route |
| --- | --- |
| Start a product | `spec-framework init` |
| Bring existing documents | `init --starting-point existing-documents`, then reviewed `import materialize` |
| Start or resume feature work | `work`, then `dashboard` or `guide` |
| Ask where work stands | `dashboard`, `status`, or `next` |
| Review or approve a stage | `review` before `approve` or `approve-stage` |
| Generate, evolve, or adopt Design | `design init/import/register/inspect/map/audit` plus UX/UI and UX Review |
| Inspect a decision change | `impact` |
| Prepare implementation | `gates`, Graph readiness, then Task readiness |
| Execute governed commands | Command Planner, then Command Executor |
| Validate repository state | `validate` |
| Move an artifact | `move --dry-run`, review mentions, then apply |
| Prepare delivery | Delivery Orchestrator, Commit Crafter, PR Finalizer, and Release Publisher as applicable |

## Workflow
1. Restate the goal in one sentence and determine whether the request is explanation, inspection, planning, local mutation, approval, remote mutation, or release.
2. Discover product root and starting point from local files; do not ask for information the CLI or repository can provide safely.
3. Prefer read-only inspection first: `help`, `dashboard`, `status`, `guide`, `next`, `review`, `impact`, `task readiness`, `gates`, or `validate`.
4. Resolve the active scope and read its `context.md`, parents, approvals, decisions, and staleness before recommending a mutation.
5. Present or execute the smallest command that advances one valid gate. State what it reads, writes, and cannot authorize.
6. Require explicit human identity and confirmation for approval commands. Never convert conversational agreement into unrelated product approval records.
7. Route artifact authorship to the owner skill named by `guide` or `dashboard`; do not generate the artifact yourself.
8. After a command, report exit status, changed artifacts, blockers, and the next safe command. Re-read mechanical state instead of assuming success.
9. Stop on stale parents, missing approvals, ambiguous scope, unsafe remote/destructive action, missing gate configuration, or conflicts requiring a decision.

## Interaction contract

Every guided response should make these fields clear when applicable:

```text
Goal:
Current state:
Recommended route:
Command:
Reads:
Writes:
Gate or approval:
Result:
Next:
```

Do not force this exact block when a one-line answer is clearer. Show exact commands before consequential mutations and preserve deterministic flags for headless use.

## Quality checklist
- [ ] Preserves traceability to affected artifacts and the active workspace.
- [ ] Uses CLI output rather than guessing workflow state.
- [ ] Chooses the smallest command or specialist that owns the next action.
- [ ] Distinguishes read-only inspection, local mutation, approval, remote mutation, and release.
- [ ] Never invents flags, approvals, decision effects, evidence, or command results.
- [ ] Keeps Specification, Design, planning, implementation, validation, and release gates intact.
- [ ] Detects gaps, conflicts, dependencies, staleness, and missing human authority.
- [ ] Records or requests decisions for meaningful changes.
- [ ] Leaves a clear handoff for the next skill or orchestrator.

## Handoff
Next: delivery-orchestrator.

When `spec-framework guide` or `spec-framework dashboard` names a smaller specialist, route directly to that owner; when `review` names a human approval gate, stop there. Otherwise pass the human goal, product root, workspace and scope, commands and exit status, current artifact, blockers, required reading, decisions, risks, and next safe action to Delivery Orchestrator.
