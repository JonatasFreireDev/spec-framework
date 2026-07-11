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
| [FDR-006](FDR-006-failure-routing-and-regression.md) | Failure routing and permanent regression fixes | approved | EV-008 |
| [FDR-007](FDR-007-code-review-contract.md) | Code Review operational contract | approved | EV-009 |
| [FDR-008](FDR-008-delivery-commits-and-prs.md) | Delivery commits and PR finalization | approved | EV-010 |
| [FDR-009](FDR-009-threat-modeling-baseline.md) | Threat modeling baseline | approved | EV-014 |
| [FDR-010](FDR-010-markdown-link-validation.md) | Markdown link validation | approved | EV-013 |
| [FDR-011](FDR-011-core-starter-example-boundary.md) | Framework core, product starter, and examples boundary | approved | EV-015 |
| [FDR-012](FDR-012-framework-md-english-translation.md) | FRAMEWORK.md English translation | approved | Governance baseline |
| [FDR-013](FDR-013-go-cli-and-agent-skill-installation.md) | Go CLI and multi-agent skill installation | approved | CLI portability migration |
| [FDR-014](FDR-014-generated-adopter-bootstrap.md) | Generated adopter bootstrap | approved | Adopter first-run feedback |
| [FDR-015](FDR-015-starting-points-and-source-import.md) | Starting points and source import | approved | Governance baseline |
| [FDR-016](FDR-016-delivery-closure-and-operational-workspaces.md) | Delivery closure and operational workspaces | approved | Governance baseline |
| [FDR-017](FDR-017-resumable-parallel-runtime.md) | Resumable parallel runtime | approved | Governance baseline |

## Rule

If a decision changes the framework method, skill contract, validator behavior, or delivery workflow, record it as an FDR or as an explicit amendment to `FRAMEWORK.md`.

If a decision changes product domain behavior, business rules, data, permissions, privacy, payment, security, or delivery scope, record it as a DEC under `knowledge/decisions/`.
