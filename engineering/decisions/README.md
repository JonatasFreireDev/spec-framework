# Framework Decision Records

## Purpose

This folder stores Framework Decision Records (FDRs): decisions about the framework method, agent contracts, gates, validation rules, and delivery workflow.

Product decisions do not live here. Product decisions continue to live in `knowledge/decisions/DEC-*` and must trace to a product problem, domain, feature, use case, or business rule.

## Index

| ID | Title | Status | Origin |
| --- | --- | --- | --- |
| [FDR-001](FDR-001-decision-governance.md) | Decision governance: product vs framework | approved | Governance baseline |
| [FDR-002](FDR-002-gate-commands.md) | Gate commands are product conventions | approved | EV-011 / gate governance |
| [FDR-003](FDR-003-writescope-safe-parallelism.md) | WriteScope and safe parallelism | approved | EV-012 |
| [FDR-004](FDR-004-qa-independence.md) | QA independence | approved | EV-011 |
| [FDR-005](FDR-005-code-runner-contract.md) | Code Runner operational contract | approved | EV-007 |

## Rule

If a decision changes the framework method, skill contract, validator behavior, or delivery workflow, record it as an FDR or as an explicit amendment to `FRAMEWORK.md`.

If a decision changes product domain behavior, business rules, data, permissions, privacy, payment, security, or delivery scope, record it as a DEC under `knowledge/decisions/`.
