# Framework Validators

## Purpose

This folder stores local validation tools for the Product Engineering Framework. Validators check whether documentation artifacts are structurally ready to move through the framework gates.

## Validator

Run the framework validator from the repository root:

```bash
node engineering/validators/framework-validator.mjs --write-report
```

The validator checks:

- required use-case artifact bundles;
- execution graph JSON shape and dependencies;
- `context.md` required metadata;
- stale `product/...` paths outside `FRAMEWORK.md`;
- decision index paths in `.product/decisions.json`;
- visual Mermaid standards for flowcharts;
- local Markdown links;
- template snapshots.

## Output

When `--write-report` is provided, the validator writes:

```text
audits/framework-validation-report.md
```

## Responsible Skill

Primary owner: Documentation Orchestrator.

Supporting skills: Audit Orchestrator, Dependency Analyzer, Gap Finder, Product Historian.
