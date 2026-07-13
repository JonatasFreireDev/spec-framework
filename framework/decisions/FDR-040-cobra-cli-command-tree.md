# FDR-040: Cobra CLI Command Tree

## Snapshot

| Field | Value |
| --- | --- |
| ID | `FDR-040` |
| Status | `proposed` |
| Origin EV | `CLI architecture evolution` |
| Date | `2026-07-13` |
| Owner | `Framework Guide` |

## Context

The CLI dispatches top-level commands through one manually maintained switch and each command family owns its own standard-library flag parser. This duplicates command discovery and help behavior, makes the command hierarchy harder to inspect, and provides no shared command-tree extension point.

The product manifest, runtime pin, and command flags are explicit framework contracts. The CLI does not currently have user-level settings, layered environment configuration, or configuration-file discovery that would justify a general configuration registry.

## Decision

| Boundary | Contract |
| --- | --- |
| Command tree | Cobra owns the root command, stable top-level command registration, command discovery, root and leaf help routing, and unknown-command handling. |
| Compatibility | Existing command handlers remain execution adapters during migration. They retain their current flags, positional parsing, output, and exit codes until an individual command family is deliberately migrated. |
| Audit-only guard | The root command applies the existing audit-only mutation guard before leaf execution, preserving exemptions for `init`, `upgrade`, and `version`. |
| Configuration | Do not add Viper. Product manifests and explicit flags remain authoritative; no documented layered user configuration requirement exists. |
| Future migration | New commands use Cobra directly. Existing command families may migrate their local flags to Cobra incrementally with compatibility tests. |

## Consequences

| Type | Consequence |
| --- | --- |
| Positive | The CLI has one inspectable command tree and standard help behavior. |
| Positive | The refactor preserves existing automation contracts while enabling incremental command-family cleanup. |
| Positive | The dependency graph avoids Viper and its configuration precedence where no such policy is needed. |
| Negative | Standard-library flag parsing and detailed per-flag help remain temporarily inside legacy execution adapters. |
| Follow-up | Migrate command-family flags to Cobra when their help, aliases, completion, or shared persistent flags need redesign. |

## FRAMEWORK.md Amendments

| Section | Amendment |
| --- | --- |
| `15. How To Use With Codex` | Define Cobra as the CLI command-tree implementation and the constraint for adding Viper. |

## Related Artifacts

| Artifact | Link |
| --- | --- |
| CLI root | [../../internal/cli/cobra.go](../../internal/cli/cobra.go) |
| Existing CLI handlers | [../../internal/cli/app.go](../../internal/cli/app.go) |
| Command guidance | [../../README.md](../../README.md) |
| CLI portability | [FDR-013](FDR-013-go-cli-and-agent-skill-installation.md) |
| Canonical method | [../../FRAMEWORK.md](../../FRAMEWORK.md) |

## Supersedes

- N/A
