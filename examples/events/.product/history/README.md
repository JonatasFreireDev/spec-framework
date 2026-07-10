# Product History

This folder stores product history and approval records.

Approval records are one JSON file per approval. The canonical naming pattern is:

```text
approval-<artifact-id>-<status-granted>-<hash8>.json
```

Each approval record must include:

- `artifact_id`
- `path`
- `content_hash`
- `status_granted`
- `approved_by`
- `approved_at`
- `notes`

Agents must not create, edit, or repair approval records unless a human explicitly approves a migration that names approval-record generation as a deliverable. Missing or inconsistent records are blockers to report, not something to silently fix.
