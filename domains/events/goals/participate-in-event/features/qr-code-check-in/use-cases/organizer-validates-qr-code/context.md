# Context: Organizer Validates QR Code

```yaml
id: UC-002
type: use_case
name: Organizer Validates QR Code
status: proposed
owner_skill: use-case
slug: organizer-validates-qr-code
rigor_tier: L
parents:
  - FT-001
children:
  - SPEC-002
  - DES-002
  - PLAN-002
  - GRAPH-002
  - TASKSET-002
  - TEST-002
  - QA-002
  - SEC-002
  - ANA-002
  - AUD-002
depends_on:
  - UC-001
used_by:
  - RELEASE-001
related:
  - domains/events/goals/participate-in-event/features/qr-code-check-in/use-cases/attendee-checks-in-with-qr-code/context.md
documents:
  canonical: use-case.md
  specification: specification.md
  design: design.md
  implementation_plan: implementation-plan.md
  execution_graph: execution-graph.json
  tasks: tasks.md
  tests: tests.md
  qa_evidence: qa-evidence.md
  security_review: security-review.md
  analytics: analytics.md
  audit: audit.md
delivery:
  level: L1
  priority: P0
  depends_on:
    - DOMAIN-users
    - DOMAIN-events
    - DEC-001
    - DEC-002
  rationale: Organizer validation closes the walking skeleton for event attendance by turning attendee QR proof into server-authoritative check-in.
open_questions:
  - Should organizers be able to scan when temporarily offline, or must validation always be online?
  - Which organizer roles can validate check-in for an event?
decisions:
  - DEC-001
  - DEC-002
  - DEC-008
```

## Purpose

This context gives agents the local map for the organizer validation use case. It links the behavior, implementation contract, UX, planning, task graph, tests, analytics, and audit evidence.

## Rigor Tier

| Field | Value |
| --- | --- |
| Tier | L |
| Trigger checklist | auth, permissions, roles, token privacy, check-in state mutation |
| Human approval | Approved by EV-003 policy rollout |
| Rationale | Organizer validation authorizes a server-side attendance write and handles sensitive event attendance signals. |

## Required Reading

| Artifact | Link |
| --- | --- |
| Framework | [FRAMEWORK.md](../../../../../../../../FRAMEWORK.md) |
| Domain context | [domains/events/context.md](../../../../../../context.md) |
| Goal context | [participate-in-event/context.md](../../../../context.md) |
| Feature context | [qr-code-check-in/context.md](../../context.md) |
| DEC-001 QR expiration duration | [DEC-001](../../../../../../../../knowledge/decisions/DEC-001-qr-expiration-duration.md) |
| DEC-002 QR token strategy | [DEC-002](../../../../../../../../knowledge/decisions/DEC-002-qr-token-strategy.md) |

## Local Documents

| Document | Link |
| --- | --- |
| Use Case | [use-case.md](use-case.md) |
| Specification | [specification.md](specification.md) |
| Design | [design.md](design.md) |
| Implementation Plan | [implementation-plan.md](implementation-plan.md) |
| Execution Graph | [execution-graph.json](execution-graph.json) |
| Tasks | [tasks.md](tasks.md) |
| Tests | [tests.md](tests.md) |
| QA Evidence | [qa-evidence.md](qa-evidence.md) |
| Security Review | [security-review.md](security-review.md) |
| Analytics | [analytics.md](analytics.md) |
| Audit | [audit.md](audit.md) |

## Handoff

Next recommended skill: Specification AI.

Do not generate application code from this folder. Generate implementation tasks only after Specification and Design approval.
