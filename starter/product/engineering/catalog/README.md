# Technical Entity Catalog

Owner skill: `technical-landscape`.

`catalog.yaml` indexes systems, applications, components, repositories, data
stores, interfaces, and deployments by stable ID. Create records only from
observed evidence or explicit hypotheses. Folder position does not define graph
relations.

Each entity entry is an index from its stable ID to a relative YAML file; entity
fields must not be embedded in `catalog.yaml`. For example,
`SYS-PRODUCT-001: systems/product.yaml`. The referenced file must declare at
least `schema_version`, the same `id`, the category-compatible `type`, and a
non-empty `status`.

Each relation uses a stable `REL-*` ID, an extensible type, and existing source
and target entity IDs. The graph contract does not infer relations from paths.
