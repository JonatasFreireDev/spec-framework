# FDR-007: Code Review operational contract

## Snapshot

| Field | Value |
| --- | --- |
| ID | `FDR-007` |
| Status | `approved` |
| Origin EV | `EV-009` |
| Date | `2026-07-10` |
| Owner | `Code Review` |

## Context

`FRAMEWORK.md` requires review before `validated`, but the repository did not have an operational Code Review skill or canonical review artifact.

Code Runner implements, Bug Fixer repairs, QA independently verifies, and Security Review checks risk controls. Code Review fills the read-only implementation review gate between implementation evidence and validation/release readiness.

## Decision

Create `.codex/skills/code-review/SKILL.md` and `framework/template/code-review-template.md`.

Code Review is read-only and uses three lenses:

- Completeness: all Specification, task, acceptance, and evidence requirements are delivered.
- Adherence: implementation follows approved architecture, contracts, data, permissions, and non-goals.
- Quality: code is maintainable, minimal, tested, safe around errors, and free of obvious dead code or unsafe behavior.

Findings are classified as `blocker`, `required_fix`, or `note`. Blockers and required fixes must include route and owner using FDR-006.

Validated or released executable artifacts require approved Code Review with a passing verdict and no blocking findings.

## Consequences

| Type | Consequence |
| --- | --- |
| Positive | The framework's review gate now has a concrete owner and artifact. |
| Positive | Review findings route cleanly into bug-fixer, code-runner, QA, or Product Historian. |
| Negative | Validation requires one more evidence artifact for executable work. |
| Follow-up | Future PR finalization should link PR review surfaces back to `code-review.md`. |

## FRAMEWORK.md Amendments

| Section | Amendment |
| --- | --- |
| `9. Skills` | Expand Code Review from named role to read-only review contract. |
| `10. Orquestradores` | Add review routing through FDR-006. |
| `11. Gates De Aprovacao` | Require approved Code Review for `validated+` executable artifacts. |

## Related Artifacts

| Artifact | Link |
| --- | --- |
| Code Review skill | [Code Review skill](../skills/code-review/SKILL.md) |
| Code Review template | [Code Review template](../template/code-review-template.md) |
| Failure routing FDR | [FDR-006-failure-routing-and-regression.md](FDR-006-failure-routing-and-regression.md) |
| FRAMEWORK | [../../FRAMEWORK.md](../../FRAMEWORK.md) |

## Supersedes

- `N/A`
