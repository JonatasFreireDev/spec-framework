# FDR-003: WriteScope and safe parallelism

## Snapshot

| Field | Value |
| --- | --- |
| ID | `FDR-003` |
| Status | `approved` |
| Origin EV | `EV-012` |
| Date | `2026-07-10` |
| Owner | `Task Generator` |

## Context

The Execution Graph promises safe parallelism, but parallel execution is unsafe unless each task declares which files, paths, modules, generated indexes, locales, local databases, or other shared resources it can touch.

The templates already include `writeScope`, but the validator did not enforce how it is used.

## Decision

Every Execution Graph node declares:

- `writeScope`: paths or modules the task may write.
- `sharedResources`: optional list of generated or shared resources that serialize parallel work.

Two nodes are parallel when neither depends on the other through the DAG. Parallel nodes must not have overlapping `writeScope` entries and must not declare the same `sharedResources` entry.

Path overlap is prefix based: `src/` overlaps `src/foo.ts`.

Rollout has two phases:

| Phase | Validator behavior |
| --- | --- |
| Phase A | Missing or overlapping `writeScope` and `sharedResources` produce warnings. |
| Phase B | After existing graphs are clean, the same findings may be promoted to errors by a future approved framework evolution. |

## Consequences

| Type | Consequence |
| --- | --- |
| Positive | Parallel task execution becomes mechanically auditable. |
| Positive | Agents get a concrete write boundary before implementation. |
| Negative | Task generation requires more upfront precision. |
| Follow-up | Future EV may promote warnings to errors after migration. |

## FRAMEWORK.md Amendments

| Section | Amendment |
| --- | --- |
| `8. Execution Graph` | Define `writeScope`, `sharedResources`, prefix overlap, and Phase A warning rollout. |

## Related Artifacts

| Artifact | Link |
| --- | --- |
| Execution graph template | [execution graph template](../template/execution-graph-template.json) |
| Validator | [../validators/framework-validator.mjs](../validators/framework-validator.mjs) |
| FRAMEWORK | [../../FRAMEWORK.md](../../FRAMEWORK.md) |

## Supersedes

- `N/A`
