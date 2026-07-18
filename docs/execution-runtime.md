# Execution Runtime Contract

This document contains the operational contract for resumable, parallel, and evidence-backed task execution. The method and lifecycle remain canonical in `FRAMEWORK.md`; this file owns runtime mechanics.

## Scope

The runtime coordinates workspaces, task ownership, command plans, execution waves, isolated worktrees, local integration, and recovery. It does not grant product approval, replace task write scopes, or spawn agents.

Harness-native agent delegation is distinct from process execution. The runtime may persist a supervised envelope that a parent agent passes to a native subagent, but only the active harness creates, waits for, interrupts, or collects that subagent. This keeps the CLI portable across Codex, Cursor, Claude, and environments without native delegation.

## Configuration authority

Runtime behavior is determined by the pinned product manifest, approved product
artifacts, and explicit command flags. It must not read or inherit an ambient
user configuration file for agent selection, task behavior, gates, approvals,
or external access. Operator-local caches may retain only non-semantic details
such as paths or event retention.

## Workspace and ownership

- Use `.product/workspaces/WORK-NNN/` for concurrent focus; never invent a global active feature.
- A workspace records identity, state, handoffs, checkpoints, command plans, and evidence.
- Its optional event ledger is append-only operational evidence at `events/`. It redacts secret-like detail keys and retains the newest 500 events; replay and observation are read-only and never infer an approval or product state.
- Task ownership is a renewable lease with heartbeat and expiry. A lease does not grant approval or permission beyond the task's `writeScope`.
- Isolated tasks use `.worktrees/WORK-NNN/TK-NNN` when the graph/runtime contract requires isolation.
- Resume from `state.json`, the latest checkpoint, and the latest handoff. Legacy `WORK-NNN.json` is read-only until explicit migration.

## Graph and scheduling

The Execution Graph is a DAG. The scheduler calculates deterministic conflict-free waves from dependencies, `writeScope`, `sharedResources`, capabilities, leases, priority, and capacity. It does not execute tasks or spawn agents.

Parallel tasks require no dependency path between them and must not overlap in `writeScope` or contend for an undeclared shared resource. Path overlap is prefix-based. A conflict becomes a dependency or requires task merging.

Graph lifecycle is `draft -> proposed -> materialized -> approved`. `materialized` is a Graph-specific state: it means canonical task files and the generated index exist; it is not a general artifact lifecycle state.

Scheduling is planning only. `schedule activate <wave>` is a separate, explicit operation and requires `--isolate`, a named agent, and `--yes`; it rechecks readiness against the same approved graph, acquires leases, and then creates only that wave's worktrees. A partial activation releases its newly acquired leases but preserves created worktrees for diagnosis and explicit cleanup.

## Command execution

Command plans store direct argv rather than shell strings. The command executor initially permits only R0 read-only and R1 local-temporary operations. It refuses stale inputs, scope escapes, conflicting resources, unsupported risk levels, and attempts beyond the configured limit. Human approval is required for remote, destructive, security-sensitive, or otherwise consequential operations.

## Evidence and integration

Implemented tasks record the required immutable working-tree evidence and current diff hash. Code Review and task QA review the same diff hash before Commit Crafter creates local commits. Validated task commits are integrated locally in DAG order; conflicts stop for human resolution, followed by Integrated QA where applicable.

Runtime commands include `runtime`, `resume`, `handoff`, `checkpoint`, `lease`, `commands`, `schedule`, `integrate`, and `reviews import`. The latter imports a local JSON array into immutable provider-neutral evidence under `.product/reviews/findings/`; it only proposes a route and cannot resolve external threads, change code, or create approvals. Use the installed CLI help as the authority for exact flags and syntax.

`runtime status` makes one read-only observation. `runtime watch` repeats the same local observation at an explicit interval (and supports a bounded `--count` for automation); neither command writes checkpoints, lease heartbeats, or events.

`runtime reconcile` is also read-only. It reports expired leases, missing graphs or command evidence, orphaned task worktrees, and incomplete integration event hashes. It never deletes a worktree, frees a claim, replays a command, or edits a product artifact.

Operational memory is optional and has shared (`memory/shared.md`) and task-local (`memory/tasks/<task>.md`) tiers. Its compact form uses `- source: [label](path)` and `- risk: ...` lines. Inspection is read-only; explicit compaction only removes duplicate lines and refuses approval history or flagged contradictions.

ACP dispatch is experimental and disabled by default. A run requires explicit enablement and per-run acknowledgement, then claims exactly one ready task for its named agent. It records a local transcript hash and releases that lease at the end. It cannot invoke Git delivery commands and has no approval, review-resolution, push, merge, or release authority.

Engineering baseline delegation does not use ACP process execution. Its `engineering-specialist` envelopes are activated by an Engineering Orchestrator handoff whose execution mode is `delegated`. Dispatch Orchestrator owns assignment and observation; Subagent Return Reviewer validates the compact return before CLI persistence. The envelope pins the specialist role, handoff input hash, dependency returns, minimal-context policy, phase, and product-relative write scope. Assignment capacity is claimed atomically per workspace, so concurrent CLI invocations cannot exceed `max_parallel` or duplicate an active specialist role. Returns require evidence and SHA-256 hashes for outputs inside that scope. The default handoff mode is `sequential`; delegated mode has a declared sequential fallback and a bounded `max_parallel` value.

Extensions are versioned manifests discovered without execution. A capability is usable only when that manifest declares it and the product has a matching versioned record under `.product/extensions/`; discovery itself grants no trust or authority.

## Local project-status server

`spec-framework server` is a local-only, user-facing status surface. It binds
only to `127.0.0.1`, opens the browser by default, and does not expose a remote
API. It is deliberately separate from `dashboard`, which is the technical
workspace view.

Run `spec-framework server start` to start it, `spec-framework server status`
to print its local URL, and `spec-framework server stop` from another terminal
to request graceful shutdown. The foreground process also stops on `Ctrl+C`.
The server stores a short-lived local descriptor under `.product/server.json`;
it is removed after normal shutdown and never contains approval data.

When the status surface enables review actions, it delegates approval and
rejection to the same lifecycle engine as the CLI. A rejection requires a
human identity and revision rationale, writes immutable history, and moves an
eligible artifact to `rejected`, including an approved artifact that must be
reopened. A revised rejected artifact may then be directly approved through a
new recorded human review; the runtime never silently edits product scope.

## Owning skills

- `execution-graph`: defines and validates the DAG and graph lifecycle.
- `execution-scheduler`: calculates waves and conflicts without executing tasks.
- `command-planner`: creates an approved command plan.
- `command-executor`: executes only an approved, current, scoped command plan.
- `delivery-orchestrator`: routes work through persisted state and handoffs.
- `integration-orchestrator`: integrates verified local commits and requires integrated validation.
