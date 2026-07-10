# Decisions

## Purpose

Store approved or proposed decisions that change scope, behavior, sequencing, architecture, permissions, privacy, security, external dependencies, or delivery commitments.

## When To Use

Create a decision whenever the framework says an approval gate is required or when a choice would be expensive to reverse.

## Expected Files

- `DEC-XXX-short-title.md`: one decision per file using `knowledge/templates/decision-template.md`.
- `README.md`: folder purpose and usage.

## Responsible Skill

Primary owner: Product Historian AI.

Supporting skills: Impact Analysis AI, Security Review AI, Documentation Writer AI.

## Next Step

After a decision is approved, update `.product/decisions.json` and link the decision from affected `context.md`, specification, plan, or audit files.
