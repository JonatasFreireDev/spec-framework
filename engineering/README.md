# Engineering

## Purpose

Engineering stores implementation-facing guidance that is derived from approved specifications, designs, implementation plans, execution graphs, and tasks.

## When To Use

Use this folder for engineering conventions, environment notes, testing strategy, architecture constraints, deployment assumptions, and implementation handoff guidance. Do not use it to create application code or bypass specification gates.

## Expected Files

- `README.md`: engineering purpose and operating rules.
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
