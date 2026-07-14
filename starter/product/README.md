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
| `engineering/` | Product Engineering System and operational evidence. |
| `audits/` | Readiness, security, QA, and consistency findings. |

`BOOTSTRAP.md` and `context.md` are authoritative for sequencing; this README is navigation only.

For CI navigation checks, run `python tools/check-links.py .`. It verifies local Markdown files, section anchors, links that escape the product root, and references to registered tasks, bugs, decisions, requirements, evidence, and other artifacts that are not linked.
