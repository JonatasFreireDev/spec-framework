# FDR-042: CLI Lifecycle Commands

## Snapshot

| Field | Value |
| --- | --- |
| ID | `FDR-042` |
| Status | `proposed` |
| Origin EV | `CLI installation usability` |
| Date | `2026-07-13` |
| Owner | `Framework Maintainer` |

## Context

The release bootstrap scripts install a checksum-verified binary but also launch product initialization. Reusing those scripts to update the CLI can therefore start an unrelated `init` flow. Updating the CLI, upgrading an adopter product, removing the binary, and purging framework-owned caches currently require different undocumented filesystem operations.

## Decision

| Boundary | Contract |
| --- | --- |
| Installation | `scripts/install.ps1` and `scripts/install.sh` install a checksum-verified released binary, write a versioned ownership manifest beside it, and do not run `init`. Product initialization is always explicit. |
| Compatibility | Legacy `scripts/init.ps1` and `scripts/init.sh` remain deprecated wrappers for installation and no longer initialize a product automatically. |
| CLI update | Top-level `spec-framework update` checks or installs a released CLI binary. `--check` is read-only, `--version` selects a release, and replacement requires `--yes`. |
| Product upgrade | `spec-framework upgrade` remains exclusively responsible for an adopter's pinned runtime, manifest, and selected dispatchers. It never updates the executing CLI binary. |
| Supply chain | Update downloads only the official release archive and `checksums.txt`, limits response and extraction sizes, verifies SHA-256, executes the candidate's `version` smoke check, stages beside the installed executable, and rolls back failed replacement. Authenticated GitHub requests may use `GITHUB_TOKEN`. |
| CLI removal | Top-level `spec-framework uninstall` previews the binary, ownership status, manifest, and PATH entry and requires `--yes`. Managed paths are staged transactionally before deletion. It never searches for or removes product repositories. |
| Purge | `uninstall --purge` additionally removes the versioned runtime cache and only the namespaced `spec-framework` dispatcher directory for Codex, Cursor, and Claude. Other skills and harness content remain untouched. |
| Windows | Update records recoverable sidecar state before scheduling replacement. Self-removal renames the running executable and schedules a hidden PowerShell cleanup after process exit, including removal of the installer-managed user PATH entry. |
| Audit-only | `update`, `uninstall`, `upgrade`, `init`, and `version` are outside the product mutation guard. Their own confirmation and preservation contracts still apply. |

## Consequences

| Type | Consequence |
| --- | --- |
| Positive | Installation, CLI update, product upgrade, and removal have distinct memorable commands. |
| Positive | Updating the CLI cannot accidentally initialize a product. |
| Positive | Checksum or replacement failure preserves the current executable. |
| Negative | Windows uninstall requires a short-lived PowerShell helper after the CLI process exits. |
| Negative | Existing automation that relied on install-and-init must invoke `spec-framework init` explicitly. |
| Follow-up | Add signed release provenance before treating checksums from the same release channel as protection against channel compromise. |
| Follow-up | Consider package-manager distribution only after lifecycle commands and signed release provenance are stable. |

## FRAMEWORK.md Amendments

| Section | Amendment |
| --- | --- |
| `15. How Agents Use The Framework` | Distinguish installation, CLI update, adopter upgrade, and uninstall/purge while preserving product content. |

## Related Artifacts

| Artifact | Link |
| --- | --- |
| Canonical method | [../../FRAMEWORK.md](../../FRAMEWORK.md) |
| Installation guide | [../install.md](../install.md) |
| CLI lifecycle implementation | [../../internal/clifecycle/lifecycle.go](../../internal/clifecycle/lifecycle.go) |
| CLI command tree | [../../internal/cli/cobra.go](../../internal/cli/cobra.go) |
| Cobra decision | [FDR-040](FDR-040-cobra-cli-command-tree.md) |
| Generated bootstrap | [FDR-014](FDR-014-generated-adopter-bootstrap.md) |

## Supersedes

- N/A; this decision specializes the installation portion of FDR-014 while preserving its generated bootstrap and adopter-preservation contracts.
