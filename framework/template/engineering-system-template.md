# Engineering System: [product]

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
| ID | `ENGSYS-001` |
| Status | `draft` |
| Version | `0.1.0` |
| Origin mode | `generate | evolve | adopt` |
| Mechanical catalog | [engineering-system.yaml](engineering-system.yaml) |
| Owner skill | `engineering-system` |

## Scope

[Describe the product and repository boundaries covered by this system.]

## Architecture

| Area | Contract | Evidence | Maturity |
| --- | --- | --- | --- |
| System context | [architecture/system-context.md](architecture/system-context.md) | `[code/config/decision path]` | `baseline | mapped | governed | verified | operated` |
| Modules | [architecture/modules.md](architecture/modules.md) | `[code/test path]` | `baseline | mapped | governed | verified | operated` |
| Data ownership | `[path or Not configured]` | `[evidence]` | `[maturity]` |
| Integrations | `[path or Not configured]` | `[evidence]` | `[maturity]` |
| Deployment | `[path or Not configured]` | `[evidence]` | `[maturity]` |

## Standards And Quality

| Concern | Contract | Gate or evidence |
| --- | --- | --- |
| Quality system | [quality/quality-system.md](quality/quality-system.md) | `[policy/evidence]` |
| Quality model | [quality/quality-model.md](quality/quality-model.md) | `[gate/evidence]` |
| Test strategy | [quality/test-strategy.md](quality/test-strategy.md) | `[gate/evidence]` |
| Fitness functions | `[quality/fitness-functions.yaml or Not configured]` | `[command/evidence]` |
| Engineering standards | `[standards path or Not configured]` | `[gate/evidence]` |

## Decisions And Deviations

| Decision or deviation | Scope | Status | Evidence |
| --- | --- | --- | --- |
| `[DEC-* or deviation]` | `[scope]` | `[approved/proposed/open]` | `[path]` |

## Operations

| Capability | Runbook | Evidence | Maturity |
| --- | --- | --- | --- |
| Deploy and rollback | `[path or Not configured]` | `[evidence]` | `[maturity]` |
| Incident response | `[path or Not configured]` | `[evidence]` | `[maturity]` |

## Consumers

| Use case or system | Pinned version | Deviations |
| --- | --- | --- |
| `[artifact path]` | `[version]` | `[none or links]` |

## Handoff

Next: `technical-discovery` for delivery-specific work.

## ✅ Agent Verification Checklist

- [ ] Scope, version, origin, mechanical catalog, and architecture boundaries are consistent.
- [ ] Module, data, integration, standards, quality, security, and operational ownership are explicit.
- [ ] Decisions, deviations, migrations, consumers, and compatibility expectations are traceable.
- [ ] The handoff identifies required downstream pins, evidence, and unresolved system gaps.
