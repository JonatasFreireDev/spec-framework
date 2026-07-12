# FDR-021: Design Sources and Visual Artifact Protocol

## Snapshot

| Field | Value |
| --- | --- |
| ID | `FDR-021` |
| Status | `approved` |
| Origin EV | Human-approved design workflow evolution |
| Date | `2026-07-12` |
| Owner | `UX/UI AI` |

## Context

The framework models Design as a required planning contract, but it does not distinguish designs created by an agent, evolved from an existing interface, or adopted from an external canonical source. It also lacks a portable contract for versioned images, Figma or Penpot exports, wireframes, mockups, prototypes, requirement coverage, and fidelity evidence.

## Decision

Keep `design.md` as the canonical use-case Design contract and add two independent dimensions:

- `origin_mode`: `generate`, `evolve`, or `adopt`.
- `maturity`: `contract`, `wireframe`, `mockup`, or `prototype`.

External sources may be authoritative for presentation while approved decisions and the Specification remain authoritative for product behavior, security, privacy, and business rules.

| Concern | Contract |
| --- | --- |
| Source authority | `behavioral`, `visual_canonical`, `reference`, or `inspiration`. |
| Fidelity | `strict`, `balanced`, or `exploratory`. |
| Storage | Source manifests and snapshots live under `product/design/`; `design.md` remains beside the use case. |
| Versioning | Local files use SHA-256; remote tools use an immutable version identifier and may include local snapshots. |
| Staleness | Source manifests are derivation inputs. A source version change invalidates dependent Design and downstream planning. |
| Adapters | Images are built in. Impeccable, Figma, and Penpot are optional adapters behind a common contract. |
| Prototype boundary | Design prototypes live under `product/design/` and are explicitly non-production artifacts. |
| Approval | Import, generation, audit, and synchronization never approve Design or create product approval records. |
| Review independence | UX Review is read-only and separate from UX/UI generation. |

## Source Precedence

```text
Approved decisions
-> security, privacy, and accessibility requirements
-> approved Specification
-> canonical visual source
-> approved design system
-> references
-> inspiration
```

Conflicts are reported and routed for human resolution. They are never reconciled silently.

## Consequences

| Type | Consequence |
| --- | --- |
| Positive | Existing descriptive designs remain valid as `generate/contract`. |
| Positive | Figma, Penpot, images, websites, and generated artifacts share one normalized handoff. |
| Positive | Requirement coverage and visual fidelity become mechanically inspectable. |
| Negative | The product gains manifests, snapshots, and additional validation rules. |
| Negative | Remote adapters require credentials and can be unavailable. |
| Mitigation | The local image adapter is the baseline and all remote adapters are optional. |

## Rollout

1. Extend the framework method, template, UX/UI skill, and add independent UX Review.
2. Add source/use-case manifests and local image import.
3. Add requirement mapping, validation, staleness, dashboard, guide, and migration.
4. Add optional Impeccable, Figma, and Penpot adapters.
5. Add visual audit and fidelity evidence.

## Related Artifacts

| Artifact | Link |
| --- | --- |
| Framework method | [FRAMEWORK.md](../../FRAMEWORK.md) |
| UX/UI skill | [UX/UI Skill](../skills/ux-ui/SKILL.md) |
| Design template | [Design template](../template/design-template.md) |

## Supersedes

- N/A
