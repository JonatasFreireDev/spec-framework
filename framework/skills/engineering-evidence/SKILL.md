---
name: engineering-evidence
description: "Engineering Evidence Skill. Use when an agent needs to inventory, normalize, assess, or audit engineering evidence, maturity claims, coverage, and staleness across the shared Engineering System."
---

# Engineering Evidence Skill

## Layer
Validation

## Responsibility
Own the engineering evidence inventory, evidence records, coverage maps, maturity assessments, verification-run references, and gap or staleness reports. It does not fabricate evidence, execute application changes, or approve maturity.

## Operating modes
- create: establish the first evidence inventory and baseline assessments.
- update: refresh evidence references and affected maturity findings.
- audit: detect missing, unresolved, stale, duplicated, or scope-mismatched evidence.
- explain: summarize what current evidence proves and does not prove.

## Inputs
Technical graph; standards; operations; Quality System; code and tests; CI and runtime references; decisions; current maturity claims.

## Outputs
`engineering/evidence/inventory.md`; optional mechanical catalog and `EVID-*` records; coverage map; maturity assessments; gap and staleness reports.

## Required reading
- [`engineering-systems.md`](../../docs/engineering-systems.md) and [`engineering-catalog-and-standards.md`](../../docs/engineering-catalog-and-standards.md).
- the framework root's `FRAMEWORK.md`
- Relevant parent context.md files.
- This skill owns every generation resource in `assets/` and the read-only evidence inventory scripts in `scripts/`.
- Approved decisions are discovered through the active product root's `.product/decisions.json`.

## Discovery and challenge

Follow the shared [Discovery And Challenge contract](../discovery-and-challenge.md) when evidence limitations, freshness policy, maturity interpretation, or ownership requires a material choice.

## Workflow
1. Run the platform-appropriate evidence inventory script and inspect referenced local, CI, runtime, and approved external sources.
2. Normalize each material source into a stable `EVID-*` reference without copying volatile logs into canonical contracts.
3. Map evidence to technical entities, standards, operations, quality areas, decisions, and maturity claims.
4. Evaluate source existence, observation time, freshness policy, scope, and what the evidence actually proves.
5. Mark unsupported or stale claims as findings; never promote maturity or synthesize approval from evidence volume.
6. Update coverage, maturity, gap, and staleness reports and identify the exact owner required to resolve each gap.
7. Run Engineering System validation when requested and return findings to `engineering-orchestrator`.

## Quality checklist
- [ ] Preserves traceability to affected artifacts.
- [ ] Uses the correct template and naming conventions.
- [ ] Evidence references are safe, resolvable, scoped, and freshness-aware.
- [ ] Observed evidence, inference, and hypothesis remain distinguishable.
- [ ] Every maturity above baseline has sufficient current evidence.
- [ ] Volatile outputs are referenced rather than copied into stable contracts.
- [ ] Detects gaps, conflicts, and dependencies.
- [ ] Records or requests decisions for meaningful changes.
- [ ] Leaves a clear handoff for the next skill or orchestrator.

## Handoff
Next: `engineering-orchestrator`.

Pass forward inventory, coverage, maturity findings, stale or missing evidence, source limitations, blockers, and required owners.
