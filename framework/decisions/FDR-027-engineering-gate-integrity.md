# FDR-027: Engineering Gate Integrity

## Snapshot

| Field | Value |
| --- | --- |
| ID | `FDR-027` |
| Status | `approved` |
| Origin EV | `Governance baseline` |
| Date | `2026-07-12` |
| Owner | `Framework Maintainer` |

## Context

FDR-026 introduced the Engineering System, Engineering Proposal, and Engineering Review. A subsequent framework gap audit found that operational routing still trusted Markdown status without checking current approval records, Engineering Review did not prove that its proposal hash matched current content, Engineering Proposal pins were not validated against the declared system, the mechanical YAML contract used line-oriented parsing, and dashboard output omitted Engineering System blockers.

## Decision

| Boundary | Contract |
| --- | --- |
| Approval integrity | Operational workflow gates require both an approved-or-later artifact status and a current hash-matching approval record. Status prose alone never advances work. |
| Review freshness | A passed Engineering Review must record the SHA-256 hash of the current Engineering Proposal. Any proposal change makes the review stale and blocks planning. |
| Consumer pin | Applicable Engineering Proposals pin the declared Engineering System id/version or explicitly declare `Not configured` only when no configured system exists. A configured system consumed by proposed-or-later work must be approved with current evidence. |
| Mechanical parsing | Engineering System catalogs and use-case contexts are parsed as YAML through a structured parser. Valid YAML collection syntax must not change trigger or evidence semantics. |
| Catalog coherence | Human context and mechanical catalog agree on id, status, version, and origin mode. Mismatches block validation. |
| Dashboard | Engineering System validation blockers are visible in dashboard output and block a misleading ready presentation. |
| Migration | Draft artifacts remain editable without approval evidence. Existing approved-or-later artifacts that lack current evidence become blocked and require human re-approval; no approval record is generated automatically. |

## Consequences

| Type | Consequence |
| --- | --- |
| Positive | Workflow navigation can no longer be advanced by manually editing status text. |
| Positive | Engineering Review and Engineering System consumption become content-addressed and stale-safe. |
| Positive | YAML syntax is interpreted consistently instead of through indentation heuristics. |
| Negative | Existing products with manually promoted statuses or stale reviews may become blocked. |
| Negative | The released CLI gains a small YAML runtime dependency. |
| Follow-up | Apply the same explicit consumer-staleness pattern to other shared product systems when introduced. |

## FRAMEWORK.md Amendments

| Section | Amendment |
| --- | --- |
| `7. Implementation Plan` | Require current proposal hash and current approval evidence before planning. |
| `11. Approval Gates` | Clarify that operational navigation uses current approval records, not status alone. |
| `15. How To Use With Codex` | Surface Engineering System blockers through dashboard and validation. |

## Related Artifacts

| Artifact | Link |
| --- | --- |
| Canonical method | [../../FRAMEWORK.md](../../FRAMEWORK.md) |
| Engineering System decision | [FDR-026](FDR-026-canonical-product-engineering-system.md) |
| Engineering Review skill | [../skills/engineering-review/SKILL.md](../skills/engineering-review/SKILL.md) |

## Supersedes

- N/A; this decision amends FDR-026.
