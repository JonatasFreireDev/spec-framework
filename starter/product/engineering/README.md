# Product Engineering

This area contains the versioned Engineering System aggregate and specialist-owned technical landscape, standards, operations, quality, and evidence contracts.

## Start here

- Read `context.md`, `engineering-system.md`, and `engineering-system.yaml`.
- Route complete baseline creation and evolution through `engineering-orchestrator`.
- Use `technical-landscape`, `engineering-standards`, `operations-baseline`, and `engineering-evidence` for their respective subtrees; use `engineering-system` only for aggregate and Quality System consolidation.
- Base maturity claims on real code and operational evidence.
- Model systems, applications, components, repositories, interfaces, data stores, and deployments as stable graph entities under `catalog/`.
- Compose verifiable rules through versioned profiles under `standards/`; never weaken inherited required rules without a governed exception.
- Materialize optional entity, standard, environment, deployment, and runbook records only when evidence or explicit hypotheses require them.
- Keep product gates in `../knowledge/conventions/gates.md`.
- Inspect with `spec-framework engineering-system inspect` and validate with `spec-framework engineering-system validate`.

Framework validators and maintenance history do not belong in product engineering content.
