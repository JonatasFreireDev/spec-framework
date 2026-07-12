---
name: fdr
description: Create or update a Framework Decision Record in framework/decisions/. Use when a change alters the framework method, a skill contract, validator behavior, gates, or the delivery workflow — or when the user asks to record a framework decision (FDR).
---

# FDR Skill

## Purpose

Framework Decision Records capture decisions about the framework method itself: agent contracts, gates, validation rules, and delivery workflow. This skill creates them in the right place, in the canonical format, and keeps the index in sync.

## FDR or DEC? Decide first

- **FDR** (`framework/decisions/FDR-*`): the decision changes how the framework works — method, skill contracts, validator behavior, gates, delivery workflow, repository boundaries.
- **DEC** (active product root's `knowledge/decisions/DEC-*`, e.g. `examples/events/knowledge/decisions/`): the decision changes product behavior — business rules, data, permissions, privacy, payments, scope. Use `framework/template/decision-template.md` and index it in the product's `.product/decisions.json`.

Recording a framework decision in a product decision log (or vice versa) is itself a defect. When a change produces both kinds, write both records and cross-link them.

## Creating an FDR

### 1. Allocate ID and file

Take the next sequential number from the index in `framework/decisions/README.md`. Filename: `FDR-NNN-<kebab-case-title>.md` (three-digit, e.g. `FDR-013-analytics-skill-ownership.md`).

### 2. Write the record

Follow the structure of the existing FDRs exactly:

```markdown
# FDR-NNN: <Title>

## Snapshot

| Field | Value |
| --- | --- |
| ID | `FDR-NNN` |
| Status | `proposed` |
| Origin EV | `EV-NNN` or `Governance baseline` |
| Date | `YYYY-MM-DD` |
| Owner | `<owning skill, e.g. Documentation Orchestrator>` |

## Context

<Why the current state is insufficient. Facts, not advocacy.>

## Decision

<What is decided, stated so an agent can act on it without reading the discussion. Use tables for boundaries and responsibilities.>

## Consequences

| Type | Consequence |
| --- | --- |
| Positive | ... |
| Negative | ... |
| Follow-up | ... |

## FRAMEWORK.md Amendments

| Section | Amendment |
| --- | --- |
| `<section>` | <what must change there> or `N/A` |

## Related Artifacts

| Artifact | Link |
| --- | --- |
| ... | [relative link](../../...) |

## Supersedes

- `FDR-NNN` or `N/A`
```

Notes:

- `Origin EV` references the evolution record that motivated the decision; use `Governance baseline` when none exists.
- New FDRs start as `proposed`. Only the user approves; set `approved` only after explicit approval, and do not fabricate approval records.
- If the FDR supersedes an earlier one, also mark the old record as superseded rather than deleting it.

### 3. Sync surfaces

- Add a row to the index table in `framework/decisions/README.md` (ID, linked title, status, origin).
- Apply the listed `FRAMEWORK.md` amendments in the same change, or record them as explicit follow-ups.
- If the decision changes a skill contract, update the affected `framework/skills/*/SKILL.md` (use the `sync-framework-assets` checklist).

## Verification

Run the `verify` skill — the test suite includes markdown link validation (FDR-010), so broken relative links in the new record will fail. Confirm the index row, the record's status, and the amendments are consistent with each other.
