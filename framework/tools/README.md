# Framework Tools

## Purpose

This folder contains framework utilities that are installed into adopter repositories under `.spec-framework/tools/`.

## Installed Tools

| Tool | Purpose |
| --- | --- |
| `move-artifact.mjs` | Move an artifact safely while rewriting resolvable Markdown links and JSON paths. |

The bootstrap process also installs `scripts/validate-product.mjs` as `.spec-framework/tools/validate-product.mjs`.

Framework development tests remain in `framework/tests/` and are not copied into adopter repositories.
