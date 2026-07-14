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
| [FDR-018](FDR-018-feature-to-task-operational-closure.md) | Feature-to-Task operational closure | approved | Governance baseline |
| [FDR-019](FDR-019-executable-product-decision-effects.md) | Executable product decision effects | approved | Governance baseline |
| [FDR-020](FDR-020-consolidated-cli-and-guided-migrations.md) | Consolidated CLI and guided migrations | approved | Governance baseline |
| [FDR-021](FDR-021-design-sources-and-visual-artifacts.md) | Design sources and visual artifact protocol | approved | Human-approved design workflow evolution |
| [FDR-022](FDR-022-conversational-cli-guidance.md) | Conversational CLI guidance | approved | Human-approved CLI usability evolution |
| [FDR-023](FDR-023-canonical-product-design-system.md) | Canonical product Design System | approved | Human-approved Design System evolution |
| [FDR-024](FDR-024-supervised-adapter-management.md) | Supervised adapter management | approved | Human-approved adapter usability evolution |
| [FDR-025](FDR-025-external-runtime-and-manifest-only-activation.md) | External runtime and manifest-only activation | proposed | Governance baseline |
| [FDR-026](FDR-026-canonical-product-engineering-system.md) | Canonical Product Engineering System | approved | Governance baseline |
| [FDR-027](FDR-027-engineering-gate-integrity.md) | Engineering Gate Integrity | approved | Governance baseline |
| [FDR-028](FDR-028-composite-engineering-approval-and-structured-gates.md) | Composite Engineering Approval And Structured Gates | approved | Governance baseline |
| [FDR-029](FDR-029-canonical-vision-companion-contracts.md) | Canonical Vision Companion Contracts | proposed | Governance baseline |
| [FDR-030](FDR-030-foundation-approval-registry.md) | Foundation Approval Registry | proposed | Governance baseline |
| [FDR-031](FDR-031-feature-scoped-bootstrap.md) | Feature-Scoped Bootstrap | proposed | Adopter first-run feedback |
| [FDR-032](FDR-032-existing-implementation-assessment.md) | Existing Implementation Assessment | proposed | Adopter first-run review |
| [FDR-033](FDR-033-code-first-existing-product-baseline.md) | Code-First Existing Product Baseline | proposed | Adopter starting-point review |
| [FDR-034](FDR-034-existing-documents-import-gate.md) | Existing Documents Import Gate | proposed | Adopter starting-point review |
| [FDR-035](FDR-035-audit-only-mutation-guard.md) | Audit-Only Mutation Guard | proposed | Adopter starting-point review |
| [FDR-036](FDR-036-registry-relationship-preservation.md) | Registry Relationship Preservation | proposed | Starting-point gate review |
| [FDR-037](FDR-037-engineering-quality-system.md) | Engineering Quality System | proposed | Governance baseline |
| [FDR-038](FDR-038-guide-first-dispatch.md) | Guide-First Dispatch | proposed | Governance baseline |
| [FDR-039](FDR-039-domain-modeling-reference-and-warnings.md) | Domain Modeling Reference And Warnings | proposed | Adopter domain-modeling feedback |
| [FDR-040](FDR-040-cobra-cli-command-tree.md) | Cobra CLI Command Tree | proposed | CLI architecture evolution |
| [FDR-041](FDR-041-declarative-initialization-contracts.md) | Declarative Initialization Contracts | proposed | CLI architecture evolution |
| [FDR-042](FDR-042-cli-lifecycle-commands.md) | CLI Lifecycle Commands | proposed | CLI installation usability |
| [FDR-043](FDR-043-native-discovery-and-challenge.md) | Native Discovery And Challenge | proposed | Planning skill interaction feedback |
| [FDR-044](FDR-044-lean-product-readmes.md) | Lean Product READMEs | proposed | Adopter starter context review |
| [FDR-045](FDR-045-canonical-method-compression.md) | Canonical Method Compression | proposed | Agent context efficiency review |

## Rule

If a decision changes the framework method, skill contract, validator behavior, or delivery workflow, record it as an FDR or as an explicit amendment to `FRAMEWORK.md`.

If a decision changes product domain behavior, business rules, data, permissions, privacy, payment, security, or delivery scope, record it as a DEC under `knowledge/decisions/`.
