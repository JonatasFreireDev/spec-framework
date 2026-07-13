# FDR-034: Existing Documents Import Gate

## Snapshot

| Field | Value |
| --- | --- |
| ID | `FDR-034` |
| Status | `proposed` |
| Origin EV | `Adopter starting-point review` |
| Date | `2026-07-12` |
| Owner | `Existing Product Import Orchestrator` |

## Context

The `existing-documents` starting point already creates an immutable import run and requires explicit materialization approval, but workspace creation does not verify that the latest declared run was reviewed and materialized. Generated bootstrap guidance also jumps toward Foundation approval before clearly naming mapping review and draft materialization as the current gate.

## Decision

| Boundary | Contract |
| --- | --- |
| Entry contract | The latest import run declared in `.product/framework.json` is the starting-point contract; no additional summary artifact is introduced. |
| Review | Inventory, conflicts, mapping selections, targets, source references, and draft content are reviewed before mutation. |
| Materialization | `spec-framework import materialize` requires an explicit human identity and creates selected canonical targets as `draft` only. |
| Approval boundary | Import materialization approval authorizes the selected draft writes. It does not approve the resulting product artifacts. |
| Workspace gate | `spec-framework work` is blocked until the latest declared run has status `materialized`, explicit materialization approval metadata, and at least one materialized path. |
| Integrity | Reported materialized paths must exactly match selected mappings and exist under the product root. New runs record hashes of approved draft content; legacy runs without hashes require the materialized file to remain identical until the gate is satisfied. |
| Foundation | Materialized drafts continue through their normal owners, parents, validation, and individual approval gates. The import itself does not choose a proportional Foundation route silently. |
| Re-import | A newer run declared in the manifest becomes the active gate; an older materialized run cannot satisfy it. |

## Consequences

| Type | Consequence |
| --- | --- |
| Positive | Existing documents cannot be bypassed after init and cannot become approved truth through import alone. |
| Positive | The first actionable command matches the actual import state. |
| Negative | Workspace creation waits for at least one reviewed mapping to be materialized. |
| Follow-up | Add a read-only `import status` command if adopters need richer navigation across multiple runs. |

## FRAMEWORK.md Amendments

| Section | Amendment |
| --- | --- |
| `15. How To Use With Codex` | Make latest-run materialization the explicit existing-documents entry gate. |

## Related Artifacts

| Artifact | Link |
| --- | --- |
| Starting points and import | [FDR-015](FDR-015-starting-points-and-source-import.md) |
| Artifact Importer | [artifact-importer](../skills/artifact-importer/SKILL.md) |
| Import orchestrator | [existing-product-import-orchestrator](../skills/existing-product-import-orchestrator/SKILL.md) |
| Canonical method | [../../FRAMEWORK.md](../../FRAMEWORK.md) |

## Supersedes

- N/A; this decision adds a mechanical workspace gate to FDR-015.
