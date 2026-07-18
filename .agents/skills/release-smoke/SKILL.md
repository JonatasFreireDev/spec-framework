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
4. Set isolated cache and agent-home overrides, confirm the empty fixture has no implementation roots, then run `spec-framework init <tmp>/repo --agents codex,cursor,claude --no-code-roots --yes` against a repository that already contains a client-owned marker.
5. Verify only `product/` was added, `product/.product/framework.json` uses manifest-only activation, `product/BOOTSTRAP.md` exists, the versioned cache is complete, and the three user-scoped dispatchers exist outside the repository.
6. Verify `.spec-framework`, repository-local agent trees, root guides, and generated CI workflows are absent.
7. Run `spec-framework validate` and `spec-framework skill path code-runner` from the generated repository.
8. Add a product-owned marker, run `upgrade --yes`, and confirm all product and client-owned markers are preserved.
9. Exercise `move --dry-run` and `move` on a temporary artifact.
10. Exercise every starting point in an isolated repository:
   - `new-product`: full Foundation remains registered.
   - `existing-product`: Product Baseline and Strategy are active; consolidated Foundation artifacts are excluded.
   - `existing-documents`: the pinned run blocks `work` before materialization and permits it after selected drafts are materially complete.
   - `existing-feature`: Feature Brief approval unlocks only its declared Target Feature, and registry regeneration preserves the existing Goal parent.
   - `existing-implementation`: Assessment approval alone still blocks `work`; full Foundation approvals unlock it.
   - `audit-only`: read-only validation leaves the product tree unchanged and representative mutations are refused.
11. For import smoke, confirm materialization hashes are recorded and a later legitimate draft refinement does not invalidate the entry gate.
12. For release archives, verify `checksums.txt` before extraction.
13. Clean temporary artifacts.

Any crash, missing asset, overwritten product file, version mismatch, invalid CI pin, or checksum failure blocks release.
