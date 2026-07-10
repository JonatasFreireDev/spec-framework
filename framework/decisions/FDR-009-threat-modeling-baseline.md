# FDR-009: Threat modeling baseline

## Snapshot

| Field | Value |
| --- | --- |
| ID | `FDR-009` |
| Status | `approved` |
| Origin EV | `EV-014` |
| Date | `2026-07-10` |
| Owner | `Security Review AI / Threat Modeler AI` |

## Context

`FRAMEWORK.md` defines Security Review as a validation gate for executable artifacts, and `.codex/skills/security-review/SKILL.md` evaluates one artifact bundle at a time. That reactive gate is necessary, but it does not create product-wide or domain-level security memory.

The framework needs a proactive security layer that records recurring product rules, known threat scenarios, residual risks, and mitigations before individual use cases reach validation.

## Decision

Add Threat Modeler AI as the proactive security-planning skill. It owns product-wide and domain-level threat modeling, security baselines, and the living threat register.

Security baselines live under `knowledge/conventions/` because they are product-adopter conventions that downstream Security Review must read. Concrete threat tracking lives under `audits/security/`, starting with `audits/security/threat-register.md`.

Security Review remains the artifact-level validation gate. It must read the relevant security baseline and active threat register entries before issuing a verdict.

## Consequences

| Type | Consequence |
| --- | --- |
| Positive | Security context becomes reusable across domains and no longer depends only on late artifact review. |
| Positive | Security Review can validate against product-specific rules, not only a generic checklist. |
| Negative | Products must maintain one more living artifact: the baseline and threat register can go stale if not reviewed after scope changes. |
| Negative | Threat modeling may surface decision gaps earlier and slow risky work until humans approve residual risk. |
| Follow-up | Future validators may check that high-risk domains or Tier L use cases link to baseline and threat register entries. |

## FRAMEWORK.md Amendments

| Section | Amendment |
| --- | --- |
| `6. Specification Driven Development` | Clarify that Security Review reads the product security baseline and threat register. |
| `9. Skills` | Add Threat Modeler AI under Engineering and Validation. |
| `10. Orchestrators` | Add proactive threat modeling before Security Review. |
| `13. Auditoria` | Clarify that security audits may maintain a living threat register. |

## Related Artifacts

| Artifact | Link |
| --- | --- |
| FRAMEWORK | [../../FRAMEWORK.md](../../FRAMEWORK.md) |
| Threat Modeler skill | [Threat Modeler skill](../skills/threat-modeler/SKILL.md) |
| Security baseline convention | [../../knowledge/conventions/security-baseline.md](../../examples/events/knowledge/conventions/security-baseline.md) |
| Threat register | [../../audits/security/threat-register.md](../../examples/events/audits/security/threat-register.md) |

## Supersedes

- `N/A`
