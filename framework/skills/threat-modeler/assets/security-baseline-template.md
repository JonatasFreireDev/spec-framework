# Security Baseline: [product or domain name]

## 🧾 Generation And Agent Self-Check

> Complete this section when materializing the artifact. Keep unresolved items explicit in the relevant scope, findings, risks, or handoff section.

| Field | Value |
| --- | --- |
| Generated on | `YYYY-MM-DD` |
| Purpose | `[decision, evidence, contract, or handoff this artifact supports]` |
| Use when | `[workflow stage, trigger, or condition]` |
| Prepared by | `[owning skill, role, or accountable person]` |
| Scope covered | `[artifact, product area, use case, or review boundary]` |
| Required inputs and evidence | `[links to approved parents, documents, code, decisions, or observations]` |
| Ready when | `[artifact-specific completion, evidence, and gate conditions]` |
| Current status | `[status allowed by this artifact's owning workflow]` |


## Snapshot

| Field | Value |
| --- | --- |
| ID | `[SECBASE-XXX]` |
| Status | `[draft | proposed | approved]` |
| Scope | `[product | domain | goal | feature family]` |
| Owner skill | Threat Modeler AI |
| Governed by | `FRAMEWORK.md security policy or local DEC-*` |
| Last reviewed | `YYYY-MM-DD` |
| Next review trigger | `[scope change / release / incident / dependency change]` |

## Navigation

| Artifact | Link |
| --- | --- |
| Framework baseline convention | `[knowledge/conventions/security-baseline.md]` |
| Threat register | `[audits/security/threat-register.md]` |
| Related domain context | `[context.md or N/A]` |
| Related decisions | `[DEC-* links]` |

## Scope

| Field | Value |
| --- | --- |
| In scope | `[actors, domains, features, data, integrations]` |
| Out of scope | `[explicit exclusions]` |
| Assumptions | `[verified assumptions]` |
| Open questions | `[questions that block stronger guarantees]` |

## Actors And Roles

| Actor | Trust Level | Allowed Actions | Explicitly Forbidden |
| --- | --- | --- | --- |
| `[actor]` | `[trusted/semi-trusted/untrusted/service]` | `[actions]` | `[forbidden actions]` |

## Data Classes

| Data | Classification | Storage | Retention | Exposure Rules |
| --- | --- | --- | --- | --- |
| `[data]` | `[public/internal/sensitive/PII/secret]` | `[location]` | `[period]` | `[UI/API/log/analytics rules]` |

## Trust Boundaries

| Boundary | Risk | Required Control | Evidence Expected |
| --- | --- | --- | --- |
| `[client -> API]` | `[risk]` | `[control]` | `[test/log/review]` |

## Required Controls

| Control | Applies When | Required Evidence | Blocks Validation If Missing |
| --- | --- | --- | --- |
| Server-side authorization | `[condition]` | `[test/log/code review]` | `[yes/no]` |
| Idempotency/replay prevention | `[condition]` | `[test/log/design decision]` | `[yes/no]` |
| Sensitive data redaction | `[condition]` | `[log/analytics review]` | `[yes/no]` |

## Domain Threats

| Threat ID | Scenario | Impact | Mitigation | Register Link |
| --- | --- | --- | --- | --- |
| `[THR-XXX]` | `[attack scenario]` | `[impact]` | `[mitigation]` | `[threat-register.md#...]` |

## Residual Risk Policy

| Risk Type | Approval Needed | Approval Artifact | Notes |
| --- | --- | --- | --- |
| `[risk type]` | `[yes/no]` | `[DEC/approval record/PR review]` | `[notes]` |

## Handoff To Security Review

Security Review must verify that each applicable required control has evidence and that active threat-register entries are mitigated, explicitly accepted, or route-blocked before validation.

## ✅ Agent Verification Checklist

- [ ] Scope, actors, roles, data classes, trust boundaries, and threat sources are complete.
- [ ] Required controls have owners, verification evidence, and applicable decisions.
- [ ] Domain threats and residual-risk policy define review and escalation triggers.
- [ ] The handoff tells Security Review which baseline, threats, exceptions, and evidence to inspect.
