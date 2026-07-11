# Delivery Closure

## Canonical Flow

```text
Domain -> Domain Evolution -> Feature Selection -> Use Cases
-> Specification Contracts -> Design -> Technical Discovery -> Architecture Gate
-> Plan -> Graph -> Tasks -> Code Runner -> Code Review -> QA
-> Commit Crafter -> PR Finalizer
```

## Select Work

```bash
spec-framework work
spec-framework work --feature FT-001 --domain events --goal manage-event --created-by "Product Owner"
spec-framework work --feature FT-001 --use-case send-invitation --created-by "Product Owner"
spec-framework status --work WORK-001
spec-framework next --work WORK-001
```

Calling `work` without a selector lists registered features. Runtime v2 workspaces are independent directories under `.product/workspaces/`, with state, handoffs, checkpoints, command plans, and evidence; multiple humans or agents do not share mutable global focus.

## Approve

```bash
spec-framework approve --artifact domains/events/context.md --grant approved --approved-by "Product Owner"
spec-framework approve --artifact domains/events/context.md --grant approved --approved-by "Product Owner" --yes
```

The first command previews the transition, resulting hash, and parent blockers. `--yes` atomically updates artifact status, registry status, and the approval record. Invalid status jumps and unapproved parents are blocked.

## Implementation Readiness

```bash
spec-framework gates
```

Applicable `TBD` commands block Code Runner. Replace them with real commands or explicit `N/A` plus rationale.

## Graph Runtime

```bash
spec-framework graph ready --graph <execution-graph.json>
spec-framework graph claim --graph <execution-graph.json> --task TK-001 --agent codex
spec-framework graph release --task TK-001 --agent codex
spec-framework graph complete --graph <execution-graph.json> --task TK-001 --agent codex
```

Claims are exclusive. Parallel claims with overlapping `writeScope` or `sharedResources` are rejected. Completion requires the owning agent and changes only operational graph state; it does not approve artifacts or execute code automatically.

## Working-tree Evidence

`implemented` records branch, base commit, changed paths, tests, gates, and a normalized diff hash. Code Review and QA independently approve that same hash. Any change makes both stale. Commit Crafter runs only after both pass; `validated` then requires commits and the remaining evidence.
