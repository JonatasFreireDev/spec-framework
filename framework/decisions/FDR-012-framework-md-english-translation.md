# FDR-012: FRAMEWORK.md English translation

## Snapshot

| Field | Value |
| --- | --- |
| ID | `FDR-012` |
| Status | `approved` |
| Origin EV | `governance baseline` |
| Date | `2026-07-10` |
| Owner | `Documentation Orchestrator` |

## Context

`FRAMEWORK.md` was written in Portuguese while every other framework-core document (`README.md`, `AGENTS.md`, `framework/adoption.md`, `framework/decisions/FDR-*`, `framework/template/`, `.codex/skills/`) is written in English. The framework is a reusable, adoptable artifact intended for repositories and agents beyond this lab; a mixed-language core document adds friction for adopters and for Codex when it reads `FRAMEWORK.md` alongside the rest of the English-language corpus.

## Decision

`FRAMEWORK.md` is translated to English in full. Section numbering and order are preserved exactly (`1.` through `17.`); only prose, headings, and in-example strings change language. Code blocks that are structural (YAML field names, JSON keys, file paths, shell commands) are left unchanged; Portuguese prose embedded inside example values (for example the `rationale` and `open_questions` sample strings in section 5) is translated along with the surrounding text.

Existing FDR-001 through FDR-011 amendment tables that cite a Portuguese section title (for example `4. Estrutura De Pastas`, `15. Como Usar Com Codex`) are historical records of the amendment made at that time and are not rewritten. Section numbers are the stable identifier across the rename; readers resolve an old FDR's amendment entry by section number, not by matching the Portuguese title text verbatim.

## Consequences

| Type | Consequence |
| --- | --- |
| Positive | The framework core is language-consistent, lowering the barrier for non-Portuguese-speaking adopters and agents. |
| Positive | Section numbers stayed stable, so cross-references by number in tooling, tasks, and prior FDRs still resolve. |
| Negative | FDR-001 through FDR-011 amendment tables now reference section titles in a language that no longer matches `FRAMEWORK.md` verbatim; readers must resolve by section number. |
| Follow-up | Future FDRs that amend `FRAMEWORK.md` must cite the English section title. |

## FRAMEWORK.md Amendments

| Section | Amendment |
| --- | --- |
| `Whole document` | Translated from Portuguese to English; section numbers and structure unchanged. |

## Related Artifacts

| Artifact | Link |
| --- | --- |
| FRAMEWORK | [../../FRAMEWORK.md](../../FRAMEWORK.md) |
| Prior amendments | [FDR-001](FDR-001-decision-governance.md), [FDR-002](FDR-002-gate-commands.md), [FDR-004](FDR-004-qa-independence.md), [FDR-005](FDR-005-code-runner-contract.md), [FDR-006](FDR-006-failure-routing-and-regression.md), [FDR-007](FDR-007-code-review-contract.md), [FDR-008](FDR-008-delivery-commits-and-prs.md), [FDR-011](FDR-011-core-starter-example-boundary.md) |

## Supersedes

- `N/A`
