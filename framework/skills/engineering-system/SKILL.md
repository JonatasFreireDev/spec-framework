---
name: engineering-system
description: "Engineering System Skill. Use when an agent needs to consolidate, version, validate, migrate, or explain the shared Engineering System aggregate and Engineering Quality System after specialist engineering contracts are ready."
---

# Engineering System Skill

## Layer
Planning

## Responsibility
Own the aggregate `engineering/engineering-system.md`, `engineering-system.yaml`, and Engineering Quality System contracts. Consolidate specialist-owned technical landscape, standards, operations, and evidence without authoring those contracts, making product decisions, planning a specific delivery, or writing application code.

## Operating modes
- create: consolidate the first evidence-backed Engineering System after specialist outputs exist.
- update: revise the aggregate and quality contracts while preserving approved decisions and compatibility.
- adopt: register an authoritative versioned external engineering source.
- audit: find stale evidence, missing ownership, ungoverned boundaries, and unsupported maturity claims.
- explain: summarize the system, maturity, consumers, and gaps.

## Inputs
Engineering Orchestrator handoff; specialist-owned technical landscape, standards, operations, and evidence contracts; product context; Quality System evidence; approved decisions.

## Outputs
`engineering/engineering-system.md`; `engineering/engineering-system.yaml`; Engineering Quality System contracts; consolidated area maturity and references; version, compatibility, migration, and approval-impact notes.

## Required reading
- [`engineering-systems.md`](../../docs/engineering-systems.md) for shared Engineering System and Quality System versioning, migration, and approval boundaries.
- [`engineering-catalog-and-standards.md`](../../docs/engineering-catalog-and-standards.md) for the scalable entity graph, standards inheritance, exceptions, and materialization boundary.
- the framework root's `FRAMEWORK.md`
- Relevant parent context.md files.
- This skill owns only the aggregate Engineering System and Engineering Quality System generation resources in `assets/`: Engineering System, quality system and model, test strategy, and fitness functions.
- Approved decisions are discovered through the active product root's `.product/decisions.json`; resolve each registered `path` from its declared domain root (`knowledge/decisions/`, `design/decisions/`, or `engineering/decisions/`).

## Discovery and challenge

Follow the shared [Discovery And Challenge contract](../discovery-and-challenge.md) before substantive creation or material revision.


Use scripts/invoke-cli.ps1 on Windows or scripts/invoke-cli.sh on macOS/Linux for the CLI operation in this skill's reviewed scope. The wrapper never adds --yes or an approver identity.

## Workflow
1. Require a current `engineering-orchestrator` handoff and read the specialist-owned technical landscape, standards, operations, and evidence contracts.
2. Stop when a required specialist contract is absent, inconsistent, stale, or reports a blocking gap; return the exact scope to the orchestrator.
3. Define the covered product boundary and choose `generate`, `evolve`, or `adopt` without redefining specialist-owned content.
4. Maintain `engineering/quality/quality-system.md` and `quality-system.yaml` together with the quality model, test strategy, and fitness functions. Keep commands canonical in `knowledge/conventions/gates.md`.
5. Consolidate area references and maturity only when the evidence skill supports each claim; maturity never implies approval.
6. Link approved `DEC-*` records and surface unresolved decision candidates without creating decisions or approvals.
7. Update the human and mechanical aggregate together, applying semantic versioning, compatibility, migration, consumer, and stale-approval notes.
8. Run `engineering-system validate` and return the current composite hash and blockers to `engineering-orchestrator`.
9. Stop before editing specialist catalogs, application code, product decisions, or approval records.

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
- [ ] Every specialist-owned area is referenced rather than duplicated.
- [ ] Aggregate maturity agrees with the evidence assessment and specialist contracts.
- [ ] Version, compatibility, consumers, migration, and approval staleness are explicit.
- [ ] Detects gaps, conflicts, and dependencies.
- [ ] Records or requests decisions for meaningful changes.
- [ ] Leaves a clear handoff for the next skill or orchestrator.

## Handoff
Next: `engineering-orchestrator` for completeness review and the human approval gate.

Pass forward the system id/version, composite hash, quality contract, specialist references, maturity, decisions, blockers, compatibility, and required approval.
