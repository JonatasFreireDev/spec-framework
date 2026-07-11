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
- the framework root's `FRAMEWORK.md`
- Relevant parent context.md files.
- Relevant templates in framework/template/.
- Approved product decisions in the active product root's `knowledge/decisions/` and `.product/decisions.json`.

## Workflow
1. Read every applicable specification contract and Design.
2. Inspect stable architecture under the active product root's `engineering/` and the real code/test tree.
3. Map each `REQ-*` to existing modules, APIs, data owners, tests, conventions, and likely paths.
4. Identify new dependencies, migrations, permission boundaries, shared resources, and rollout constraints.
5. Set the Architecture Gate to `Decision required` with candidates, or `Not required` with concrete rationale.
6. Stop when the codebase cannot be inspected or an architecture decision is unresolved.
7. Hand the approved report to Implementation Planner.

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
Next: implementation-planner after the Architecture Gate is resolved.

Pass forward requirement mappings, code paths, architecture evidence, decisions, dependencies, risks, migrations, tests, and rollout constraints.
