# FDR-015: Starting Points and Source Import

## Snapshot

| Field | Value |
| --- | --- |
| ID | `FDR-015` |
| Status | `approved` |
| Origin EV | `Governance baseline` |
| Date | `2026-07-11` |
| Owner | `Documentation Orchestrator` |

## Context

Adopter repositories may start from a blank product, existing product documents, an existing feature, implementation, or an audit need. Removing orchestrators at installation time confuses available automation with required method inputs and can weaken traceability. Existing documents also need a controlled path into the canonical product graph without silently becoming approved product truth.

## Decision

All framework skills and orchestrators remain installed. `init` records a starting point and generates contextual bootstrap guidance. For `existing-documents`, sources are copied into `product/knowledge/imports/sources/` and an immutable import run inventories them before any canonical product artifact is created.

| Boundary | Rule |
| --- | --- |
| Capability installation | Agent selection controls target formats; it does not remove framework capabilities. |
| Source documents | Sources are evidence, not approved product artifacts. |
| Analysis | Inventory, candidates, conflicts, and mappings are produced read-only against canonical product artifacts. |
| Materialization | Requires explicit human approval and creates `draft` artifacts only. |
| Conflicts | Ambiguities, duplicates, and contradictions are reported and never resolved silently. |
| Gates | Starting point does not waive parent artifacts, rigor tiers, approvals, or validation gates. |
| Re-runs | Source hashes make changed inputs and stale mappings mechanically detectable. |

Canonical import runs live at `product/knowledge/imports/runs/IMPORT-NNN/` and contain `inventory.json`, `import-plan.json`, `mapping.json`, `conflicts.md`, and `import-report.md`.

## Consequences

| Type | Consequence |
| --- | --- |
| Positive | Existing product knowledge can enter the framework with source-level traceability. |
| Positive | Adopters retain every workflow while receiving a relevant first-run path. |
| Negative | Import adds schemas, validation rules, and a human reconciliation gate. |
| Follow-up | Add assisted extraction for richer document formats without changing the approval boundary. |

## FRAMEWORK.md Amendments

| Section | Amendment |
| --- | --- |
| `4. Folder Structure` | Add the canonical imports area and import-run artifacts. |
| `10. Orchestrators` | Add Existing Product Import Orchestrator. |
| `15. How To Use With Codex` | Describe starting-point selection and source import. |

## Related Artifacts

| Artifact | Link |
| --- | --- |
| Go CLI installation decision | [FDR-013](FDR-013-go-cli-and-agent-skill-installation.md) |
| Generated bootstrap decision | [FDR-014](FDR-014-generated-adopter-bootstrap.md) |
| Import specialist | [artifact-importer](../skills/artifact-importer/SKILL.md) |
| Import orchestrator | [existing-product-import-orchestrator](../skills/existing-product-import-orchestrator/SKILL.md) |

## Supersedes

- `N/A`
