# Domains

## Purpose

Domains organize product knowledge by coherent business or product area. A domain owns goals, features, use cases, and their traceability to strategy.

## When To Use

Use this folder when creating or auditing product areas such as `events`, `users`, `friendship`, or `payments`. Do not create global features outside a domain.

## Expected Files

- `<domain>/context.md`: domain context, dependencies, children, and handoff.
- `<domain>/domain.md`: domain definition, boundaries, rules, metrics, and risks.
- `<domain>/decisions.md`: domain-local decision index when useful.
- `<domain>/goals/<goal>/`: user goals inside the domain.

## Responsible Skill

Primary owner: Domain Architect AI.

Supporting skills: User Goal AI, Feature AI, Use Case AI, Specification AI.

## Next Step

For a new domain, create `context.md` and `domain.md`, then model user goals under `goals/`.
