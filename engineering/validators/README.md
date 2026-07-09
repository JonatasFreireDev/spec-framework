# Framework Validators

## Purpose

This folder stores local validation tools for the Product Engineering Framework. Validators check whether documentation artifacts are structurally ready to move through the framework gates.

## Validator

Run the framework validator from the repository root:

```bash
node engineering/validators/framework-validator.mjs --write-registry --write-report
```

The validator checks:

- required use-case artifact bundles;
- approval gates between Specification, Design, Implementation Plan, Execution Graph, and Tasks;
- execution graph JSON shape and dependencies;
- `context.md` required metadata;
- `.product/artifacts.json` registry consistency;
- stale `product/...` paths outside `FRAMEWORK.md`;
- decision index paths in `.product/decisions.json`;
- visual Mermaid standards for flowcharts;
- Mermaid progress state assignments using `done`, `current`, `pending`, and `blocked`;
- Mermaid semantic state bindings against `.product/artifacts.json`;
- local Markdown links;
- template snapshots.

## Output

When `--write-report` is provided, the validator writes:

```text
audits/framework-validation-report.md
```

When `--write-registry` is provided, the validator writes:

```text
.product/artifacts.json
```

## Responsible Skill

Primary owner: Documentation Orchestrator.

Supporting skills: Audit Orchestrator, Dependency Analyzer, Gap Finder, Product Historian.
