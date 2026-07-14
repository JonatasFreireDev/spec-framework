---
name: technical-discovery
description: "Technical Discovery Skill. Use when Codex needs to map an approved specification and design to the existing codebase, architecture, modules, data ownership, tests, and change risks before implementation planning."
---

# Technical Discovery Skill

## Layer
Planning

## Responsibility
Owns delivery-specific `technical-discovery.md`. It maps requirements to existing engineering reality and declares whether an architecture decision is required; it does not make that decision or write implementation code.

## Operating modes
- create: inspect the codebase and create the first discovery report.
- update: refresh mappings after specification, design, or codebase changes.
- audit: detect stale paths, missing requirements, and unapproved architecture changes.
- explain: summarize affected modules, risks, and planning prerequisites.

## Inputs
Approved specification contracts; approved design or `Not applicable`; stable engineering architecture; codebase; tests; conventions; dependencies.

## Outputs
`technical-discovery.md`; requirement-to-module map; probable change surface; architecture gate verdict; risks and plan inputs.

## Required reading
- [`engineering-systems.md`](../../docs/engineering-systems.md) for the stable engineering baseline and migration boundaries.
- the framework root's `FRAMEWORK.md`
- Relevant parent context.md files.
- Relevant templates in framework/template/.
- Approved decisions are discovered through the active product root's `.product/decisions.json`; resolve each registered `path` from its declared domain root (`knowledge/decisions/`, `design/decisions/`, or `engineering/decisions/`).

## Discovery and challenge

Follow the shared [Discovery And Challenge contract](../discovery-and-challenge.md) before substantive creation or material revision.

## Workflow
1. Read every applicable specification contract and Design.
2. Inspect stable architecture under the active product root's `engineering/` and the real code/test tree.
3. Map each `REQ-*` to existing modules, APIs, data owners, tests, conventions, and likely paths.
4. Identify new dependencies, migrations, permission boundaries, shared resources, and rollout constraints.
5. Set the Architecture Gate to `Decision required` with candidates, or `Not required` with concrete rationale. A referenced DEC must exist, be indexed, approved with a current hash-matching approval record, and apply to this use-case scope.
6. Stop when the codebase cannot be inspected or an architecture decision is unresolved.
7. Hand the approved report to Engineering Proposal when the Architecture Gate is resolved.

## Quality checklist
- [ ] Preserves traceability to affected artifacts.
- [ ] Uses the correct template and naming conventions.
- [ ] Every applicable requirement maps to real engineering evidence or an explicit gap.
- [ ] Stable architecture knowledge is linked, not duplicated.
- [ ] Architecture decision requirement is explicit.
- [ ] Detects gaps, conflicts, and dependencies.
- [ ] Records or requests decisions for meaningful changes.
- [ ] Leaves a clear handoff for the next skill or orchestrator.

## Handoff
Next: engineering-proposal after the Architecture Gate is resolved.

Pass forward requirement mappings, code paths, architecture evidence, decisions, dependencies, risks, migrations, tests, and rollout constraints.
