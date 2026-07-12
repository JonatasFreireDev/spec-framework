# FDR-023: Canonical Product Design System

## Snapshot

| Field | Value |
| --- | --- |
| ID | `FDR-023` |
| Status | `approved` |
| Origin EV | Human-approved Design System evolution |
| Date | `2026-07-12` |
| Owner | `Design System Skill` |

## Context

Design and UX Review consume a design system, but the framework has no canonical artifact, ownership, token contract, validation, versioning, or impact rules for one. This allows each use case to invent local visual rules and makes external libraries difficult to adopt safely.

## Decision

Products may maintain one canonical Design System under `design/system/`. `context.md` owns identity and status, `design-system.md` is the human contract, and `tokens/tokens.json` is the mechanical token source.

| Concern | Contract |
| --- | --- |
| Identity | Stable `DSYS-NNN`; the default starter identity is `DSYS-001`. |
| Ownership | The specialist `design-system` owns canonical content; UX Review remains independent. |
| Origin | `generate`, `evolve`, or `adopt`. |
| Structure | Foundations, tokens, components, patterns, sources, and evidence live below `design/system/`. |
| Tokens | Tool-independent JSON with primitive, semantic, and component aliases. |
| Approval | Approved and later statuses require normal product approval records. No migration creates them. |
| Consumption | Use-case Design pins the Design System id/version and records tokens, components, patterns, and deviations. |
| Impact | Token removal, rename, semantic change, component removal, or external-system replacement requires impact analysis and a product decision when it changes product commitments. |
| Adapters | Figma, Penpot, Storybook, Impeccable, and implementation libraries remain optional sources/adapters. |

## Consequences

| Type | Consequence |
| --- | --- |
| Positive | Shared visual and interaction rules become versioned and reviewable. |
| Positive | Agents can validate token references and Design consumers mechanically. |
| Negative | Products with UI gain an additional shared artifact to govern. |
| Follow-up | Live provider synchronization and code generation remain outside the first release. |

## FRAMEWORK.md Amendments

| Section | Amendment |
| --- | --- |
| `3` | Define Design System as a shared product artifact. |
| `4` | Add the canonical `design/system/` structure. |
| `6.2` | Require Designs to pin an approved system when one is declared. |
| `9` | Add Design System Skill. |

## Related Artifacts

| Artifact | Link |
| --- | --- |
| Visual sources | [FDR-021](FDR-021-design-sources-and-visual-artifacts.md) |
| Framework method | [FRAMEWORK.md](../../FRAMEWORK.md) |
| Design System Skill | [Design System Skill](../skills/design-system/SKILL.md) |

## Supersedes

- N/A
