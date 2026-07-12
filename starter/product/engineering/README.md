# Product Engineering

## Purpose

Store the versioned product Engineering System: stable architecture, ownership, standards, quality attributes, fitness functions, runbooks, and evidence.

Framework validators and framework method decisions should be installed or referenced from the framework core, not mixed with product engineering decisions.

## Expected Areas

| Area | Purpose |
| --- | --- |
| `architecture/` | Product architecture notes and ADR links. |
| `runbooks/` | Product operational runbooks. |
| `conventions/` | Product-specific engineering conventions. |
| `quality/` | Product quality model and configured fitness functions. |
| `evidence/` | Evidence supporting maturity and operational claims. |

## Next Step

Complete `engineering-system.md` and `engineering-system.yaml` from real code and operational evidence. Define product-specific gates in `knowledge/conventions/gates.md`; do not claim maturity or create architecture decisions from placeholder content.

Inspect with `spec-framework engineering-system inspect` and validate with `spec-framework engineering-system validate`. Use `spec-framework engineering-system triggers` before declaring delivery-specific triggers in a use-case context.
