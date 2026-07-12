# Codex Skills

## Purpose

This folder contains repository-local skills that operationalize the framework. Each child folder is a Codex skill with a required `SKILL.md`.

## When To Use

Use these skills when creating, updating, auditing, explaining, routing, or handing off framework artifacts. Specialist skills own canonical artifact content. Orchestrator skills own flow, gates, sequencing, and handoff.

## Path Resolution

Skills are written to work in both this framework repository and adopter repositories.

| Path kind | In this repository | In adopter repositories |
| --- | --- | --- |
| Framework method assets | `framework/skills/`, `framework/template/`, `framework/decisions/`, `framework/validators/`, `framework/tools/`, `framework/tests/` | Versioned external user cache resolved from `product/.product/framework.json` |
| Active product artifacts | `examples/events/` | `product/` |
| Product decisions and state | `examples/events/knowledge/decisions/`, `examples/events/.product/` | `product/knowledge/decisions/`, `product/.product/` |
| Product gates and conventions | `examples/events/knowledge/conventions/` | `product/knowledge/conventions/` |

When a skill mentions a product-relative path such as `knowledge/conventions/gates.md`, `.product/decisions.json`, `domains/`, `audits/`, or `releases/`, resolve it under the active product root. When it mentions framework assets such as `FRAMEWORK.md`, `templates/`, `validators/`, or FDRs, resolve it under the framework root.

## Expected Files

- `<skill-name>/SKILL.md`: one skill per folder.
- Specialist skills: problem, vision, strategy, domain, goal, journey, feature, use case, Design System, UX/UI, specification, technical discovery, implementation planning, graph, task, code runner, bug fixer, QA, code review, security review, threat modeler, commit crafter, PR finalizer, audit, documentation, history, and artifact import.
- Orchestrator skills: product, domain evolution, existing product import, new feature, audit, evolution, documentation, release, delivery, execution scheduling, and integration.
- Guidance skill: Framework Guide translates human goals into current CLI state, the smallest safe command, and the correct specialist or approval handoff without authoring artifacts.

Runtime v2 also includes the `command-planner` and `command-executor` operational skills. The planner owns immutable argv-based plans; the executor is restricted to local R0/R1 plans.

Use `framework-guide` as the conversational entry point when the person does not know which command, workspace, artifact owner, or gate applies. It must read CLI state first and route runtime execution back through Command Planner and Command Executor.

Delivery Orchestrator reports through the consolidated dashboard model. Product Historian owns review of guided legacy decision migration; the migration tool only updates `.product/decisions.json` metadata and never edits DEC content or approvals.

## Responsible Skill

Primary owner: Documentation Orchestrator.

## Next Step

Read the smallest skill that owns the next artifact. If a request spans multiple artifacts, start with the matching orchestrator and stop at approval gates.

All `Next:` handoffs use canonical skill folder names. Numbered legacy filenames are invalid and the validator reports them.
