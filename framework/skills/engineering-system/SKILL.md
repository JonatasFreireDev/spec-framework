---
name: engineering-system
description: "Engineering System Skill. Use when an agent needs to create, adopt, evolve, version, or audit shared product architecture, standards, quality attributes, fitness functions, operations, and engineering evidence in the Spec Framework workflow."
---

# Engineering System Skill

## Layer
Planning

## Responsibility
Own the shared `engineering/engineering-system.md` and `engineering-system.yaml` contracts, technical entity graph, standards, operations catalogs, quality policy, and evidence index. Maintain stable product engineering knowledge without making product decisions, planning a specific delivery, or writing application code.

## Operating modes
- create: establish the first evidence-backed Engineering System.
- update: revise stable contracts while preserving approved decisions and compatibility.
- adopt: register an authoritative versioned external engineering source.
- audit: find stale evidence, missing ownership, ungoverned boundaries, and unsupported maturity claims.
- explain: summarize the system, maturity, consumers, and gaps.

## Inputs
Product context; real code and test tree; deployment and operations configuration; approved decisions; conventions; runbooks; dependency evidence.

## Outputs
`engineering/engineering-system.md`; `engineering/engineering-system.yaml`; root technical catalog; on-demand system, application, component, repository, interface, data-store, and deployment records; standards catalog, profiles, standards, and exceptions; architecture, quality, operations, runbook, and evidence links; decision candidates; version and compatibility notes.

## Required reading
- [`engineering-systems.md`](../../docs/engineering-systems.md) for shared Engineering System and Quality System versioning, migration, and approval boundaries.
- [`engineering-catalog-and-standards.md`](../../docs/engineering-catalog-and-standards.md) for the scalable entity graph, standards inheritance, exceptions, and materialization boundary.
- the framework root's `FRAMEWORK.md`
- Relevant parent context.md files.
- This skill owns every generation resource in `assets/`: Engineering System, technical catalog and entity, standards catalog, standard, profile, exception, quality system and model, test strategy, fitness functions, operations catalog, and evidence inventory.
- Approved decisions are discovered through the active product root's `.product/decisions.json`; resolve each registered `path` from its declared domain root (`knowledge/decisions/`, `design/decisions/`, or `engineering/decisions/`).

## Discovery and challenge

Follow the shared [Discovery And Challenge contract](../discovery-and-challenge.md) before substantive creation or material revision.


Use scripts/invoke-cli.ps1 on Windows or scripts/invoke-cli.sh on macOS/Linux for the CLI operation in this skill's reviewed scope. The wrapper never adds --yes or an approver identity.

Use `scripts/inventory-engineering-evidence.ps1` on Windows or `scripts/inventory-engineering-evidence.sh` on macOS/Linux before authoring the shared baseline. It inventories deterministic engineering evidence and can call `engineering-system validate`; it never changes product or code files.

## Workflow
1. Inspect the real product code, tests, configuration, environments, operations evidence, and existing engineering documents. If no code exists, create explicit hypothesis contracts, pending decisions, and intended constraints rather than pretending evidence exists.
2. Define the covered product and repository boundaries and choose `generate`, `evolve`, or `adopt`.
3. Inventory systems, applications, components, repositories, ownership, data stores, interfaces, deployments, standards, quality attributes, environments, runbooks, and consumers with evidence paths. Model stable IDs and graph relations instead of inferring them from folders.
4. Initialize only root catalogs. Materialize entity records, standards, profiles, exceptions, environments, deployments, and runbooks on demand from evidence or explicit hypotheses.
5. Compose standards through profiles selected by entity type, capability, or explicit assignment. A narrower contract may add constraints but requires a governed exception or decision to weaken an inherited required rule.
6. Maintain `engineering/quality/quality-system.md` and `quality-system.yaml` together with the quality model, test strategy, and fitness functions. Keep commands canonical in `knowledge/conventions/gates.md`.
7. Declare maturity per area only when its required evidence exists; maturity never implies approval.
8. Link approved `DEC-*` records for governed choices and record candidates for missing structural, data, security, standards, or operational decisions.
9. Update the human contract and every affected mechanical catalog together, applying semantic versioning and compatibility notes.
10. Stop before inventing architecture, editing application code, or creating approval records.

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
- [ ] Repository, application, component, and deployment relations use stable IDs and support monorepo and polyrepo shapes.
- [ ] Standards declare version, obligation level, applicability, verification, evidence, and exception policy.
- [ ] Profiles do not contain cycles and consumers cannot silently weaken inherited required standards.
- [ ] Optional records are materialized only when evidence or explicit hypotheses require them.
- [ ] Detects gaps, conflicts, and dependencies.
- [ ] Records or requests decisions for meaningful changes.
- [ ] Leaves a clear handoff for the next skill or orchestrator.

## Handoff
Next: `technical-discovery` for delivery-specific work.

Pass forward the pinned system id/version, architecture and ownership contracts, standards, quality attributes, decisions, evidence, gaps, and compatibility constraints.
