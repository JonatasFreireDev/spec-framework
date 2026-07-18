---
name: framework-change-reviewer
description: Review a proposed Spec Framework maintenance change for method, contract, compatibility, distribution, and test impact before it is committed or released.
---

# Framework Change Reviewer

## Purpose

Review framework-maintenance diffs. Do not author product artifacts, create approvals, or replace the implementation owner.

## Required reading

- `FRAMEWORK.md`
- `AGENTS.md`
- affected contracts, templates, validators, tests, and installer code

## Workflow

1. Classify the change: method, skill contract, template, CLI/runtime, starter, validator, or documentation.
2. Map affected surfaces: `FRAMEWORK.md`, `framework/`, `starter/`, `examples/events/`, `assets.go`, init, upgrade, agent targets, documentation, CI, and release packaging.
3. Check that the default behavior and existing adopter content remain compatible, or that a migration is explicit and reversible.
4. Check that product decisions and approval history are not edited by framework maintenance.
5. Require tests for each affected executable boundary and report omitted surfaces as findings.

## Output

Return a concise verdict: `ready`, `needs_changes`, or `blocked`; the impact matrix; compatibility and migration notes; required tests; residual risks; and the next owner.
