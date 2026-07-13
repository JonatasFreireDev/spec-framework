# FDR-037: Engineering Quality System

## Snapshot

| Field | Value |
| --- | --- |
| ID | `FDR-037` |
| Status | `proposed` |
| Origin EV | `Governance baseline` |
| Date | `2026-07-13` |
| Owner | `Engineering System` |

## Context

The framework already defines delivery-specific `tests.md`, independent `qa-evidence.md`, QA independence, an Engineering Quality Model, and fitness functions. It does not yet provide one versioned shared contract for test levels, risk-based coverage, environments, evidence policy, flaky tests, exceptions, and the boundary between test planning and independent QA.

Without that contract, use cases can repeat product-wide policy, Engineering System quality maturity can remain underspecified, and QA can verify evidence without a stable shared quality baseline.

## Decision

Introduce the **Engineering Quality System** as the canonical shared quality contract under `engineering/quality/`.

| Concern | Canonical owner | Boundary |
| --- | --- | --- |
| Shared quality policy and capability | Engineering Quality System | Defines product-wide expectations, test levels, risk coverage, environments, evidence, exceptions, and maturity. |
| Delivery-specific validation plan | `tests.md` | Pins the Engineering System version and applies the shared policy to one use case. |
| Test implementation | Delivery tasks and Code Runner | Implements tests but does not grant QA approval. |
| Independent verification | QA and `qa-evidence.md` | Re-runs applicable gates, verifies the pinned policy and real evidence, and remains read-only. |
| Security judgment | Security Review | Remains a specialized validation gate and is not absorbed by the Quality System. |

The canonical package is:

- `engineering/quality/quality-system.md`: human contract;
- `engineering/quality/quality-system.yaml`: mechanical catalog;
- `engineering/quality/quality-model.md`: product quality attributes and required evidence;
- `engineering/quality/test-strategy.md`: shared test levels, risk coverage, environments, and data policy;
- `engineering/quality/fitness-functions.yaml`: configured mechanical checks.

The Engineering Quality System is part of the Engineering System composite approval hash. Its maturity records available evidence and never grants approval. Gate commands remain product conventions in `knowledge/conventions/gates.md`.

Engineering System approval first rejects an invalid configured Quality System, then updates `status` atomically in `engineering/context.md`, `engineering-system.md`, `engineering-system.yaml`, `quality/quality-system.md`, and `quality/quality-system.yaml` before calculating the composite hash. Human and mechanical capability maturity must agree, and non-baseline evidence must be safe and resolvable. Proposed-or-later `tests.md` must structurally apply the pinned canonical policy, select configured environment/data/platform values, and declare `None` or consumable `QEX-*` records; approved-or-later QA evidence must record `passed`, not `N/A`, for policy, environment/data, and flaky-test/exception checks.

Quality exceptions are mechanically consumable only while `open`, before their ISO `expiry_or_review` date, and when their scope is `product` or a safe `domains/...` path containing the consumer. Closed or expired exceptions remain history but do not authorize deviations.

Compatibility is additive. Existing adopter repositories without the new contracts continue to use the legacy `quality-model.md` area until they explicitly run `spec-framework engineering-system migrate`. That migration previews changes with `--dry-run`, preserves the legacy quality model and existing quality files, materializes only missing Quality System contracts before atomically replacing the catalog, rolls back generated files on failure, and never creates approval records. Because the Engineering System composite hash changes, an approved system requires explicit human re-approval. General `upgrade` refreshes only the external runtime and pinned manifest; it never mutates adopter-owned quality policy.

A future move to a top-level `quality/` pillar requires a separate FDR supported by evidence that recurring quality ownership materially extends beyond engineering boundaries.

## Consequences

| Type | Consequence |
| --- | --- |
| Positive | Product-wide test and evidence policy becomes reusable, versioned, and traceable. |
| Positive | `tests.md` and QA evidence gain an explicit shared baseline without weakening QA independence. |
| Negative | Engineering System evolution requires maintaining additional human and mechanical contracts. |
| Follow-up | Evaluate a top-level Quality pillar only after adopter evidence shows independent non-engineering ownership. |

## FRAMEWORK.md Amendments

| Section | Amendment |
| --- | --- |
| `3. Conceptual Model` | Define Engineering Quality System and its relationship to tests, QA, and Security Review. |
| `4. Folder Structure` | Add the canonical quality contract package under `engineering/quality/`. |
| `6. Specification Driven Development` | Require delivery test plans to apply a pinned shared policy when configured. |

## Related Artifacts

| Artifact | Link |
| --- | --- |
| QA independence | [FDR-004](FDR-004-qa-independence.md) |
| Engineering System | [FDR-026](FDR-026-canonical-product-engineering-system.md) |
| Quality System template | [Quality System template](../template/quality-system-template.md) |
| QA skill | [QA skill](../skills/qa/SKILL.md) |
| Framework | [FRAMEWORK.md](../../FRAMEWORK.md) |

## Supersedes

- `N/A`
