# Engineering

## Purpose

Engineering stores implementation-facing guidance that is derived from approved specifications, designs, implementation plans, execution graphs, and tasks.

## When To Use

Use this folder for engineering conventions, environment notes, testing strategy, architecture constraints, deployment assumptions, and implementation handoff guidance. Do not use it to create application code or bypass specification gates.

## Expected Files

- `README.md`: engineering purpose and operating rules.
- `validators/`: local framework validation scripts.
- `move-artifact.mjs`: safe artifact move tool that rewrites Markdown links and JSON paths.
- Future convention files when approved, such as `testing.md`, `architecture-notes.md`, `observability.md`, or `release-checklist.md`.

## Responsible Skill

Primary owner: Implementation Planner AI.

Supporting skills: Code Runner AI, QA AI, Code Review AI, Security Review AI, Documentation Writer AI.

## Operating Rules

- No code implementation starts without an approved Specification.
- Use cases with UI require approved `design.md` before `implementation-plan.md`.
- Tasks must be generated from an Execution Graph.
- Engineering decisions that change architecture, security, data, or external dependencies require a decision record.

## Next Step

When a specification and design are approved, create or update the use-case `implementation-plan.md` and `execution-graph.json` before generating tasks.

## Validation

Run the framework validator before major documentation handoffs:

```bash
node engineering/validators/framework-validator.mjs --write-registry --write-report
```

The registry is written to `.product/artifacts.json` and the report is written to `audits/framework-validation-report.md`.

Move an artifact with:

```bash
node engineering/move-artifact.mjs --from domains/old/path --to domains/new/path --dry-run
node engineering/move-artifact.mjs --from domains/old/path --to domains/new/path
```

The move tool rewrites Markdown links and JSON paths it can resolve mechanically, then reports free-text mentions for human review.

Run engineering tool tests with:

```bash
node engineering/tests/run-tests.mjs
```

The GitHub Actions workflow at `.github/workflows/framework-validation.yml` runs the same checks on pushes and pull requests:

```bash
node --check engineering/validators/framework-validator.mjs
node --check engineering/move-artifact.mjs
node --check engineering/tests/run-tests.mjs
node engineering/tests/run-tests.mjs
node engineering/validators/framework-validator.mjs
```
