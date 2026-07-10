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

Copy [_template-domain](_template-domain/README.md) to the first product domain slug after problem, vision, and strategy are coherent enough to guide scope.

## Template Chain

```text
_template-domain -> _template-goal -> _template-feature -> _template-use-case
```
