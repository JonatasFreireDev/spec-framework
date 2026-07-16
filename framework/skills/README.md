# Codex Skills

## Purpose

This folder contains repository-local skills that operationalize the framework. Each child folder is a Codex skill with a required `SKILL.md`.

## When To Use

Use these skills when creating, updating, auditing, explaining, routing, or handing off framework artifacts. Specialist skills own canonical artifact content. Orchestrator skills own flow, gates, sequencing, and handoff.

Definition and planning skills must follow [Discovery And Challenge](discovery-and-challenge.md). They inspect evidence before asking, use the harness-native structured question tool when available, compare meaningful options, recommend a path, warn about material risks, and stop on unanswered blocking questions.

## Path Resolution

Skills are written to work in both this framework repository and adopter repositories.

| Path kind | In this repository | In adopter repositories |
| --- | --- | --- |
| Framework method assets | `framework/skills/`, `framework/template/`, `framework/validators/`, `framework/tools/`, `framework/tests/` | Versioned external user cache resolved from `product/.product/framework.json` |
| Active product artifacts | `examples/events/` | `product/` |
| Product decisions and state | `examples/events/.product/` and its indexed decision roots | `product/.product/` and its indexed decision roots |
| Product gates and conventions | `examples/events/knowledge/conventions/` | `product/knowledge/conventions/` |

When a skill mentions a product-relative path such as `knowledge/conventions/gates.md`, `.product/decisions.json`, `domains/`, `audits/`, or `releases/`, resolve it under the active product root. When it mentions framework assets such as `FRAMEWORK.md`, `templates/`, or `validators/`, resolve it under the framework root.

## Expected Files

- `<skill-name>/SKILL.md`: one skill per folder.
- Specialist skills: problem, vision, strategy, domain, goal, journey, feature, use case, Design System, UX/UI, specification, Engineering System, technical discovery, engineering proposal, engineering review, implementation planning, graph, task, code runner, bug fixer, QA, code review, security review, threat modeler, commit crafter, PR finalizer, audit, documentation, history, artifact import, and subagent return review.
- Orchestrator skills: product, domain evolution, existing product import, new feature, audit, evolution, documentation, release, delivery, execution scheduling, integration, and dispatch.
- Guidance skill: Framework Guide translates human goals into current CLI state, the smallest safe command, and the correct specialist or approval handoff without authoring artifacts.

Runtime v2 also includes the `command-planner` and `command-executor` operational skills. The planner owns immutable argv-based plans; the executor is restricted to local R0/R1 plans.

## Shared Runtime Contracts

Skills and orchestrators must read the shared contract that matches their responsibility:

- [`execution-runtime.md`](../../docs/execution-runtime.md): workspaces, leases, graph scheduling, command plans, and integration.
- [`engineering-systems.md`](../../docs/engineering-systems.md): Engineering System and Engineering Quality System versioning, migration, evidence, and approval boundaries.
- [`lifecycle-and-approvals.md`](../../docs/lifecycle-and-approvals.md): lifecycle states, approval records, staleness, authority, and failure routing.

Use `framework-guide` as the default conversational entry point when no verified specialist route exists. A direct specialist route requires current CLI guidance or an explicit human request that names both the specialist and concrete scope. Persisted handoffs/checkpoints must first be revalidated with `dashboard`, `status`, `next`, or `guide`; a skill name without scope is only a hint. Framework Guide activates product operations only from a valid `product/.product/framework.json`, routes bootstrap/init before activation when explicitly requested, resolves pinned specialist contracts with `spec-framework skill path`, reads CLI state first, and routes governed runtime execution back through Command Planner and Command Executor.

Delivery Orchestrator reports through the consolidated dashboard model. Product Historian owns review of guided legacy decision migration; the migration tool only updates `.product/decisions.json` metadata and never edits DEC content or approvals.

Vision AI owns a three-artifact package with exclusive content boundaries: `vision.md` contains direction and strategic boundaries, `principles.md` contains product decision rules and trade-offs, and `north-star.md` contains the outcome metric and guardrails. The artifacts link to one another instead of duplicating canonical content.

Domain Architect AI reads the pinned runtime's `examples/events/` before first-domain modeling. That worked example is the canonical reference for business-area boundaries, non-ownership, cross-domain contracts, and the first Domain -> User Goal -> Feature -> Use Case walking skeleton.

## Responsible Skill

Primary owner: Documentation Orchestrator.

## Next Step

Read the smallest skill that owns the next artifact. If a request spans multiple artifacts, start with the matching orchestrator and stop at approval gates.

All `Next:` handoffs use canonical skill folder names. Numbered legacy filenames are invalid and the validator reports them.
