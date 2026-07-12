# FDR-022: Conversational CLI Guidance

## Snapshot

| Field | Value |
| --- | --- |
| ID | `FDR-022` |
| Status | `approved` |
| Origin EV | Human-approved CLI usability evolution |
| Date | `2026-07-12` |
| Owner | `Delivery Orchestrator` |

## Context

The CLI exposes deterministic commands and a consolidated dashboard, but a person still needs to translate goals such as “start a feature”, “use these Figma screens”, or “what should I do next?” into framework terminology and safe command sequences. Command Planner and Command Executor govern approved runtime commands; they are not a conversational entry point.

## Decision

Ship `framework-guide` as an orchestration skill that translates human intent into read-first CLI navigation, explains gates, and routes artifact creation to the existing owner skill.

| Concern | Contract |
| --- | --- |
| Authority | Guidance never grants approval, authors specialist artifacts, or broadens command authority. |
| State | CLI `dashboard`, `guide`, `status`, `review`, readiness, impact, and validation output are mechanical sources of truth. |
| Mutation | Recommend or execute the smallest command that advances one valid gate. |
| Approvals | Require explicit human identity and confirmation through existing approval commands. |
| Runtime | Governed execution still routes through Command Planner and Command Executor. |
| Portability | The skill ships with framework assets to every supported agent target. |

## Consequences

| Type | Consequence |
| --- | --- |
| Positive | New users can operate the framework without memorizing commands or artifact ownership. |
| Positive | Agents share one intent-routing and explanation contract. |
| Negative | Guidance quality depends on current CLI help and dashboard coverage. |
| Mitigation | The skill must inspect current mechanical state and never invent flags or results. |

## Related Artifacts

| Artifact | Link |
| --- | --- |
| Guided migrations and dashboard | [FDR-020](FDR-020-consolidated-cli-and-guided-migrations.md) |
| Runtime authority | [FDR-017](FDR-017-resumable-parallel-runtime.md) |
| Framework Guide | [Framework Guide Skill](../skills/framework-guide/SKILL.md) |

## Supersedes

- N/A
