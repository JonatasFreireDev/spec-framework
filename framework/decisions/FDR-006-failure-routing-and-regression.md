# FDR-006: Failure routing and permanent regression fixes

## Snapshot

| Field | Value |
| --- | --- |
| ID | `FDR-006` |
| Status | `approved` |
| Origin EV | `EV-008` |
| Date | `2026-07-10` |
| Owner | `Release Orchestrator` |

## Context

QA Evidence and Security Review can block validation, but blockers need an explicit owner and route. Without routing, a red gate can stall indefinitely or be bypassed by an agent that silently edits the wrong artifact.

The framework now has `code-runner` for implementation and QA as an independent read-only verifier. EV-008 adds `bug-fixer` for confirmed defects and a routing matrix for failures.

## Decision

Use a failure routing matrix:

| Finding Type | Route |
| --- | --- |
| Defect, regression, security bug with clear expected behavior, production error | `bug-fixer` |
| Missing test, hollow test, missing negative/permission coverage | `qa` or test owner |
| Incomplete implementation or code outside the task contract | `code-runner` |
| Missing product decision or ambiguous business/security rule | `product-historian` plus human approval |
| Framework method gap | FDR / Evolution Orchestrator |

Every defect that escaped into QA, Security Review, audit, or production must get a permanent regression test before the fix is considered complete.

After any code change, the workflow re-enters QA. A red QA or Security Review gate cannot be skipped.

Automated fix attempts are capped at three per gate or finding. After the third failed attempt, the owner escalates to the human with reproduction, logs, attempted fixes, and remaining hypothesis.

## Consequences

| Type | Consequence |
| --- | --- |
| Positive | Blockers get a clear next owner instead of becoming vague follow-ups. |
| Positive | Escaped defects become permanent regression coverage. |
| Positive | QA independence remains intact because fixes return to QA rather than self-approving. |
| Negative | More findings require structured route metadata. |
| Follow-up | Future EVs should add code-review and PR finalization routing. |

## FRAMEWORK.md Amendments

| Section | Amendment |
| --- | --- |
| `10. Orquestradores` | Add failure routing and mandatory QA re-entry after code changes. |
| `11. Gates De Aprovacao` | Red gates cannot be bypassed; blockers require route/owner before approved validation artifacts can advance. |

## Related Artifacts

| Artifact | Link |
| --- | --- |
| Bug Fixer skill | [Bug Fixer skill](../skills/bug-fixer/SKILL.md) |
| QA skill | [QA skill](../skills/qa/SKILL.md) |
| Security Review skill | [Security Review skill](../skills/security-review/SKILL.md) |
| FRAMEWORK | [../../FRAMEWORK.md](../../FRAMEWORK.md) |

## Supersedes

- `N/A`
