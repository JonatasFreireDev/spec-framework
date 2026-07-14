# Security Baseline: [product or domain name]

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
