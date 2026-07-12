---
name: release-smoke
description: End-to-end smoke test of the versioned Go release binary, including init, validate, upgrade, skill targets, archives, and checksums.
---

# Release Smoke Skill

## Purpose

Prove that release archives work without Node.js or a Go toolchain on the adopter machine.

## Procedure

1. Run the verify skill.
2. Build or unpack the candidate binary in a clean temporary directory.
3. Run `spec-framework version` and `spec-framework help`.
4. Set isolated cache and agent-home overrides, then run `spec-framework init <tmp>/repo --agents codex,cursor,claude --yes` against a repository that already contains a client-owned marker.
5. Verify only `product/` was added, `product/.product/framework.json` uses manifest-only activation, `product/BOOTSTRAP.md` exists, the versioned cache is complete, and the three user-scoped dispatchers exist outside the repository.
6. Verify `.spec-framework`, repository-local agent trees, root guides, and generated CI workflows are absent.
7. Run `spec-framework validate` and `spec-framework skill path code-runner` from the generated repository.
8. Add a product-owned marker, run `upgrade --yes`, and confirm all product and client-owned markers are preserved.
9. Exercise `move --dry-run` and `move` on a temporary artifact.
10. For release archives, verify `checksums.txt` before extraction.
11. Clean temporary artifacts.

Any crash, missing asset, overwritten product file, version mismatch, invalid CI pin, or checksum failure blocks release.
