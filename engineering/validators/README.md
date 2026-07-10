# Framework Validators

## Purpose

This folder stores local validation tools for the Product Engineering Framework. Validators check whether documentation artifacts are structurally ready to move through the framework gates.

## Validator

Run the framework validator from the repository root:

```bash
node engineering/validators/framework-validator.mjs --write-registry --write-report
```

Run validator and move-tool fixture tests with:

```bash
node engineering/tests/run-tests.mjs
```

The GitHub Actions workflow `.github/workflows/framework-validation.yml` runs syntax checks, engineering tests, and the validator without write flags so CI verifies the committed framework state without mutating generated reports.

The validator checks:

- required use-case artifact bundles;
- traceability between parent, child, source, graph, and task artifacts;
- status policy between parent and child artifacts;
- approval gates between Specification, Design, Implementation Plan, Execution Graph, and Tasks;
- approval records in `.product/history/` for `approved` and later statuses;
- derived staleness through `.product/derivations.json`;
- validation gates for approved tests, QA evidence, Security Review, and audit before `validated` or `released`;
- Code Review gate and review verdict quality before `validated` or `released`;
- traceable commit hashes/URLs for implemented tasks and traceable PR references for validated tasks;
- Delivery Level and Priority metadata for executable framework artifacts;
- execution graph JSON shape and dependencies;
- `writeScope` and `sharedResources` safety for parallel graph nodes, currently as Phase A warnings;
- `context.md` required metadata;
- `.product/artifacts.json` registry consistency;
- stale `product/...` paths outside framework adoption documentation;
- decision index paths in `.product/decisions.json`;
- decision references against `.product/decisions.json`, including approved delivery dependencies;
- visual Mermaid standards for flowcharts;
- Mermaid progress state assignments using `done`, `current`, `pending`, and `blocked`;
- Mermaid semantic state bindings against `.product/artifacts.json`;
- local Markdown links, including real links inside templates;
- template snapshots.
- identity policy metadata in `.product/ids.json`;
- immutable `slug` metadata in `context.md`.
- concrete QA evidence for approved or later `qa-evidence.md` artifacts.
- task handoff fields that reference repository-local skills.
- blocker/high/required_fix findings in approved QA evidence, Code Review, Security Review, and audit artifacts include route and owner.

## Output

When `--write-report` is provided, the validator writes:

```text
audits/framework-validation-report.md
audits/readiness/framework-readiness.md
```

When `--write-registry` is provided, the validator writes:

```text
.product/artifacts.json
```

## Responsible Skill

Primary owner: Documentation Orchestrator.

Supporting skills: Audit Orchestrator, Dependency Analyzer, Gap Finder, QA AI, Security Review AI, Product Historian.
