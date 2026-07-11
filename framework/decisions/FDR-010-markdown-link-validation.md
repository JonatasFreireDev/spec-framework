# FDR-010: Markdown link validation

## Snapshot

| Field | Value |
| --- | --- |
| ID | `FDR-010` |
| Status | `approved` |
| Origin EV | `EV-013` |
| Date | `2026-07-10` |
| Owner | `Documentation Orchestrator / framework-validator` |

## Context

The framework depends on navigable documentation. The move tool can rewrite Markdown links during artifact relocation, but broken relative links in normal documents, reports, and templates could still pass review if no mechanical gate checked them.

Deep product paths such as `domains/<domain>/goals/<goal>/features/<feature>/use-cases/<use-case>/...` make relative links easy to break during moves, generation, or template updates.

## Decision

The framework validator must validate Markdown inline links in all `.md` files in the repository, including templates.

Relative Markdown links must resolve to an existing file or directory inside the repository. Broken relative links are errors. Links that point outside the repository are warnings. External URLs, mail links, phone links, internal anchors, images, and explicit template placeholders such as `[path]`, `TBD`, or `N/A` are ignored by this check.

## Consequences

| Type | Consequence |
| --- | --- |
| Positive | Documentation navigation becomes a CI-verifiable property. |
| Positive | Generated templates are pressured to use real relative links when the target is part of the framework. |
| Positive | The move tool and validator now form a closed loop: move rewrites resolvable links; validator catches leftovers. |
| Negative | Some intentional draft links must be represented as placeholders or planned outputs instead of broken Markdown links. |
| Follow-up | Future work may add optional anchor validation for intra-file heading fragments. |

## FRAMEWORK.md Amendments

| Section | Amendment |
| --- | --- |
| `13. Auditoria` | Link integrity is part of consistency auditing. |

## Related Artifacts

| Artifact | Link |
| --- | --- |
| FRAMEWORK | [../../FRAMEWORK.md](../../FRAMEWORK.md) |
| Validator | [../../internal/validator/](../../internal/validator/) |
| Move tool | [../../internal/moveartifact/](../../internal/moveartifact/) |
| Engineering tests | [../../internal/](../../internal/) |

## Supersedes

- `N/A`
