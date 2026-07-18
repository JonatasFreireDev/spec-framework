# Product

Product-owned specifications, decisions, plans, and evidence live here. Framework skills and templates remain in the versioned user cache pinned by `.product/framework.json`.

## Start here

1. Read `BOOTSTRAP.md` for the selected starting point and current gate.
2. Read the nearest `context.md` before changing an artifact.
3. Use `spec-framework guide` or `spec-framework dashboard` for the next safe action.

## Areas

| Area | Owns |
| --- | --- |
| `foundation/` | Product direction and starting-point contracts. |
| `domains/` | Domains, goals, features, use cases, and delivery artifacts. |
| `knowledge/` | Product rules, decisions, conventions, and imports. |
| `design/` | Shared and use-case design sources. |
| `engineering/` | Product Engineering System, technical entity graph, standards, quality, operations, and evidence. |
| `audits/` | Readiness, security, QA, and consistency findings. |

`BOOTSTRAP.md` and `context.md` are authoritative for sequencing; this README is navigation only.

For CI navigation checks, run `python tools/check-links.py .` and `spec-framework decisions check --strict`. The first verifies local Markdown files and anchors; the second verifies indexed decisions, domain/path coherence, approvals, references, and canonical decision links.
