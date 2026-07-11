# FDR-014: Generated adopter bootstrap

## Snapshot

| Field | Value |
| --- | --- |
| ID | `FDR-014` |
| Status | `approved` |
| Origin EV | `Adopter first-run feedback` |
| Date | `2026-07-11` |
| Owner | `Documentation Orchestrator` |

## Context

`spec-framework init` creates a structurally valid repository, but the copied starter README still speaks as if it were inside the framework source repository. A first-time adopter sees hundreds of intentional placeholders without a concise explanation of the next gate. Development builds also generate a CI workflow pinned to the nonexistent release `vdev`.

## Decision

After installing framework assets, `init` will generate two adopter-owned guides:

| File | Contract |
| --- | --- |
| `README.md` | Identifies the repository as a newly initialized product, shows installed agent integrations, documents the stable CLI commands, and links only to files that exist in the adopter repository. |
| `BOOTSTRAP.md` | Provides an ordered checklist from product identity through Problem, Vision, Strategy, first Domain, gates, and validation. It explicitly distinguishes structural validation from product readiness. |

`upgrade` must never overwrite either guide because adopters own them after initialization.

For a development build, generated CI installs the CLI from the repository's `master` Go source instead of requesting `vdev`. Versioned release builds continue to generate a workflow that downloads the pinned release binary.

## Consequences

| Type | Consequence |
| --- | --- |
| Positive | A fresh repository has a clear first-run path without reading framework-source documentation. |
| Positive | Development builds no longer generate a guaranteed-broken `vdev` download. |
| Positive | Upgrade preserves adopter edits to onboarding documentation. |
| Negative | README and bootstrap generation become part of the CLI contract and require fixtures. |
| Follow-up | Future interactive questions may optionally prefill product identity and stack-specific gates. |

## FRAMEWORK.md Amendments

| Section | Amendment |
| --- | --- |
| `4. Folder Structure` | Document generated adopter guides at repository root. |
| `15. How To Use With Codex` | Point new adopters to `BOOTSTRAP.md` before creating downstream artifacts. |

## Related Artifacts

| Artifact | Link |
| --- | --- |
| Installer | [../../internal/install/install.go](../../internal/install/install.go) |
| Installer tests | [../../internal/install/install_test.go](../../internal/install/install_test.go) |
| Adoption guide | [../adoption.md](../adoption.md) |
| Go CLI decision | [FDR-013](FDR-013-go-cli-and-agent-skill-installation.md) |

## Supersedes

- `N/A`
