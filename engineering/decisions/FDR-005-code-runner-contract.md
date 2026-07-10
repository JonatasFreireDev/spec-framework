# FDR-005: Code Runner operational contract

## Snapshot

| Field | Value |
| --- | --- |
| ID | `FDR-005` |
| Status | `approved` |
| Origin EV | `EV-007` |
| Date | `2026-07-10` |
| Owner | `Code Runner` |

## Context

`FRAMEWORK.md` names Code Runner AI as the actor that implements tasks, and existing task files already hand off to `code-runner`. Until EV-007, the repository did not contain an operational skill contract for that role.

The framework now has prerequisites that Code Runner can consume: per-task files, `writeScope`, `sharedResources`, product-specific gates in `knowledge/conventions/gates.md`, and independent QA rules.

## Decision

Create `.codex/skills/code-runner/SKILL.md` as the operational implementation skill.

Code Runner works on exactly one task per invocation, uses TDD, respects `writeScope`, reads gates from `knowledge/conventions/gates.md`, stops on missing decisions or specification gaps, runs gates until green or reports blockers, and never commits, pushes, merges, or creates approval records.

## Consequences

| Type | Consequence |
| --- | --- |
| Positive | Tasks that already route to `code-runner` now have an executable owner. |
| Positive | Implementation remains bounded by Specification, graph, task contract, and writeScope. |
| Positive | QA remains independent because Code Runner does not mark QA evidence as passed. |
| Negative | Code Runner may stop frequently when task scope, gates, or source sections are incomplete. |
| Follow-up | Future EVs should add bug-fixer, code-review, commit-crafter, and pr-finalizer contracts. |

## FRAMEWORK.md Amendments

| Section | Amendment |
| --- | --- |
| `9. Skills` | Expand Code Runner AI from a named role into a bounded implementation contract. |
| `10. Orquestradores` | Clarify implementation handoff into QA after code changes. |
| `11. Gates De Aprovacao` | Clarify that Code Runner does not commit or create approval records. |

## Related Artifacts

| Artifact | Link |
| --- | --- |
| Code Runner skill | [../../.codex/skills/code-runner/SKILL.md](../../.codex/skills/code-runner/SKILL.md) |
| Gates convention | [../../knowledge/conventions/gates.md](../../knowledge/conventions/gates.md) |
| FRAMEWORK | [../../FRAMEWORK.md](../../FRAMEWORK.md) |

## Supersedes

- `N/A`
