# Product Commit Convention

## Snapshot

| Field | Value |
| --- | --- |
| Status | `draft` |
| Owner | Product adopter |
| Purpose | Define how product commits should be grouped and referenced from task evidence. |

## Rule

Use atomic commits that can be traced back to task files under `domains/.../tasks/` inside the product root.

## Message Format

```text
<verb> <product concern>
```

Examples:

```text
Add event check-in validation
Fix profile privacy rule
Update payment retry tests
```

## Evidence

When a task reaches `implemented`, record branch, commits, and code paths in the task file. When it reaches `validated`, add PR and test evidence.
