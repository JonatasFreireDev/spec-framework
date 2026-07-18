# Engineering Systems Contract

This document contains the detailed contract for the shared Engineering System and Engineering Quality System. `FRAMEWORK.md` defines when these systems apply; this file defines their shared boundaries, versioning, evidence, and migration rules.

## Engineering System

The Engineering System versions stable architecture, module and data ownership, integrations, standards, quality attributes, operations, and evidence under `engineering/`. `engineering-system.md` is the human contract and `engineering-system.yaml` is the mechanical catalog. Origin is `generate`, `evolve`, or `adopt`; maturity records evidence and never grants approval.

Its scalable entity and standards model is defined in [`engineering-catalog-and-standards.md`](engineering-catalog-and-standards.md). Systems, applications, components, repositories, interfaces, data stores, and deployments are independent graph entities. Root catalogs are initialized once; entity records are created only from evidence or explicit hypotheses.

The standards system composes versioned rules through profiles. Consumers pin applicable profiles and standards. A narrower scope may add constraints but may not silently weaken inherited contracts. Exceptions require a stable ID, exact scope, owner, rationale, residual risk, mitigation, expiry or review date, re-entry gate, status, and approval where required.

The shared baseline is coordinated by `engineering-orchestrator`. It routes `technical-landscape`, `engineering-standards`, `operations-baseline`, and `engineering-evidence`, then asks `engineering-system` to consolidate the aggregate and Quality System. Specialist contracts remain independently owned even though the Engineering System composite approval hash governs their use together.

The orchestrator may execute these owners sequentially or delegate them through harness-native subagents. Delegation is phased rather than unrestricted: Technical Landscape first; Standards and Operations concurrently only while their scopes do not overlap; Evidence after both; and Engineering System aggregation last. Each delegated return is compact and hash-verifiable so the parent context retains conclusions and blockers instead of full specialist exploration.

Engineering System approval hashes its complete contract surface deterministically. A change to an approved shared contract makes its approval stale and requires human re-approval. Specification and approved product decisions remain authoritative when contracts conflict.

## Engineering Quality System

The Engineering Quality System is the shared quality contract under `engineering/quality/`. It covers quality attributes, test levels, risk-based coverage, environments, data, fitness functions, evidence, flaky tests, exceptions, and maturity. `quality-system.md` is human-readable, `quality-system.yaml` is mechanical, and supporting contracts define the quality model, test strategy, and fitness functions.

It defines policy and capability, not delivery approval. Use-case tests pin and apply it; tasks implement coverage; QA verifies acceptance criteria, configured gates, and real evidence; Security Review remains separate. Gate commands are canonical in `knowledge/conventions/gates.md`.

Maturity cannot waive gates or residual risk. Exceptions require scope, owner, rationale, mitigation, expiry or review date, re-entry gate, and open status. A Quality System migration is additive, previewable, rollback-safe, preserves adopter content, creates no approval evidence, and requires re-approval when the approved Quality System contract or its containing Engineering System composite hash changes.

## Legacy compatibility

Legacy products without the new contract remain compatible until they explicitly migrate their Engineering System or Quality System. Migration must be previewable, preserve product content and approval history, and identify the new approval or evidence required after the composite contract changes.

## Owning skills

- `engineering-orchestrator`: controls completeness, sequencing, blockers, and the human approval handoff.
- `technical-landscape`: owns the technical entity graph, topology, entities, relations, and boundary assessments.
- `engineering-standards`: owns standards, profiles, resolutions, conformance, and governed exceptions.
- `operations-baseline`: owns environments, deployments, runbooks, service levels, release, rollback, and continuity contracts.
- `engineering-evidence`: owns evidence inventory, coverage, maturity, gap, and staleness assessments.
- `engineering-system`: consolidates, versions, validates, and migrates the aggregate and Engineering Quality System.
- `technical-discovery`: maps delivery requirements to the codebase and stable engineering baseline.
- `engineering-proposal`: translates approved delivery contracts into an intended technical solution.
- `engineering-review`: independently reviews the proposal and its alignment with the pinned systems.
- `qa`: verifies Quality System application, configured gates, evidence, and exceptions.
