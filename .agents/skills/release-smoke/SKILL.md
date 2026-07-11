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
4. Run `spec-framework init --target <tmp>/product --agents codex,cursor,claude --yes`.
5. Verify `.agents/skills`, `.cursor/skills`, `.claude/skills`, `.spec-framework`, `README.md`, `BOOTSTRAP.md`, manifests, and CI workflow.
6. Run `spec-framework validate` from the generated repository.
7. Add a product-owned marker, run `upgrade --yes`, and confirm the marker plus generated README/bootstrap are preserved.
8. Exercise `move --dry-run` and `move` on a temporary artifact.
9. For release archives, verify `checksums.txt` before extraction.
10. Clean temporary artifacts.

Any crash, missing asset, overwritten product file, version mismatch, invalid CI pin, or checksum failure blocks release.
