# Commit Convention

## Snapshot

| Field | Value |
| --- | --- |
| Status | `approved` |
| Governed by | [FDR-008](../../engineering/decisions/FDR-008-delivery-commits-and-prs.md) |
| Owner | Commit Crafter |

## Rule

Commits are local delivery evidence. They should be atomic, reviewable, and traceable to tasks or framework evolution work.

## Message Format

Use concise imperative messages:

```text
<verb> <scope or artifact>
```

Examples:

```text
Add code review gate and validation
Update QA evidence template routing
Fix QR check-in duplicate handling
```

## Grouping

| Concern | Commit separately when practical |
| --- | --- |
| Framework docs | Yes |
| Skills | Yes |
| Validator/tooling | Yes |
| Templates/conventions | Yes |
| Product docs | Yes |
| Application code | Yes, grouped by task/use case |

## Safety

- Never commit secrets, credentials, tokens, private keys, `.env` files, or local-only artifacts.
- Do not stage unrelated user changes.
- Do not push from Commit Crafter.
- Record concrete commit hashes in task files only when task status requires them and the user requested the artifact update.
