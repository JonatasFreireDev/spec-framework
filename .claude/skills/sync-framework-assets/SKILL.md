---
name: sync-framework-assets
description: Consistency checklist for changes to framework assets (framework/skills, framework/template, framework/validators, framework/tools, framework/decisions, FRAMEWORK.md, AGENTS.md, scripts/). Use whenever a framework asset is added, renamed, moved, or has its contract changed, to keep the starter, the worked example, the docs, and the npm package in sync.
---

# Sync Framework Assets Skill

## Purpose

This repository ships the same assets through several surfaces: the repository lab itself, the `starter/` skeleton, the worked example in `examples/events/`, the npm package, and the `.spec-framework/` tree that `installFrameworkAssets` builds inside adopter repositories. Most regressions here are not bugs in one file — they are one surface changing while its replicas silently drift. This skill is the checklist that prevents that.

## How assets flow

`scripts/init-product.mjs` copies `starter/` into the target, then `installFrameworkAssets` in `scripts/framework-assets.mjs` overwrites `.spec-framework/` with the live assets:

| Source in this repo | Installed in adopter repo |
| --- | --- |
| `FRAMEWORK.md` | `.spec-framework/FRAMEWORK.md` |
| `framework/decisions/` | `.spec-framework/decisions/` |
| `framework/skills/` | `.spec-framework/skills/` (and mirrored to `.codex/skills/`) |
| `framework/template/` | `.spec-framework/templates/` (note the rename: `template` → `templates`) |
| `framework/validators/` | `.spec-framework/validators/` |
| `framework/tools/move-artifact.mjs` | `.spec-framework/tools/move-artifact.mjs` |
| `scripts/validate-product.mjs` | `.spec-framework/tools/validate-product.mjs` |
| `framework/AGENTS.framework.md` | `.spec-framework/AGENTS.framework.md` |

`starter/.spec-framework/` holds only placeholder READMEs and `manifest.json`; real content arrives via the copy above. Do not hand-copy assets into `starter/.spec-framework/`.

## Checklist by change type

### Adding, renaming, or removing a skill in `framework/skills/`

- Register it in `framework/skills/README.md` (specialist vs orchestrator list).
- Check `FRAMEWORK.md` and `AGENTS.md` for skill listings or flow references that must mention it.
- Write paths using the dual-root contract: "the framework root's ..." / "the active product root's ..." — never hardcode `framework/` or `.spec-framework/` inside a SKILL.md body.
- If the skill owns an artifact, confirm a matching template exists in `framework/template/` and the artifact appears in the canonical flow documentation.

### Changing a template in `framework/template/`

- Update every artifact in `examples/events/` that instantiates the template, or `npm run validate` will flag the drift.
- Update the matching placeholder structure in `starter/product/` if the skeleton mirrors that artifact.
- Update the owning skill's SKILL.md if the template's sections or contract changed.
- Remember adopters see this directory as `.spec-framework/templates/` — docs addressed to adopters must use that name.

### Changing a validator or `scripts/validate-product.mjs`

- Run `npm run validate` — `examples/events/` must pass under the new rule; fix the example rather than weakening the rule.
- Check whether `framework/tests/run-tests.mjs` asserts on the old behavior.
- New rules that constrain product artifacts usually need a matching statement in `FRAMEWORK.md` or `AGENTS.md`; a rule that exists only in code is undiscoverable for agents.

### Adding a new shipped file or directory

Three places must agree, and the test suite checks them:

1. The `files` list in `package.json`.
2. The copy list in `installFrameworkAssets` (`scripts/framework-assets.mjs`).
3. The asset-inclusion tests in `framework/tests/run-tests.mjs`.

### Changing `FRAMEWORK.md`, `AGENTS.md`, or `framework/AGENTS.framework.md`

- `AGENTS.md` addresses this repository; `framework/AGENTS.framework.md` is the copy shipped to adopters. A rule that applies to both must be stated in both, with paths resolved for each audience.
- Method decisions that motivated the change belong in `framework/decisions/FDR-*`, not in product decision logs.

## Always finish with verification

After syncing, run the `verify` skill (`npm run check`, `npm test`, `npm run validate`, plus `npm run pack:dry` for packaging changes). Report which surfaces you touched and which you checked and found already consistent — silence about a surface is indistinguishable from having forgotten it.
