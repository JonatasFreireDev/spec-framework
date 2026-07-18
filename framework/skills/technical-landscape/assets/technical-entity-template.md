# [Entity type]: [name]

## Snapshot

| Field | Value |
| --- | --- |
| ID | `[SYS|APP|CMP|REPO|DATA|IFACE|DEPLOY]-*` |
| Type | `system | application | component | repository | data-store | interface | deployment` |
| Status | `draft` |
| Mechanical record | [entity.yaml](entity.yaml) |

## Responsibility

[Stable responsibility and explicit boundaries.]

## Relations

| Relation | Target ID | Evidence |
| --- | --- | --- |
| `[contains | depends_on | exposes | stores_in | deployed_as]` | `[stable ID]` | `[path]` |

## Ownership And Evidence

| Concern | Owner | Evidence |
| --- | --- | --- |
| Technical ownership | `[owner or Unassigned]` | `[path]` |

## Standards

| Profile or standard | Version | Exception |
| --- | --- | --- |
| `[PROFILE-* or STD-*]` | `[version]` | `[None or STDEX-*]` |
