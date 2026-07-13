# FDR-036: Registry Relationship Preservation

## Snapshot

| Field | Value |
| --- | --- |
| ID | `FDR-036` |
| Status | `proposed` |
| Origin EV | `Starting-point gate review` |
| Date | `2026-07-12` |
| Owner | `Framework Validator` |

## Context

Registry regeneration reads scalar identity and status from canonical Markdown but previously ignored structured `parents`, `children`, `depends_on`, decisions, and delivery dependencies in companion `context.md` YAML. Running `validate --write-registry` could therefore remove Domain, User Goal, Feature, and Use Case relationships that approval gates depend on.

## Decision

| Boundary | Contract |
| --- | --- |
| Structured parsing | Registry generation parses the fenced YAML document in the canonical companion `context.md` with the repository YAML parser. |
| Relationships | `parents`, `children`, `depends_on`, `decisions`, and `delivery.depends_on` survive registry regeneration. |
| Canonical override | An explicit `Parent IDs` table field on the canonical artifact takes precedence over companion parents. |
| Starting-point additions | Feature Brief, Implementation Assessment, and Product Baseline parents are added after preserving existing relationships. |
| Failure safety | Missing or invalid companion YAML does not invent relationships; normal context validation reports malformed contracts. |

## Consequences

| Type | Consequence |
| --- | --- |
| Positive | Rebuilding the registry no longer weakens parent approval gates. |
| Positive | Generated registry state remains consistent with canonical context contracts. |
| Negative | Registry generation now depends on structured YAML parsing rather than scalar regex extraction alone. |
| Follow-up | Add parity fixtures for every canonical artifact type with companion relationships. |

## FRAMEWORK.md Amendments

| Section | Amendment |
| --- | --- |
| `13. Recommended Governance Rules` | Registry regeneration preserves structured relationships from companion contexts. |

## Related Artifacts

| Artifact | Link |
| --- | --- |
| Foundation registry | [FDR-030](FDR-030-foundation-approval-registry.md) |
| Feature-scoped bootstrap | [FDR-031](FDR-031-feature-scoped-bootstrap.md) |
| Canonical method | [../../FRAMEWORK.md](../../FRAMEWORK.md) |

## Supersedes

- N/A.
