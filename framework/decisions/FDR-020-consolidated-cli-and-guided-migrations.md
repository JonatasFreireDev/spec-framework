# FDR-020: Consolidated CLI and Guided Migrations

## Snapshot

| Field | Value |
| --- | --- |
| ID | `FDR-020` |
| Status | `approved` |
| Origin EV | `Governance baseline` |
| Date | `2026-07-12` |
| Owner | `Delivery Orchestrator` |

## Context

The CLI exposes safe commands, but users must combine status, guide, graph, readiness, decisions, and runtime records mentally. Legacy decisions also need typed scoped metadata without unsafe bulk inference or manual JSON editing.

## Decision

| Concern | Contract |
| --- | --- |
| Consolidated view | One dashboard presents the canonical Feature-to-Task stages, current gate, blockers, graph/task summary, decisions, leases, checkpoints, and next actions. |
| Output modes | Human output is scannable and color-aware; `--json` is stable for agents and CI. `status --graph` reuses the same model. |
| Interaction | Bubble Tea/Huh may collect navigation and migration choices only when explicitly requested or when an existing interactive command has no headless confirmation flags. |
| Headless safety | Every interactive operation has deterministic flags, preview/dry-run, exit codes, and no TTY requirement. |
| Decision migration | Migration detects legacy entries, proposes type/scope/effects, previews every change, creates a backup, writes atomically, and never changes DEC content/status or creates approvals. |
| Inference | Type may be inferred from legacy scope; scope paths come from affected artifacts. Ambiguous values remain explicit review items. Effects default to empty arrays. |
| Compatibility | Existing commands and decision indexes remain readable. Migration is opt-in and idempotent. |

## Consequences

| Type | Consequence |
| --- | --- |
| Positive | Humans and agents share one workflow snapshot and next-action model. |
| Positive | Legacy decisions can adopt FDR-019 safely without hand-editing JSON. |
| Negative | CLI presentation and migration schemas require golden/JSON compatibility tests. |
| Follow-up | A full-screen long-running agent monitor remains future work. |

## FRAMEWORK.md Amendments

| Section | Amendment |
| --- | --- |
| `15` | Add dashboard, status graph, and decision migration workflows. |

## Related Artifacts

| Artifact | Link |
| --- | --- |
| Decision effects | [FDR-019](FDR-019-executable-product-decision-effects.md) |
| Runtime | [FDR-017](FDR-017-resumable-parallel-runtime.md) |

## Supersedes

- N/A
