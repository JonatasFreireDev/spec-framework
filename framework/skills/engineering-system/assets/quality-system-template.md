# Engineering Quality System: [product]

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
| Engineering System | `[ENGSYS-XXX @ semver]` |
| Status | `[draft | proposed | approved]` |
| Mechanical catalog | [quality-system.yaml](quality-system.yaml) |
| Quality model | [quality-model.md](quality-model.md) |
| Test strategy | [test-strategy.md](test-strategy.md) |
| Fitness functions | [fitness-functions.yaml](fitness-functions.yaml) |
| Owner skill | `engineering-system` |

## Scope

[Define the repositories, surfaces, environments, platforms, and delivery types governed by this system.]

## Principles

- Quality requirements must be observable and traceable to `REQ-*` or `AC-*` identifiers.
- Coverage is proportional to user, business, operational, security, and change risk.
- Test implementation does not grant independent QA approval.
- Maturity describes evidence; it does not approve risk or waive gates.

## Capability Model

| Area | Policy | Required evidence | Maturity |
| --- | --- | --- | --- |
| Behavioral | [test-strategy.md](test-strategy.md) | `[test output]` | `[baseline | mapped | governed | verified | operated]` |
| Accessibility | [test-strategy.md](test-strategy.md) | `[audit/test/screenshot]` | `[maturity]` |
| Security and privacy | `[security baseline and strategy]` | `[test/review/log]` | `[maturity]` |
| Performance and reliability | [quality-model.md](quality-model.md) | `[benchmark/load/failure evidence]` | `[maturity]` |
| Observability | [quality-model.md](quality-model.md) | `[log/metric/trace assertion]` | `[maturity]` |

## Risk And Coverage Policy

| Risk | Minimum coverage | Required independent evidence |
| --- | --- | --- |
| Low | `[unit/component or explicit review method]` | `[gate output or limitation]` |
| Medium | `[unit plus integration/contract]` | `[gate output and acceptance mapping]` |
| High | `[negative, permission, integration, E2E, resilience as applicable]` | `[QA evidence plus specialized reviews]` |

## Environments And Test Data

| Environment or data class | Purpose | Ownership | Constraints | Evidence |
| --- | --- | --- | --- | --- |
| `[local/CI/staging/etc.]` | `[purpose]` | `[owner]` | `[privacy/isolation/reset]` | `[path/log]` |

## Evidence Policy

| Evidence | Required when | Source |
| --- | --- | --- |
| Gate output | A configured gate applies | `knowledge/conventions/gates.md` |
| Acceptance mapping | Every delivery | `tests.md` and `qa-evidence.md` |
| Visual/accessibility evidence | Delivery has a visual surface | Screenshot, audit, or test artifact |
| Fix verification | A defect or regression was found | Re-run output and permanent regression test |

## Flaky Tests And Exceptions

Every known flaky test or policy exception must record scope (`product` or a safe product-relative `domains/...` path), owner, rationale, residual risk, mitigation, a valid future expiry or review date, status, and re-entry gate. Only an `open`, unexpired exception whose scope contains the consuming use case may authorize a deviation. Closed and expired records remain auditable but cannot be consumed. An exception never changes an acceptance criterion or grants approval.

| ID | Scope | Owner | Rationale | Residual risk | Mitigation | Expires/review | Re-entry gate | Status |
| --- | --- | --- | --- | --- | --- | --- | --- | --- |
| `[QEX-001]` | `[scope]` | `[owner]` | `[why temporarily necessary]` | `[risk]` | `[mitigation]` | `[date]` | `[gate]` | `[open/closed]` |

## Consumers

| Use case | Pinned Engineering System | Deviations |
| --- | --- | --- |
| `[artifact path]` | `[ENGSYS-XXX @ semver]` | `[none or QEX-*]` |

## Handoff

Next: `qa` for delivery-specific test planning or independent evidence, depending on delivery state.

## ✅ Agent Verification Checklist

- [ ] Scope, principles, capability levels, environments, data, and mechanical catalog agree.
- [ ] Risk triggers map to required test levels, evidence, gates, and fitness functions.
- [ ] Flaky-test, exception, expiry, ownership, and escalation policies are explicit.
- [ ] Consumers, versions, migrations, deviations, and handoff requirements are traceable.
