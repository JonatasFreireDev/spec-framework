# Product Gates

## Snapshot

| Field | Value |
| --- | --- |
| Status | `draft` |
| Owner | Product adopter |
| Purpose | Declare the real product commands used by Code Runner, QA, CI, and release checks. |

## Rule

Replace placeholder commands before any task can move to `implemented`, `validated`, or `released` with strong evidence.

## Gate Catalog

| ID | Command | When Runs | Blocks Status From | Evidence Expected | Notes |
| --- | --- | --- | --- | --- | --- |
| `GATE-TYPECHECK` | `TBD by product adopter` | Before implementation is marked complete. | `implemented` | Command output or CI log. | Remove if the product has no typecheck. |
| `GATE-LINT` | `TBD by product adopter` | Before implementation and during QA. | `implemented` | Command output or CI log. | Include formatting if applicable. |
| `GATE-TEST` | `TBD by product adopter` | During implementation and QA. | `validated` | Test log or CI URL. | Include unit/integration/e2e scope. |
| `GATE-DATABASE` | `TBD by product adopter` | When migrations, policies, or seed data change. | `validated` | Migration or schema test output. | Required for database work. |
| `GATE-VISUAL` | `TBD by product adopter` | When a task changes UI or visual states. | `validated` | Screenshot, CI artifact, and accessibility notes. | Required for user-visible surfaces. |
| `GATE-SECURITY` | `TBD by product adopter` | When auth, permissions, PII, uploads, payment, public surfaces, or secrets are involved. | `validated` | Security scan, review notes, or explicit limitation. | Complements Security Review. |

## Evidence Rules

| Rule | Detail |
| --- | --- |
| No checkbox-only QA | QA evidence must include real output, CI URL, screenshot, or explicit limitation. |
| No hidden failures | Failing gates must route to owner and required fix. |
| Product-specific | Do not edit `.spec-framework/` to encode product commands. |

