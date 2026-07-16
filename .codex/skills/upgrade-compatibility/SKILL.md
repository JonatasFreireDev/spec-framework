---
name: upgrade-compatibility
description: Verify that a Spec Framework distribution change can initialize new products and upgrade existing products without overwriting adopter-owned content or approvals.
---

# Upgrade Compatibility

## Purpose

Protect adopter-owned product content during distribution changes. Do not change an adopter product or approval record merely to make an upgrade pass.

## Required reading

- `FRAMEWORK.md`
- `AGENTS.md`
- `internal/install/`
- affected starter assets and upgrade tests

## Workflow

1. Identify changed embedded assets and whether they are framework runtime assets or product-owned starter assets.
2. Exercise a fresh `init` for each affected starting point and agent target.
3. Exercise `upgrade` over an initialized fixture containing changed product files, decisions, and approval history.
4. Verify upgrade refreshes only the runtime, manifest, and selected dispatchers; it must not replay initialization or overwrite adopter content.
5. For removed assets, verify no stale reference remains and document whether the removal is safe, retained, or requires an explicit migration.

## Output

Report init and upgrade evidence, preserved paths, changed managed paths, compatibility status, rollback path, and any human migration decision required.
