# Technical Gates

## Snapshot

| Field | Value |
| --- | --- |
| Status | `placeholder` |
| Governed by | `FRAMEWORK.md configured gate policy` |
| Owner | Product adopter |
| Purpose | Declare product-specific technical gates for implementation and QA. |

## Rule

The framework is stack-agnostic. Each adopting product must replace these placeholders with real commands before implementation and validation gates can produce strong evidence.

Code Runner and QA read this file. They must not hardcode project gate commands in their skill contracts.

## Gate Catalog

| ID | Command | When runs | Blocks status from | Evidence expected | Notes |
| --- | --- | --- | --- | --- | --- |
| `GATE-TYPECHECK` | `TBD by product adopter` | Before implementation is marked complete and during QA verification. | `implemented` | Command output or CI log. | Placeholder. Example command might be a product typecheck. |
| `GATE-LINT` | `TBD by product adopter` | Before implementation is marked complete and during QA verification. | `implemented` | Command output or CI log. | Placeholder. Example command might be a product lint check. |
| `GATE-TEST` | `TBD by product adopter` | During implementation, regression verification, and QA. | `validated` | Test log, CI URL, or captured local output. | Placeholder. Include unit, integration, and e2e commands as product needs. |
| `GATE-DATABASE` | `TBD by product adopter` | When tasks touch migrations, policies, seed data, or local database state. | `validated` | Migration/test output. | Placeholder. Use only when the product has database gates. |
| `GATE-VISUAL` | `TBD by product adopter` | When a delivery has UI or user-visible states. | `validated` | Screenshot path or CI artifact plus accessibility notes. | Placeholder. Required only for visual surface. |

## Adoption Checklist

- [ ] Replace each `TBD by product adopter` command with the product command or remove gates that do not apply.
- [ ] Add gates for framework-specific project needs, such as security scans, generated client checks, localization checks, or schema drift checks.
- [ ] Keep gate IDs stable after tasks reference them.
- [ ] Record real output in `qa-evidence.md`; do not replace command output with a checkbox-only verdict.
