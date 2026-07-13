# FDR-041: Declarative Initialization Contracts

## Snapshot

| Field | Value |
| --- | --- |
| ID | `FDR-041` |
| Status | `proposed` |
| Origin EV | `CLI architecture evolution` |
| Date | `2026-07-13` |
| Owner | `Framework Maintainer` |

## Context

The Go installer copies the complete `starter/product/` tree for every starting point and then applies entry-specific files, prose substitutions, registry rewrites, and import behavior through hard-coded branches. Adding or reviewing a starting point therefore requires reconstructing its materialized product from installer control flow, starter assets, generated bootstrap profiles, and tests. The static copy also prevents the physical product tree from evolving proportionally without another installer rewrite.

## Decision

| Boundary | Contract |
| --- | --- |
| Contract source | Every supported starting point has one versioned JSON contract under `framework/init/contracts/`; reusable asset groups are declared in `framework/init/catalog.json`. |
| Selection | `init --starting-point` resolves exactly one contract whose id and bootstrap profile match the parsed starting point. |
| Composition | Contracts select asset sets, explicit directories including intentionally empty paths, entry-specific template files and deterministic patches, initial artifact-registry transformations, and typed post-materialization actions. |
| Planning | The CLI strictly parses contracts, expands all selected assets and directories in memory, applies replacements and patches, constructs the registry, and rejects unknown fields, missing sources, file/directory collisions, ambiguous anchors, invalid paths, duplicate artifacts, missing artifact files, and broken parent or target relationships before writing. |
| Execution | A valid plan is written under a temporary directory inside the target repository. Guides, manifest, registry, runtime preparation, dispatchers, and typed initialization actions must succeed before the staged `product/` is atomically published. Failure removes staging and does not leave a partial `product/`. |
| Safety | Contract targets are confined to `product/`. Contracts contain data only and cannot execute shell commands, load plugins, or name arbitrary Go functions. |
| Typed actions | Runtime-dependent behavior remains a closed set implemented by the CLI. `existing-documents` may select `create-import-run`; contracts cannot invent executable behavior or select arbitrary functions. |
| Compatibility | Initial contracts reproduce existing starting-point outputs and preserve reference skeletons required by FDR-031 and FDR-033. Future physical reductions are contract changes only when they preserve the applicable method or receive a separate framework decision. |
| Upgrade | `upgrade` never replays initialization composition, patches, registry transformations, or actions over adopter-owned content. It continues to refresh only the runtime, dispatcher, and manifest. |
| Existing product | `init` never overwrites an existing `product/`; `--force` does not relax adopter-content preservation. |
| Configuration boundary | These versioned framework-owned contracts configure materialization only. They do not introduce Viper, user-level configuration precedence, or adopter-editable executable configuration. |

## Consequences

| Type | Consequence |
| --- | --- |
| Positive | Starting-point output becomes inspectable and reviewable without following entry-specific Go branches. |
| Positive | Adding or changing a starting point is primarily a validated asset-contract change. |
| Positive | Shared assets remain canonical while entry-specific registries and files can vary proportionally. |
| Positive | Invalid contracts and artifact graphs fail before product publication, without leaving a partially initialized tree. |
| Positive | Contracts can preserve meaningful empty directories instead of relying on incidental placeholder files. |
| Negative | Starter file moves must now update catalog or contract references in addition to embedded asset tests. |
| Negative | Deterministic text patches still depend on stable anchors until affected starter prose is promoted to dedicated templates. |
| Follow-up | Replace prose patches with dedicated rendered templates if starting-point-specific guidance grows materially. |
| Follow-up | Persist a deterministic contract/catalog/input digest in the product manifest when reproducible initialization audits require it. |
| Follow-up | Extract typed actions into a dedicated Strategy Registry if more runtime-dependent actions are introduced. |
| Follow-up | Share starting-point invariants with validator and workflow packages if contract evolution creates duplicated semantic rules. |

## FRAMEWORK.md Amendments

| Section | Amendment |
| --- | --- |
| `15. How Agents Use The Framework` | Define agent-independent declarative initialization contracts, strict planning, atomic materialization, typed actions, compatibility, and upgrade preservation. |

## Related Artifacts

| Artifact | Link |
| --- | --- |
| Canonical method | [../../FRAMEWORK.md](../../FRAMEWORK.md) |
| Contract catalog | [../init/catalog.json](../init/catalog.json) |
| Contract schema | [../init/schema.json](../init/schema.json) |
| Starting-point contracts | [../init/contracts/new-product.json](../init/contracts/new-product.json) |
| Installer | [../../internal/install/install.go](../../internal/install/install.go) |
| External runtime | [FDR-025](FDR-025-external-runtime-and-manifest-only-activation.md) |
| Cobra configuration boundary | [FDR-040](FDR-040-cobra-cli-command-tree.md) |

## Supersedes

- N/A
