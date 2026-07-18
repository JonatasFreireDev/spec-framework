---
name: technical-landscape
description: "Technical Landscape Skill. Use when an agent needs to discover, create, evolve, or audit the shared graph of systems, applications, components, repositories, interfaces, data stores, deployments, and their relations."
---

# Technical Landscape Skill

## Layer
Planning

## Responsibility
Own the shared system context, module map, technical entity catalog, topology, entity records, and boundary assessments under `engineering/`. It does not define standards, operational procedures, delivery-specific changes, or application code.

## Operating modes
- create: establish the first evidence-backed technical graph.
- update: evolve entities and relations while preserving stable IDs.
- adopt: map authoritative external architecture sources into local references.
- audit: find omitted roots, boundaries, ownership, relations, and stale evidence.
- explain: summarize the topology and its gaps.

## Inputs
Product Landscape; declared code roots; repositories; manifests; source and test trees; configuration; deployment evidence; approved decisions.

## Outputs
`engineering/architecture/system-context.md`; `engineering/architecture/modules.md`; `engineering/catalog/catalog.yaml`; `engineering/architecture/topology.yaml`; on-demand entity YAML and Markdown records; technical landscape and boundary assessments; relation findings.

## Required reading
- [`engineering-catalog-and-standards.md`](../../docs/engineering-catalog-and-standards.md).
- the framework root's `FRAMEWORK.md`
- Relevant parent context.md files.
- This skill owns every generation resource in `assets/` and the read-only inventory scripts in `scripts/`.
- Approved decisions are discovered through the active product root's `.product/decisions.json`.

## Discovery and challenge

Follow the shared [Discovery And Challenge contract](../discovery-and-challenge.md) before material boundary, ownership, identity, split, merge, or adoption choices.

## Delegated execution

When `engineering-orchestrator` supplies an `engineering-specialist` dispatch envelope through `dispatch-orchestrator`, treat it as a minimal-context subagent assignment. Verify the role, input hash, dependencies, and write scope; read only the declared inputs plus evidence discovered within the code roots; write only under `engineering/architecture/` and `engineering/catalog/`; and return a compact summary, blockers, evidence, decision candidates, and product-relative output hashes to `subagent-return-reviewer`. Do not request or retain the parent conversation.

## Workflow
1. Run the platform-appropriate technical landscape inventory script and inspect every declared code root, not only the first repository or dominant application.
2. Identify systems, applications, components, repositories, interfaces, data stores, deployments, ownership, capabilities, and observable relations.
3. Distinguish observed evidence from inference and explicit hypotheses; do not convert folder nesting into ownership automatically.
4. Reuse stable IDs when meaning is unchanged. Record move, split, merge, replacement, or compatibility implications when identity changes.
5. Update `catalog.yaml` and `topology.yaml`; create entity records only when evidence or explicit hypotheses justify them.
6. Validate entity identity, type, safe references, relation endpoints, and coverage of every declared code root.
7. Record decision candidates for unresolved boundaries, ownership, data, integration, or repository structure.
8. Return the current graph and gaps to `engineering-orchestrator` without changing aggregate status or approvals.

## Quality checklist
- [ ] Preserves traceability to affected artifacts.
- [ ] Uses the correct template and naming conventions.
- [ ] Covers monorepo, polyrepo, shared-component, and multi-deployment shapes without path-derived ownership.
- [ ] Every entity and relation has a stable ID and evidence or explicit hypothesis.
- [ ] Every declared code root is represented or reported as a gap.
- [ ] Detects gaps, conflicts, and dependencies.
- [ ] Records or requests decisions for meaningful changes.
- [ ] Leaves a clear handoff for the next skill or orchestrator.

## Handoff
Next: `engineering-orchestrator`.

Pass forward catalogs, topology, entity records, source evidence, assumptions, missing coverage, decision candidates, and compatibility risks.
