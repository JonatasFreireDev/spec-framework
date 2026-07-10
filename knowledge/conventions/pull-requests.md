# Pull Request Convention

## Snapshot

| Field | Value |
| --- | --- |
| Status | `approved` |
| Governed by | [FDR-008](../../engineering/decisions/FDR-008-delivery-commits-and-prs.md) |
| Owner | PR Finalizer |

## Rule

Pull requests are validation and review surfaces. A PR should carry code, documentation, and evidence links together in the monorepo.

## Required PR Body

```markdown
## Summary

## Tasks

## Evidence

## Gates

## QA

## Code Review

## Security Review

## Risks / Rollback
```

## Required Links

| Link | Required when |
| --- | --- |
| Task files | Always for task delivery. |
| QA Evidence | Validation or release PR. |
| Code Review | Validation or release PR. |
| Security Review | Sensitive/executable delivery. |
| Gate logs or CI URL | When gates run. |
| Screenshots | When delivery has UI. |

## Boundaries

- PR Finalizer may prepare a PR body without opening a PR.
- PR Finalizer may open a PR only when requested and provider tooling is available.
- PR Finalizer does not merge.
- Red QA, Code Review, or Security Review gates block non-draft PR finalization.
