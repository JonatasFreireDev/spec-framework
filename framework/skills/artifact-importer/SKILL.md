---
name: artifact-importer
description: "Artifact Importer Skill. Use when Codex needs to inventory and normalize existing product documents into an approved, traceable import plan without silently creating canonical product truth."
---

# Artifact Importer Skill

## Layer
Discovery

## Responsibility
Owns one import run under the active product root's `knowledge/imports/runs/`. It inventories sources, proposes candidates and mappings, and reports conflicts; it does not approve or silently overwrite canonical product artifacts.

## Operating modes
- create: inventory sources and create the first import plan.
- update: re-run analysis and mark mappings affected by changed source hashes.
- audit: find missing sources, duplicate candidates, conflicts, and stale mappings.
- explain: summarize the import plan and its unresolved decisions.

## Inputs
Source documents; product context; existing Domains, User Goals, Features, glossary, rules, and decisions.

## Outputs
`inventory.json`; `import-plan.json`; `mapping.json`; `traceability.json`; `conflicts.md`; `import-report.md` in one canonical import run.

## Required reading
- the framework root's `FRAMEWORK.md`
- Relevant parent context.md files.
- Relevant templates in framework/template/.
- `framework/template/import-traceability-template.json` when normalizing the traceability ledger.
- Approved decisions are discovered through the active product root's `.product/decisions.json`; resolve each registered `path` from its declared domain root (`knowledge/decisions/`, `design/decisions/`, or `engineering/decisions/`).

## Workflow
1. Copy or reference sources without modifying their content and compute a SHA-256 hash for each source.
2. Inventory source path, format, size, and hash.
3. Read each imported source and update `traceability.json` with review status, section-level evidence, extracted claims, candidate ids, mapped targets, and gaps. Every source must end as `reviewed`, `partially_mapped`, `mapped`, or `not_applicable`; never leave the reason implicit.
4. Extract candidate Domains, User Goals, Features, rules, decisions, priorities, and dependencies with section-level evidence.
5. Compare candidates with existing artifacts and the glossary.
6. Record duplicates, contradictions, ambiguous parents, and open questions; never resolve them silently.
7. Propose source-to-artifact mappings in `draft` and leave `materialization_approved` false.
8. Stop for explicit human approval before creating canonical product artifacts.
9. When approved, materialize only selected mappings as `draft`, preserving `source_documents` traceability and never creating approval records.
10. Use `spec-framework import materialize --run IMPORT-NNN --approved-by <human> --yes` for mechanical materialization; do not edit the approval fields manually.

## Quality checklist
- [ ] Preserves traceability to affected artifacts.
- [ ] Uses the correct template and naming conventions.
- [ ] Every source has a stable path and SHA-256 hash.
- [ ] Every candidate cites at least one source location.
- [ ] Every imported source has one traceability entry and an explicit coverage status.
- [ ] Unmapped content is recorded as a gap or open question, never silently dropped.
- [ ] Existing artifacts are proposed as merge targets, never overwritten.
- [ ] Detects gaps, conflicts, and dependencies.
- [ ] Records or requests decisions for meaningful changes.
- [ ] Leaves a clear handoff for the next skill or orchestrator.

## Handoff
Next: Existing Product Import Orchestrator.

Pass forward the complete import run, source hashes, candidates, mappings, conflicts, open questions, and the explicit materialization decision.
