# Engineering Catalog And Standards

## Purpose

This contract defines how `engineering/` grows from a product baseline into a
multi-system, multi-application, and multi-repository technical knowledge graph.
It also records the implementation plan for adopting that model across the
framework.

The Engineering System must not assume that one product equals one repository
or one deployable application. Repositories, systems, applications, components,
interfaces, data stores, and deployments are independent entities connected by
stable identifiers.

Ownership is deliberately split. `technical-landscape` owns the graph,
`engineering-standards` owns rules and profiles, `operations-baseline` owns
operational contracts, `engineering-evidence` owns evidence and maturity, and
`engineering-system` owns only the versioned aggregate and Quality System.
`engineering-orchestrator` sequences these owners and stops at human approval.
The aggregate `engineering-system.yaml` declares `owner_skill` for every area,
and specialist root catalogs repeat their owner. Missing owner metadata remains
valid for legacy adopters; once declared, mismatched ownership is a blocker.

## Canonical Structure

```text
engineering/
├── context.md
├── engineering-system.md
├── engineering-system.yaml
├── architecture/
│   ├── system-context.md
│   ├── modules.md
│   └── topology.yaml
├── catalog/
│   ├── catalog.yaml
│   ├── systems/
│   ├── applications/
│   ├── components/
│   ├── repositories/
│   ├── data-stores/
│   ├── interfaces/
│   └── deployments/
├── standards/
│   ├── standards.yaml
│   ├── profiles/
│   ├── catalog/
│   └── exceptions/
├── quality/
├── operations/
│   ├── operations.yaml
│   ├── environments/
│   ├── deployments/
│   └── runbooks/
├── evidence/
│   └── inventory.md
└── decisions/
```

Only the root catalogs and README contracts are initialized. Entity records,
profiles, standards, exceptions, environments, and runbooks are materialized
on demand from evidence or explicit hypotheses; empty placeholder forests are
not required.

## Entity Graph

| Entity | Responsibility |
| --- | --- |
| System | Technical capability composed from one or more applications or components. |
| Application | Executable or user-facing unit delivered as a whole. |
| Component | Service, module, package, worker, frontend, or reusable library. |
| Repository | Physical source-control boundary containing one or more components. |
| Interface | API, event, queue, file, webhook, or other cross-component contract. |
| Data store | Database, cache, bucket, index, or other persistent source. |
| Deployment | Mapping from components and artifacts to an environment and release strategy. |

Relations live in mechanical catalogs and use stable IDs. Folder position is
for navigation and never determines ownership or containment by itself. This
supports monorepos, polyrepos, shared components, multiple deployables per
repository, and applications assembled from several repositories.

`catalog/catalog.yaml` is an index, not an entity store. Every category maps a
stable entity ID to one relative YAML file:

```yaml
entities:
  systems:
    SYS-PRODUCT-001: systems/product.yaml
```

Embedded entity mappings are incompatible with this contract. Each referenced
file must exist below `engineering/`, parse as YAML, and declare at least
`schema_version: 1`, an `id` equal to the catalog key, the category-compatible
`type`, and a non-empty `status`.

`architecture/topology.yaml` uses the same seven entity categories as the
catalog. Every catalog entity must appear in the topology, every topology ID
must be indexed by the catalog, and every topology relation endpoint must exist.

Each relation declares a stable `REL-*` ID, an extensible relation `type`, and
existing `source` and `target` entity IDs. Evidence may be attached when the
relation is observed rather than hypothetical. The framework validates graph
integrity without imposing a closed relation vocabulary.

## Contract Inheritance

```text
Engineering System
→ System
→ Application
→ Component
→ delivery-specific Engineering Proposal
```

A more specific contract may add constraints. It may not silently weaken or
replace an inherited contract. A divergence requires a governed exception,
deviation, or approved `DEC-*` record.

## Standards System

`engineering/standards/` is a versioned catalog of verifiable technical
rules. Standards are independent entities grouped into profiles and selected
by entity type, capability, or explicit assignment.

Canonical categories are architecture, code, API, events, data, dependencies,
security, observability, testing, and delivery. A standard declares:

- stable ID, semantic version, status, category, and obligation level;
- applicable entity types and capabilities;
- individually identifiable rules;
- verification methods and required evidence;
- exception policy and compatibility notes.

Obligation levels are `required`, `recommended`, `experimental`, and
`deprecated`. Required standards block a governed gate unless conformity or
an open, approved, unexpired, in-scope `STDEX-*` exception is recorded.

Profiles compose standards for common shapes such as web applications, HTTP
APIs, workers, event consumers, shared libraries, and product defaults.
Consumers pin profile and standard versions. Profiles may extend other
profiles, but cycles are invalid.

`standards/standards.yaml` is also strictly indexed: `PROFILE-*`, `STD-*`, and
`STDEX-*` keys map to relative YAML records and never contain embedded fields.
Profiles and standards require `schema_version: 1`, matching identity, semantic
version, and non-empty status. Standards additionally require a supported
category and obligation level plus at least one verifiable rule.

`operations/operations.yaml` maps `ENV-*` and `DEPLOY-*` IDs to relative YAML
records and `RUNBOOK-*` IDs to relative Markdown files. Embedded records are
invalid. Environment and deployment YAML must declare `schema_version: 1`, a
matching ID, and non-empty status; deployments also reference an indexed
environment, a `DEPLOY-*` technical entity, valid application/component IDs,
and indexed runbooks.

Standards define technical rules. The Quality System defines quality policy and
coverage. The security baseline defines threats, trust, and governed controls.
Standards reference these contracts rather than duplicating them.

## Evidence And Maturity

Every maturity above `baseline` requires resolvable evidence. Evidence may
reference repository paths, tests, CI, commands, runtime observations, or
approved external sources. `engineering/evidence/inventory.md` indexes the
evidence used by catalogs without copying volatile output into contracts.

The Engineering System composite approval hash covers every file under
`engineering/`. Adding or changing a catalog, profile, standard, exception,
entity, operation, or evidence contract makes existing approval stale.

## Compatibility And Migration

The model is additive. Existing adopters with only `engineering-system.md`,
`engineering-system.yaml`, architecture, and quality remain valid. Upgrade
must not create or overwrite adopter-owned entity records. New starter
repositories receive empty root catalogs; existing products opt into them by a
previewable migration or normal skill-driven evolution followed by approval.

## Skill Flow

```text
engineering-orchestrator
→ technical-landscape
→ engineering-standards
→ operations-baseline
→ engineering-evidence
→ engineering-system aggregation and validation
→ engineering-orchestrator readiness review
→ human approval
→ domain-architect or technical-discovery
```

The orchestrator may revisit only affected specialists during evolution, but it
must revalidate downstream contracts and the final composite hash before asking
for approval.
