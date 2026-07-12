# FDR-025: External runtime and manifest-only activation

## Snapshot

| Field | Value |
| --- | --- |
| ID | `FDR-025` |
| Status | `proposed` |
| Origin EV | Governance baseline |
| Date | `2026-07-12` |
| Owner | Framework Maintainer |

## Context

The adopter installation currently copies framework assets, generated skill trees, guides, and CI files into every product repository. This makes initialization visibly larger than the product documentation it is intended to create, allows generic skill names to collide with other projects, and couples an adopter repository to generated copies of framework-owned content. The product already has a canonical adoption record at `product/.product/framework.json`, but commands and skills do not use it as the exclusive activation and version-resolution boundary.

## Decision

| Boundary | Contract |
| --- | --- |
| Activation | A repository activates Spec Framework only when `product/.product/framework.json` exists, is valid, identifies `spec-framework`, pins a concrete version, and declares `activation.mode: manifest-only`. A user mention or keyword is never an activation signal. |
| Adopter writes | Initialization adds only `product/` to the adopter repository. Framework bootstrap and documentary commands do not create `.spec-framework/`, repository-local agent skill trees, root guides, or CI workflows. |
| Runtime | The released CLI embeds framework assets and materializes the pinned version under the operating system's user cache. Multiple versions may coexist. `SPEC_FRAMEWORK_CACHE` is the test and administration override. |
| Skills | Each selected harness receives one user-scoped, namespaced `spec-framework` dispatcher. It checks the canonical manifest before resolving a specialized skill from the pinned runtime. Specialized framework skills are not copied into adopter repositories. |
| Commands | Commands run from the repository root and resolve product and framework roots from the manifest. Explicit roots remain available for framework development and diagnostics. |
| Configuration | Initialization choices are persisted in the product manifest and product-owned metadata so later commands follow the selected starting point and agents. |
| Upgrade | Upgrade refreshes the external versioned runtime and the product manifest while preserving adopter-owned content. It does not copy framework assets into the repository. |
| Code writes | Documentation and bootstrap operations remain inside `product/`. Explicit implementation work may write outside `product/` only under the approved task `writeScope`. |
| Compatibility | Legacy `.spec-framework/` repositories remain migration inputs. Migration must preview removals and must not create approval records or delete adopter content implicitly. |

## Consequences

| Type | Consequence |
| --- | --- |
| Positive | Adopter repositories contain product artifacts instead of duplicated framework runtime assets. |
| Positive | Projects can pin different framework versions without global skill-contract collisions. |
| Positive | The embedded binary provides an offline runtime after the initial CLI bootstrap. |
| Negative | A fresh machine needs the CLI bootstrap before commands can resolve the cached runtime. |
| Negative | Agent harnesses need a supported user-scoped dispatcher location. |
| Follow-up | Add signed/checksummed one-command bootstrap scripts for Windows, Linux, and macOS publication. |
| Follow-up | Provide an explicit legacy migration command before removing legacy read compatibility. |

## FRAMEWORK.md Amendments

| Section | Amendment |
| --- | --- |
| `4. Folder Structure` | Make `product/` the only framework-created adopter directory and move method/runtime assets to the versioned user cache. |
| `15. How To Use With Codex` | Define manifest-only activation, user-scoped dispatch, root command discovery, and the one-command initialization target. |

## Related Artifacts

| Artifact | Link |
| --- | --- |
| Canonical method | [../../FRAMEWORK.md](../../FRAMEWORK.md) |
| Adoption guide | [../adoption.md](../adoption.md) |
| Installation guide | [../install.md](../install.md) |
| Prior CLI decision | [FDR-013](FDR-013-go-cli-and-agent-skill-installation.md) |
| Prior bootstrap decision | [FDR-014](FDR-014-generated-adopter-bootstrap.md) |

## Supersedes

- The repository-local asset and generated-skill installation portions of `FDR-013` and `FDR-014`; their CLI portability, guided bootstrap, and preservation requirements remain in force.
