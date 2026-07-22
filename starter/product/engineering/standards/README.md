# Engineering Standards

Owner skill: `engineering-standards`.

`standards.yaml` indexes versioned profiles, standards, and governed
exceptions. Standards apply by entity type, capability, or explicit assignment.
Required inherited standards cannot be silently weakened.

Each `PROFILE-*`, `STD-*`, or `STDEX-*` key maps to a relative YAML record.
Embedded records are invalid. Referenced profiles and standards require a
matching ID, semantic version, and non-empty status.
