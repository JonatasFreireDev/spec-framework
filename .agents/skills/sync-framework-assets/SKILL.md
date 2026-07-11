---
name: sync-framework-assets
description: Keep embedded Go assets, starter content, framework docs, generated agent skill trees, examples, tests, and release packaging synchronized.
---

# Sync Framework Assets Skill

## Asset flow

The root `assets.go` embeds `FRAMEWORK.md`, `starter/`, framework decisions, skills, templates, and adopter instructions. `internal/install` copies those assets into `.spec-framework/` and renders selected agent trees.

## Checklist

- When changing framework skills, update `framework/skills/README.md`, neighboring handoffs, templates, and relevant FDRs.
- When changing starter content, check generated README/bootstrap ownership and upgrade preservation.
- When adding shipped assets, update `assets.go`, installer fixtures, adoption documentation, and release smoke coverage.
- Keep Codex extensions such as `agents/openai.yaml` out of Cursor and Claude trees.
- Keep `.claude/skills` and `.agents/skills` maintenance copies synchronized.
- When changing CLI distribution, update GoReleaser, CI, installation docs, manifests, and versioned workflow generation.

Finish with the verify and release-smoke skills.
