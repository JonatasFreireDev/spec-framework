---
name: engineering-system
description: "Engineering System Skill. Use when Codex needs to create, adopt, evolve, version, or audit shared product architecture, standards, quality attributes, fitness functions, operations, and engineering evidence in the Spec Framework workflow."
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
- the framework root's `FRAMEWORK.md`
- Relevant parent context.md files.
- Relevant templates in framework/template/.
- Approved product decisions in the active product root's `knowledge/decisions/` and `.product/decisions.json`.

## Workflow
1. Inspect the real product code, tests, configuration, environments, operations evidence, and existing engineering documents.
2. Define the covered product and repository boundaries and choose `generate`, `evolve`, or `adopt`.
3. Inventory modules, ownership, data, integrations, standards, quality attributes, gates, runbooks, and consumers with evidence paths.
4. Declare maturity per area only when its required evidence exists; maturity never implies approval.
5. Link approved `DEC-*` records for governed choices and record candidates for missing structural, data, security, or operational decisions.
6. Update the human contract and mechanical catalog together, applying semantic versioning and compatibility notes.
7. Stop before inventing architecture, editing application code, or creating approval records.

## Quality checklist
- [ ] Preserves traceability to affected artifacts.
- [ ] Uses the correct template and naming conventions.
- [ ] Human and mechanical contracts agree on identity, version, origin, scope, and maturity.
- [ ] Every maturity claim points to concrete evidence.
- [ ] Stable knowledge is not duplicated into delivery proposals.
- [ ] Detects gaps, conflicts, and dependencies.
- [ ] Records or requests decisions for meaningful changes.
- [ ] Leaves a clear handoff for the next skill or orchestrator.

## Handoff
Next: `technical-discovery` for delivery-specific work.

Pass forward the pinned system id/version, architecture and ownership contracts, standards, quality attributes, decisions, evidence, gaps, and compatibility constraints.
