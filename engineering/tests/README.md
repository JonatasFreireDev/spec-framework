# Engineering Tests

## Purpose

This folder contains local fixture-based tests for framework engineering tools. The tests create temporary repositories, run the real scripts, and remove the fixtures afterward.

## Current Coverage

| Tool | Coverage |
| --- | --- |
| `engineering/validators/framework-validator.mjs` | Approval-record enforcement, derived staleness blocking, Phase A writeScope warnings, task handoff skill references, concrete QA evidence enforcement, and blocker route/owner enforcement. |
| `engineering/move-artifact.mjs` | Folder move, Markdown link rewrite, JSON path rewrite, and free-text mention reporting. |

## Run

```bash
node engineering/tests/run-tests.mjs
```

Run these tests before changing validator gates, identity policy, staleness behavior, approval-record behavior, or artifact movement behavior.

## Next Step

Add fixtures for task-file validation, code-evidence gates, rigor-tier gates, Mermaid semantic bindings, and future Phase B writeScope errors as those areas evolve.
