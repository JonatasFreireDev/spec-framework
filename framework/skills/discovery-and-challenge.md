# Discovery And Challenge Contract

This contract applies to framework skills that define product intent, scope, experience, architecture, or implementation plans.

## Required behavior

Before drafting or materially updating an artifact, the skill must:

1. Inspect available repository and CLI evidence first. Do not ask the human for facts that can be discovered safely.
2. Identify missing information, assumptions, conflicts, dependencies, risks, and meaningful choices.
3. Use the harness-native structured question tool when it is available. Do not replace an available native tool with assumptions or a question buried in prose.
4. Ask one to three focused questions per round. Explain why each answer changes the artifact.
5. For a meaningful choice, present two or three concrete options, state trade-offs, recommend one option, and allow a free-form answer.
6. Distinguish blocking questions from useful but non-blocking questions. Continue with an explicitly labeled assumption only when the question is non-blocking and the harness supports that flow.
7. Challenge the requested direction when evidence exposes scope, dependency, usability, security, operability, reversibility, approval, or delivery risk. State the likely consequence and a safer alternative.
8. Summarize resolved answers, retained assumptions, rejected alternatives, warnings, and open questions in the artifact or its `context.md`. Do not store raw conversation transcripts.
9. Stop before finalizing or handing off when a blocking question remains unanswered. A conversational answer never grants formal product approval.

## Harness capability

The canonical capability is `native_user_question`. Harness adapters map it to their default structured question mechanism, such as `request_user_input` in Codex or `AskUserQuestion` in Claude Code. Cursor uses its native user-question mechanism when exposed. If the harness exposes no structured question tool, ask a concise explicit question in conversation and record that fallback in the handoff.

## Mode behavior

| Mode | Requirement |
| --- | --- |
| `create` | Discovery is mandatory before the first substantive draft. |
| `update` or `evolve` | Discovery is mandatory when scope, behavior, risk, priority, or an approved decision may change. |
| `audit`, `compare`, or `explain` | Ask only when human input is required to resolve ambiguity or choose a next action. Otherwise report findings directly. |

## Minimum handoff

The handoff must make clear:

- questions resolved;
- assumptions still in force;
- alternatives considered and the recommendation;
- risks and warnings;
- blocking open questions;
- the next human or skill decision.
