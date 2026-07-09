# Audits

## Purpose

Audits report whether framework artifacts are complete, consistent, traceable, and safe to use as inputs for downstream work.

## When To Use

Use audits before approval gates, before implementation planning, before task generation, and before release. Audits should analyze existing artifacts and findings; they should not silently create product scope.

## Expected Files

- `README.md`: audit purpose and operating rules.
- `readiness/`: readiness reports for use cases or releases.
- Future folders for `gaps/`, `conflicts/`, `dependencies/`, `impact/`, `security/`, or `ux/` when reports become numerous.

## Responsible Skill

Primary owner: Audit Orchestrator.

Specialist owners: Gap Finder AI, Conflict AI, Dependency AI, Impact Analysis AI, QA AI, Security Review AI, UX Review AI.

## Report Shape

Each audit should include:

- Verdict: `approved`, `approved_with_notes`, or `blocked`.
- Findings with severity and evidence.
- Required fixes.
- Suggested improvements.
- Residual risk.
- Next recommended skill.

## Visual Reporting Standard

Reports should be easy to scan in GitHub and in Codex. Prefer:

- status icons such as `✅`, `🟡`, `🔴`, and `➖`;
- summary tables before long prose;
- Mermaid diagrams for flows, dependencies, and gates;
- finding matrices with severity, evidence, impact, required fix, and owner;
- decision tables when human approval is needed.

## Next Step

Run readiness checks on any use case before it feeds implementation planning or release preparation.
