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
- A downstream artifact must not advance to `proposed`, `approved`, or a later status while its required parent gate is incomplete.
- No `design.md` should be generated before a Specification exists.
- No `implementation-plan.md` should be generated until `design.md` is approved or explicitly marked `Not applicable`.
- No task file or generated `tasks.md` index should be generated until `execution-graph.json` exists and is valid.
- Do not implement application code until the relevant Specification, Design, Implementation Plan, Execution Graph, and Tasks are approved or the user explicitly asks for a draft/prototype exception.
- Artifacts with status `approved`, `in_progress`, `implemented`, `validated`, or `released` must have a matching approval record in `.product/history/`.
- Agents must not create, edit, or repair approval records unless the user explicitly approves a migration that names approval-record generation as a deliverable. If approval records are missing or inconsistent, report the blocker and stop.
- Staleness is derived by the validator from `.product/derivations.json`; agents must not set `stale` as an artifact status. If a derived artifact is stale, report it and regenerate or request re-approval through the appropriate approved flow.
- A task must not move to `implemented` unless its task file contains structured `Branch`, `Commits`, and `Code paths`.
- A task must not move to `validated` or `released` unless its task file contains `PR`, passing `Test status`, and concrete evidence such as gate logs, CI URL, screenshots, or QA evidence.

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

## Tasks

Tasks are canonical one-file artifacts under `tasks/<task-id>.md`.

- Edit task status, contract, delivery metadata, code links, and evidence only in the task file.
- Treat `tasks.md` as a generated index from `execution-graph.json` and `tasks/*.md`; do not edit it by hand except when regenerating the index.
- Every `execution-graph.json` node must include `path` pointing to its task file.
- If graph snapshot fields such as `title` or `type` disagree with the task file, stop and reconcile the source artifact before continuing.

## Rigor Tiers

Every use case context must declare `rigor_tier: S | M | L | N/A`.

- Tier S requires specification, tasks, and tests; design, analytics, and audit may be `Not applicable`.
- Tier M adds design, implementation plan, and execution graph.
- Tier L adds analytics, audit, QA evidence, and Security Review.
- Auth, permissions, roles, payments, PII, uploads, UGC, public surfaces, or migrations touching RLS/policies force Tier L.
- Do not lower or change a tier without a matching approval record for the use case.

## Identity And Moves

Every important folder-backed artifact must declare `slug` in `context.md`.

- The slug must match the folder name and remain stable when the human-readable title changes.
- Do not allocate new IDs by incrementing `.product/ids.json`; IDs are scoped by parent and references should include paths when ambiguity is possible.
- Use `node engineering/move-artifact.mjs --from <old-path> --to <new-path>` when moving an artifact folder or file.
- The move tool rewrites Markdown links and JSON paths; review its reported free-text mentions manually.

## Code Links And Evidence

The framework uses a monorepo delivery convention: product documentation and product code live in the same product repository. This `spec-framework` repository remains a template/lab.

- Use repository-relative paths for code paths whenever the code lives in the same repository.
- Keep branch, commits, PR, code paths, gate logs, CI URL, screenshots, and QA evidence in structured task or QA evidence fields.
- Do not mark task code as `implemented` with only prose notes.
- Do not mark task code as `validated` without passing test evidence and a PR/equivalent review surface.

## Reporting

After a generation or orchestration task, report:

- files created;
- files modified;
- files still incomplete;
- decisions needing human approval;
- recommended next steps;
- validation performed.

For larger documentation bootstraps or audits, save the report under `audits/`.

Reports should be visual and scannable:

- use status icons such as `✅`, `🟡`, `🔴`, and `➖`;
- use summary tables for status, files, findings, decisions, and next steps;
- use Mermaid diagrams for flows, gates, dependencies, and artifact chains;
- keep prose concise and use tables for comparison or review surfaces;
- include a final result section with verdict, blockers, and next owner.

Mermaid flow diagrams should show current progress when they represent a framework sequence:

- use `done` for approved or completed prior steps;
- use `current` for the step being created, reviewed, or executed now;
- use `pending` for future steps;
- use `blocked` for steps stopped by missing decisions, gaps, conflicts, or dependencies.

The skill that changes an artifact status must update the related `context.md` and directly related Mermaid flow. The Documentation Orchestrator owns final synchronization across reports, templates, indexes, and context files. The Audit Orchestrator verifies that Mermaid visual state matches artifact status during audits.

## Boundaries

Do not change the architecture defined in `FRAMEWORK.md` without asking for approval.

Do not implement application code as part of documentation bootstrap, planning, readiness, or framework maintenance tasks.

Do not invent product scope silently. If a needed detail is missing, mark it as an assumption, open question, or decision candidate.
