# FDR-017: Resumable Parallel Runtime

## Snapshot

| Field | Value |
| --- | --- |
| ID | `FDR-017` |
| Status | `approved` |
| Origin EV | `Governance baseline` |
| Date | `2026-07-11` |
| Owner | `Delivery Orchestrator` |

## Context

FDR-016 makes graphs claimable and delivery navigable, but execution state still depends on a flat workspace, permanent claims, a shared checkout, manually selected commands, and conversational handoffs. Safe parallel agents require isolated Git worktrees, expiring ownership, resumable checkpoints, shell-free command plans, deterministic waves, and governed integration.

## Decision

Adopt runtime v2 under `.product/`:

```text
workspaces/WORK-NNN/{workspace.json,state.json,handoffs/,checkpoints/,command-plans/,evidence/}
claims/TK-NNN.json
scheduler/waves/
integrations/INTEGRATION-NNN.json
```

| Concern | Contract |
| --- | --- |
| Orchestration | `delivery-orchestrator` routes only; it does not author or execute. |
| Scheduling | `execution-scheduler` calculates safe waves; it does not spawn agents initially. |
| Commands | `command-planner` produces argv-based plans from approved gates/tasks/runbooks. |
| Execution | `command-executor` runs validated R0/R1 commands only, with timeout and sanitized environment. |
| Integration | `integration-orchestrator` plans ordered local integration and requires Integrated QA. |
| Claims | One JSON lease per task; 30-minute default, heartbeat every 5 minutes, max three attempts. |
| Isolation | One Git branch/worktree per task under ignored `.worktrees/WORK/TASK`. |
| Resume | Checkpoints hash inputs/outputs/base commit. Changed inputs make downstream state stale. |
| Risk | R0 read-only; R1 local temporary; R2 persistent local requires approval; R3 remote and R4 production/destructive remain disabled initially. |
| Integration | Default strategy is cherry-pick plan; conflicts are never resolved automatically. |
| Versioning | Checkpoints, handoffs, plan summaries, claims, and integration records are versionable; worktrees, locks, raw logs, and secrets are not. |
| Compatibility | v2 reads flat v1 workspaces; writes v2 only. Migration is explicit and supports dry-run. |

## Consequences

| Type | Consequence |
| --- | --- |
| Positive | A new agent can reconstruct scope, staleness, ownership, evidence, and next action without chat history. |
| Positive | Independent tasks can run in isolated worktrees with deterministic resource conflict checks. |
| Negative | Git/runtime state and cross-platform process control become framework responsibilities. |
| Follow-up | R2, automatic agent spawning, R3/R4, and remote merge remain future gated evolutions. |

## FRAMEWORK.md Amendments

| Section | Amendment |
| --- | --- |
| `8` | Add leases, scheduler, worktrees, command plans, checkpoints, and integration. |
| `9-10` | Register five runtime skills and their authority boundaries. |
| `15` | Add resume, checkpoint, handoff, heartbeat, worktree, commands, schedule, and integrate commands. |

## Related Artifacts

| Artifact | Link |
| --- | --- |
| Delivery closure | [FDR-016](FDR-016-delivery-closure-and-operational-workspaces.md) |
| Framework method | [FRAMEWORK.md](../../FRAMEWORK.md) |

## Supersedes

- Runtime storage and permanent-lock portions of FDR-016; delivery method clauses remain active.
