# FDR-035: Audit-Only Mutation Guard

## Snapshot

| Field | Value |
| --- | --- |
| ID | `FDR-035` |
| Status | `proposed` |
| Origin EV | `Adopter starting-point review` |
| Date | `2026-07-12` |
| Owner | `Framework Guide` |

## Context

The `audit-only` starting point currently changes guidance but does not mechanically prevent commands from creating workspaces, approvals, imports, registry rewrites, migrations, or delivery state. Its generated bootstrap also recommends `validate --write-registry --write-report`, contradicting the stated read-only intent.

## Decision

| Boundary | Contract |
| --- | --- |
| Activation | The guard activates only when the canonical manifest declares `starting_point: audit-only`. |
| Allowed | Help, version, read-only validation, inspection, status, readiness, review, impact, dashboard, skill resolution, and other commands that do not write product or delivery state. |
| Blocked | Artifact moves, imports, approvals, workspace creation, registry/report writes, design or system initialization, non-dry-run migrations, graph/task execution mutations, decision migration application, and runtime delivery mutations. |
| Validation | `validate` is allowed without `--write-registry` or `--write-report`; audit results remain terminal output unless a human explicitly changes starting point. |
| Framework maintenance | `init` establishes the audit manifest and `upgrade` may maintain the pinned framework runtime. Neither grants permission to mutate adopter product truth. |
| Transition | Continuing from findings into product delivery requires an explicit manifest starting-point change through a future supported transition, not a hidden bypass. |

## Consequences

| Type | Consequence |
| --- | --- |
| Positive | Audit intent is mechanically enforced instead of relying on agent restraint. |
| Positive | Approval history, registry, workspaces, and product artifacts cannot drift during an audit-only session. |
| Negative | Persisting audit reports under `product/` is blocked while audit-only remains active. |
| Follow-up | Add an explicit `starting-point transition` command with preview and confirmation. |

## FRAMEWORK.md Amendments

| Section | Amendment |
| --- | --- |
| `15. How To Use With Codex` | Define the audit-only read/write boundary and transition requirement. |

## Related Artifacts

| Artifact | Link |
| --- | --- |
| Starting points | [FDR-015](FDR-015-starting-points-and-source-import.md) |
| Framework Guide | [framework-guide](../skills/framework-guide/SKILL.md) |
| Canonical method | [../../FRAMEWORK.md](../../FRAMEWORK.md) |

## Supersedes

- N/A; this decision makes the audit-only intent in FDR-015 mechanical.
