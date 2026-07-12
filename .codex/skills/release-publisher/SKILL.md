---
name: release-publisher
description: "Coordinate this repository's approved release from final diff review through versioning, commit, push, tag-triggered publication, checksum verification, and post-publication smoke testing."
---

# Release Publisher

## Purpose

Finish one approved release of this repository without bypassing gates or silently performing remote mutations.

## Required reading

- `FRAMEWORK.md`.
- `.codex/skills/release-smoke/SKILL.md`.
- `.codex/skills/verify/SKILL.md`.
- `framework/skills/release-orchestrator/SKILL.md`.
- `framework/skills/commit-crafter/SKILL.md`.
- `framework/skills/pr-finalizer/SKILL.md`.
- `.github/workflows/release.yml` and the release packaging configuration.
- The current Git status, branch, remotes, tags, release notes, and applicable conventions.

## Authority boundaries

- Read-only inspection, local builds, tests, checksums, and temporary smoke repositories are allowed within the requested release scope.
- Require explicit user approval for the exact version when it has not already been approved.
- Require explicit user authorization before commit, push, tag creation/push, release publication, merge, or deployment.
- Never overwrite or move a published tag.
- Never publish a duplicate release when tag-triggered automation already owns publication.
- Never stage unrelated user changes or include secrets and sensitive local artifacts.

## Workflow

1. Confirm repository, branch, remote, release scope, and the approved readiness verdict.
2. Inspect the complete diff and separate unrelated changes.
3. Run the complete verify gates, validator, candidate smoke, and required cross-builds.
4. Compare the proposed version with existing tags and change severity; obtain approval for the exact version.
5. Prepare release notes with changes, migration guidance, breaking changes, limitations, installation, and rollback.
6. With explicit authorization, use Commit Crafter to stage and create scoped atomic commits.
7. Use PR Finalizer when repository policy requires a PR; do not merge without separate authorization.
8. Re-check branch, remote, commits, and gates. Push only after explicit approval.
9. Show the exact tag and target commit. Create and push it only after explicit approval.
10. Observe the tag-triggered release workflow until completion.
11. Download published archives and `checksums.txt` into a clean temporary directory.
12. Verify checksums before extraction, version output, filenames, and the full OS/architecture matrix.
13. Run Release Smoke with the downloaded binary: version, help, manifest-only init, validate, upgrade preservation, external dispatcher/runtime checks, move dry-run/apply, import/materialization when applicable, and cleanup.
14. Report release URL, tag, commit, CI, assets, checksums, smoke evidence, limitations, rollback, and final verdict.
15. If anything fails, stop and route a corrective release. Do not conceal or rewrite publication history.

## Required release matrix

| OS | Architectures |
| --- | --- |
| Windows | `amd64`, `arm64` |
| Linux | `amd64`, `arm64` |
| macOS | `amd64`, `arm64` |

## Completion checklist

- [ ] Release readiness is approved.
- [ ] Diff and staged scope were reviewed.
- [ ] Secrets and sensitive local artifacts were excluded.
- [ ] All mechanical gates passed.
- [ ] Exact version, branch, commit, tag, and remote are recorded.
- [ ] Required remote mutations have explicit authority.
- [ ] CI release completed successfully.
- [ ] All expected assets exist.
- [ ] Published checksums passed before extraction.
- [ ] Smoke used downloaded release artifacts, not only local builds.
- [ ] Known limitations and rollback are documented.
- [ ] Final verdict is `released` or `blocked`.

## Handoff

On success, hand off to the local documentation/synchronization workflow with the tag, commit, release URL, CI URL, asset inventory, checksum results, smoke evidence, limitations, and rollback.

On failure, hand off to Release Orchestrator and the appropriate fixing skill with exact evidence and a proposed corrective-release path.
