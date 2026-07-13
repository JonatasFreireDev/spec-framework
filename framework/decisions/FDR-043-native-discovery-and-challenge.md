# FDR-043: Native Discovery And Challenge

## Snapshot

| Field | Value |
| --- | --- |
| ID | `FDR-043` |
| Status | `proposed` |
| Origin EV | `Planning skill interaction feedback` |
| Date | `2026-07-13` |
| Owner | `Framework Maintainer` |

## Context

Definition and planning skills require agents to identify gaps and open questions, but their workflows can proceed directly from inspection to drafting. They do not require the harness-native question mechanism, compared alternatives, recommendations, or proactive risk warnings. Agents can therefore fill templates from assumptions without eliciting the human choices that materially shape the artifact.

## Decision

| Boundary | Contract |
| --- | --- |
| Shared contract | Framework-owned definition and planning skills reference `framework/skills/discovery-and-challenge.md`. |
| Native capability | `native_user_question` is the harness-neutral capability. Each dispatcher maps it to the default structured question tool exposed by Codex, Claude Code, or Cursor. |
| Discovery | `create` requires discovery before a substantive draft. `update` and `evolve` require it when scope, behavior, risk, priority, or an approved decision may change. |
| Question quality | Ask one to three focused questions per round; meaningful choices include two or three options, trade-offs, a recommendation, and a free-form path. |
| Challenge | Skills proactively warn about material scope, dependency, usability, security, operability, reversibility, approval, and delivery risks and propose a safer alternative. |
| Evidence first | Skills inspect repository and CLI evidence before asking and never ask humans for safely discoverable facts. |
| Blocking | A skill cannot finalize or hand off a substantive artifact while a blocking question remains unanswered. Conversational answers do not grant formal approval. |
| Persistence | Store resolved choices, assumptions, alternatives, warnings, and open questions in the artifact or context; never persist raw conversation transcripts. |
| Fallback | When no structured question tool is exposed, ask explicitly in conversation and record the fallback in the handoff. |
| Validation | The validator reports an error when a governed skill omits the shared contract reference or its `Discovery and challenge` section. |

## Consequences

| Type | Consequence |
| --- | --- |
| Positive | Planning becomes collaborative and surfaces human choices before artifacts harden assumptions. |
| Positive | Harness-specific tool names remain outside canonical specialist contracts. |
| Positive | Skills consistently propose alternatives, recommendations, and material risk warnings. |
| Negative | Creation and material updates may require additional interaction rounds. |
| Negative | Repository validation proves contract presence, not that a remote harness actually emitted tool telemetry. |
| Follow-up | Add harness execution telemetry only if tool-call attestation becomes a product requirement. |

## FRAMEWORK.md Amendments

| Section | Amendment |
| --- | --- |
| `9. Skills` | Define mandatory evidence-first discovery, native structured questions, option comparison, recommendations, challenge behavior, and blocking-question stops. |

## Related Artifacts

| Artifact | Link |
| --- | --- |
| Canonical method | [../../FRAMEWORK.md](../../FRAMEWORK.md) |
| Shared interaction contract | [../skills/discovery-and-challenge.md](../skills/discovery-and-challenge.md) |
| Skill catalog | [../skills/README.md](../skills/README.md) |
| Dispatcher | [../../internal/dispatcher/dispatcher.go](../../internal/dispatcher/dispatcher.go) |
| Validator | [../../internal/validator/rules.go](../../internal/validator/rules.go) |

## Supersedes

- N/A
