# FDR-024: Supervised Adapter Management

## Snapshot

| Field | Value |
| --- | --- |
| ID | `FDR-024` |
| Status | `approved` |
| Origin EV | Human-approved adapter usability evolution |
| Date | `2026-07-12` |
| Owner | `Framework Guide` |

## Context

Visual adapters are optional, but users must currently discover runtimes, installation commands, versions, and harness paths outside the framework. Automatic installation would create an unsafe supply-chain and authority expansion.

## Decision

Add supervised adapter management with read-only `list`, `status`, and `doctor`, plus explicit version-pinned `install` and `update` operations.

| Concern | Contract |
| --- | --- |
| Optionality | No adapter is required for core framework validation. |
| Preview | Install/update print the exact provider, package, version, argv, cwd, and detected runtime before mutation. |
| Confirmation | Mutation requires `--yes` and an explicit version; no conversational inference authorizes it. |
| Execution | Commands use direct argv in the repository root, not a shell string. |
| State | Status is detected from supported harness skill paths; credentials are never persisted. |
| Removal | Not provided until the upstream provider has a documented, reversible removal contract. |
| Failure | Provider exit code and output are preserved; the framework never reports a failed install as ready. |

## Consequences

| Type | Consequence |
| --- | --- |
| Positive | Users can diagnose and install optional adapters without memorizing provider commands. |
| Positive | External execution remains visible, pinned, and explicitly authorized. |
| Negative | Installation still depends on Node, npx, network access, and upstream interactive behavior. |

## FRAMEWORK.md Amendments

| Section | Amendment |
| --- | --- |
| `15` | Document supervised adapter discovery and installation. |

## Related Artifacts

| Artifact | Link |
| --- | --- |
| Visual adapter protocol | [FDR-021](FDR-021-design-sources-and-visual-artifacts.md) |
| Framework Guide | [Framework Guide](../skills/framework-guide/SKILL.md) |

## Supersedes

- N/A
