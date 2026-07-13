# FDR-038: Guide-First Dispatch

## Snapshot

| Field | Value |
| --- | --- |
| ID | `FDR-038` |
| Status | `proposed` |
| Origin EV | `Governance baseline` |
| Date | `2026-07-13` |
| Owner | `Framework Guide` |

## Context

The framework defines Framework Guide as the conversational entry point, but the installed dispatcher currently resolves any named specialist directly. Adopter instructions also recommend Framework Guide only when the agent already recognizes that routing is unclear. An agent can therefore select a plausible but incorrect specialist before inspecting the workspace, artifact owner, current gate, or persisted handoff.

Always routing through Framework Guide would remove that ambiguity but would add unnecessary work to resumable flows whose route is already established mechanically.

## Decision

The installed dispatcher uses **Guide-first dispatch** for framework-governed product operations.

| Situation | Route |
| --- | --- |
| No verified route | Resolve and follow `framework-guide` before selecting a specialist. |
| Current CLI route | A current-session `guide`, `dashboard`, `status`, or `next` result that names the workspace, concrete feature or use-case scope, current gate, and owner skill may route directly to that skill. |
| Persisted route | A handoff or checkpoint identifies where to resume but is not direct-route evidence by itself. Revalidate it with `dashboard`, `status`, `next`, or `guide`; only the current CLI result may route directly. |
| Explicit human route | A human request that names both the specialist and the concrete artifact or workspace scope may route directly after manifest and scope validation. |
| Incomplete hint | A skill name, keyword, or remembered chat instruction without concrete scope is not a verified route. Use Framework Guide. |

Guide-first dispatch is an agent routing contract, not a new CLI approval gate. It does not block direct diagnostic commands, alter product state, create approval evidence, or require migration of product artifacts. `init` and `upgrade` install or refresh the user-scoped dispatcher for every agent selected in the manifest or command. When a direct route becomes stale, ambiguous, or inconsistent with current mechanical state, the agent returns to Framework Guide.

## Consequences

| Type | Consequence |
| --- | --- |
| Positive | Agents inspect mechanical state before choosing among specialist contracts. |
| Positive | Valid resumptions require only a read-only CLI revalidation, not an unnecessary guidance pass. |
| Negative | Ambiguous first requests require one additional skill resolution and state inspection. |
| Follow-up | Measure whether a future machine-readable routing envelope can safely remove the read-only revalidation step. |

## FRAMEWORK.md Amendments

| Section | Amendment |
| --- | --- |
| `10. Orchestrators` | Define Framework Guide as the default route when no verified specialist route exists. |
| `15. How To Use With Codex` | Define dispatcher routing evidence and direct-route bypass conditions. |

## Related Artifacts

| Artifact | Link |
| --- | --- |
| Conversational CLI guidance | [FDR-022](FDR-022-conversational-cli-guidance.md) |
| Manifest-only activation | [FDR-025](FDR-025-external-runtime-and-manifest-only-activation.md) |
| Framework Guide | [Framework Guide Skill](../skills/framework-guide/SKILL.md) |
| Canonical method | [FRAMEWORK.md](../../FRAMEWORK.md) |

## Supersedes

- N/A.
