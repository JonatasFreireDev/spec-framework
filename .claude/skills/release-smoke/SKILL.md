---
name: release-smoke
description: End-to-end smoke test of the packaged spec-framework CLI — pack the tarball, install it in a clean consumer, init a product, validate it, and exercise upgrade. Use before a release, after changes to scripts/init-product.mjs, upgrade-product.mjs, framework-assets.mjs, or the package.json files/bin fields.
---

# Release Smoke Skill

## Purpose

Unit tests cover the pieces; this skill proves the shipped artifact works: the npm tarball installs, the CLI runs from a consumer project, `init` assembles a complete adopter repository, `validate` runs against it, and `upgrade` refreshes assets in place. Run it before any release and after any change to the packaging or bootstrap path.

## Procedure

Work in a temporary directory outside the repository (use the session scratchpad). `<repo>` is the spec-framework repository root.

### 1. Pack

```bash
cd <repo>
npm run check && npm test
npm pack
```

Note the tarball name (`spec-framework-<version>.tgz`). Inspect the file list in the pack output: it must include `framework/skills/`, `framework/template/`, `framework/validators/`, `framework/decisions/`, `framework/tools/move-artifact.mjs`, `scripts/`, `starter/`, `FRAMEWORK.md`, and `AGENTS.md` — and must not include `examples/`, `.claude/`, or `.codex/`.

### 2. Install in a clean consumer

```bash
mkdir <tmp>/consumer && cd <tmp>/consumer
npm init -y
npm install <repo>/spec-framework-<version>.tgz --no-save
npx spec-framework help
```

`help` must print usage without errors.

### 3. Init a product

```bash
npx spec-framework init --target ../my-product
```

Then verify in `<tmp>/my-product`:

- `product/` matches the starter skeleton (foundation, domains, knowledge, audits, design, releases, engineering).
- `.spec-framework/` is fully populated: `FRAMEWORK.md`, `skills/` (all skills, not just the placeholder README), `templates/`, `validators/`, `decisions/`, `tools/move-artifact.mjs`, `tools/validate-product.mjs`, `AGENTS.framework.md`.
- `.codex/skills/` mirrors `framework/skills/` (default; disabled only when init was run with the mirror off).
- `.github/workflows/framework-validation.yml` exists.
- `.spec-framework/manifest.json` and `product/.product/framework.json` carry the current git short SHA as `version` and a populated `installed_assets` list.

### 4. Validate the fresh product

```bash
cd ../my-product
npx spec-framework validate
# equivalently: node .spec-framework/tools/validate-product.mjs
```

Gate: the validator must run to completion and produce a report. Findings about unfilled foundation content are expected on a fresh skeleton and are not failures; a crash, path-resolution error, or missing-asset error is.

### 5. Upgrade in place

Make a marker to prove refresh semantics, then upgrade:

```bash
# from <tmp>/consumer
npx spec-framework upgrade --target ../my-product
```

Verify: `.spec-framework/` assets are refreshed (version in `manifest.json` updated), and product-owned content under `product/` is untouched. `upgrade` must refuse to run against a target missing `product/` or `.spec-framework/`.

### 6. Clean up

Remove the tarball from the repo root and the temporary consumer/product directories.

## Reporting

Report a table of the six steps with pass/fail and evidence (command output excerpts). Any failure blocks release; route packaging failures to the three-surface check in the `sync-framework-assets` skill (package.json `files`, `installFrameworkAssets` copy list, asset-inclusion tests).
