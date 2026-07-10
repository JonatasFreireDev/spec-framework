# FDR-008: Delivery commits and PR finalization

## Snapshot

| Field | Value |
| --- | --- |
| ID | `FDR-008` |
| Status | `approved` |
| Origin EV | `EV-010` |
| Date | `2026-07-10` |
| Owner | `Release Orchestrator` |

## Context

DEC-004 established the monorepo code-link convention: implemented tasks require branch, commits, and code paths; validated tasks require PR and concrete validation evidence.

The framework had validator checks for these fields, but no delivery skills to package commits or prepare PRs according to the framework gates.

## Decision

Create two delivery skills:

| Skill | Contract |
| --- | --- |
| `commit-crafter` | Creates local atomic commits by concern when explicitly asked, follows commit conventions, checks staged content for secrets, and never pushes. |
| `pr-finalizer` | Verifies hard preconditions, prepares or opens a PR, links evidence, records the PR back to task files when appropriate, and never merges. |

Delivery conventions live in:

- `knowledge/conventions/commits.md`
- `knowledge/conventions/pull-requests.md`

Commit and PR references in task evidence must be concrete and traceable.

## Consequences

| Type | Consequence |
| --- | --- |
| Positive | Delivery evidence becomes operational instead of only declarative. |
| Positive | Commits and PRs remain under explicit human/requested control. |
| Positive | PRs link back to tasks, QA Evidence, Code Review, Security Review, and gates. |
| Negative | Delivery steps add friction when evidence is incomplete. |
| Follow-up | Future release orchestration can consume PR evidence for release notes and deployment readiness. |

## FRAMEWORK.md Amendments

| Section | Amendment |
| --- | --- |
| `9. Skills` | Add Commit Crafter and PR Finalizer. |
| `10. Orquestradores` | Add delivery orchestration from local commits to PR finalization. |
| `11. Gates De Aprovacao` | Clarify implemented/validated evidence and no-push/no-merge boundaries. |

## Related Artifacts

| Artifact | Link |
| --- | --- |
| Commit Crafter skill | [../../.codex/skills/commit-crafter/SKILL.md](../../.codex/skills/commit-crafter/SKILL.md) |
| PR Finalizer skill | [../../.codex/skills/pr-finalizer/SKILL.md](../../.codex/skills/pr-finalizer/SKILL.md) |
| Commit convention | [../../knowledge/conventions/commits.md](../../examples/events/knowledge/conventions/commits.md) |
| Pull request convention | [../../knowledge/conventions/pull-requests.md](../../examples/events/knowledge/conventions/pull-requests.md) |
| FRAMEWORK | [../../FRAMEWORK.md](../../FRAMEWORK.md) |

## Supersedes

- `N/A`
