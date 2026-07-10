# Security Baseline

## Snapshot

| Field | Value |
| --- | --- |
| Status | `template baseline` |
| Governed by | [FDR-009](../../engineering/decisions/FDR-009-threat-modeling-baseline.md) |
| Owner | Threat Modeler AI |
| Consumers | Security Review AI, QA AI, Code Runner AI, Audit Orchestrator |
| Purpose | Declare product-specific security rules that must be reused across domains and delivery artifacts. |

## Rule

Each adopting product should keep a security baseline in this folder or split it into domain-specific files linked from this index. Security Review must read the relevant baseline before validating an executable artifact.

This baseline is not a checklist substitute. It records product-specific rules, assumptions, boundaries, and evidence expectations that complement generic security review.

## Baseline Sections

| Section | Required Content | Consumer |
| --- | --- | --- |
| Actors and roles | Canonical users, admins, service actors, external systems, and trust level. | Specification, Security Review |
| Data classes | PII, sensitive business data, public data, derived data, retention rules. | Specification, QA, Security Review |
| Trust boundaries | Client/server boundary, database boundary, third-party boundary, job/queue boundary. | Design, Implementation Plan |
| Authorization rules | Server-side ownership, role checks, escalation limits, default deny behavior. | Code Runner, Code Review |
| Abuse controls | Rate limits, replay prevention, enumeration prevention, anti-spam/fraud controls. | QA, Threat Modeler |
| Logging and analytics | What must be logged, what must never be logged, audit trails, redaction rules. | QA, Security Review |
| Operational controls | Migrations, feature flags, rollback, monitoring, alerting, incident evidence. | Release Orchestrator |
| Residual risk policy | Which risks require human approval and where that approval is recorded. | Security Review, Product Historian |

## Domain Baseline Index

| Domain | Baseline Link | Status | Notes |
| --- | --- | --- | --- |
| `_example-domain` | `TBD by product adopter` | `placeholder` | Example domain has no product-specific security rules. |
| `events` | `TBD by product adopter` | `placeholder` | Event flows should define attendee, organizer, QR/check-in, replay, and privacy rules before real implementation. |

## Security Review Reading Contract

Before issuing a Security Review verdict, read:

1. This file.
2. Any domain-specific baseline linked above.
3. Active entries in [../../audits/security/threat-register.md](../../audits/security/threat-register.md) that affect the artifact.
4. Approved product decisions related to permissions, privacy, data retention, payments, uploads, UGC, public surfaces, or admin behavior.

## Maintenance Triggers

Update the baseline when:

- a domain introduces a new actor, role, permission, or trust boundary;
- a use case handles PII, payments, upload, UGC, public surfaces, or admin operations;
- a threat register entry becomes a reusable rule;
- a residual risk is accepted by a human decision;
- a stack dependency or deployment topology changes security expectations.

## Open Questions

| Question | Owner | Needed Before |
| --- | --- | --- |
| Which product-specific auth provider, database policy model, and deployment surface should future adopters document here? | Product adopter | First executable implementation |
| Should domain-specific baselines be separate files by default or sections in this index for small products? | Threat Modeler AI | First real product adoption |
