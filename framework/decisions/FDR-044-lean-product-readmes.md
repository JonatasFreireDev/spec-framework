# FDR-044: Lean Product READMEs

## Snapshot

| Field | Value |
| --- | --- |
| ID | `FDR-044` |
| Status | `proposed` |
| Origin EV | `Adopter starter context review` |
| Date | `2026-07-13` |
| Owner | `Documentation Orchestrator` |

## Context

The starter contains 37 product READMEs with more than 500 lines. Many leaf READMEs exist only to preserve an otherwise empty directory or repeat instructions already owned by `BOOTSTRAP.md`, `context.md`, a specialist skill, or a framework template. Agents beginning product work can encounter generic folder prose that competes with the selected starting point and increases context without adding current product state.

## Decision

| Boundary | Contract |
| --- | --- |
| Authority | `BOOTSTRAP.md` owns starting-point sequence; `context.md` owns local state and handoff; skills own workflow; templates own artifact structure. |
| README role | Product READMEs are concise navigation aids only. They may state purpose, ownership boundary, authoritative files, and the next safe entry point. |
| Leaf directories | Do not ship a README solely to preserve an empty directory. Declarative init contracts or the owning skill create a directory when a real artifact needs it. |
| Retention | Keep READMEs at major product areas, governed import/decision boundaries, and copyable template entry points where local navigation prevents misuse. |
| Duplication | Do not duplicate gates, full workflows, templates, status, or starting-point-specific sequencing in a README. |
| Research templates | Promote reusable interview and research structures to `framework/template/`; product evidence remains under Foundation when created. |
| Compatibility | Upgrade does not delete adopter-owned READMEs or directories. The reduction applies only to new initialization output. |

## Consequences

| Type | Consequence |
| --- | --- |
| Positive | New products expose less generic context and clearer sources of authority. |
| Positive | Empty areas are created proportionally instead of appearing complete because a placeholder README exists. |
| Positive | Reusable research structures become canonical framework templates. |
| Negative | Users browsing an unmaterialized leaf path no longer see a local explanation until that area is created. |
| Negative | Tests and initialization contracts must not assume placeholder READMEs preserve directories. |

## FRAMEWORK.md Amendments

| Section | Amendment |
| --- | --- |
| `3. Context System` | Define README as navigation only and prohibit it from competing with Bootstrap, context, skills, or templates. |

## Related Artifacts

| Artifact | Link |
| --- | --- |
| Canonical method | [../../FRAMEWORK.md](../../FRAMEWORK.md) |
| Product starter | [../../starter/product/README.md](../../starter/product/README.md) |
| Initialization catalog | [../init/catalog.json](../init/catalog.json) |
| Framework templates | [../template/README.md](../template/README.md) |
| Declarative initialization | [FDR-041](FDR-041-declarative-initialization-contracts.md) |

## Supersedes

- N/A
