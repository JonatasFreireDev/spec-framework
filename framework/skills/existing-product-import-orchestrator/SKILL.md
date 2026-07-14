---
name: existing-product-import-orchestrator
description: "Existing Product Import Orchestrator Skill. Use when Codex needs to coordinate multiple existing documents into a reviewed Domain, User Goal, and Feature graph without bypassing framework gates."
---

# Existing Product Import Orchestrator Skill

## Layer
Governance

## Responsibility
Coordinates inventory, classification, reconciliation, approval, draft materialization, and downstream handoff. It does not author candidate content itself or approve product artifacts.

## Operating modes
- create: coordinate a new import run.
- update: coordinate re-analysis after source or mapping changes.
- audit: verify run completeness, conflict resolution, and approval boundaries.
- explain: present the proposed graph and next human decisions.

## Inputs
Source documents; product context; Artifact Importer outputs; existing product graph; approved decisions.

## Outputs
Sequenced import run; approval request; routed draft artifacts; synchronization handoff.

## Required reading
- the framework root's `FRAMEWORK.md`
- Relevant parent context.md files.
- Relevant templates in framework/template/.
- Approved decisions are discovered through the active product root's `.product/decisions.json`; resolve each registered `path` from its declared domain root (`knowledge/decisions/`, `design/decisions/`, or `engineering/decisions/`).

## Workflow
1. Ask Artifact Importer to inventory and classify all sources.
2. Review candidates against the existing product graph.
3. Route terminology and contradiction findings through Conflict Finder and Product Historian when decisions are required.
4. Present the proposed Domain → User Goal → Feature graph and unresolved questions.
5. Stop at the materialization gate until the user explicitly approves selected mappings.
6. Ask Artifact Importer to create approved targets as `draft` only.
7. Route missing foundation context to Product Orchestrator and ready feature candidates to New Feature Orchestrator.
8. Ask Documentation Orchestrator to synchronize contexts, indexes, mappings, and reports.
9. Treat the latest run declared in `.product/framework.json` as the active workspace gate; never use an older materialized run to satisfy a newer import.

## Quality checklist
- [ ] Preserves traceability to affected artifacts.
- [ ] Uses the correct template and naming conventions.
- [ ] Keeps analysis separate from materialization.
- [ ] Does not treat source prose as an approval record.
- [ ] Distinguishes approval to materialize selected drafts from approval of the resulting product artifacts.
- [ ] Detects gaps, conflicts, and dependencies.
- [ ] Records or requests decisions for meaningful changes.
- [ ] Leaves a clear handoff for the next skill or orchestrator.

## Handoff
Next: Product Orchestrator, New Feature Orchestrator, or Documentation Orchestrator according to the approved mapping.

Pass forward approved selections, draft artifacts, findings, open questions, decisions, dependencies, risks, and required follow-up work.
