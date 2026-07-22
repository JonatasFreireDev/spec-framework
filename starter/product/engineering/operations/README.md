# Engineering Operations

Owner skill: `operations-baseline`.

Index environments, deployments, and runbooks in `operations.yaml`. Create
records on demand and link them to stable catalog entity IDs.

Each `ENV-*` and `DEPLOY-*` key maps to a relative YAML record; each
`RUNBOOK-*` key maps to a relative Markdown file. Embedded records are invalid.
Environment and deployment records require `schema_version: 1`, matching ID,
and non-empty status.
