---
name: operations-baseline
description: "Operations Baseline Skill. Use when an agent needs to create, evolve, or audit shared environments, deployments, release and rollback strategies, observability, service levels, continuity, and runbooks."
---

# Operations Baseline Skill

## Layer
Planning

## Responsibility
Own the shared operations catalog and on-demand environment, deployment, runbook, service-level, observability, continuity, release, and rollback contracts. It does not execute deployments, respond to incidents, or own delivery-specific rollout approval.

## Operating modes
- create: establish the first evidence-backed operations baseline.
- update: evolve operations contracts after platform or topology changes.
- adopt: reference authoritative external operational sources.
- audit: find missing runbooks, unsafe recovery assumptions, stale ownership, and unsupported service claims.
- explain: summarize operational capability and gaps.

## Inputs
Technical graph; deployment configuration; infrastructure; runtime and observability evidence; incident knowledge; quality and security constraints; approved decisions.

## Outputs
`engineering/operations/operations.yaml`; environment and deployment records; runbooks; service-level, observability, continuity, release, and rollback contracts.

## Required reading
- [`engineering-systems.md`](../../docs/engineering-systems.md) and [`engineering-catalog-and-standards.md`](../../docs/engineering-catalog-and-standards.md).
- the framework root's `FRAMEWORK.md`
- Relevant parent context.md files.
- This skill owns every generation resource in `assets/`.
- Approved decisions are discovered through the active product root's `.product/decisions.json`.

## Discovery and challenge

Follow the shared [Discovery And Challenge contract](../discovery-and-challenge.md) before material service-level, recovery, data-loss, rollout, or ownership choices.

## Delegated execution

When `engineering-orchestrator` supplies an `engineering-specialist` dispatch envelope through `dispatch-orchestrator`, treat it as a minimal-context subagent assignment. Require the returned Technical Landscape dependency, verify the input hash and write scope, read the graph plus operational evidence, write only under `engineering/operations/`, and return a compact summary, blockers, evidence, decision candidates, and product-relative output hashes to `subagent-return-reviewer`. Do not request or retain the parent conversation.

## Workflow
1. Read the technical graph and inspect deployment, infrastructure, observability, incident, continuity, and runbook evidence.
2. Inventory real environments, deployable units, release paths, rollback mechanisms, dependencies, service objectives, and operational owners.
3. Distinguish observed capability from intended or hypothetical capability; do not claim recovery, monitoring, or rollback without evidence.
4. Update the root operations catalog and create detailed records only for evidenced or explicitly proposed operational surfaces. Keep `operations.yaml` indexed: `ENV-*` and `DEPLOY-*` keys map to relative YAML files and `RUNBOOK-*` keys map to relative Markdown files, never to embedded fields. Environment and deployment records must declare `schema_version`, matching `id`, and non-empty `status`.
5. Link operational deployment records to technical `DEPLOY-*`, application, component, interface, and data-store IDs.
6. Record gaps and decision candidates for ownership, recovery objectives, data loss, observability, continuity, and irreversible rollout risk.
7. Validate safe references and return the baseline to `engineering-orchestrator` without performing operational changes.

## Quality checklist
- [ ] Preserves traceability to affected artifacts.
- [ ] Uses the correct template and naming conventions.
- [ ] Environments, deployments, dependencies, owners, release, and rollback are explicit.
- [ ] Runbooks have triggers, safe actions, validation, escalation, and recovery boundaries.
- [ ] Service-level and continuity claims are supported by evidence or marked as hypotheses.
- [ ] Detects gaps, conflicts, and dependencies.
- [ ] Records or requests decisions for meaningful changes.
- [ ] Leaves a clear handoff for the next skill or orchestrator.

## Handoff
Next: `engineering-orchestrator`.

Pass forward catalogs, environment and deployment records, runbooks, service objectives, evidence, operational gaps, decisions, and risks.
