---
name: framework-guide
description: "Framework Guide Skill. Default entry point for Spec Framework product operations when no verified specialist route exists; use it to translate a person's goal into the safest next CLI action, explain current workflow state and gates, and choose between bootstrap/init/import/work/design/review/migration/upgrade/release routes."
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
Human goal; current working directory; canonical `product/.product/framework.json` when present; repository and product roots; the active product root's `BOOTSTRAP.md`; workspace id when present; CLI help, status, guide, dashboard, review, readiness, impact, migration, runtime, skill-resolution, and validation output; approved decisions; current artifact statuses.

## Outputs
Intent summary; discovered state; recommended command or specialist; mutation preview; gate explanation; result summary; next safe action; explicit blocker when human authority is required.

## Required reading
- The pinned runtime contracts: [`execution-runtime.md`](../../docs/execution-runtime.md), [`engineering-systems.md`](../../docs/engineering-systems.md), and [`lifecycle-and-approvals.md`](../../docs/lifecycle-and-approvals.md).
- The pinned runtime's `AGENTS.framework.md` common agent rules.
- the framework root's `FRAMEWORK.md`
- The active product root's `BOOTSTRAP.md` when it exists.
- Relevant parent and local `context.md` files after the scope is known.
- Approved decisions are discovered through the active product root's `.product/decisions.json`; resolve each registered `path` from its declared domain root (`knowledge/decisions/`, `design/decisions/`, or `engineering/decisions/`).
- CLI help for commands whose flags or behavior are not already evidenced in the current session.

## Activation boundary

1. Before routing any product operation, look for `product/.product/framework.json` in the current repository.
2. Activate only when that file is valid, identifies `framework: spec-framework`, pins a concrete version, and declares `activation.mode: manifest-only`.
3. Never activate because the user mentions Spec Framework, a prompt contains matching keywords, or another similarly named file exists.
4. When the canonical manifest is absent, the only permitted Spec Framework routes are explaining or running the explicit bootstrap or `init` requested by the user. Do not load specialized framework contracts first.
5. When the manifest is invalid, stop product operations and report the exact validation problem. Do not infer or repair adoption metadata silently.
6. After activation, resolve specialized contracts with `spec-framework skill path <skill-name>` and read the returned versioned `SKILL.md` completely.

## Dispatch boundary

1. Framework Guide is the default first route for framework-governed product operations.
2. A specialist may be resolved directly only from one of these verified routes:
   - current-session `guide`, `dashboard`, `status`, or `next` output that names the workspace, concrete feature or use-case scope, current gate, and owner skill;
   - an explicit human request that names both the specialist and the concrete artifact or workspace scope.
3. A persisted handoff or checkpoint identifies the workspace to resume but is not direct-route evidence by itself. Revalidate it with `dashboard`, `status`, `next`, or `guide` first.
4. A skill name, keyword, or remembered chat instruction without concrete scope is only a hint. Resolve Framework Guide first.
5. Before using a direct route, validate the manifest, scope, ownership, gate, and staleness against current mechanical state.
6. If direct-route evidence is missing, stale, ambiguous, or conflicting, return to Framework Guide.
7. This boundary does not intercept direct diagnostic CLI commands and does not grant approval or mutation authority.

## Intent routing

| Human intent | First route |
| --- | --- |
| Install without an existing CLI | Download and inspect `scripts/install.ps1` or `scripts/install.sh`, then run the chosen bootstrap; installation does not initialize a product, and piping a remote script remains an explicit user choice |
| Update or remove the installed CLI | Use `spec-framework update [--check | --version <version>]` or preview `spec-framework uninstall [--purge]`; mutations require `--yes` and never remove product repositories |
| Start a product interactively | `spec-framework init <repository-path>` |
| Start a product headlessly | `spec-framework init <repository-path> --agents <agents> --yes`; declare known sibling implementation roots with `--code-roots web:web,api:api` |
| Bring existing documents | `init --starting-point existing-documents`, review the latest run, then `import materialize`; resulting artifacts remain draft |
| Adopt a code-first operating product | `init --starting-point existing-product`, approve `foundation/product-baseline.md`, then approve Strategy |
| Deliver one bounded existing feature | `init --starting-point existing-feature`, complete and individually approve `foundation/feature-brief.md`, then `work` |
| Adopt an existing implementation | `init --starting-point existing-implementation`, approve `knowledge/assessments/implementation-assessment.md`, then derive and approve the full Foundation |
| Resolve a specialized skill | `spec-framework skill path <skill-name>` after manifest activation |
| Upgrade the pinned external runtime | Inspect current manifest and version, then `spec-framework upgrade --yes` |
| Migrate a legacy local runtime | `spec-framework migrate external-runtime --dry-run`, review preserved legacy paths, then rerun with `--yes` |
| Start or resume feature work | `work`, then `dashboard` or `guide` |
| Ask where work stands | `dashboard`, `status`, or `next` |
| Review or approve artifacts | `approve` for one artifact; `approve-batch` for files, IDs, Foundation, or eligible artifacts through a stage; `review`/`approve-stage` remains the workspace-stage compatibility route |
| Generate, evolve, or adopt Design | `design init/import/register/inspect/map/audit` plus UX/UI and UX Review |
| Create, adopt, evolve, or inspect a shared Design System | `design-system init/inspect/validate/migrate` plus Design System Skill |
| Diagnose or install an optional visual adapter | `adapters list/status/doctor`, then version-pinned `install/update --yes` after preview |
| Inspect a decision change | `impact` |
| Prepare implementation | `gates`, Graph readiness, then Task readiness |
| Establish, inspect, or migrate shared engineering contracts | `engineering-system inspect/validate/migrate` plus Engineering System, then Technical Discovery for a delivery |
| Review a proposed technical solution | Engineering Proposal, then independent Engineering Review |
| Execute governed commands | Command Planner, then Command Executor |
| Validate repository state | `validate` |
| Move an artifact | `move --dry-run`, review mentions, then apply |
| Prepare delivery | Delivery Orchestrator, Commit Crafter, PR Finalizer, and Release Publisher as applicable |

## Discovery and challenge

Follow the shared [Discovery And Challenge contract](../discovery-and-challenge.md) when intent, scope, starting point, route, or a meaningful human choice remains ambiguous after read-only inspection. Do not ask for state the CLI can discover.

## Workflow

For a deterministic workspace snapshot, run `scripts/inspect-workspace.ps1` on Windows or `scripts/inspect-workspace.sh` on macOS/Linux instead of reconstructing the `guide`, `status`, and `dashboard` sequence.

For a declared implementation layout, run `scripts/inspect-code-roots.ps1` on Windows or `scripts/inspect-code-roots.sh` on macOS/Linux. The script inventories deterministic code-root evidence and optionally calls `spec-framework validate`; interpret its output before routing domain work.

When the intent concerns a decision or ADR, first run `spec-framework decisions check --product-root product --json`. Use the indexed `path` and reported domain as the routing evidence. `--strict` is the CI gate; `--fix-links` requires a reviewed preview followed by `--yes` and never creates approval records.
1. Restate the goal in one sentence and determine whether the request is explanation, inspection, planning, local mutation, approval, remote mutation, or release.
2. Apply the activation and dispatch boundaries. If active, discover the repository root, product root, pinned version, agents, and starting point from the canonical manifest; do not ask for information the CLI or repository can provide safely.
3. Prefer read-only inspection first: `help`, `dashboard`, `status`, `guide`, `next`, `review`, `impact`, `task readiness`, `gates`, `validate`, `engineering-system inspect`, `engineering-system triggers`, `skill path`, or migration `--dry-run`.
4. Resolve the active scope and read its `context.md`, parents, approvals, decisions, and staleness before recommending a mutation.
5. Present or execute the smallest command that advances one valid gate. State what it reads, writes, and cannot authorize.
6. For `existing-feature`, route through the registered Feature Brief instead of the full product Foundation. Escalate to Problem, Vision, Product Principles, North Star, and Strategy when the requested work reveals broad or uncertain product direction.
7. For `existing-implementation`, route through the registered Implementation Assessment before full Foundation. Treat code, tests, configuration, and history as evidence rather than approved product intent.
8. For `existing-product`, consolidate evidenced current state in Product Baseline and keep future Strategy separate. Escalate to full Foundation when current audience, value, or direction is uncertain.
9. Before any delivery Domain, Goal, Feature, Use Case, or Specification, read `knowledge/assessments/product-landscape.md`, `engineering/engineering-system.md`, and `design/system/design-system.md`. For code roots declared in the manifest, inventory every root comprehensively and separate observed evidence from inferred intent. For no-code products, establish these as explicit draft hypotheses and identify the intended stack and official scaffold command before proposing implementation creation.
10. Keep implementation projects as semantic sibling roots beside `product/` (`web/`, `api/`, `worker/`, `mobile/`, `infrastructure/`, or `library/`); never place application code inside `product/`.
11. For `audit-only`, use terminal-output inspection commands without write flags. Do not create reports, registry changes, approvals, workspaces, migrations, or delivery state; request an explicit starting-point transition before product work.
12. Require explicit human identity and confirmation for approval commands. For batch approval, preview the exact artifact list and hashes with `approve-batch` before applying `--yes`; ask the human which scope to approve and never convert conversational agreement into unrelated product approval records.
13. Route artifact authorship to the owner skill named by `guide` or `dashboard`; do not generate the artifact yourself.
14. After a command, report exit status, changed artifacts, blockers, and the next safe command. Re-read mechanical state instead of assuming success.
15. Stop on stale parents, missing approvals, ambiguous scope, unsafe remote/destructive action, missing gate configuration, or conflicts requiring a decision.

## Runtime and repository cleanliness

- `init` may add `product/` to a new or existing repository but must refuse an existing `product/` unless the user explicitly selects a supported migration or force route.
- Framework method assets belong in the versioned user cache. Selected harnesses receive only the user-scoped, namespaced dispatcher.
- Treat creation of `.spec-framework/`, repository-local framework skill trees, root bootstrap guides, or generated CI workflows as a regression in the manifest-only model.
- `upgrade` updates the external runtime and pinned product manifest while preserving adopter-owned product and code content.
- Legacy migration previews paths that remain for manual review; it does not delete `.spec-framework/`, local agent trees, product content, or approval history.
- Documentation and bootstrap work writes only within `product/`. Implementation work outside `product/` requires the approved task `writeScope`.

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
- [ ] Enforces manifest-only activation and resolves the pinned skill version.
- [ ] Requires Guide-first routing unless a verified direct route names the specialist and concrete scope.
- [ ] Keeps framework runtime assets and dispatchers outside the adopter repository.
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
