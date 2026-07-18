---
name: engineering-orchestrator
description: "Engineering Orchestrator Skill. Use when an agent needs to route creation, adoption, evolution, or audit of the shared Engineering System across technical landscape, standards, operations, evidence, aggregation, and human approval."
---

# Engineering Orchestrator Skill

## Layer
Governance

## Responsibility
Coordinate the complete shared engineering baseline across specialist owners, gates, sequencing, and handoffs. It never authors specialist catalogs, changes application code, or grants approval.

## Operating modes
- create: route the first complete engineering baseline.
- update: route affected specialists after evidence or contracts change.
- audit: find missing, stale, conflicting, or unapproved engineering surfaces.
- explain: summarize baseline coverage, blockers, and next ownership.

## Inputs
Product Landscape; declared code roots; product and engineering contexts; current Engineering System; approvals; decisions; specialist findings.

## Outputs
Persisted engineering-baseline handoff; readiness verdict; ordered specialist routes; blockers; human approval checkpoint.

## Required reading
- [`engineering-systems.md`](../../docs/engineering-systems.md) and [`engineering-catalog-and-standards.md`](../../docs/engineering-catalog-and-standards.md).
- [`execution-runtime.md`](../../docs/execution-runtime.md) and [`lifecycle-and-approvals.md`](../../docs/lifecycle-and-approvals.md).
- the framework root's `FRAMEWORK.md`
- Relevant parent context.md files.
- This skill owns `assets/engineering-baseline-handoff-template.json` and `assets/engineering-readiness-report-template.md` for persisted routing and readiness reporting.
- Approved decisions are discovered through the active product root's `.product/decisions.json`; resolve each registered path from its declared decision domain root.

## Discovery and challenge

Follow the shared [Discovery And Challenge contract](../discovery-and-challenge.md) when scope, ownership, compatibility, or a blocking route requires human choice.

## Workflow
1. Revalidate the product manifest, Product Landscape, code roots, current Engineering System hash, approvals, and persisted workspace state.
2. Classify the route as create, update, adopt, or audit and identify affected specialist contracts without assuming one repository, application, or deployment.
3. Route `technical-landscape` first when entity or relation coverage is absent or stale.
4. Route `engineering-standards` after the entity graph can express applicability and inheritance.
5. Route `operations-baseline` after applications, components, deployments, interfaces, and data stores are identifiable.
6. Route `engineering-evidence` to resolve evidence, maturity, coverage, and staleness for every claimed area.
7. Route `engineering-system` only after specialist outputs are internally consistent, so it can consolidate identity, version, quality, compatibility, and the composite contract.
8. Re-run Engineering System validation, report uncovered scope and decision candidates, and stop at the human approval gate.
9. After current approval, route product-wide bootstrap to `domain-architect` or delivery-specific work to `technical-discovery` according to the initiating scope.

## Quality checklist
- [ ] Preserves traceability to affected artifacts.
- [ ] Uses the correct template and naming conventions.
- [ ] Every declared code root and engineering area has an owner or explicit gap.
- [ ] Specialist outputs agree on stable IDs, applicability, maturity, evidence, and compatibility.
- [ ] No specialist is bypassed merely because a previous aggregate contract exists.
- [ ] Approval remains human and applies to the current composite hash.
- [ ] Detects gaps, conflicts, and dependencies.
- [ ] Records or requests decisions for meaningful changes.
- [ ] Leaves a clear handoff for the next skill or orchestrator.

## Handoff
Next: `technical-landscape`.

This is the default first route. For an update or resumed baseline, route only the affected specialist or `engineering-system` consolidation described by the current readiness state. Pass forward scope, code roots, current hashes, completed specialist contracts, evidence, decisions, blockers, and the post-approval route.
