# Agent Instructions

## Purpose

This repository is a Product Engineering Framework. Documentation is the execution infrastructure. Agents must use the framework artifacts, project skills, templates, decisions, and context files before creating downstream work.

## Source Of Truth

Always read `FRAMEWORK.md` first.

Use the repository root as the framework root. The canonical structure shown as `product/` in `FRAMEWORK.md` is represented by this repository root.

## Project Skills

Use project-local skills from `.codex/skills/`.

Each skill owns one step of the framework. Specialist skills create or update canonical artifact content. Orchestrator skills control sequencing, gates, handoffs, and readiness.

When a task maps to a skill:

1. Read the matching `.codex/skills/<skill>/SKILL.md`.
2. Read relevant parent and local `context.md` files.
3. Read the matching template in `knowledge/templates/` when creating or normalizing an artifact.
4. Read approved decisions in `knowledge/decisions/` and `.product/decisions.json` when relevant.
5. Create or update only the artifact owned by the skill unless the user asks for an orchestrated flow.

## Canonical Flow

Follow this sequence:

```text
Problem -> Vision -> Strategy -> Domain -> User Goal -> Feature -> Use Case -> Specification -> Design -> Implementation Plan -> Execution Graph -> Tasks -> Code -> Validation -> Audit
```

For feature delivery, follow:

```text
Feature -> Use Cases -> Specification -> Design -> Implementation Plan -> Execution Graph -> Tasks
```

## Approval Gates

Do not skip gates.

- No downstream work should be generated from an unapproved or incomplete parent artifact unless the output is explicitly marked `draft`.
- No `design.md` should be generated before a Specification exists.
- No `implementation-plan.md` should be generated until `design.md` is approved or explicitly marked `Not applicable`.
- No `tasks.md` should be generated until `execution-graph.json` exists and is valid.
- Do not implement application code until the relevant Specification, Design, Implementation Plan, Execution Graph, and Tasks are approved or the user explicitly asks for a draft/prototype exception.

## Delivery And Priority

Every executable artifact must include:

```yaml
delivery:
  level: L0 | L1 | L2 | L3 | L4 | L5
  priority: P0 | P1 | P2 | P3
  depends_on:
    - artifact-id-or-path
  rationale: Why this level and priority are assigned.
```

Changing Delivery Level or Priority is an approval-gated decision.

## Decisions

Record meaningful decisions in `knowledge/decisions/` using `knowledge/templates/decision-template.md`.

A decision is required when work changes:

- architecture;
- scope;
- security, privacy, permissions, or payment behavior;
- core business rules;
- external dependencies;
- delivery commitments;
- hard-to-reverse strategy.

Approved decisions should also be indexed in `.product/decisions.json`.

## Context Files

Every important object should have a `context.md`.

Use context files to preserve:

- ID, type, name, status, and owner skill;
- parents and children;
- dependencies and related artifacts;
- documents owned by the object;
- Delivery Level and Priority;
- open questions;
- decisions;
- next recommended skill.

## Templates

Use `knowledge/templates/` as the starting structure for new artifacts. Replace placeholders with useful content. Do not create title-only documents.

## Reporting

After a generation or orchestration task, report:

- files created;
- files modified;
- files still incomplete;
- decisions needing human approval;
- recommended next steps;
- validation performed.

For larger documentation bootstraps or audits, save the report under `audits/`.

## Boundaries

Do not change the architecture defined in `FRAMEWORK.md` without asking for approval.

Do not implement application code as part of documentation bootstrap, planning, readiness, or framework maintenance tasks.

Do not invent product scope silently. If a needed detail is missing, mark it as an assumption, open question, or decision candidate.
