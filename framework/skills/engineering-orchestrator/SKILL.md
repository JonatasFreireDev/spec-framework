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
- Compact delegated returns use the `engineering-specialist-return-template.json` owned by `subagent-return-reviewer`.
- Approved decisions are discovered through the active product root's `.product/decisions.json`; resolve each registered path from its declared decision domain root.

## Discovery and challenge

Follow the shared [Discovery And Challenge contract](../discovery-and-challenge.md) when scope, ownership, compatibility, or a blocking route requires human choice.

## Execution modes

- `sequential` is the compatible default. The current agent follows one specialist contract at a time and persists the same handoff state.
- `delegated` uses harness-native subagents when that capability is available. The orchestrator remains in the parent context and gives each subagent only its dispatch envelope, the pinned specialist skill, the engineering handoff, declared code roots, and files required by that specialist.
- When selecting `delegated`, set `max_parallel` to the supported bounded capacity, normally `2` for the disjoint Standards and Operations phase. Keep it at `1` when concurrency is not safe.
- If delegated execution is requested but the harness cannot create subagents, record the unavailable capability and use the handoff's `fallback`. Never simulate delegation by loading all specialist context into the parent window.
- Route envelope creation and observation through `dispatch-orchestrator`, and return validation through `subagent-return-reviewer`. The CLI persists and validates their runtime contracts but never starts an agent. Agent creation, waiting, interruption, and collection use the current harness's native capability.

## Workflow
1. Revalidate the product manifest, Product Landscape, code roots, current Engineering System hash, approvals, and persisted workspace state.
   Stop when code-root discovery is `cli-fallback` or `needs-agent-review`; route back to Framework Guide for agent-led discovery and manifest correction before assigning specialists.
2. Classify the route as create, update, adopt, or audit and identify affected specialist contracts without assuming one repository, application, or deployment.
3. Persist the execution mode, minimal-context policy, fallback, maximum parallelism, specialist phases, dependencies, and non-overlapping write scopes in the engineering handoff.
4. In delegated mode, route supervised envelope creation to `dispatch-orchestrator` before spawning specialists. Do not inherit the full conversation; pass only the envelope and its `required_reading`. In sequential mode, follow the same phases without creating subagents.
5. Route `technical-landscape` alone in phase 1 when entity or relation coverage is absent or stale.
6. Route `engineering-standards` and `operations-baseline` in phase 2 after the technical graph returns. They may run concurrently because their canonical write scopes do not overlap.
7. Route `engineering-evidence` in phase 3 after standards and operations return, so its coverage and maturity assessment sees their current contracts.
8. Route `engineering-system` in phase 4 only after all specialist outputs are internally consistent, so it can consolidate identity, version, quality, compatibility, and the composite contract.
9. Keep the handoff immutable while a phase has active assignments. Route every compact result through `subagent-return-reviewer`, which verifies the dispatch ID, input hash, dependency returns, write scope, output hashes, evidence, blockers, and decision candidates before the CLI records it. Update handoff state once the phase is collected. Reject stale, escaped, incomplete, or conflicting returns and route only the affected specialist again.
10. Re-run Engineering System validation, report uncovered scope and decision candidates, and stop at the human approval gate.
11. After current approval, route product-wide bootstrap to `domain-architect` or delivery-specific work to `technical-discovery` according to the initiating scope.

## Quality checklist
- [ ] Preserves traceability to affected artifacts.
- [ ] Uses the correct template and naming conventions.
- [ ] Every declared code root and engineering area has an owner or explicit gap.
- [ ] Specialist outputs agree on stable IDs, applicability, maturity, evidence, and compatibility.
- [ ] No specialist is bypassed merely because a previous aggregate contract exists.
- [ ] Delegated specialists receive minimal context and cannot write outside their owned scope.
- [ ] Parallel assignments have returned dependencies, disjoint write scopes, and bounded concurrency.
- [ ] The parent reconciles compact returns rather than importing subagent conversation history.
- [ ] Approval remains human and applies to the current composite hash.
- [ ] Detects gaps, conflicts, and dependencies.
- [ ] Records or requests decisions for meaningful changes.
- [ ] Leaves a clear handoff for the next skill or orchestrator.

## Handoff
Next: `technical-landscape`.

This is the default first route. For an update or resumed baseline, route only the affected specialist or `engineering-system` consolidation described by the current readiness state. Pass forward scope, code roots, current hashes, completed specialist contracts, evidence, decisions, blockers, and the post-approval route.
