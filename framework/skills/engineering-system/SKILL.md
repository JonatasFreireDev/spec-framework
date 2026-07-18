---
name: engineering-system
description: "Engineering System Skill. Use when an agent needs to create, adopt, evolve, version, or audit shared product architecture, standards, quality attributes, fitness functions, operations, and engineering evidence in the Spec Framework workflow."
---

# Engineering System Skill

## Layer
Planning

## Responsibility
Own the shared `engineering/engineering-system.md` and `engineering-system.yaml` contracts. Maintain stable product engineering knowledge without making product decisions, planning a specific delivery, or writing application code.

## Operating modes
- create: establish the first evidence-backed Engineering System.
- update: revise stable contracts while preserving approved decisions and compatibility.
- adopt: register an authoritative versioned external engineering source.
- audit: find stale evidence, missing ownership, ungoverned boundaries, and unsupported maturity claims.
- explain: summarize the system, maturity, consumers, and gaps.

## Inputs
Product context; real code and test tree; deployment and operations configuration; approved decisions; conventions; runbooks; dependency evidence.

## Outputs
`engineering/engineering-system.md`; `engineering/engineering-system.yaml`; architecture, standards, quality, runbook, and evidence links; decision candidates; version and compatibility notes.

## Required reading
- [`engineering-systems.md`](../../docs/engineering-systems.md) for shared Engineering System and Quality System versioning, migration, and approval boundaries.
- the framework root's `FRAMEWORK.md`
- Relevant parent context.md files.
- This skill owns its generation resources: the engineering-system, quality-system, quality-model, test-strategy, and fitness-function resources in `assets/`.
- Approved decisions are discovered through the active product root's `.product/decisions.json`; resolve each registered `path` from its declared domain root (`knowledge/decisions/`, `design/decisions/`, or `engineering/decisions/`).

## Discovery and challenge

Follow the shared [Discovery And Challenge contract](../discovery-and-challenge.md) before substantive creation or material revision.

## Workflow
1. Inspect the real product code, tests, configuration, environments, operations evidence, and existing engineering documents.
2. Define the covered product and repository boundaries and choose `generate`, `evolve`, or `adopt`.
3. Inventory modules, ownership, data, integrations, standards, quality attributes, test strategy, gates, environments, test data, runbooks, and consumers with evidence paths.
4. Maintain `engineering/quality/quality-system.md` and `quality-system.yaml` together with the quality model, test strategy, and fitness functions. Keep commands canonical in `knowledge/conventions/gates.md`.
5. Declare maturity per area only when its required evidence exists; maturity never implies approval.
6. Link approved `DEC-*` records for governed choices and record candidates for missing structural, data, security, or operational decisions.
7. Update the human contract and mechanical catalogs together, applying semantic versioning and compatibility notes.
8. Stop before inventing architecture, editing application code, or creating approval records.

## Quality checklist
- [ ] Preserves traceability to affected artifacts.
- [ ] Uses the correct template and naming conventions.
- [ ] Human and mechanical contracts agree on identity, version, origin, scope, and maturity.
- [ ] Every maturity claim points to concrete evidence.
- [ ] Human and mechanical capability maturity agree, and non-baseline evidence references are safe and resolvable.
- [ ] Quality policy separates shared expectations, delivery-specific test planning, test implementation, independent QA, and Security Review.
- [ ] Quality exceptions identify scope, owner, rationale, residual risk, mitigation, expiry or review date, re-entry gate, and status.
- [ ] Only open, unexpired, in-scope exceptions are passed to delivery consumers.
- [ ] Stable knowledge is not duplicated into delivery proposals.
- [ ] Detects gaps, conflicts, and dependencies.
- [ ] Records or requests decisions for meaningful changes.
- [ ] Leaves a clear handoff for the next skill or orchestrator.

## Handoff
Next: `technical-discovery` for delivery-specific work.

Pass forward the pinned system id/version, architecture and ownership contracts, standards, quality attributes, decisions, evidence, gaps, and compatibility constraints.
