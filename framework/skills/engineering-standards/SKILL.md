---
name: engineering-standards
description: "Engineering Standards Skill. Use when an agent needs to create, evolve, resolve, or audit versioned engineering standards, profiles, applicability, verifiable rules, and governed exceptions."
---

# Engineering Standards Skill

## Layer
Planning

## Responsibility
Own the standards catalog, profiles, individual standards, effective resolutions, conformance reports, and `STDEX-*` exceptions. It does not own the technical graph, Quality System, product decisions, or delivery approval.

## Operating modes
- create: establish the first evidence-backed standards catalog and profiles.
- update: version rules and profiles while declaring compatibility impact.
- adopt: register authoritative external standards with pinned versions.
- audit: find cycles, silent weakening, unverifiable rules, gaps, and expired exceptions.
- explain: summarize effective standards for a scope.

## Inputs
Technical entity graph; engineering conventions; code and CI evidence; quality and security contracts; approved decisions; compatibility constraints.

## Outputs
`engineering/standards/standards.yaml`; `PROFILE-*`; `STD-*` YAML and Markdown; `STDEX-*`; profile resolutions; conformance and gap reports.

## Required reading
- [`engineering-systems.md`](../../docs/engineering-systems.md) and [`engineering-catalog-and-standards.md`](../../docs/engineering-catalog-and-standards.md).
- the framework root's `FRAMEWORK.md`
- Relevant parent context.md files.
- This skill owns every generation resource in `assets/`.
- Approved decisions are discovered through the active product root's `.product/decisions.json`.

## Discovery and challenge

Follow the shared [Discovery And Challenge contract](../discovery-and-challenge.md) before selecting obligation levels, compatibility policy, profile inheritance, or exception boundaries.

## Delegated execution

When `engineering-orchestrator` supplies an `engineering-specialist` dispatch envelope through `dispatch-orchestrator`, treat it as a minimal-context subagent assignment. Require the returned Technical Landscape dependency, verify the input hash and write scope, read the graph plus standards evidence, write only under `engineering/standards/`, and return a compact summary, blockers, evidence, decision candidates, and product-relative output hashes to `subagent-return-reviewer`. Do not request or retain the parent conversation.

## Workflow
1. Read the technical graph, existing conventions, quality and security contracts, and real enforcement evidence.
2. Define standards only for stable technical rules with an explicit applicability boundary and verification method.
3. Assign semantic versions, categories, obligation levels, individually identifiable rules, compatibility notes, and exception policy.
4. Compose profiles by entity type, capability, or explicit assignment; reject inheritance cycles and silent weakening of required rules.
5. Resolve effective standards for affected entities and record uncovered or conflicting applicability.
6. Create `STDEX-*` only with exact scope, owner, rationale, residual risk, mitigation, expiry, re-entry gate, and status; never self-approve it.
7. Validate all references, versions, rules, profiles, and exception expiry.
8. Return standards coverage and gaps to `engineering-orchestrator`.

## Quality checklist
- [ ] Preserves traceability to affected artifacts.
- [ ] Uses the correct template and naming conventions.
- [ ] Rules are verifiable and do not duplicate Quality System or security policy.
- [ ] Profiles are acyclic, version-pinned, and cannot silently weaken inherited requirements.
- [ ] Exceptions are scoped, temporary, risk-explicit, and approval-bound.
- [ ] Detects gaps, conflicts, and dependencies.
- [ ] Records or requests decisions for meaningful changes.
- [ ] Leaves a clear handoff for the next skill or orchestrator.

## Handoff
Next: `engineering-orchestrator`.

Pass forward standards, profiles, resolved applicability, exceptions, verification evidence, compatibility impacts, and blockers.
