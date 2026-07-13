# Domains

## Purpose

Store product domains. A domain contains user goals, goals contain features, and features contain use cases.

## Expected Shape

```text
domains/
  <domain>/
    context.md
    domain.md
    goals/
```

## Next Step

Before creating the first real domain, read the pinned framework runtime's `examples/events/domains/README.md` and `examples/events/domains/events/domain.md`. They are the canonical reference for domain modeling; the example is not copied into product scope.

Copy [_template-domain](_template-domain/README.md) to the first product domain slug after the applicable Foundation contract is approved and scope is coherent enough to guide it. For `existing-feature`, the approved Feature Brief is that contract.

Choose a stable business-area slug such as `users`, `task-management`, or `payments`; do not use the product name, a navigation label, or a catch-all domain. Define what the domain does not own, record its dependencies, and materialize one walking skeleton before creating a workspace:

```text
Domain -> User Goal -> Feature -> Use Case
```

Authentication and identity normally belong to a `users`/identity boundary, not inside an unrelated business domain. If a different boundary is intentional, document the decision and cross-domain contract.

## Anti-patterns

- A single domain named after the product that owns unrelated capabilities.
- Treating a dashboard, screen, or sidebar item as a domain boundary.
- Putting authentication in a business domain without an explicit identity boundary.
- Stopping after `domain.md` without a goal, feature, and use case.

## Template Chain

```text
_template-domain -> _template-goal -> _template-feature -> _template-use-case
```
