---
name: sync-framework-assets
description: Keep embedded Go assets, starter content, framework docs, generated agent skill trees, examples, tests, and release packaging synchronized.
---

# Sync Framework Assets Skill

## Asset flow

The root `assets.go` embeds `FRAMEWORK.md`, `starter/product/`, framework skills, templates, and adopter instructions. `internal/install` copies only product-owned starter content into the adopter repository; framework assets are materialized in the versioned user cache and selected harnesses receive a user-scoped dispatcher.

## Checklist

- When changing framework skills, update `FRAMEWORK.md`, `framework/skills/README.md`, neighboring handoffs, templates, validators, and tests as applicable.
- When changing starter content, check generated README/bootstrap ownership and upgrade preservation.
- When adding shipped assets, update `assets.go`, installer fixtures, adoption documentation, and release smoke coverage.
- Keep `.agents/skills/` as the single repository-maintenance skill tree; do not recreate target-specific maintenance copies.
- Keep the dispatcher namespaced and gated exclusively by `product/.product/framework.json`.
- Verify that initialization creates no repository-local agent trees or `.spec-framework/` directory.
- When changing CLI distribution, update GoReleaser, CI, installation docs, manifests, and versioned workflow generation.

Finish with the verify and release-smoke skills.
