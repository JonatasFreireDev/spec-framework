# Spec Framework Assets

## Purpose

This folder contains the installed framework assets used to operate on `product/`.

It explains how to do Specification Driven Development. It is not product scope.

## Expected Installed Assets

| Area | Purpose |
| --- | --- |
| `FRAMEWORK.md` | Method contract copied from the framework release. |
| `AGENTS.framework.md` | Agent instructions for using the installed framework without mixing product scope into `.spec-framework/`. |
| `manifest.json` | Installed framework version and asset map. |
| `skills/` | Operational agent skills. |
| `templates/` | Reusable artifact templates. |
| `validators/` | Mechanical validation gates. |
| `tools/` | Move, bootstrap, and upgrade tools. |

## Product Boundary

Product artifacts live in `../product/`.

Product decisions live in `../product/knowledge/decisions/`.

Framework method evolution is incorporated directly into `FRAMEWORK.md`, the affected contracts, validators, and tests. Git history is the maintenance record.
