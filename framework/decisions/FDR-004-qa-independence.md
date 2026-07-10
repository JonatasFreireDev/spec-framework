# FDR-004: QA independence

## Snapshot

| Field | Value |
| --- | --- |
| ID | `FDR-004` |
| Status | `approved` |
| Origin EV | `EV-011` |
| Date | `2026-07-10` |
| Owner | `QA` |

## Context

QA Evidence exists as a gate, but QA must not become a checkbox reviewer that trusts task declarations. The gate is valuable only when QA independently verifies acceptance criteria, real gate output, regressions, security controls, and visual evidence when applicable.

## Decision

QA is an independent, read-only verifier.

QA must:

- read `knowledge/conventions/gates.md`;
- re-run applicable gates when possible;
- record real command output, logs, CI URLs, screenshots, or explicit limitations;
- hunt for hollow tests, missing negative or permission cases, scope drift, and Specification divergence;
- require proportional visual evidence only for deliveries with UI;
- check basic accessibility for UI deliveries;
- return blockers through routing instead of fixing code.

QA must not:

- edit application code;
- create, edit, or repair approval records;
- mark evidence as real when the gate was not run;
- declare `validated` when required evidence is placeholder-only.

## Consequences

| Type | Consequence |
| --- | --- |
| Positive | Validation evidence becomes independently trustworthy. |
| Positive | UI and accessibility risks become visible before release. |
| Negative | QA may take longer because gates are re-executed. |
| Follow-up | Future bug-fixer and code-runner skills should consume QA findings through a routing matrix. |

## FRAMEWORK.md Amendments

| Section | Amendment |
| --- | --- |
| `10. Orquestradores` | QA blockers route back to implementation/fix skills in future EVs instead of being repaired by QA. |
| `11. Gates De Aprovacao` | `validated+` requires non-placeholder QA evidence and real gate output or documented limitations. |

## Related Artifacts

| Artifact | Link |
| --- | --- |
| QA skill | [QA skill](../skills/qa/SKILL.md) |
| QA evidence template | [QA evidence template](../template/qa-evidence-template.md) |
| FRAMEWORK | [../../FRAMEWORK.md](../../FRAMEWORK.md) |

## Supersedes

- `N/A`
