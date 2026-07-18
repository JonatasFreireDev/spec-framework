# Template Authoring

## Purpose

Templates are resources owned by the skill that generates each artifact. They live in that skill's `assets/` directory, keeping generation instructions and reusable structure together. This reference defines their shared editorial contract without changing the architecture defined in `FRAMEWORK.md`.

## When To Use

Use the generating skill's own templates whenever a new canonical artifact is created or an existing artifact must be normalized. A template is not a substitute for product thinking; it is the checklist that keeps the output traceable, auditable, and ready for the next skill.

## Template Ownership

- `context-template.md`: baseline for every `context.md`.
- `problem-template.md`, `vision-template.md`, `strategy-template.md`: foundation artifacts.
- `interview-note-template.md`, `research-summary-template.md`: reusable problem-evidence notes stored in adopter-owned Foundation paths when needed.
- `feature-brief-template.md`: bounded `existing-feature` entry contract with one Target Feature.
- `product-baseline-template.md`: code-first current-state contract for `existing-product`.
- `implementation-assessment-template.md`: observed implementation evidence before full Foundation.
- `persona-template.md`, `metric-template.md`, `roadmap-item-template.md`: strategy support artifacts.
- `domain-template.md`, `goal-template.md`, `feature-template.md`, `use-case-template.md`: product hierarchy artifacts.
- `journey-template.md`: user journey artifact.
- `specification-template.md`, `design-template.md`, `design-system-template.md`, `design-component-template.md`, `design-pattern-template.md`, `engineering-system-template.md`, `engineering-proposal-template.md`, `engineering-review-template.md`, `implementation-plan-template.md`: planning, shared Design and Engineering contracts, and independent technical review.
- `technical-landscape/assets/` owns topology, technical catalog, entity, and boundary templates; `engineering-standards/assets/` owns standards, profiles, resolution, exception, and conformance templates; `operations-baseline/assets/` owns operations, environment, deployment, runbook, service-level, release, and rollback templates; `engineering-evidence/assets/` owns evidence, coverage, maturity, gap, and staleness templates.
- `engineering-system/assets/` owns only the aggregate Engineering System and shared product quality contracts: `engineering-system-template.*`, `quality-system-template.*`, `quality-model-template.md`, `test-strategy-template.md`, and `fitness-functions-template.yaml`.
- `execution-graph-template.json`, `tasks-template.md`, `tests-template.md`: executable planning artifacts.
- `qa-evidence-template.md`, `security-review-template.md`, `security-baseline-template.md`, `threat-register-template.md`: validation evidence, security gate, and proactive threat modeling artifacts.
- `analytics-template.md`, `audit-template.md`, `readiness-report-template.md`: validation and measurement artifacts.
- `decision-template.md`, `approval-record-template.json`, `derivation-record-template.json`, `release-template.md`: human approval, derivation, decision, and release artifacts.
- `domain-evolution-template.md`: opportunity mapping, candidate comparison, slicing, and explicit feature selection.
- `technical-discovery-template.md`: requirement-to-codebase mapping and Architecture Gate.
- `execution-graph-template.json`: proposed DAG contract; task paths become mandatory after atomic materialization.
- `task-template.md`: canonical task contract and readiness handoff.
- `specification-contract-template.md`: reusable structure for modular product, behavior, UX, API, data, security, quality, observability, and rollout contracts.
- `import-traceability-template.json`, `import-plan-template.json`, `import-mapping-template.json`, `import-report-template.md`: staged source-import evidence, review, selection, and reporting contracts.
- `command-plan-template.json`, `handoff-template.json`, `runtime-workspace-template.json`, `integration-template.json`: execution-runtime planning, handoff, workspace, and integration state contracts.

## Responsibility

The skill that generates an artifact owns its template and keeps it in `assets/`. The specialist named by the artifact's `owner_skill` or owning workflow remains accountable for its content and gate. Runtime JSON templates are owned by their command, workspace, handoff, or integration workflow.

Documentation Writer owns this editorial reference and the shared `context-template.md` and `derivation-record-template.json` resources. It reviews cross-template consistency; it does not become the owner of templates produced by other skills.

## Visual Standard

Templates should produce artifacts that are easy to scan in common Markdown renderers:

- use a `🧭 Snapshot` table near the top;
- use status icons such as `✅`, `🟡`, `🔴`, and `➖` where a report or gate has a result;
- use tables for scope, decisions, risks, dependencies, owners, and acceptance;
- use Mermaid diagrams for flows, artifact chains, gates, journeys, and dependencies;
- keep prose focused on decisions, evidence, and handoff.

## Generation And Self-Check Standard

Every Markdown artifact template includes a `Generation And Agent Self-Check` section near the top. When materializing an artifact, the owning agent must:

- record the generation date in `YYYY-MM-DD` format;
- state the artifact purpose, workflow trigger, responsible owner, and covered scope;
- link required inputs and evidence instead of naming unsupported claims;
- define the artifact-specific conditions that make the artifact ready;
- use the status vocabulary defined by the specific template;
- explain unresolved, blocked, or not-applicable items in the relevant artifact section;
- complete the checklist at the end of the artifact before handing it to the next skill.

Each checklist is tailored to the artifact contract. It must verify the artifact's specialized gates, evidence, approvals, lifecycle boundaries, and delivery criteria rather than repeat a generic documentation checklist.

Machine-readable JSON and YAML templates do not include editorial guidance, icons, comments, or checklists unless their validated schema defines those fields. Their field names, enum values, timestamps, and placeholders are the executable contract; accompanying Markdown templates and owning skills provide human guidance.

## Link Standard

Generated documents should use real Markdown links for local artifacts instead of plain paths when the target exists.

Use repository-relative paths in indexes and root-level reports:

```markdown
[FRAMEWORK.md](../../FRAMEWORK.md)
```

Use sibling-relative links inside artifact bundles:

```markdown
[Context](context.md)
[Specification](specification.md)
[Implementation Plan](implementation-plan.md)
[Execution Graph](execution-graph.json)
```

When a template contains placeholders, keep the placeholder in the label and use the expected target shape:

```markdown
[`[SPEC-XXX] Specification`](specification.md)
```

Do not link to a file that does not exist yet unless the link target is an explicit placeholder such as `[planned-output.md]`, `TBD`, or `N/A`. Real relative Markdown links in templates are validated by `framework/validators/framework-validator.mjs`.

## Mermaid Progress Classes

When a Mermaid flow represents framework progress, include visual classes:

```mermaid
flowchart LR
  %% artifact: SPEC-XXX node: H %%
  A["Problem"] --> B["Vision"] --> C["Strategy"] --> D["Domain"] --> E["User Goal"] --> F["Feature"] --> G["Use Case"] --> H["Specification"] --> I["Design"] --> J["Implementation Plan"] --> K["Execution Graph"] --> L["Tasks"]

  classDef done fill:#dcfce7,stroke:#16a34a,color:#14532d;
  classDef current fill:#fef3c7,stroke:#d97706,color:#78350f,stroke-width:3px;
  classDef pending fill:#f8fafc,stroke:#94a3b8,color:#334155;
  classDef blocked fill:#fee2e2,stroke:#dc2626,color:#7f1d1d;

  class A,B,C,D,E,F,G done;
  class H current;
  class I,J,K,L pending;
```

Responsibility:

| Owner | Responsibility |
| --- | --- |
| Skill that owns the current artifact | Update the local Mermaid flow and `context.md` when artifact status changes. |
| Documentation Orchestrator | Synchronize Mermaid progress across reports, templates, indexes, and context files. |
| Audit Orchestrator | Verify visual state matches real artifact status during audits. |
| Release Orchestrator | Verify release/readiness visual flows before release approval. |

Semantic validation:

Use a Mermaid comment to bind a node to a framework artifact when the node represents a real artifact:

```text
%% artifact: UC-002 node: U %%
```

Placeholder IDs such as `SPEC-XXX` are examples. Generated documents should replace them with real artifact IDs before relying on semantic validation.

The validator maps artifact status from `.product/artifacts.json` to visual state:

| Artifact status | Expected Mermaid class |
| --- | --- |
| `approved`, `implemented`, `validated`, `released` | `done` |
| `draft`, `proposed`, `in_progress` | `current` |
| `deprecated`, `superseded` | `blocked` |
| `unknown` or missing | `pending` |

## Next Step

When creating an artifact, read the relevant parent `context.md`, copy the matching template structure into the target document, replace placeholders with concrete content, and leave the artifact in `draft` or `proposed` until human approval is recorded.

## Reference links

Every cross-artifact reference in a materialized document must be a Markdown link to the canonical file and, when applicable, its section anchor. This includes tasks, bugs, decisions, requirements, acceptance criteria, QA/security evidence, audits, releases, gaps, mappings, and any other document registered by the product. An ID by itself (for example `GAP-003`, `MAP-010`, `TASK-021`, or `DEC-014`) is not sufficient for navigation. During drafting, use the explicit placeholder form `[GAP-003](<path-to-gap-003.md>#gap-003)` and replace the path with the real product-relative path before handoff. Keep the ID in the link label so search, registry checks, and human navigation all retain the stable identity.
