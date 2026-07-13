# Template Domain

## Purpose

Copy this folder when creating the first real product domain. Rename `_template-domain` to the stable domain slug and replace all `TBD` values.

## When To Use

Use after the applicable Foundation contract has enough evidence to define a coherent business area. For `existing-feature`, this means an approved Feature Brief.

Before copying, read the pinned framework runtime's `examples/events/domains/events/domain.md`. Name the folder for the business area it owns, not for the product or a UI section. Define both ownership and non-ownership, then create the first goal -> feature -> use-case chain; a domain document alone is not a delivery slice.

## Expected Files

| File | Purpose |
| --- | --- |
| `context.md` | Domain identity, parents, children, delivery metadata, and next skill. |
| `domain.md` | Domain model, responsibilities, boundaries, and rules. |
| `decisions.md` | Domain-local decision index. |
| `goals/` | User goals owned by this domain. |

## Next Step

Rename this folder, update [context.md](context.md), define cross-domain boundaries in [domain.md](domain.md), and create the first goal from [goals/_template-goal](goals/_template-goal/README.md). Continue through its first feature and use case before creating a workspace.
