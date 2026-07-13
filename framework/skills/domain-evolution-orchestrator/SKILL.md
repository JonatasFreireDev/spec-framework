---
name: domain-evolution-orchestrator
description: "Domain Evolution Orchestrator Skill. Use when Codex needs to turn an approved domain and its goals into compared, sliced, and explicitly selected feature candidates before New Feature Orchestrator begins."
---

# Domain Evolution Orchestrator Skill

## Layer
Governance

## Responsibility
Coordinates goals, journeys, opportunities, candidate features, dependencies, impact, delivery slices, and human selection inside one domain. It does not approve a candidate or author specialist-owned artifacts.

## Operating modes
- create: coordinate a new domain evolution cycle.
- update: re-evaluate candidates after evidence or roadmap changes.
- audit: find missing goals, journeys, dependencies, slices, and selection evidence.
- explain: compare candidates and explain the recommended handoff.

## Inputs
Approved domain; goals; journeys; strategy; roadmap; evidence; metrics; decisions; existing features and dependencies.

## Outputs
`evolution/EVOLUTION-NNN/opportunity-map.md`; `candidate-features.md`; `dependency-map.md`; `impact-report.md`; `selection.md`.

## Required reading
- the framework root's `FRAMEWORK.md`
- Relevant parent context.md files.
- Relevant templates in framework/template/.
- Approved product decisions in the active product root's `knowledge/decisions/` and `.product/decisions.json`.

## Discovery and challenge

Follow the shared [Discovery And Challenge contract](../discovery-and-challenge.md) when comparing candidates or requesting human selection.

## Workflow
1. Confirm the Domain and relevant Goals are approved.
2. Route missing user intent to User Goal and missing journeys to Journey.
3. Ask Gap Finder for unmet outcomes and Dependency Analyzer for cross-domain constraints.
4. Ask Feature to draft candidates with explicit delivery slices and non-goals.
5. Ask Impact Analyzer to compare value, cost, risk, reversibility, dependencies, and Delivery Level/Priority.
6. Present candidates without silently selecting one.
7. Record the human selection, rejected/deferred candidates, rationale, and exact feature path.
8. Hand the approved selection to New Feature Orchestrator.

## Quality checklist
- [ ] Preserves traceability to affected artifacts.
- [ ] Uses the correct template and naming conventions.
- [ ] Every candidate serves an approved Goal and has an independently assessable slice.
- [ ] Dependencies, risks, deferred behavior, and rejection rationale are visible.
- [ ] Detects gaps, conflicts, and dependencies.
- [ ] Records or requests decisions for meaningful changes.
- [ ] Leaves a clear handoff for the next skill or orchestrator.

## Handoff
Next: new-feature-orchestrator.

Pass forward the selected feature path, approved slice, evidence, impacts, decisions, dependencies, risks, deferred candidates, and open questions.
