---
name: design-system
description: "Design System Skill. Use when an agent needs to create, adopt, evolve, document, or audit the canonical product Design System, including foundations, tokens, themes, components, patterns, sources, versioning, compatibility, and handoff to use-case Design."
---

# Design System Skill

## Layer
Design

## Responsibility
Own the canonical product Design System under `design/system/`. Create or evolve its human contract and structured tokens; do not approve it, silently change product behavior, or implement application components.

## Operating modes
- create: produce the first draft Design System from approved product foundations.
- update: revise an existing system while preserving approved decisions and compatibility policy.
- audit: find broken tokens, inconsistent components, missing states, accessibility gaps, and affected consumers.
- explain: summarize foundations, versions, sources, components, patterns, and migration consequences.

## Inputs
Approved Vision and Strategy; personas; brand sources; existing interfaces or component libraries; visual source manifests; accessibility conventions; approved product decisions; consuming Designs.

## Outputs
`design/system/context.md`; `design/system/design-system.md`; foundations; tokens and themes; component and pattern contracts; source references; compatibility findings; decision candidates; handoff to UX/UI.

## Required reading
- the framework root's `FRAMEWORK.md`
- Relevant product foundation and design context files.
- This skill owns its generation resources: `assets/design-system-template.md`, `assets/design-component-template.md`, and `assets/design-pattern-template.md`.
- Approved decisions are discovered through the active product root's `.product/decisions.json`; resolve each registered `path` from its declared domain root (`knowledge/decisions/`, `design/decisions/`, or `engineering/decisions/`).

## Discovery and challenge

Follow the shared [Discovery And Challenge contract](../discovery-and-challenge.md) before substantive creation or material revision.


Use scripts/invoke-cli.ps1 on Windows or scripts/invoke-cli.sh on macOS/Linux for the CLI operation in this skill's reviewed scope. The wrapper never adds --yes or an approver identity.

Use `scripts/inventory-design-evidence.ps1` on Windows or `scripts/inventory-design-evidence.sh` on macOS/Linux before authoring the shared baseline. It inventories interface and design-system evidence and can call `design-system validate`; it never changes product or code files.

## Workflow
1. Confirm whether origin is `generate`, `evolve`, or `adopt`; never promote a reference to canonical authority implicitly. When there is no interface source, establish explicit accessibility and experience hypotheses before Specifications rather than omitting the shared baseline.
2. Inventory foundations, sources, tokens, themes, components, patterns, consumers, and implementation links.
3. Create or update `context.md` and `design-system.md` using canonical templates.
4. Maintain tool-independent tokens in primitive, semantic, and component layers; resolve aliases and reject cycles.
5. Document every component's anatomy, variants, states, accessibility, tokens, content, responsive behavior, and lifecycle.
6. Record compatibility and migration consequences before removing or changing approved tokens, components, patterns, or external dependencies.
7. Route product, architecture, accessibility, or dependency changes to a human decision gate.
8. Run Design System validation and independent UX Review; never create approval records.
9. When implementation roots exist, inspect design evidence across every agent-confirmed root, including tokens, themes, components, assets, typography, interaction patterns, accessibility conventions, and duplicated or divergent systems. Do not derive a shared Design System from CLI fallback candidates.
9. Hand the approved id/version and supported tokens/components/patterns to UX/UI.

## Quality checklist
- [ ] Preserves traceability to sources, decisions, consumers, and implementation links.
- [ ] Uses the correct templates, stable identities, semantic version, and naming conventions.
- [ ] Token aliases resolve without cycles and themes preserve semantic contracts.
- [ ] Components cover required states, accessibility, content, and responsive behavior.
- [ ] Breaking changes include impact, migration, deprecation, and rollback notes.
- [ ] Detects gaps, conflicts, dependencies, and stale consumers.
- [ ] Records or requests decisions for meaningful changes.
- [ ] Leaves a clear handoff for the next skill or orchestrator.

## Handoff
Next: ux-ui.

Pass forward Design System id/version, status, source versions, foundations, supported tokens/components/patterns, deviations, decisions, compatibility risks, consumers, and required follow-up work.
